/*
   This file is part of go-palletone.
   go-palletone is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-palletone is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with go-palletone.  If not, see <http://www.gnu.org/licenses/>.
*/
/*
 * @author PalletOne core developer Albert·Gou <dev@pallet.one>
 * @date 2018
 */

package mediatorplugin

import (
	"fmt"

	"github.com/dedis/kyber"
	"github.com/dedis/kyber/share/dkg/pedersen"
	"github.com/dedis/kyber/share/vss/pedersen"
	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/event"
	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/dag/modules"
)

func GenInitPair(suite vss.Suite) (kyber.Scalar, kyber.Point) {
	sc := suite.Scalar().Pick(suite.RandomStream())

	return sc, suite.Point().Mul(sc, nil)
}

func (mp *MediatorPlugin) BroadcastVSSDeals() {
	for localMed, dkg := range mp.dkgs {
		deals, err := dkg.Deals()
		if err != nil {
			log.Error(err.Error())
		}

		for index, deal := range deals {
			event := VSSDealEvent{
				DstIndex: index,
				Deal:     deal,
			}

			go mp.vssDealFeed.Send(event)
		}

		go mp.processResponseBuf(&dkgVerifier{localMed, localMed})
	}
}

func (mp *MediatorPlugin) SubscribeVSSDealEvent(ch chan<- VSSDealEvent) event.Subscription {
	return mp.vssDealScope.Track(mp.vssDealFeed.Subscribe(ch))
}

func (mp *MediatorPlugin) ToProcessDeal(deal *VSSDealEvent) error {
	select {
	case <-mp.quit:
		return errTerminated
	case mp.toProcessDealCh <- deal:
		return nil
	}
}

func (mp *MediatorPlugin) processDealLoop() {
	for {
		select {
		case <-mp.quit:
			return
		case deal := <-mp.toProcessDealCh:
			go mp.processVSSDeal(deal)
		}
	}
}

func (mp *MediatorPlugin) getLocalActiveMediatorDKG(add common.Address) *dkg.DistKeyGenerator {
	if !mp.IsLocalActiveMediator(add) {
		log.Error(fmt.Sprintf("The following mediator is not local active mediator: %v", add.String()))
		return nil
	}

	dkg, ok := mp.dkgs[add]
	if !ok || dkg == nil {
		log.Error(fmt.Sprintf("The following mediator`s dkg is not existed: %v", add.String()))
		return nil
	}

	return dkg
}

func (mp *MediatorPlugin) processVSSDeal(dealEvent *VSSDealEvent) {
	dag := mp.getDag()
	localMed := dag.GetActiveMediatorAddr(dealEvent.DstIndex)

	dkgr := mp.getLocalActiveMediatorDKG(localMed)
	if dkgr == nil {
		return
	}

	deal := dealEvent.Deal

	resp, err := dkgr.ProcessDeal(deal)
	if err != nil {
		log.Error(err.Error())
		return
	}

	vrfrMed := dag.GetActiveMediatorAddr(int(deal.Index))
	go mp.processResponseBuf(&dkgVerifier{localMed, vrfrMed})

	if resp.Response.Status != vss.StatusApproval {
		log.Error(fmt.Sprintf("DKG: own deal gave a complaint: %v", localMed.String()))
		return
	}

	respEvent := VSSResponseEvent{
		Resp: resp,
	}
	go mp.vssResponseFeed.Send(respEvent)
}

func (mp *MediatorPlugin) ToProcessResponse(resp *VSSResponseEvent) error {
	select {
	case <-mp.quit:
		return errTerminated
	case mp.toProcessResponseCh <- resp:
		return nil
	}
}

func (mp *MediatorPlugin) processResponseLoop() {
	for {
		select {
		case <-mp.quit:
			return
		case resp := <-mp.toProcessResponseCh:
			go mp.processVSSResponse(resp)
		}
	}
}

func (mp *MediatorPlugin) initRespBuf(localMed common.Address) {
	aSize := mp.getDag().GetActiveMediatorCount()
	mp.respBuf[localMed] = make(map[common.Address]chan *dkg.Response, aSize)

	for i := 0; i < aSize; i++ {
		vrfrMed := mp.getDag().GetActiveMediatorAddr(i)
		mp.respBuf[localMed][vrfrMed] = make(chan *dkg.Response, aSize-1)
	}
}

func (mp *MediatorPlugin) processResponseBuf(dvp *dkgVerifier) {
	localMed := dvp.localMed
	vrfrMed := dvp.vrfrMed
	dkg := mp.getLocalActiveMediatorDKG(localMed)
	if dkg == nil {
		return
	}

	respCount := 0
	aSize := mp.getDag().GetActiveMediatorCount()
	for resp := range mp.respBuf[localMed][vrfrMed] {
		respCount++

		jstf, err := dkg.ProcessResponse(resp)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		if jstf != nil {
			log.Error(fmt.Sprintf("DKG: wrong Process Response: %v", localMed.String()))
			continue
		}

		if respCount+1 >= aSize-1 {
			log.Debug(fmt.Sprintf("%v 's DKG certifing is %v, vrfrMed: %v, count: %v",
				localMed.Str(), dkg.Certified(), vrfrMed.Str(), respCount))
		}
	}
}

func (mp *MediatorPlugin) processVSSResponse(respEvent *VSSResponseEvent) {
	resp := respEvent.Resp
	lams := mp.GetLocalActiveMediators()
	for _, localMed := range lams {
		dag := mp.getDag()

		//ignore the message from myself
		srcIndex := resp.Response.Index
		srcMed := dag.GetActiveMediatorAddr(int(srcIndex))
		if srcMed == localMed {
			continue
		}

		dkg := mp.getLocalActiveMediatorDKG(localMed)
		if dkg == nil {
			continue
		}

		vrfrMed := dag.GetActiveMediatorAddr(int(resp.Index))
		mp.respBuf[localMed][vrfrMed] <- resp
	}
}

func (mp *MediatorPlugin) SubscribeVSSResponseEvent(ch chan<- VSSResponseEvent) event.Subscription {
	return mp.vssResponseScope.Track(mp.vssResponseFeed.Subscribe(ch))
}

func (mp *MediatorPlugin) ToUnitTBLSSign(peer string, unit *modules.Unit) error {
	op := &toBLSSigned{
		origin: peer,
		unit:   unit,
	}

	select {
	case <-mp.quit:
		return errTerminated
	case mp.toBLSSigned <- op:
		return nil
	}
}

func (mp *MediatorPlugin) unitBLSSignLoop() {
	for {
		select {
		// Mediator Plugin terminating, abort operation
		case <-mp.quit:
			return
		case op := <-mp.toBLSSigned:
			//			PushUnit(mp.ptn.Dag(), op.unit)
			go mp.unitBLSSign(op)
		}
	}
}

func (mp *MediatorPlugin) unitBLSSign(toBLSSigned *toBLSSigned) {
	//todo

}
