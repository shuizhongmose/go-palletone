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
 *
 */

package dag

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dedis/kyber/sign/bls"
	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/hexutil"
	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/contracts/syscontract"
	"github.com/palletone/go-palletone/core"

	"github.com/palletone/go-palletone/dag/modules"
)

func (dag *Dag) validateUnit(unit *modules.Unit) error {
	author := unit.Author()
	if !dag.IsActiveMediator(author) && !dag.IsPrecedingMediator(author) {
		errStr := fmt.Sprintf("The author(%v) of unit(%v) is not mediator!",
			author.Str(), unit.UnitHash.TerminalString())
		log.Debugf(errStr)

		return fmt.Errorf(errStr)
	}

	err := dag.validate.ValidateUnitExceptGroupSig(unit)
	if err != nil {
		return err
	}

	return nil
}

func (dag *Dag) validateUnitHeader(nextUnit *modules.Unit) error {
	pHash := nextUnit.ParentHash()[0]
	headHash, idx, _ := dag.stablePropRep.GetNewestUnit(nextUnit.Number().AssetID)
	if pHash != headHash {
		// todo 出现分叉, 调用本方法之前未处理分叉
		errStr := fmt.Sprintf("unit(%v) on the forked chain: parentHash(%v) not equal headUnitHash(%v)",
			nextUnit.UnitHash.TerminalString(), pHash.TerminalString(), headHash.TerminalString())
		log.Debugf(errStr)
		return fmt.Errorf(errStr)
	}

	if idx.Index+1 != nextUnit.NumberU64() {
		errStr := fmt.Sprintf("invalidated unit(%v)'s height number!, last height:%d, next unit height:%d",
			nextUnit.UnitHash.TerminalString(), idx.Index, nextUnit.NumberU64())
		log.Debugf(errStr)
		return fmt.Errorf(errStr)
	}

	return nil
}

// @author Albert·Gou
func (d *Dag) ValidateUnitExceptGroupSig(unit *modules.Unit) error {
	unitState := d.validate.ValidateUnitExceptGroupSig(unit)
	return unitState
}

// author Albert·Gou
func (d *Dag) IsActiveMediator(add common.Address) bool {
	return d.GetGlobalProp().IsActiveMediator(add)
}

func (d *Dag) IsPrecedingMediator(add common.Address) bool {
	return d.GetGlobalProp().IsPrecedingMediator(add)
}

func (dag *Dag) InitPropertyDB(genesis *core.Genesis, unit *modules.Unit) error {
	//  全局属性不是交易，不需要放在Unit中
	// @author Albert·Gou
	gp := modules.InitGlobalProp(genesis)
	if err := dag.stablePropRep.StoreGlobalProp(gp); err != nil {
		return err
	}

	//  动态全局属性不是交易，不需要放在Unit中
	// @author Albert·Gou
	dgp := modules.InitDynGlobalProp(unit)
	if err := dag.stablePropRep.StoreDynGlobalProp(dgp); err != nil {
		return err
	}
	//dag.stablePropRep.SetNewestUnit(unit.Header())

	//  初始化mediator调度器，并存在数据库
	// @author Albert·Gou
	ms := modules.InitMediatorSchl(gp, dgp)
	if err := dag.stablePropRep.StoreMediatorSchl(ms); err != nil {
		return err
	}

	return nil
}

func (dag *Dag) InitStateDB(genesis *core.Genesis, unit *modules.Unit) error {
	// Create initial mediators
	list := make(map[string]bool, len(genesis.InitialMediatorCandidates))
	for _, imc := range genesis.InitialMediatorCandidates {
		// 存储 mediator info
		err := imc.Validate()
		if err != nil {
			log.Debugf(err.Error())
			panic(err.Error())
		}

		mi := modules.NewMediatorInfo()
		*mi.MediatorInfoBase = *imc.MediatorInfoBase

		addr, _ := common.StringToAddress(mi.AddStr)
		err = dag.stableStateRep.StoreMediatorInfo(addr, mi)
		if err != nil {
			log.Debugf(err.Error())
			panic(err.Error())
		}

		list[mi.AddStr] = true
	}

	// 存储 initMediatorCandidates
	imcB, err := json.Marshal(list)
	if err != nil {
		log.Debugf(err.Error())
		return err
	}

	ws := &modules.ContractWriteSet{
		IsDelete: false,
		Key:      modules.MediatorList,
		Value:    imcB,
	}

	version := &modules.StateVersion{
		Height:  unit.Number(),
		TxIndex: ^uint32(0),
	}

	err = dag.stableStateRep.SaveContractState(syscontract.DepositContractAddress.Bytes(), ws, version)
	if err != nil {
		log.Debugf(err.Error())
		return err
	}

	return nil
}

