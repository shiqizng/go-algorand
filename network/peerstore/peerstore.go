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
	"path/filepath"
	"time"

	ds "github.com/ipfs/go-datastore"
	leveldb "github.com/ipfs/go-ds-leveldb"
	sqliteds "github.com/ipfs/go-ds-sql/sqlite"
	"github.com/libp2p/go-libp2p/core/peer"
	libp2p "github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/record"
	"github.com/libp2p/go-libp2p/p2p/host/peerstore/pstoreds"
	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

func dbStore(t, path string) ds.Batching {
	var store ds.Batching
	var err error
	switch t {
	case "kv":
		// empty path creates an in-memory datastore. set path to a directory to persist
		store, err = leveldb.NewDatastore(path, nil)
		if err != nil {
			log.Fatal(err)
		}
	case "sqlite":
		opts := &sqliteds.Options{
			DSN: filepath.Join(path, "peerstore.sqlite"),
		}
		store, err = opts.Create()
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("unknown datastore type")
	}
	return store
}

// NewPeerStore creates a new peerstore backed by a datastore.
func NewPeerStore(ctx context.Context, storeType string, path string) (libp2p.Peerstore, error) {
	datastore := dbStore(storeType, path)
	return pstoreds.NewPeerstore(ctx, datastore, pstoreds.DefaultOpts())
}

// AddrBook implements libp2p.AddrBook.; pstoreds.addr_book also implements a certifiedAddrBook
type AddrBook struct {
	libp2p.AddrBook
	libp2p.CertifiedAddrBook
}

// CertAddrBook implements libp2p.CertifiedAddrBook.
type CertAddrBook struct {
	Records map[peer.ID]*record.Envelope
}

// ConsumePeerRecord implements libp2p.CertifiedAddrBook.
func (cab CertAddrBook) ConsumePeerRecord(s *record.Envelope, ttl time.Duration) (accepted bool, err error) {
	rec, _ := s.Record()
	prec := rec.(*peer.PeerRecord)
	cab.Records[prec.PeerID] = s
	return true, nil
}

// GetPeerRecord implements libp2p.CertifiedAddrBook.
func (cab CertAddrBook) GetPeerRecord(p peer.ID) *record.Envelope {
	return cab.Records[p]
	return nil
}

type Metrics struct {
	libp2p.Metrics
	PeersCounter map[peer.ID]int
}

func (m Metrics) RecordCount(p peer.ID) int {
	m.PeersCounter[p]++
	return m.PeersCounter[p]
}

type AlgoPeerStore struct {
	libp2p.Peerstore
	AddrBook
	Metrics
}

func NewAlgoPeerstore(ctx context.Context, storeType string, path string) (AlgoPeerStore, error) {
	datastore := dbStore(storeType, path)
	store, err := pstoreds.NewPeerstore(ctx, datastore, pstoreds.DefaultOpts())
	m := Metrics{store.Metrics, make(map[peer.ID]int)}
	records := make(map[peer.ID]*record.Envelope)
	cb := CertAddrBook{Records: records}
	db := dbStore(storeType, path)
	addrBook, err := pstoreds.NewAddrBook(context.Background(), db, pstoreds.DefaultOpts())
	ab := AddrBook{
		CertifiedAddrBook: cb,
		AddrBook:          addrBook,
	}
	aps := AlgoPeerStore{store, ab, m}
	return aps, err
}
