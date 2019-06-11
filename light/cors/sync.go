// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package cors

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/common/p2p/discover"
	"github.com/palletone/go-palletone/dag/modules"
	"github.com/palletone/go-palletone/ptn/downloader"
)

// dataPack is a data message returned by a peer for some query.
type dataPack interface {
	PeerId() string
	Items() int
	Stats() string
}

// headerPack is a batch of block headers returned by a peer.
type headerPack struct {
	peerId  string
	headers []*modules.Header
}

func (p *headerPack) PeerId() string { return p.peerId }
func (p *headerPack) Items() int     { return len(p.headers) }
func (p *headerPack) Stats() string  { return fmt.Sprintf("%d", len(p.headers)) }

func (pm *ProtocolManager) StartCorsSync() (string, error) {
	mainchain, err := pm.dag.GetMainChain()
	if mainchain == nil || err != nil {
		log.Debug("Cors ProtocolManager StartCorsSync", "GetMainChain err", err)
		return err.Error(), err
	}
	pm.mclock.Lock()
	pm.mainchain = mainchain
	pm.mclock.Unlock()

	pm.mclock.RLock()
	for _, peer := range mainchain.Peers {
		node, err := discover.ParseNode(peer)
		if err != nil {
			log.Debugf("Cors ProtocolManager StartCorsSync invalid pnode: %v", err)
			continue
		}
		pm.server.corss.AddPeer(node)
	}
	pm.mclock.RUnlock()

	go func() {
		time.Sleep(time.Duration(3) * time.Second)
		if pm.peers.Len() >= pm.mainchainpeers()/2+1 {
			pm.PullSync()
			pm.PushSync()
		}
	}()

	return "OK", nil
}

func (pm *ProtocolManager) PushSync() {
	log.Debug("Enter Cors ProtocolManager PushSync")
	defer log.Debug("End Cors ProtocolManager PushSync")
	if atomic.LoadUint32(&pm.corsSync) == 0 {
		atomic.StoreUint32(&pm.corsSync, 1)
		index, _ := pm.pushSync()
		log.Info("Cors Push Sync OK", "index", index)
		atomic.StoreUint32(&pm.corsSync, 0)
	}
}

func (pm *ProtocolManager) pushSync() (uint64, error) {
	var (
		bytes   common.StorageSize
		headers []*modules.Header
		index   uint64
		flag    int
	)

	pheader, err := pm.fetchHeader()
	if err != nil {
		log.Debug("Cors ProtocolManager", "pushSync fetchHeader err", err)
		return 0, err
	}

	flag = 0
	if pheader.Number.Index <= fsMinFullBlocks {
		index = 0
	} else {
		index = pheader.Number.Index - fsMinFullBlocks
	}

	log.Debug("Cors ProtocolManager", "pheader.index", pheader.Number.Index, "push index", index, "pushSync fetchHeader header", pheader)

	number := &modules.ChainIndex{pm.assetId, index}
	for {
		bytes = 0
		headers = []*modules.Header{}

		for bytes < softResponseLimit && len(headers) < MaxHeaderFetch {
			bytes += estHeaderRlpSize
			number.Index = index
			header, err := pm.dag.GetHeaderByNumber(number)
			if err != nil {
				if len(headers) == MaxHeaderFetch {
					index--
					break
				} else {
					flag = 1
				}
				break
			}
			headers = append(headers, header)
			index++
		}

		rand.Seed(time.Now().UnixNano())
		peers := pm.peers.AllPeers()
		x := rand.Intn(len(peers))
		p := peers[x]
		log.Info("Cors ProtocolManager", "pushSync SendHeaders len(headers)", len(headers), "index", index)
		if len(headers) == 0 {
			header := modules.Header{}
			number := modules.ChainIndex{pm.assetId, 0}
			header.Number = &number
			headers = append(headers, &header)
		}
		p.SendHeaders(headers)
		if flag == 1 {
			break
		} else {
			time.Sleep(waitPushSync)
		}
	}
	return index, nil
}

