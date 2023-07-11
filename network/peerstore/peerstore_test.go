package peerstore

import (
	"crypto/rand"
	"testing"
	"time"

	libp2p_crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
)

func TestPeerstoreKV(t *testing.T) {
	ps, err := NewPeerStore("kv", "")
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
	ps, err := NewPeerStore("sqlite", dir)
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