func (dag *Dag) IsSynced() bool {
	gp := dag.GetGlobalProp()
	dgp := dag.GetDynGlobalProp()

	//nowFine := time.Now()
	//now := time.Unix(nowFine.Add(500*time.Millisecond).Unix(), 0)
	now := time.Now()
	nextSlotTime := dag.stablePropRep.GetSlotTime(gp, dgp, 1)

	if nextSlotTime.Before(now) {
		return false
	}

	return true
}

// author Albert·Gou
func (d *Dag) ChainThreshold() int {
	return d.GetGlobalProp().ChainThreshold()
}

func (d *Dag) PrecedingThreshold() int {
	return d.GetGlobalProp().PrecedingThreshold()
}

func (d *Dag) UnitIrreversibleTime() time.Duration {
	gp := d.GetGlobalProp()
	it := uint(gp.ChainThreshold()) * uint(gp.ChainParameters.MediatorInterval)
	return time.Duration(it) * time.Second
}

func (d *Dag) IsIrreversibleUnit(hash common.Hash) bool {
	exist, err := d.stableUnitRep.IsHeaderExist(hash)
	if err != nil {
		log.Errorf("IsHeaderExist execute error:%s", err.Error())
		return false
	}

	return exist
}

func (d *Dag) GetIrreversibleUnit(id modules.AssetId) (*modules.ChainIndex, error) {
	_, idx, err := d.stablePropRep.GetNewestUnit(id)
	return idx, err
}

func (d *Dag) VerifyUnitGroupSign(unitHash common.Hash, groupSign []byte) error {
	unit, err := d.GetUnitByHash(unitHash)
	if err != nil {
		log.Debug(err.Error())
		return err
	}

	pubKey, err := unit.GroupPubKey()
	if err != nil {
		log.Debug(err.Error())
		return err
	}

	err = bls.Verify(core.Suite, pubKey, unitHash[:], groupSign)
	if err == nil {
		//log.Debug("the group signature: " + hexutil.Encode(groupSign) +
		//	" of the Unit that hash: " + unitHash.Hex() + " is verified through!")
	} else {
		log.Debug("the group signature: " + hexutil.Encode(groupSign) + " of the Unit that hash: " +
			unitHash.Hex() + " is verified that an error has occurred: " + err.Error())
		return err
	}

	return nil
}

func (dag *Dag) IsConsecutiveMediator(nextMediator common.Address) bool {
	dgp := dag.GetDynGlobalProp()

	if !dgp.IsShuffledSchedule && nextMediator.Equal(dgp.LastMediator) {
		return true
	}

	return false
}

// 计算最近64个生产slots的mediator参与度，不包括当前unit
// Calculate the percent of unit production slots that were missed in the
// past 64 units, not including the current unit.
func (dag *Dag) MediatorParticipationRate() uint32 {
	popCount := func(x uint64) uint8 {
		m := []uint64{
			0x5555555555555555,
			0x3333333333333333,
			0x0F0F0F0F0F0F0F0F,
			0x00FF00FF00FF00FF,
			0x0000FFFF0000FFFF,
			0x00000000FFFFFFFF,
		}

		var i, w uint8
		for i, w = 0, 1; i < 6; i, w = i+1, w+w {
			x = (x & m[i]) + ((x >> w) & m[i])
		}

		return uint8(x)
	}

	recentSlotsFilled := dag.GetDynGlobalProp().RecentSlotsFilled
	participationRate := core.PalletOne100Percent * int(popCount(recentSlotsFilled)) / 64

	return uint32(participationRate)
}
