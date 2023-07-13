package peerstore

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	libp2p_crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	libp2p "github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/record"
	"github.com/libp2p/go-libp2p/p2p/host/peerstore/pstoreds"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeerstoreKV(t *testing.T) {
	ps, err := NewPeerStore(context.Background(), "kv", "")
	defer ps.Close()
	require.NoError(t, err)

	// add kv pair
	ps.Put("0", "foo", "bar")
	v, err := ps.Get("0", "foo")
	require.NoError(t, err)
	require.Equal(t, "bar", v)

	// add peers
	// ephemeral key for peer ID
	privKey, _, err := libp2p_crypto.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)
	peerID, err := peer.IDFromPrivateKey(privKey)
	require.NoError(t, err)
	privKey2, _, err := libp2p_crypto.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)
	peerID2, err := peer.IDFromPrivateKey(privKey2)
	require.NoError(t, err)

	ps.AddAddr(peerID, ma.StringCast("/ip4/1.2.3.4/tcp/1234"), time.Hour)
	ps.AddAddrs(peerID2, []ma.Multiaddr{
		ma.StringCast("/ip4/1.2.3.4/tcp/1234"),
		ma.StringCast("/ip4/1.2.3.4/tcp/1111"),
		ma.StringCast("/ip4/1.2.3.4/tcp/3456"),
	}, time.Hour)
	peers := ps.PeersWithAddrs()
	info := ps.PeerInfo(peerID)
	require.Equal(t, 2, len(peers))
	require.Equal(t, peerID, info.ID)
	require.Equal(t, 1, len(info.Addrs))
	require.Equal(t, "/ip4/1.2.3.4/tcp/1234", info.Addrs[0].String())

	// remove peer only removes the keys, not the addresses
	ps.RemovePeer(peerID)
	peers = ps.PeersWithAddrs()
	require.Equal(t, 2, len(peers))
	// remove address
	ps.ClearAddrs(peerID)
	peers = ps.PeersWithAddrs()
	require.Equal(t, 1, len(peers))

	// built-in metrics. how do read this data?
	ps.RecordLatency(peerID, 5*time.Second)

}

func TestPeerstoreSQLite(t *testing.T) {
	dir := t.TempDir()
	ps, err := NewPeerStore(context.Background(), "sqlite", dir)
	defer ps.Close()
	require.NoError(t, err)

	// add kv pair
	ps.Put("0", "foo", "bar")
	v, err := ps.Get("0", "foo")
	require.NoError(t, err)
	require.Equal(t, "bar", v)

	// add peers
	// ephemeral key for peer ID
	privKey, _, err := libp2p_crypto.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)
	peerID, err := peer.IDFromPrivateKey(privKey)
	require.NoError(t, err)
	privKey2, _, err := libp2p_crypto.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)
	peerID2, err := peer.IDFromPrivateKey(privKey2)
	require.NoError(t, err)

	// address no longer valid after an hour
	ps.AddAddr(peerID, ma.StringCast("/ip4/1.2.3.4/tcp/1234"), time.Hour)
	ps.AddAddrs(peerID2, []ma.Multiaddr{
		ma.StringCast("/ip4/1.2.3.4/tcp/1234"),
		ma.StringCast("/ip4/1.2.3.4/tcp/1111"),
		ma.StringCast("/ip4/1.2.3.4/tcp/3456"),
	}, time.Hour)
	peers := ps.PeersWithAddrs()
	info := ps.PeerInfo(peerID)
	require.Equal(t, 2, len(peers))
	require.Equal(t, peerID, info.ID)
	require.Equal(t, 1, len(info.Addrs))
	require.Equal(t, "/ip4/1.2.3.4/tcp/1234", info.Addrs[0].String())

	// remove peer only removes the keys, not the addresses
	ps.RemovePeer(peerID)
	peers = ps.PeersWithAddrs()
	require.Equal(t, 2, len(peers))
	// remove address
	ps.ClearAddrs(peerID)
	peers = ps.PeersWithAddrs()
	require.Equal(t, 1, len(peers))
}

func TestCertifiedAddrBook(t *testing.T) {
	records := make(map[peer.ID]*record.Envelope)
	cb := CertAddrBook{Records: records}
	dir := t.TempDir()
	db := dbStore("kv", dir)
	addrBook, err := pstoreds.NewAddrBook(context.Background(), db, pstoreds.DefaultOpts())
	require.NoError(t, err)
	ab := AddrBook{
		CertifiedAddrBook: cb,
		AddrBook:          addrBook,
	}
	// check ab is type libp2p.AddrBook
	var _ libp2p.AddrBook = (*AddrBook)(nil)
	// check ab is type libp2p.CertifiedAddrBook
	_, ok := libp2p.GetCertifiedAddrBook(&ab)
	assert.True(t, ok)

	// create a signed record
	privKey, _, err := libp2p_crypto.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)
	peerID, err := peer.IDFromPrivateKey(privKey)
	require.NoError(t, err)

	addr := ma.StringCast("/ip4/1.2.3.4/tcp/1234")
	rec := peer.NewPeerRecord()
	rec.PeerID = peerID
	rec.Addrs = []ma.Multiaddr{addr}
	signed, err := record.Seal(rec, privKey)
	if err != nil {
		t.Fatalf("error generating peer record: %s", err)
	}

	accepted, err := ab.CertifiedAddrBook.ConsumePeerRecord(signed, time.Second)
	require.True(t, accepted)
	require.NoError(t, err)

	// get sealed record
	env := ab.CertifiedAddrBook.GetPeerRecord(peerID)
	envrec, _ := env.Record()
	peerRec := envrec.(*peer.PeerRecord)
	require.NotNil(t, peerRec)
	require.Equal(t, peerID, peerRec.PeerID)
}

func TestAlgoPeerStore(t *testing.T) {
	dir := t.TempDir()
	ps, err := NewAlgoPeerstore(context.Background(), "kv", dir)
	defer ps.Close()
	require.NoError(t, err)

	// add kv pair
	ps.Put("0", "foo", "bar")
	v, err := ps.Get("0", "foo")
	require.NoError(t, err)
	require.Equal(t, "bar", v)

	// peer ID
	privKey, _, err := libp2p_crypto.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)
	peerID, err := peer.IDFromPrivateKey(privKey)
	require.NoError(t, err)

	ps.RecordCount(peerID)
	cnt := ps.RecordCount(peerID)
	require.Equal(t, 2, cnt)

	addr := ma.StringCast("/ip4/1.2.3.4/tcp/1234")
	rec := peer.NewPeerRecord()
	rec.PeerID = peerID
	rec.Addrs = []ma.Multiaddr{addr}
	signed, err := record.Seal(rec, privKey)
	if err != nil {
		t.Fatalf("error generating peer record: %s", err)
	}

	// check ab is type libp2p.CertifiedAddrBook
	_, ok := libp2p.GetCertifiedAddrBook(&ps.AddrBook)
	assert.True(t, ok)

	accepted, err := ps.CertifiedAddrBook.ConsumePeerRecord(signed, time.Second)
	require.True(t, accepted)
	require.NoError(t, err)

	// get sealed record
	env := ps.CertifiedAddrBook.GetPeerRecord(peerID)
	envrec, _ := env.Record()
	peerRec := envrec.(*peer.PeerRecord)
	require.NotNil(t, peerRec)
	require.Equal(t, peerID, peerRec.PeerID)

}
