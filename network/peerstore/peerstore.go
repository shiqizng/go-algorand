// Copyright (C) 2019-2023 Algorand, Inc.
// This file is part of go-algorand
//
// go-algorand is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// go-algorand is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with go-algorand.  If not, see <https://www.gnu.org/licenses/>.

package peerstore

import (
	"context"
	"log"

	"github.com/algorand/go-algorand/network"
	ds "github.com/ipfs/go-datastore"
	leveldb "github.com/ipfs/go-ds-leveldb"
	libp2p "github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/host/peerstore/pstoreds"
)

func leveldbStore() (ds.Batching, func()) {
	store, err := leveldb.NewDatastore("", nil)
	if err != nil {
		log.Fatal(err)
	}
	closer := func() {
		store.Close()
	}
	return store, closer
}

type Peerstore struct {
	ps libp2p.Peerstore
	ds ds.Batching
}

func NewPeerStore() (Peerstore, error) {
	datastore, _ := leveldbStore()
	peerstore, _ := pstoreds.NewPeerstore(context.Background(), datastore, pstoreds.DefaultOpts())
	return Peerstore{
		ps: peerstore,
		ds: datastore,
	}, nil
}

func (ps Peerstore) AddPeer(peer *network.Peer) (string, error) {
	return "", nil
}

func (ps Peerstore) RemovePeer(peer *network.Peer) (string, error) {
	return "", nil
}

func (ps Peerstore) GetPeer(peer *network.Peer) (string, error) {
	return "", nil
}

func (ps Peerstore) Close() {
	ps.ds.Close()
	ps.ps.Close()
}
