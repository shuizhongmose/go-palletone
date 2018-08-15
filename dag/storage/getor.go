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

package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"reflect"
	"unsafe"

	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/rlp"
	config "github.com/palletone/go-palletone/dag/dagconfig"
	"github.com/palletone/go-palletone/dag/modules"
)

// DatabaseReader wraps the Get method of a backing data store.
type DatabaseReader interface {
	Get(key []byte) (value []byte, err error)
}

// @author Albert·Gou
func Retrieve(key string, v interface{}) error {
	//rv := reflect.ValueOf(v)
	//if rv.Kind() != reflect.Ptr || rv.IsNil() {
	//	return errors.New("an invalid argument, the argument must be a non-nil pointer")
	//}

	data, err := Get([]byte(key))
	if err != nil {
		return err
	}

	err = rlp.DecodeBytes(data, v)
	if err != nil {
		return err
	}

	return nil
}

// get bytes
func Get(key []byte) ([]byte, error) {
	if Dbconn == nil {
		Dbconn = ReNewDbConn(config.DefaultConfig.DbPath)
	}
	// return Dbconn.Get(key)
	b, err := Dbconn.Get(key)
	return b, err
}

// get string
func GetString(key []byte) (string, error) {
	if Dbconn == nil {
		Dbconn = ReNewDbConn(config.DefaultConfig.DbPath)
	}
	if re, err := Dbconn.Get(key); err != nil {
		return "", err
	} else {
		return *(*string)(unsafe.Pointer(&re)), nil
	}
}

// get prefix: return maps
func GetPrefix(prefix []byte) map[string][]byte {
	if Dbconn != nil {
		return getprefix(Dbconn, prefix)
	} else {
		db := ReNewDbConn(config.DefaultConfig.DbPath)
		if db == nil {
			return nil
		}

		Dbconn = db
		return getprefix(db, prefix)
	}

}

// get prefix
func getprefix(db DatabaseReader, prefix []byte) map[string][]byte {
	iter := Dbconn.NewIteratorWithPrefix(prefix)
	result := make(map[string][]byte)
	for iter.Next() {
		key := iter.Key()
		value := make([]byte, 0)
		// 请注意： 直接赋值取得iter.Value()的最后一个指针
		result[string(key)] = append(value, iter.Value()...)
	}
	return result
}

func GetUnit(hash common.Hash) *modules.Unit {
	unit_bytes, err := Get(append(UNIT_PREFIX, hash.Bytes()...))
	log.Println(err)
	var unit modules.Unit
	json.Unmarshal(unit_bytes, &unit)

	return &unit
}
func GetUnitFormIndex(height uint64, asset modules.IDType16) *modules.Unit {
	key := fmt.Sprintf("%s_%s_%d", UNIT_NUMBER_PREFIX, asset.String(), height)
	hash, err := Get([]byte(key))
	if err != nil {
		return nil
	}
	var h common.Hash
	h.SetBytes(hash)
	return GetUnit(h)
}

func GetHeader(hash common.Hash, index uint64) *modules.Header {

	encNum := encodeBlockNumber(index)
	key := append(HEADER_PREFIX, encNum...)
	header_bytes, err := Get(append(key, hash.Bytes()...))
	// rlp  to  Header struct
	log.Println(err)
	header := new(modules.Header)
	if err := rlp.Decode(bytes.NewReader(header_bytes), &header); err != nil {
		log.Println("Invalid unit header rlp:", err)
		return nil
	}

	return header
}

func GetHeaderRlp(db DatabaseReader, hash common.Hash, index uint64) rlp.RawValue {
	encNum := encodeBlockNumber(index)
	key := append(HEADER_PREFIX, encNum...)
	header_bytes, err := db.Get(append(key, hash.Bytes()...))
	// rlp  to  Header struct
	log.Println(err)
	return header_bytes
}

