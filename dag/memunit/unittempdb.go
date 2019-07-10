/*
 *
 *    This file is part of go-palletone.
 *    go-palletone is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU General Public License as published by
 *    the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *    go-palletone is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU General Public License for more details.
 *    You should have received a copy of the GNU General Public License
 *    along with go-palletone.  If not, see <http://www.gnu.org/licenses/>.
 * /
 *
 *  * @author PalletOne core developer <dev@pallet.one>
 *  * @date 2018-2019
 *
 */

package memunit

import (
	"github.com/palletone/go-palletone/common/ptndb"
	"github.com/palletone/go-palletone/dag/common"
	"github.com/palletone/go-palletone/dag/modules"
	"github.com/palletone/go-palletone/dag/palletcache"
	"github.com/palletone/go-palletone/validator"
)

type UnitTempDb struct {
	Tempdb    *Tempdb
	UnitRep   common.IUnitRepository
	UtxoRep   common.IUtxoRepository
	StateRep  common.IStateRepository
	PropRep   common.IPropRepository
	Validator validator.Validator
	Unit      *modules.Unit
}

func NewUnitTempDb(db ptndb.Database, newestUnit *modules.Unit, cache palletcache.ICache) *UnitTempDb {
	tempdb, _ := NewTempdb(db)
	trep := common.NewUnitRepository4Db(tempdb)
	tutxoRep := common.NewUtxoRepository4Db(tempdb)
	tstateRep := common.NewStateRepository4Db(tempdb)
	tpropRep := common.NewPropRepository4Db(tempdb)
	v := validator.NewValidate(trep, tutxoRep, tstateRep, tpropRep, cache)
	return &UnitTempDb{
		Tempdb:    tempdb,
		UnitRep:   trep,
		UtxoRep:   tutxoRep,
		StateRep:  tstateRep,
		PropRep:   tpropRep,
		Validator: v,
		Unit:      newestUnit,
	}
}
