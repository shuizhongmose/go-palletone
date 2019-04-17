// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package les

import (
	"context"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/crypto"
	"github.com/palletone/go-palletone/dag/modules"
)

var sha3_nil = crypto.HashResult(nil)

func GetHeaderByNumber(ctx context.Context, odr OdrBackend, number uint64) (*modules.Header, error) {
	return nil, nil
	/*
		db := odr.Database()
		hash := core.GetCanonicalHash(db, number)
		if (hash != common.Hash{}) {
			// if there is a canonical hash, there is a header too
			header := core.GetHeader(db, hash, number)
			if header == nil {
				panic("Canonical hash present but header not found")
			}
			return header, nil
		}

		var (
			chtCount, sectionHeadNum uint64
			sectionHead              common.Hash
		)
		if odr.ChtIndexer() != nil {
			chtCount, sectionHeadNum, sectionHead = odr.ChtIndexer().Sections()
			canonicalHash := core.GetCanonicalHash(db, sectionHeadNum)
			// if the CHT was injected as a trusted checkpoint, we have no canonical hash yet so we accept zero hash too
			for chtCount > 0 && canonicalHash != sectionHead && canonicalHash != (common.Hash{}) {
				chtCount--
				if chtCount > 0 {
					sectionHeadNum = chtCount*CHTFrequencyClient - 1
					sectionHead = odr.ChtIndexer().SectionHead(chtCount - 1)
					canonicalHash = core.GetCanonicalHash(db, sectionHeadNum)
				}
			}
		}
		if number >= chtCount*CHTFrequencyClient {
			return nil, ErrNoTrustedCht
		}
		r := &ChtRequest{ChtRoot: GetChtRoot(db, chtCount-1, sectionHead), ChtNum: chtCount - 1, BlockNum: number}
		if err := odr.Retrieve(ctx, r); err != nil {
			return nil, err
		}
		return r.Header, nil
	*/
}

func GetCanonicalHash(ctx context.Context, odr OdrBackend, number uint64) (common.Hash, error) {
	return common.Hash{}, nil
	//hash := core.GetCanonicalHash(odr.Database(), number)
	//if (hash != common.Hash{}) {
	//	return hash, nil
	//}
	//header, err := GetHeaderByNumber(ctx, odr, number)
	//if header != nil {
	//	return header.Hash(), nil
	//}
	//return common.Hash{}, err
}

// GetBodyRLP retrieves the block body (transactions and uncles) in RLP encoding.
func GetBodyRLP(ctx context.Context, odr OdrBackend, hash common.Hash, number uint64) (rlp.RawValue, error) {
	//if data := core.GetBodyRLP(odr.Database(), hash, number); data != nil {
	//	return data, nil
	//}
	//r := &BlockRequest{Hash: hash, Number: number}
	//if err := odr.Retrieve(ctx, r); err != nil {
	//	return nil, err
	//} else {
	//	return r.Rlp, nil
	//}
	return nil, nil
}

// GetBody retrieves the block body (transactons, uncles) corresponding to the
// hash.
func GetBody(ctx context.Context, odr OdrBackend, hash common.Hash, number uint64) (*modules.Transactions, error) {
	return nil, nil
	//data, err := GetBodyRLP(ctx, odr, hash, number)
	//if err != nil {
	//	return nil, err
	//}
	//body := new(types.Body)
	//if err := rlp.Decode(bytes.NewReader(data), body); err != nil {
	//	return nil, err
	//}
	//return body, nil
}

// GetBlock retrieves an entire block corresponding to the hash, assembling it
// back from the stored header and body.
func GetBlock(ctx context.Context, odr OdrBackend, hash common.Hash, number uint64) (*modules.Unit, error) {
	// Retrieve the block header and body contents
	//header := core.GetHeader(odr.Database(), hash, number)
	//if header == nil {
	//	return nil, ErrNoHeader
	//}
	//body, err := GetBody(ctx, odr, hash, number)
	//if err != nil {
	//	return nil, err
	//}
	//// Reassemble the block and return
	//return types.NewBlockWithHeader(header).WithBody(body.Transactions, body.Uncles), nil
	return nil, nil
}