func GetHeaderFormIndex(height uint64, asset modules.IDType16) *modules.Header {
	unit := GetUnitFormIndex(height, asset)
	return unit.UnitHeader
}

// GetTxLookupEntry
func GetTxLookupEntry(db DatabaseReader, hash common.Hash) (common.Hash, uint64, uint64) {
	data, _ := Get(append(LookupPrefix, hash.Bytes()...))
	if len(data) == 0 {
		return common.Hash{}, 0, 0
	}
	var entry modules.TxLookupEntry
	if err := rlp.DecodeBytes(data, &entry); err != nil {
		return common.Hash{}, 0, 0
	}
	return entry.UnitHash, entry.UnitIndex, entry.Index

}

// GetTransaction retrieves a specific transaction from the database , along with its added positional metadata
// p2p 同步区块 分为同步header 和body。 GetBody可以省掉节点包装交易块的过程。
func GetTransaction(hash common.Hash) (*modules.Transaction, common.Hash, uint64, uint64) {
	unitHash, unitNumber, txIndex := GetTxLookupEntry(Dbconn, hash)
	if unitHash != (common.Hash{}) {
		body, _ := GetBody(unitHash)
		if body == nil || len(body) <= int(txIndex) {
			return nil, common.Hash{}, 0, 0
		}
		tx, err := gettrasaction(body[txIndex])
		if err == nil {
			return tx, unitHash, unitNumber, txIndex
		}
	}
	tx, err := gettrasaction(hash)
	if err != nil {
		return nil, unitHash, unitNumber, txIndex
	}
	return tx, unitHash, unitNumber, txIndex
}

// gettrasaction can get a transaction by hash.
func gettrasaction(hash common.Hash) (*modules.Transaction, error) {
	if hash == (common.Hash{}) {
		return nil, errors.New("hash is not exist.")
	}
	data, err := Get(append(TRANSACTION_PREFIX, hash.Bytes()...))
	if err != nil {
		return nil, err
	}
	tx := new(modules.Transaction)
	if err := rlp.DecodeBytes(data, &tx); err != nil {
		return nil, err
	}
	return tx, nil
}

// GetContract can get a Contract by the contract hash
func GetContract(id common.Hash) (*modules.Contract, error) {
	if common.EmptyHash(id) {
		return nil, errors.New("the filed not defined")
	}
	con_bytes, err := Get(append(CONTRACT_PTEFIX, id[:]...))
	if err != nil {
		log.Println("err:", err)
		return nil, err
	}
	contract := new(modules.Contract)
	err = rlp.DecodeBytes(con_bytes, contract)
	if err != nil {
		log.Println("err:", err)
		return nil, err
	}
	return contract, nil
}

func GetContractRlp(id common.Hash) (rlp.RawValue, error) {
	if common.EmptyHash(id) {
		return nil, errors.New("the filed not defined")
	}
	con_bytes, err := Get(append(CONTRACT_PTEFIX, id[:]...))
	if err != nil {
		return nil, err
	}
	return con_bytes, nil
}

// Get contract key's value
func GetContractKeyValue(id common.Hash, key string) (interface{}, error) {
	var val interface{}
	if common.EmptyHash(id) {
		return nil, errors.New("the filed not defined")
	}
	con_bytes, err := Get(append(CONTRACT_PTEFIX, id[:]...))
	if err != nil {
		return nil, err
	}
	contract := new(modules.Contract)
	err = rlp.DecodeBytes(con_bytes, contract)
	if err != nil {
		log.Println("err:", err)
		return nil, err
	}
	obj := reflect.ValueOf(contract)
	myref := obj.Elem()
	typeOftype := myref.Type()

	for i := 0; i < myref.NumField(); i++ {
		filed := myref.Field(i)
		if typeOftype.Field(i).Name == key {
			val = filed.Interface()
			log.Println(i, ". ", typeOftype.Field(i).Name, " ", filed.Type(), "=: ", filed.Interface())
			break
		} else if i == myref.NumField()-1 {
			val = nil
		}
	}
	return val, nil
}

