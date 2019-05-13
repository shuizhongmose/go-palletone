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

package storage

import (
	"bytes"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/common/ptndb"
	"github.com/palletone/go-palletone/core"
	"github.com/palletone/go-palletone/dag/constants"
	"github.com/palletone/go-palletone/dag/modules"
)

func mediatorKey(address common.Address) []byte {
	key := append(constants.MEDIATOR_INFO_PREFIX, address.Bytes21()...)
	//key := append(constants.MEDIATOR_INFO_PREFIX, address.Str()...)

	return key
}

func StoreMediator(db ptndb.Database, med *core.Mediator) error {
	mi := modules.MediatorToInfo(med)

	return StoreMediatorInfo(db, med.Address, mi)
}

func StoreMediatorInfo(db ptndb.Database, add common.Address, mi *modules.MediatorInfo) error {
	log.Debugf("Store Mediator Info %v:", mi.AddStr)

	err := storeToJson(db, mediatorKey(add), mi)
	if err != nil {
		log.Debugf("Store mediator error:%v", err.Error())
		return err
	}

	return nil
}

func RetrieveMediatorInfo(db ptndb.Database, address common.Address) (*modules.MediatorInfo, error) {
	mi := modules.NewMediatorInfo()

	err := readFromJson(db, mediatorKey(address), mi)
	if err != nil {
		log.Errorf("Retrieve mediator error: %v", err.Error())
		return nil, err
	}

	return mi, nil
}

func RetrieveMediator(db ptndb.Database, address common.Address) (*core.Mediator, error) {
	mi, err := RetrieveMediatorInfo(db, address)
	if mi == nil || err != nil {
		return nil, err
	}

	med := mi.InfoToMediator()
	//med.Address = address

	return med, nil
}

func GetMediatorCount(db ptndb.Database) int {
	mc := getCountByPrefix(db, constants.MEDIATOR_INFO_PREFIX)

	return mc
}

// todo
func IsMediator(db ptndb.Database, address common.Address) bool {
	has, err := db.Has(mediatorKey(address))
	if err != nil {
		log.Debugf("Error in determining if it is a mediator: %v", err.Error())
	}

	return has
}

// todo
func GetMediators(db ptndb.Database) map[common.Address]bool {
	result := make(map[common.Address]bool)

	iter := db.NewIteratorWithPrefix(constants.MEDIATOR_INFO_PREFIX)
	for iter.Next() {
		key := iter.Key()
		if key == nil {
			continue
		}

		//log.Debugf("Get Mediator's key : %s", key))
		addB := bytes.TrimPrefix(key, constants.MEDIATOR_INFO_PREFIX)

		result[common.BytesToAddress(addB)] = true
		//result[core.StrToMedAdd(string(addStr))] = true
	}

	return result
}

// todo
func LookupMediator(db ptndb.Database) map[common.Address]*core.Mediator {
	result := make(map[common.Address]*core.Mediator)

	iter := db.NewIteratorWithPrefix(constants.MEDIATOR_INFO_PREFIX)
	for iter.Next() {
		key := iter.Key()
		if key == nil {
			continue
		}

		value := iter.Value()
		if value == nil {
			continue
		}

		mi := modules.NewMediatorInfo()
		err := rlp.DecodeBytes(value, mi)
		if err != nil {
			log.Debugf("Error in Decoding Bytes to MediatorInfo: %v", err.Error())
		}

		addB := bytes.TrimPrefix(key, constants.MEDIATOR_INFO_PREFIX)
		add := common.BytesToAddress(addB)
		med := mi.InfoToMediator()
		//med.Address = add

		result[add] = med
	}

	return result
}
