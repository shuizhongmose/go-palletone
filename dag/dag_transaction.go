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
	"fmt"

	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/core"
	"github.com/palletone/go-palletone/dag/modules"
	"github.com/palletone/go-palletone/tokenengine"
)

type Txo4Greedy struct {
	modules.OutPoint
	Amount uint64
}

func (txo *Txo4Greedy) GetAmount() uint64 {
	return txo.Amount
}

func newTxo4Greedy(outPoint modules.OutPoint, amount uint64) *Txo4Greedy {
	return &Txo4Greedy{
		OutPoint: outPoint,
		Amount:   amount,
	}
}

func (dag *Dag) CreateBaseTransaction(from, to common.Address, daoAmount, daoFee uint64) (*modules.Transaction, error) {
	// 1. 获取转出账户所有的utxo
	allUtxos, err := dag.GetAddrUtxos(from)
	if err != nil {
		return &modules.Transaction{}, err
	}

	// 2. 利用贪心算法得到指定额度的utxo集合
	greedyUtxos := core.Utxos{}
	for outPoint, utxo := range allUtxos {
		tg := newTxo4Greedy(outPoint, utxo.Amount)
		greedyUtxos = append(greedyUtxos, tg)
	}

	selUtxos, change, err := core.Select_utxo_Greedy(greedyUtxos, daoAmount+daoFee)
	if err != nil {
		return nil, fmt.Errorf("Select utxo err")
	}

	// 3. 构建PaymentPayload的Inputs
	pload := new(modules.PaymentPayload)
	pload.LockTime = 0

	for _, selTxo := range selUtxos {
		tg := selTxo.(*Txo4Greedy)
		txInput := modules.NewTxIn(&tg.OutPoint, []byte{})
		pload.AddTxIn(txInput)
	}

	// 4. 构建PaymentPayload的Outputs
	outAmount := map[common.Address]uint64{}
	outAmount[to] = daoAmount
	if change > 0 {
		outAmount[from] = change
	}

	for addr, amount := range outAmount {
		pkScript := tokenengine.GenerateLockScript(addr)
		txOut := modules.NewTxOut(amount, pkScript, modules.NewPTNAsset())
		pload.AddTxOut(txOut)
	}

	// 4. 构建Transaction
	tx := &modules.Transaction{
		TxMessages: make([]*modules.Message, 0),
	}
	tx.TxMessages = append(tx.TxMessages, modules.NewMessage(modules.APP_PAYMENT, pload))

	return tx, nil
}
