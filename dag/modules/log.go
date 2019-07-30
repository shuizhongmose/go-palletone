// Copyright 2018 PalletOne

package modules

import (
	"fmt"
	"io"

	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
)

//go:generate gencodec -type Log -field-override logMarshaling -out gen_log_json.go

// Log represents a contract log event. These events are generated by the LOG opcode and
// stored/indexed by the node.
type Log struct {
	// Consensus fields:
	// address of the contract that generated the event
	Address common.Address `json:"address" gencodec:"required"`
	// list of topics provided by the contract.
	Topics []common.Hash `json:"topics" gencodec:"required"`
	// supplied by the contract, usually ABI-encoded
	Data []byte `json:"data" gencodec:"required"`

	// Derived fields. These fields are filled in by the node
	// but not secured by consensus.
	// block in which the transaction was included
	UnitNumber uint64 `json:"unitNumber"`
	// UnitHash of the transaction
	TxHash common.Hash `json:"transactionHash" gencodec:"required"`
	// index of the transaction in the block
	TxIndex uint `json:"transactionIndex" gencodec:"required"`
	// UnitHash of the block in which the transaction was included
	UnitHash common.Hash `json:"unitHash"`
	// index of the log in the receipt
	Index uint `json:"logIndex" gencodec:"required"`

	// The Removed field is true if this log was reverted due to a chain reorganization.
	// You must pay attention to this field if you receive logs through a filter query.
	Removed bool `json:"removed"`
}

type logMarshaling struct {
	Data       hexutil.Bytes
	UnitNumber hexutil.Uint64
	TxIndex    hexutil.Uint
	Index      hexutil.Uint
}

type rlpLog struct {
	Address common.Address
	Topics  []common.Hash
	Data    []byte
}

type rlpStorageLog struct {
	Address    common.Address
	Topics     []common.Hash
	Data       []byte
	UnitNumber uint64
	TxHash     common.Hash
	TxIndex    uint
	UnitHash   common.Hash
	Index      uint
}

// EncodeRLP implements rlp.Encoder.
func (l *Log) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, rlpLog{Address: l.Address, Topics: l.Topics, Data: l.Data})
}

// DecodeRLP implements rlp.Decoder.
func (l *Log) DecodeRLP(s *rlp.Stream) error {
	var dec rlpLog
	err := s.Decode(&dec)
	if err == nil {
		l.Address, l.Topics, l.Data = dec.Address, dec.Topics, dec.Data
	}
	return err
}

func (l *Log) String() string {
	return fmt.Sprintf(`log: %x %x %x %x %d %x %d`, l.Address, l.Topics, l.Data, l.TxHash, l.TxIndex, l.UnitHash, l.Index)
}

// LogForStorage is a wrapper around a Log that flattens and parses the entire content of
// a log including non-consensus fields.
type LogForStorage Log

// EncodeRLP implements rlp.Encoder.
func (l *LogForStorage) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, rlpStorageLog{
		Address:    l.Address,
		Topics:     l.Topics,
		Data:       l.Data,
		UnitNumber: l.UnitNumber,
		TxHash:     l.TxHash,
		TxIndex:    l.TxIndex,
		UnitHash:   l.UnitHash,
		Index:      l.Index,
	})
}

// DecodeRLP implements rlp.Decoder.
func (l *LogForStorage) DecodeRLP(s *rlp.Stream) error {
	var dec rlpStorageLog
	err := s.Decode(&dec)
	if err == nil {
		*l = LogForStorage{
			Address:    dec.Address,
			Topics:     dec.Topics,
			Data:       dec.Data,
			UnitNumber: dec.UnitNumber,
			TxHash:     dec.TxHash,
			TxIndex:    dec.TxIndex,
			UnitHash:   dec.UnitHash,
			Index:      dec.Index,
		}
	}
	return err
}
