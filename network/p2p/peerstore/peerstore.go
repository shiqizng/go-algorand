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
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	libp2p "github.com/libp2p/go-libp2p/core/peerstore"
	mempstore "github.com/libp2p/go-libp2p/p2p/host/peerstore/pstoremem"
	"github.com/multiformats/go-multiaddr"
	"golang.org/x/exp/slices"

	"github.com/algorand/go-algorand/network"
)

// PeerStore implements Peerstore and CertifiedAddrBook.
type PeerStore struct {
	peerStoreCAB
}

// addressData: holds the information associated with each phonebook address.
type addressData struct {
	// retryAfter is the time to wait before retrying to connect to the address.
	retryAfter time.Time

	// recentConnectionTimes: is the log of connection times used to observe the maximum
	//                        connections to the address in a given time window.
	recentConnectionTimes []time.Time

	// networkNames: lists the networks to which the given address belongs.
	networkNames map[string]bool

	// role is the role that this address serves.
	role network.PhoneBookEntryRoles

	// persistent is set true for peers whose record should not be removed for the peer list
	persistent bool
}

// peerStoreCAB combines the libp2p Peerstore and CertifiedAddrBook interfaces.
type peerStoreCAB interface {
	libp2p.Peerstore
	libp2p.CertifiedAddrBook
}

// NewPeerStore creates a new peerstore backed by a datastore.
func NewPeerStore(addrInfo []*peer.AddrInfo) (*PeerStore, error) {
	ps, err := mempstore.NewPeerstore()
	if err != nil {
		return nil, fmt.Errorf("cannot initialize a peerstore: %w", err)
	}

	// initialize peerstore with addresses
	for i := 0; i < len(addrInfo); i++ {
		info := addrInfo[i]
		ps.AddAddrs(info.ID, info.Addrs, libp2p.AddressTTL)
	}
	pstore := &PeerStore{peerStoreCAB: ps}
	return pstore, nil
}

// GetAddresses returns up to N addresses, but may return fewer
func (ps *PeerStore) GetAddresses(n int, role network.PhoneBookEntryRoles) []string {
	return shuffleSelect(ps.filterRetryTime(time.Now(), role), n)
}

// UpdateRetryAfter updates the retryAfter time for the given address.
func (ps *PeerStore) UpdateRetryAfter(addr string, retryAfter time.Time) {
	info, err := peer.AddrInfoFromString(addr)
	if err != nil {
		return
	}
	metadata, err := ps.Get(info.ID, "addressData")
	if err != nil && errors.Is(err, libp2p.ErrNotFound) {
		return
	} else if err != nil {
		panic(err)
	}
	ad := metadata.(addressData)
	ad.retryAfter = retryAfter
	_ = ps.Put(info.ID, "addressData", ad)

}

// GetConnectionWaitTime returns the connection wait time for the given address.
func (ps *PeerStore) GetConnectionWaitTime(addr string) (bool, time.Duration, time.Time) {
	curTime := time.Now()
	info, err := peer.AddrInfoFromString(addr)
	if err != nil {
		return false, 0 /* not used */, curTime /* not used */
	}
	if info == nil {
		// The addr is not in this phonebook.
		// Will find the addr in a different phonebook.
		return false, 0 /* not used */, curTime /* not used */
	}
	var timeSince time.Duration
	var numElmtsToRemove int
	metadata, err := ps.Get(info.ID, "addressData")
	if err != nil {
		return false, 0 /* not used */, curTime /* not used */
	}
	ad := metadata.(addressData)
	// Remove from recentConnectionTimes the times later than ConnectionsRateLimitingWindowSeconds
	connectionsRateLimitingWindow, err := ps.Get("peerID", "connectionsRateLimitingWindow")
	if err != nil {
		return false, 0 /* not used */, curTime /* not used */
	}
	for numElmtsToRemove < len(ad.recentConnectionTimes) {
		timeSince = curTime.Sub(ad.recentConnectionTimes[numElmtsToRemove])
		if timeSince >= connectionsRateLimitingWindow.(time.Duration) {
			numElmtsToRemove++
		} else {
			break // break the loop. The rest are earlier than 1 second
		}
	}

	// Remove the expired elements from e.data[addr].recentConnectionTimes
	ps.popNElements(numElmtsToRemove, peer.ID(addr))
	// If there are max number of connections within the time window, wait
	connectionsRateLimitingCount, err := ps.Get("peerID", "connectionsRateLimitingCount")
	if err != nil {
		return false, 0 /* not used */, curTime /* not used */
	}
	numElts := len(ad.recentConnectionTimes)
	if uint(numElts) >= connectionsRateLimitingCount.(uint) {
		return true, /* true */
			connectionsRateLimitingWindow.(time.Duration) - timeSince, curTime /* not used */
	}

	// Else, there is space in connectionsRateLimitingCount. The
	// connection request of the caller will proceed
	// Update curTime, since it may have significantly changed if waited
	provisionalTime := time.Now()
	// Append the provisional time for the next connection request
	ps.appendTime(info.ID, provisionalTime)
	return true, 0 /* no wait. proceed */, provisionalTime
}

// UpdateConnectionTime updates the connection time for the given address.
func (ps *PeerStore) UpdateConnectionTime(addr string, provisionalTime time.Time) bool {
	info, err := peer.AddrInfoFromString(addr)
	if err != nil {
		return false
	}
	metadata, err := ps.Get(info.ID, "addressData")
	if err != nil {
		return false
	}
	ad := metadata.(addressData)
	entry := ad.recentConnectionTimes
	defer func() {
		_ = ps.Put(info.ID, "addressData", entry)

	}()

	// Find the provisionalTime and update it
	for indx, val := range entry {
		if provisionalTime == val {
			entry[indx] = time.Now()
			return true
		}
	}

	// Case where the time is not found: it was removed from the list.
	// This may happen when the time expires before the connection was established with the server.
	// The time should be added again.
	entry = append(entry, time.Now())

	return true
}

