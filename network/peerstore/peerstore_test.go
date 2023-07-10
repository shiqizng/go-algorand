package peerstore

import (
	"context"
	"testing"

	"github.com/libp2p/go-libp2p/p2p/host/peerstore/pstoreds"
	"github.com/stretchr/testify/require"
)

func TestPeerstore(t *testing.T) {
	ds, dsclose := leveldbStore()
	defer dsclose()
	ps, err := pstoreds.NewPeerstore(context.Background(), ds, pstoreds.DefaultOpts())
	defer ps.Close()
	require.NoError(t, err)
	ps.Put("0", "foo", "bar")
	v, err := ps.Get("0", "foo")
	require.NoError(t, err)
	require.Equal(t, "bar", v)
}
