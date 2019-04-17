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
 * @author PalletOne core developers <dev@pallet.one>
 * @date 2018
 */

package jury

import (
	"encoding/json"
	"time"

	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/crypto"
	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/core/accounts"
	"github.com/palletone/go-palletone/dag/errors"
)

func checkValid(reqEvt *AdapterRequestEvent) bool {
	hash := crypto.Hash(reqEvt.ConsultData, reqEvt.Answer)
	log.Debugf("sig: %s", common.Bytes2Hex(reqEvt.Sig))
	sig := reqEvt.Sig[:len(reqEvt.Sig)-1] // remove recovery id
	return crypto.VerifySignature(reqEvt.Pubkey, hash, sig)
}
func (p *Processor) saveSig(msgType uint32, reqEvt *AdapterRequestEvent) (firstSave bool) {
	p.locker.Lock()
	defer p.locker.Unlock()

	if _, exist := p.mtx[reqEvt.ReqId]; !exist { //todo how to process
		return true
	}
	if p.mtx[reqEvt.ReqId].adaInf == nil {
		p.mtx[reqEvt.ReqId].adaInf = make(map[uint32]*AdapterInf)
	}
	pubkeyHex := common.Bytes2Hex(reqEvt.Pubkey)
	if _, exist := p.mtx[reqEvt.ReqId].adaInf[msgType]; !exist {
		//all jury msg
		typeAdaInf := &AdapterInf{JuryMsgAll: make(map[string]*MsgSigCollect)}
		//one msg collect
		msgSigCollect := &MsgSigCollect{OneMsgAllSig: make(map[string]JuryMsgSig)}
		msgSigCollect.OneMsgAllSig[pubkeyHex] = JuryMsgSig{reqEvt.Sig, reqEvt.Answer}
		typeAdaInf.JuryMsgAll[string(reqEvt.ConsultData)] = msgSigCollect
		//
		p.mtx[reqEvt.ReqId].adaInf[msgType] = typeAdaInf
	} else {
		//
		typeAdaInf := p.mtx[reqEvt.ReqId].adaInf[msgType]
		if _, existCollect := typeAdaInf.JuryMsgAll[string(reqEvt.ConsultData)]; !existCollect { //new collect
			msgSigCollect := &MsgSigCollect{OneMsgAllSig: make(map[string]JuryMsgSig)}
			msgSigCollect.OneMsgAllSig[pubkeyHex] = JuryMsgSig{reqEvt.Sig, reqEvt.Answer}
			typeAdaInf.JuryMsgAll[string(reqEvt.ConsultData)] = msgSigCollect
		} else {
			if _, exist := typeAdaInf.JuryMsgAll[string(reqEvt.ConsultData)].OneMsgAllSig[pubkeyHex]; exist {
				return false
			}
			typeAdaInf.JuryMsgAll[string(reqEvt.ConsultData)].OneMsgAllSig[pubkeyHex] = JuryMsgSig{reqEvt.Sig, reqEvt.Answer}
		}
	}
	return true
}

func (p *Processor) checkJury(reqEvt *AdapterRequestEvent) bool {
	if _, exist := p.lockArf[reqEvt.ContractId]; !exist {
		return false
	}
	pubkeyHex := common.Bytes2Hex(reqEvt.Pubkey)
	juryAll := p.lockArf[reqEvt.ContractId]
	for i := range juryAll {
		if common.Bytes2Hex(juryAll[i].PublicKey) == pubkeyHex {
			return true
		}
	}
	return false
}

func (p *Processor) processAdapterRequestEvent(msgType uint32, reqEvt *AdapterRequestEvent) (result *AdapterEvent, err error) {
	log.Info("processAdapterRequestEvent")

	//if not this contract's jury, just return
	if !p.checkJury(reqEvt) {
		localMediators := p.ptn.GetLocalMediators()
		if len(localMediators) == 0 {
			return nil, nil
		} //mediator continue process
	}

	//check
	isValid := checkValid(reqEvt)
	if !isValid {
		return nil, errors.New("Event invalid")
	}
	//save
	firstSave := p.saveSig(msgType, reqEvt)

	//broadcast
	if firstSave { //first receive, broadcast
		go p.ptn.AdapterBroadcast(AdapterEvent{AType: AdapterEventType(msgType), Event: reqEvt})
	}

	return nil, nil
}