// ReplacePeerList replaces the peer list for the given networkName and role.
func (ps *PeerStore) ReplacePeerList(addressesThey []string, networkName string, role network.PhoneBookEntryRoles) {
	// prepare a map of items we'd like to remove.
	removeItems := make(map[peer.ID]bool, 0)
	peerIDs := ps.Peers()
	for _, pid := range peerIDs {
		data, _ := ps.Get(pid, "addressData")
		if data != nil {
			ad := data.(addressData)
			if ad.networkNames[networkName] && ad.role == role && !ad.persistent {
				removeItems[pid] = true
			}
		}

	}
	for _, addr := range addressesThey {
		info, err := peer.AddrInfoFromString(addr)
		if err != nil {
			return
		}
		data, _ := ps.Get(info.ID, "addressData")
		if data != nil {
			// we already have this.
			// Update the networkName
			ad := data.(addressData)
			ad.networkNames[networkName] = true

			// do not remove this entry
			delete(removeItems, info.ID)
		} else {
			// we don't have this item. add it.
			entry := makePhonebookEntryData(networkName, role, false)
			_ = ps.Put(info.ID, "addressData", entry)
		}
	}

	// remove items that were missing in addressesThey
	for k := range removeItems {
		ps.deletePhonebookEntry(k, networkName)
	}
}

// AddPersistentPeers stores addresses of peers which are persistent.
// i.e. they won't be replaced by ReplacePeerList calls
func (ps *PeerStore) AddPersistentPeers(dnsAddresses []string, networkName string, role network.PhoneBookEntryRoles) {

	for _, addr := range dnsAddresses {
		val, err := ps.Get("peerID", "addressData")
		if err == nil {
			// we already have this.
			// Make sure the persistence field is set to true
			data := val.(addressData)
			data.persistent = true
			_ = ps.Put("peerID", "addressData", data)

		} else {
			// we don't have this item. add it.
			ps.AddAddr("peerID", multiaddr.Cast([]byte(addr)), libp2p.PermanentAddrTTL)
			data := makePhonebookEntryData(networkName, role, true)
			_ = ps.Put("peerID", "addressData", data)
		}
	}
}

// Length returns the number of addrs in peerstore
func (ps *PeerStore) Length() int {
	return len(ps.PeersWithKeys())
}

// makePhonebookEntryData creates a new address entry for provided network name and role.
func makePhonebookEntryData(networkName string, role network.PhoneBookEntryRoles, persistent bool) addressData {
	pbData := addressData{
		networkNames:          make(map[string]bool),
		recentConnectionTimes: make([]time.Time, 0),
		role:                  role,
		persistent:            persistent,
	}
	pbData.networkNames[networkName] = true
	return pbData
}

func (ps *PeerStore) deletePhonebookEntry(peerID peer.ID, networkName string) {
	data, err := ps.Get(peerID, "addressData")
	if err != nil {
		return
	}
	ad := data.(addressData)
	delete(ad.networkNames, networkName)
	if 0 == len(ad.networkNames) {
		_ = ps.Put(peerID, "addressData", nil)
	}
}

// AppendTime adds the current time to recentConnectionTimes in
// addressData of addr
func (ps *PeerStore) appendTime(peerID peer.ID, t time.Time) {
	data, _ := ps.Get(peerID, "addressData")
	ad := data.(addressData)
	ad.recentConnectionTimes = append(ad.recentConnectionTimes, t)
	_ = ps.Put(peerID, "addressData", ad)
}

// PopEarliestTime removes the earliest time from recentConnectionTimes in
// addressData for addr
// It is expected to be later than ConnectionsRateLimitingWindow
func (ps *PeerStore) popNElements(n int, peerID peer.ID) {
	data, _ := ps.Get(peerID, "addressData")
	ad := data.(addressData)
	ad.recentConnectionTimes = ad.recentConnectionTimes[n:]
	_ = ps.Put(peerID, "addressData", ad)
}

func (ps *PeerStore) filterRetryTime(t time.Time, role network.PhoneBookEntryRoles) []string {
	o := make([]string, 0, len(ps.Peers()))
	for _, peerID := range ps.Peers() {
		data, _ := ps.Get(peerID, "addressData")
		if data != nil {
			ad := data.(addressData)
			if t.After(ad.retryAfter) && role == ad.role {
				o = append(o, peerID.String())
			}
		}
	}
	return o
}

func shuffleSelect(set []string, n int) []string {
	if n >= len(set) || n == network.GetAllAddresses {
		// return shuffled copy of everything
		out := slices.Clone(set)
		shuffleStrings(out)
		return out
	}
	// Pick random indexes from the set
	indexSample := make([]int, n)
	for i := range indexSample {
		indexSample[i] = rand.Intn(len(set)-i) + i
		for oi, ois := range indexSample[:i] {
			if ois == indexSample[i] {
				indexSample[i] = oi
			}
		}
	}
	out := make([]string, n)
	for i, index := range indexSample {
		out[i] = set[index]
	}
	return out
}

func shuffleStrings(set []string) {
	rand.Shuffle(len(set), func(i, j int) { t := set[i]; set[i] = set[j]; set[j] = t })
}