const missingNumber = uint64(0xffffffffffffffff)

func GetUnitNumber(db DatabaseReader, hash common.Hash) uint64 {
	data, _ := db.Get(append(UNIT_HASH_NUMBER_Prefix, hash.Bytes()...))
	if len(data) != 8 {
		return missingNumber
	}
	return binary.BigEndian.Uint64(data)
}

//  GetCanonicalHash get

func GetCanonicalHash(db DatabaseReader, number uint64) (common.Hash, error) {
	key := append(HEADER_PREFIX, encodeBlockNumber(number)...)
	data, err := db.Get(append(key, NumberSuffix...))
	if err != nil {
		return common.Hash{}, err
	}
	if len(data) == 0 {
		return common.Hash{}, err
	}
	return common.BytesToHash(data), nil
}
func GetHeadHeaderHash(db DatabaseReader, hash common.Hash) (common.Hash, error) {
	data, err := db.Get(HeadHeaderKey)
	if err != nil {
		return common.Hash{}, err
	}
	if len(data) != 8 {
		return common.Hash{}, errors.New("data's len is error.")
	}
	return common.BytesToHash(data), nil
}

// GetHeadUnitHash stores the head unit's hash.
func GetHeadUnitHash(db DatabaseReader, hash common.Hash) (common.Hash, error) {
	data, err := db.Get(HeadUnitKey)
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(data), nil
}

// GetHeadFastUnitHash stores the fast head unit's hash.
func GetHeadFastUnitHash(db DatabaseReader, hash common.Hash) (common.Hash, error) {
	data, err := db.Get(HeadFastKey)
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(data), nil
}

// GetTrieSyncProgress stores the fast sync trie process counter to support
// retrieving it across restarts.
func GetTrieSyncProgress(db DatabaseReader, count uint64) (uint64, error) {
	data, err := db.Get(TrieSyncKey)
	if err != nil {
		return 0, err
	}
	return new(big.Int).SetBytes(data).Uint64(), nil
}

//  dbFetchUtxoEntry
func GetUtxoEntry(db DatabaseReader, key []byte) (*modules.Utxo, error) {
	utxo := new(modules.Utxo)
	data, err := db.Get(key)
	if err != nil {
		return nil, err
	}

	if err := rlp.DecodeBytes(data, &utxo); err != nil {
		return nil, err
	}

	return utxo, nil
}

// GetAdddrTransactionsHash
func GetAddrTransactionsHash(addr string) ([]common.Hash, error) {
	data, err := Get(append(AddrTransactionsHash_Prefix, []byte(addr)...))
	if err != nil {
		return []common.Hash{}, err
	}
	hashs := make([]common.Hash, 0)
	if err := rlp.DecodeBytes(data, hashs); err != nil {
		return []common.Hash{}, err
	}
	return hashs, nil
}

// GetAddrTransactions
func GetAddrTransactions(addr string) (modules.Transactions, error) {
	data, err := Get(append(AddrTransactionsHash_Prefix, []byte(addr)...))
	if err != nil {
		return modules.Transactions{}, err
	}
	hashs := make([]common.Hash, 0)
	if err := rlp.DecodeBytes(data, hashs); err != nil {
		return modules.Transactions{}, err
	}
	txs := make(modules.Transactions, 0)
	for _, hash := range hashs {
		tx, _, _, _ := GetTransaction(hash)
		txs = append(txs, tx)
	}
	return txs, nil
}

// Get income transactions
func GetAddrOutput(addr string) ([]modules.Output, error) {
	data := GetPrefix(append(AddrOutput_Prefix, []byte(addr)...))
	outputs := make([]modules.Output, 0)
	var err error
	for _, b := range data {
		out := new(modules.Output)
		if err := rlp.DecodeBytes(b, &out); err == nil {
			outputs = append(outputs, *out)
		} else {
			err = err
		}
	}
	return outputs, err
}