func (p *Processor) AdapterFunRequest(reqId common.Hash, contractId common.Address, msgType uint32, consultContent []byte, myAnswer []byte) ([]byte, error) {
	if reqId == (common.Hash{}) {
		return nil, errors.New("AdapterFunRequest param is nil")
	}
	log.Info("AdapterFunRequest")
	//
	account := p.getLocalAccount()
	if account == nil {
		return nil, errors.New("AdapterFunRequest no local account")
	}

	//
	hash := crypto.Hash(consultContent, myAnswer)
	sig, err := p.ptn.GetKeyStore().SignHashWithPassphrase(accounts.Account{Address: account.Address}, account.Password, hash)
	if err != nil {
		return nil, errors.New("AdapterFunRequest SignHashWithPassphrase failed")
	}
	log.Debugf("sig: %s", common.Bytes2Hex(sig))
	//
	pubKey, err := p.ptn.GetKeyStore().GetPublicKey(account.Address)
	if err != nil {
		return nil, errors.New("AdapterFunRequest GetPublicKey failed")
	}
	//
	reqEvt := &AdapterRequestEvent{
		ReqId:       reqId,
		ContractId:  contractId,
		ConsultData: consultContent,
		Answer:      myAnswer,
		Sig:         sig,
		Pubkey:      pubKey,
	}

	//todo delete test
	isValid := checkValid(reqEvt)
	if !isValid {
		return nil, errors.New("Event invalid")
	}

	go p.ptn.AdapterBroadcast(AdapterEvent{AType: AdapterEventType(msgType), Event: reqEvt})

	//save
	p.saveSig(msgType, reqEvt)

	return sig, nil
}

func (p *Processor) getRusult(reqId common.Hash, msgType uint32, consultContent []byte) ([]byte, error) {
	p.locker.Lock()
	defer p.locker.Unlock()

	adaInf := p.mtx[reqId].adaInf[msgType]
	if len(adaInf.JuryMsgAll[string(consultContent)].OneMsgAllSig) >= p.contractSigNum {
		var juryMsgSigAll []JuryMsgSig
		for _, juryMsgSig := range adaInf.JuryMsgAll[string(consultContent)].OneMsgAllSig {
			juryMsgSigAll = append(juryMsgSigAll, juryMsgSig)
		}
		result, err := json.Marshal(juryMsgSigAll)
		return result, err
	}

	return nil, errors.New("Not enough")
}
func (p *Processor) AdapterFunResult(reqId common.Hash, contractId common.Address, msgType uint32, consultContent []byte, timeOut time.Duration) ([]byte, error) {
	if reqId == (common.Hash{}) {
		return nil, errors.New("AdapterFunRequest param is nil")
	}
	log.Info("AdapterFunResult")
	result, err := p.getRusult(reqId, msgType, consultContent)
	if err == nil {
		return result, nil
	}

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(timeOut)
		timeout <- true
	}()

	select {
	case <-timeout:
		result, err := p.getRusult(reqId, msgType, consultContent)
		if err == nil {
			return result, nil
		}
		log.Debug("AdapterFunRequest, time out")
		return nil, errors.New("AdapterFunRequest, time out")
	}
}

func (p *Processor) ProcessAdapterEvent(event *AdapterEvent) (result *AdapterEvent, err error) {
	if event == nil {
		return nil, errors.New("ProcessAdapterEvent, event is nil")
	}
	log.Info("ProcessAdapterEvent", "event", event.AType)

	return p.processAdapterRequestEvent(uint32(event.AType), event.Event.(*AdapterRequestEvent))
}