// requestTTL returns the current timeout allowance for a single download request
// to finish under.
func (pm *ProtocolManager) requestTTL() time.Duration {
	var (
		rtt  = time.Duration(atomic.LoadUint64(&pm.rttEstimate))
		conf = float64(atomic.LoadUint64(&pm.rttConfidence)) / 1000000.0
	)
	ttl := time.Duration(ttlScaling) * time.Duration(float64(rtt)/conf)
	if ttl > ttlLimit {
		ttl = ttlLimit
	}
	return ttl
}

func (pm *ProtocolManager) fetchHeader() (*modules.Header, error) {
	// Request the advertised remote head block and wait for the response
	rand.Seed(time.Now().UnixNano())
	peers := pm.peers.AllPeers()
	log.Debug("Cors ProtocolManager fetchHeader", "len(peers)", len(peers))
	x := rand.Intn(len(peers))
	p := peers[x]
	log.Debug("Retrieving remote all token", "peer", p.ID())
	var number modules.ChainIndex
	number.AssetID = pm.assetId
	go p.RequestCurrentHeader(number)

	ttl := pm.requestTTL()
	timeout := time.After(ttl)
	for {
		select {
		case <-pm.quitSync:
			return nil, errCancelHeaderFetch

		case packet := <-pm.headerCh:
			// Discard anything not from the origin peer
			if packet.PeerId() != p.id {
				log.Debug("Received headers from incorrect peer", "peer", packet.PeerId())
				break
			}
			// Make sure the peer actually gave something valid
			headers := packet.(*headerPack).headers
			if len(headers) != 1 {
				log.Debug("Multiple headers for single request", "headers", len(headers), "peer", p.id)
				return nil, errBadPeer
			}
			log.Debug("Remote leaf nodes", "counts", len(headers), "peer", packet.PeerId())
			return headers[0], nil

		case <-timeout:
			log.Debug("Waiting for head header timed out", "elapsed", ttl, "peer", p.id)
			return nil, errTimeout
		}
	}
	return nil, nil
}

func (pm *ProtocolManager) PullSync() {
	log.Debug("Enter Cors ProtocolManager PullSync")
	defer log.Debug("End Cors ProtocolManager PullSync")

	peer := pm.peers.BestPeer()
	if peer == nil {
		return
	}
	//TODO modify get from getMainChain
	if peer.headInfo.number.AssetID != modules.PTNCOIN {
		log.Debug("Cors PalletOne ProtocolManager PullSync", "peer assetid", peer.headInfo.number.AssetID)
		return
	}

	if atomic.LoadUint32(&pm.corsSync) == 0 {
		atomic.StoreUint32(&pm.corsSync, 1)
		pm.pullSync(peer)
		log.Info("Cors Pull Sync OK")
		atomic.StoreUint32(&pm.corsSync, 0)
	}

	if header := pm.dag.CurrentHeader(modules.PTNCOIN); header != nil {
		pm.server.SendEvents(header)
	} else {
		log.Debug("Cors PalletOne ProtocolManager PullSync ptn CurrentHeader is nil")
	}
}

func (pm *ProtocolManager) pullSync(peer *peer) {
	//peer := pm.peers.BestPeer()
	//if peer == nil {
	//	return
	//}
	//if peer.headInfo.number.AssetID != modules.PTNCOIN {
	//	log.Debug("Cors PalletOne ProtocolManager pullSync", "peer assetid", peer.headInfo.number.AssetID)
	//	return
	//}
	lheader := pm.dag.CurrentHeader(modules.PTNCOIN)
	//hash, number := peer.HeadAndNumber(modules.PTNCOIN)
	//if lheader.Number.Index >= number.Index {
	//	log.Debug("Cors PalletOne ProtocolManager pullSync is not need sync", "local index", lheader.Number.Index, "peer index", number.Index)
	//	return
	//}

	if err := pm.downloader.Synchronise(peer.id, lheader.Hash(), lheader.Number.Index, downloader.LightSync, modules.PTNCOIN); err != nil {
		log.Debug("ptn sync downloader.", "Synchronise err:", err)
		return
	}
}
