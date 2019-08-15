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
	"time"

	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/core/accounts/keystore"
	dagcommon "github.com/palletone/go-palletone/dag/common"
	"github.com/palletone/go-palletone/dag/modules"
	"github.com/palletone/go-palletone/dag/txspool"
)

// GenerateUnit, generate unit
// @author Albert·Gou
func (dag *Dag) GenerateUnit(when time.Time, producer common.Address, groupPubKey []byte,
	ks *keystore.KeyStore, txpool txspool.ITxPool) *modules.Unit {
	t0 := time.Now()
	defer func(start time.Time) {
		log.Debugf("GenerateUnit cost time: %v", time.Since(start))
	}(t0)

	// 1. 判断是否满足生产的若干条件

	// 2. 生产unit，添加交易集、时间戳、签名
	newUnit, err := dag.CreateUnit(producer, txpool, when)
	if err != nil {
		log.Debug("GenerateUnit", "error", err.Error())
		return nil
	}
	// added by yangyu, 2018.8.9
	if newUnit == nil || newUnit.IsEmpty() {
		log.Info("No unit need to be packaged for now.", "unit", newUnit)
		return nil
	}

	newUnit.UnitHeader.Time = when.Unix()
	newUnit.UnitHeader.GroupPubKey = groupPubKey
	newUnit.Hash()

	sign_unit, err1 := dagcommon.GetUnitWithSig(newUnit, ks, producer)
	if err1 != nil {
		log.Debugf("GetUnitWithSig error: %v", err)
		return nil
	}

	sign_unit.UnitSize = sign_unit.Size()
	log.Infof("Generate new unit index:[%d],hash:[%s],size:%s, parent unit[%s],txs[%d], spent time: %s",
		sign_unit.NumberU64(), sign_unit.Hash().String(), sign_unit.UnitSize.String(),
		sign_unit.UnitHeader.ParentsHash[0].String(), sign_unit.Txs.Len(), time.Since(t0).String())

	//3.将新单元添加到MemDag中
	a, b, c, d, e, err := dag.Memdag.AddUnit(sign_unit, txpool)
	if a != nil && err == nil {
		dag.unstableUnitRep = a
		dag.unstableUtxoRep = b
		dag.unstableStateRep = c
		dag.unstablePropRep = d
		dag.unstableUnitProduceRep = e
	}

	//4.PostChainEvents
	//TODO add PostChainEvents
	go func() {
		var (
			events = make([]interface{}, 0, 2)
		)
		events = append(events, modules.ChainHeadEvent{Unit: sign_unit})
		events = append(events, modules.ChainEvent{Unit: sign_unit, Hash: sign_unit.UnitHash})
		dag.PostChainEvents(events)
	}()

	return sign_unit
}
