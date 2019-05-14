package cors

import (
	"crypto/ecdsa"
	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/common/p2p"
	"github.com/palletone/go-palletone/dag/modules"
	"github.com/palletone/go-palletone/ptn"
)

type CorsServer struct {
	config          *ptn.Config
	protocolManager *ProtocolManager
	privateKey      *ecdsa.PrivateKey
	quitSync        chan struct{}
}

func NewCoresServer(ptn *ptn.PalletOne, config *ptn.Config) (*CorsServer, error) {
	quitSync := make(chan struct{})
	gasToken := config.Dag.GetGasToken()
	genesis, err := ptn.Dag().GetGenesisUnit()
	if err != nil {
		log.Error("Light PalletOne New", "get genesis err:", err)
		return nil, err
	}

	pm, err := NewCorsProtocolManager(true, newPeerSet(), config.NetworkId, gasToken,
		ptn.Dag(), ptn.EventMux(), genesis, quitSync)
	if err != nil {
		log.Error("NewlesServer NewProtocolManager", "err", err)
		return nil, err
	}

	srv := &CorsServer{
		config:          config,
		protocolManager: pm,
		quitSync:        quitSync,
	}
	pm.server = srv

	return srv, nil
}

func (s *CorsServer) Protocols() []p2p.Protocol {
	return s.protocolManager.SubProtocols
}

// Start starts the LES server
func (s *CorsServer) Start(srvr *p2p.Server) {
	s.protocolManager.Start(s.config.LightPeers)
	s.privateKey = srvr.PrivateKey
	s.protocolManager.blockLoop()
}

// Stop stops the LES service
func (s *CorsServer) Stop() {
	go func() {
		<-s.protocolManager.noMorePeers
	}()
	s.protocolManager.Stop()
}

func (pm *ProtocolManager) blockLoop() {
	pm.wg.Add(1)
	headCh := make(chan modules.ChainHeadEvent, 10)
	headSub := pm.dag.SubscribeChainHeadEvent(headCh)
	go func() {
		var lastHead *modules.Header
		for {
			select {
			case ev := <-headCh:
				peers := pm.peers.AllPeers()
				if len(peers) > 0 {
					header := ev.Unit.Header()
					hash := header.Hash()
					number := header.Number.Index
					//td := core.GetTd(pm.chainDb, hash, number)
					if lastHead == nil || (header.Number.Index > lastHead.Number.Index) {
						lastHead = header
						log.Debug("Announcing block to peers", "number", number, "hash", hash)

						announce := announceData{Hash: hash, Number: *lastHead.Number, Header: *lastHead}
						var (
							signed         bool
							signedAnnounce announceData
						)

						for _, p := range peers {
							log.Debug("Light Palletone", "ProtocolManager->blockLoop p.announceType", p.announceType)
							switch p.announceType {

							case announceTypeSimple:
								select {
								case p.announceChn <- announce:
								default:
									pm.removePeer(p.id)
								}

							case announceTypeSigned:
								if !signed {
									signedAnnounce = announce
									signedAnnounce.sign(pm.server.privateKey)
									signed = true
								}

								select {
								case p.announceChn <- signedAnnounce:
								default:
									pm.removePeer(p.id)
								}
							}
						}
					}
				}
			case <-pm.quitSync:
				headSub.Unsubscribe()
				pm.wg.Done()
				return
			}
		}
	}()
}
