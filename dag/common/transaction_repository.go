package common

import (
	"fmt"
	"github.com/palletone/go-palletone/common"
	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/core/accounts/keystore"
	"github.com/palletone/go-palletone/dag/errors"
	"github.com/palletone/go-palletone/dag/modules"
)

func GetTxSig(tx *modules.Transaction, ks *keystore.KeyStore, signer common.Address) ([]byte, error) {
	sign, err := ks.SigData(tx, signer)
	if err != nil {
		msg := fmt.Sprintf("Failed to singure transaction:%v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	return sign, nil
}

func ValidateTxSig(tx *modules.Transaction) bool {
	if tx == nil {
		return false
	}
	var sigs []modules.SignatureSet

	tmpTx := modules.Transaction{}
	//tmpTx.TxId = tx.TxId
	//if !bytes.Equal(tx.TxHash.Bytes(), tx.Hash().Bytes()){
	//	log.Error("ValidateTxSig", "transaction hash is not equal, tx req id:", tx.TxId)
	//	return false
	//}
	//todo 检查msg的有效性

	for _, msg := range tx.TxMessages {
		if msg.App == modules.APP_SIGNATURE {
			sigs = msg.Payload.(*modules.SignaturePayload).Signatures
		} else {
			tmpTx.TxMessages = append(tmpTx.TxMessages, msg)
		}
	}

	if len(sigs) > 0 {
		for i := 0; i < len(sigs); i++ {
			//fmt.Printf("sig[%v]-pubkey[%v]--tx[%v]", sigs[i].Signature, sigs[i].PubKey, tmpTx)
			if keystore.VerifyTXWithPK(sigs[i].Signature, tmpTx, sigs[i].PubKey) != true {
				log.Error("ValidateTxSig", "VerifyTXWithPK sig fail", tmpTx.RequestHash().String())
				return false
			}
		}
	}

	return true
}
