package main

import (
	"crypto/elliptic"
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	json "github.com/json-iterator/go"
	"time"
)

func NewAccount(value string) common.Address {
	return common.HexToAddress(value)
}

type Tx struct {
	From     common.Address `json:"from"`
	To       common.Address `json:"to"`
	Gas      uint           `json:"gas"`
	GasPrice uint           `json:"gas_price"`
	Value    uint           `json:"value"`
	Nonce    uint           `json:"nonce"`
	Data     string         `json:"data"`
	Time     uint64         `json:"time"`
}

type SignedTx struct {
	Tx
	Sig []byte `json:"signature"`
}

func NewTx(from, to common.Address, gas uint, gasPrice uint, value, nonce uint, data string) Tx {
	return Tx{from, to, gas, gasPrice, value, nonce, data, uint64(time.Now().Unix())}
}

func NewBaseTx(from, to common.Address, value, nonce uint, data string) Tx {
	return NewTx(from, to, TxGas, TxGasPriceDefault, value, nonce, data)
}

func NewSignedTx(tx Tx, sig []byte) SignedTx {
	return SignedTx{tx, sig}
}

func (t Tx) IsReward() bool {
	return t.Data == "reward"
}

func (t Tx) Cost(isTip1Fork bool) uint {
	if isTip1Fork {
		return t.Value + t.GasCost()
	}

	return t.Value + TxFee
}

func (t Tx) GasCost() uint {
	return t.Gas * t.GasPrice
}

func (t Tx) Hash() (Hash, error) {
	txJson, err := t.Encode()
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(txJson), nil
}

func (t Tx) Encode() ([]byte, error) {
	return json.Marshal(t)
}

// MarshalJSON is the main source of truth for encoding a TX for hash calculation from expected attributes.
//
// The logic is bit ugly and hacky but prevents infinite marshaling loops of embedded objects and allows
// the structure to change with new TIPs.
func (t Tx) MarshalJSON() ([]byte, error) {
	// Prior TIP1
	if t.Gas == 0 {
		type legacyTx struct {
			From  common.Address `json:"from"`
			To    common.Address `json:"to"`
			Value uint           `json:"value"`
			Nonce uint           `json:"nonce"`
			Data  string         `json:"data"`
			Time  uint64         `json:"time"`
		}
		return json.Marshal(legacyTx{
			From:  t.From,
			To:    t.To,
			Value: t.Value,
			Nonce: t.Nonce,
			Data:  t.Data,
			Time:  t.Time,
		})
	}

	type tip1Tx struct {
		From     common.Address `json:"from"`
		To       common.Address `json:"to"`
		Gas      uint           `json:"gas"`
		GasPrice uint           `json:"gas_price"`
		Value    uint           `json:"value"`
		Nonce    uint           `json:"nonce"`
		Data     string         `json:"data"`
		Time     uint64         `json:"time"`
	}

	return json.Marshal(tip1Tx{
		From:     t.From,
		To:       t.To,
		Gas:      t.Gas,
		GasPrice: t.GasPrice,
		Value:    t.Value,
		Nonce:    t.Nonce,
		Data:     t.Data,
		Time:     t.Time,
	})
}

// MarshalJSON is the main source of truth for encoding a TX for hash calculation (backwards compatible for TIPs).
//
// The logic is bit ugly and hacky but prevents infinite marshaling loops of embedded objects and allows
// the structure to change with new TIPs.
func (t SignedTx) MarshalJSON() ([]byte, error) {
	// Prior TIP1
	if t.Gas == 0 {
		type legacyTx struct {
			From  common.Address `json:"from"`
			To    common.Address `json:"to"`
			Value uint           `json:"value"`
			Nonce uint           `json:"nonce"`
			Data  string         `json:"data"`
			Time  uint64         `json:"time"`
			Sig   []byte         `json:"signature"`
		}
		return json.Marshal(legacyTx{
			From:  t.From,
			To:    t.To,
			Value: t.Value,
			Nonce: t.Nonce,
			Data:  t.Data,
			Time:  t.Time,
			Sig:   t.Sig,
		})
	}

	type tip1Tx struct {
		From     common.Address `json:"from"`
		To       common.Address `json:"to"`
		Gas      uint           `json:"gas"`
		GasPrice uint           `json:"gas_price"`
		Value    uint           `json:"value"`
		Nonce    uint           `json:"nonce"`
		Data     string         `json:"data"`
		Time     uint64         `json:"time"`
		Sig      []byte         `json:"signature"`
	}

	return json.Marshal(tip1Tx{
		From:     t.From,
		To:       t.To,
		Gas:      t.Gas,
		GasPrice: t.GasPrice,
		Value:    t.Value,
		Nonce:    t.Nonce,
		Data:     t.Data,
		Time:     t.Time,
		Sig:      t.Sig,
	})
}

func (t SignedTx) Hash() (Hash, error) {
	txJson, err := t.Encode()
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(txJson), nil
}

func (t SignedTx) IsAuthentic() (bool, error) {
	txHash, err := t.Tx.Hash()
	if err != nil {
		return false, err
	}

	recoveredPubKey, err := crypto.SigToPub(txHash[:], t.Sig)
	if err != nil {
		return false, err
	}

	recoveredPubKeyBytes := elliptic.Marshal(crypto.S256(), recoveredPubKey.X, recoveredPubKey.Y)
	recoveredPubKeyBytesHash := crypto.Keccak256(recoveredPubKeyBytes[1:])
	recoveredAccount := common.BytesToAddress(recoveredPubKeyBytesHash[12:])

	return recoveredAccount.Hex() == t.From.Hex(), nil
}
