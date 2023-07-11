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

	ds "github.com/ipfs/go-datastore"
	leveldb "github.com/ipfs/go-ds-leveldb"
	sqliteds "github.com/ipfs/go-ds-sql/sqlite"
	libp2p "github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/host/peerstore/pstoreds"
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

func NewPeerStore(context context.Context, storeType string, path string) (libp2p.Peerstore, error) {
	datastore := dbStore(storeType, path)
	return pstoreds.NewPeerstore(context, datastore, pstoreds.DefaultOpts())
}
