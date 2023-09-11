package eg

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/LukeEuler/dolly/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	cc "github.com/LukeEuler/eth-go/common"
	"github.com/LukeEuler/eth-go/config"
	"github.com/LukeEuler/eth-go/key"
)

func Transfer() {
	conf := config.Get()

	unSupportTxTypes := make([]string, 0, 3)
	for txType := range conf.Transfer {
		switch txType {
		default:
			unSupportTxTypes = append(unSupportTxTypes, txType)
		case "eth", "eth_mpc", "token", "NFT", "create_contract":
		}
	}
	if len(unSupportTxTypes) != 0 {
		log.Entry.Warnf("invalid transfer types: %v", unSupportTxTypes)
	}

	tfs := conf.Transfer["eth"]
	for _, tf := range tfs {
		transfer(tf)
	}

	tfs = conf.Transfer["eth_mpc"]
	for _, tf := range tfs {
		transfer_mpc(tf)
	}

	tfs = conf.Transfer["create_contract"]
	for _, tf := range tfs {
		createContract(tf)
	}

	tfs = conf.Transfer["token"]
	for _, tf := range tfs {
		transferToken(tf)
	}
}

func transfer(tf *config.Transfer) {
	if !tf.Enable {
		return
	}
	conf := config.Get()

	privateKey, ok := conf.KeyPair[tf.From]
	if !ok {
		log.Entry.Warnf("can not found key for: %s", tf.From)
		return
	}
	key, err := key.NewKeyFromHex(privateKey)
	if err != nil {
		log.Entry.Fatal(err)
	}
	initClient()

	nonce, err := getNonce(tf.From)
	if err != nil {
		log.Entry.Fatal(err)
	}

	signer := types.NewLondonSigner(big.NewInt(conf.Net.ChainID))

	amount, ok := big.NewInt(0).SetString(tf.Amount, 10)
	if !ok {
		log.Entry.Fatalf("invalid amount: %s", tf.Amount)
	}

	to := common.HexToAddress(tf.To)

	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		GasTipCap: big.NewInt(0).SetUint64(tf.MaxPriorityFeePerGas),
		GasFeeCap: big.NewInt(0).SetUint64(tf.MaxFeePerGas),
		Gas:       tf.GasLimit,
		To:        &to,
		Value:     amount,
		Data:      nil,
	})

	tx, err = types.SignTx(tx, signer, key.ToECDSA())
	if err != nil {
		log.Entry.Fatal(err)
	}

	bs, err := tx.MarshalJSON()
	if err != nil {
		log.Entry.Fatal(err)
	}
	fmt.Println("------------ tx json ------------")
	fmt.Println(string(bs))

	data, err := tx.MarshalBinary()
	if err != nil {
		log.Entry.Fatal(err)
	}

	raw := hexutil.Encode(data)
	fmt.Println("\n------------ tx raw ------------")
	fmt.Println(raw)
	var res string
	err = node.SyncCall(&res, "eth_sendRawTransaction", raw)
	if err != nil {
		log.Entry.Fatal(err)
	}

	fmt.Printf("\n%s/tx/%s\n\n", conf.Net.Show, res)
}

func transfer_mpc(tf *config.Transfer) {
	if !tf.Enable {
		return
	}
	conf := config.Get()
	initClient()

	nonce, err := getNonce(tf.From)
	if err != nil {
		log.Entry.Fatal(err)
	}

	signer := types.NewLondonSigner(big.NewInt(conf.Net.ChainID))

	amount, ok := big.NewInt(0).SetString(tf.Amount, 10)
	if !ok {
		log.Entry.Fatalf("invalid amount: %s", tf.Amount)
	}

	to := common.HexToAddress(tf.To)
	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		GasTipCap: big.NewInt(0).SetUint64(tf.MaxPriorityFeePerGas),
		GasFeeCap: big.NewInt(0).SetUint64(tf.MaxFeePerGas),
		Gas:       tf.GasLimit,
		To:        &to,
		Value:     amount,
		Data:      nil,
	})

	h := signer.Hash(tx)

	needSig := false
	if len(tf.R) == 0 || len(tf.S) == 0 {
		needSig = true
	}
	if needSig {
		fmt.Println("---------- eth mpc ----------")
		fmt.Println(hex.EncodeToString(h[:]))
		return
	}

	bs, err := hex.DecodeString(tf.R + tf.S)
	if err != nil {
		log.Entry.Fatal(err)
	}
	sig := []byte{}
	sig = append(sig, bs...)
	sig = append(sig, byte(tf.RecID))

	tx, err = tx.WithSignature(signer, sig)
	if err != nil {
		log.Entry.Fatal(err)
	}

	bs, err = tx.MarshalJSON()
	if err != nil {
		log.Entry.Fatal(err)
	}
	fmt.Println("------------ tx json ------------")
	fmt.Println(string(bs))

	data, err := tx.MarshalBinary()
	if err != nil {
		log.Entry.Fatal(err)
	}

	raw := hexutil.Encode(data)
	fmt.Println("\n------------ tx raw ------------")
	fmt.Println(raw)
	var res string
	err = node.SyncCall(&res, "eth_sendRawTransaction", raw)
	if err != nil {
		log.Entry.Fatal(err)
	}

	fmt.Printf("\n%s/tx/%s\n\n", conf.Net.Show, res)
}

func createContract(tf *config.Transfer) {
	if !tf.Enable {
		return
	}
	conf := config.Get()

	privateKey, ok := conf.KeyPair[tf.From]
	if !ok {
		log.Entry.Warnf("can not found key for: %s", tf.From)
		return
	}
	key, err := key.NewKeyFromHex(privateKey)
	if err != nil {
		log.Entry.Fatal(err)
	}
	initClient()

	nonce, err := getNonce(tf.From)
	if err != nil {
		log.Entry.Fatal(err)
	}

	signer := types.NewLondonSigner(big.NewInt(conf.Net.ChainID))

	ddd := cc.FormatHex(tf.Data)
	dataBs, err := hex.DecodeString(ddd)
	if err != nil {
		log.Entry.Fatal(err)
	}
	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		GasTipCap: big.NewInt(0).SetUint64(tf.MaxPriorityFeePerGas),
		GasFeeCap: big.NewInt(0).SetUint64(tf.MaxFeePerGas),
		Gas:       tf.GasLimit,
		Data:      dataBs,
	})

	tx, err = types.SignTx(tx, signer, key.ToECDSA())
	if err != nil {
		log.Entry.Fatal(err)
	}

	bs, err := tx.MarshalJSON()
	if err != nil {
		log.Entry.Fatal(err)
	}
	fmt.Println("------------ tx json ------------")
	fmt.Println(string(bs))

	data, err := tx.MarshalBinary()
	if err != nil {
		log.Entry.Fatal(err)
	}

	raw := hexutil.Encode(data)
	fmt.Println("\n------------ tx raw ------------")
	fmt.Println(raw)
	var res string
	err = node.SyncCall(&res, "eth_sendRawTransaction", raw)
	if err != nil {
		log.Entry.Fatal(err)
	}

	fmt.Printf("\n%s/tx/%s\n\n", conf.Net.Show, res)
}

func transferToken(tf *config.Transfer) {
	if !tf.Enable {
		return
	}
	conf := config.Get()

	privateKey, ok := conf.KeyPair[tf.From]
	if !ok {
		log.Entry.Warnf("can not found key for: %s", tf.From)
		return
	}
	key, err := key.NewKeyFromHex(privateKey)
	if err != nil {
		log.Entry.Fatal(err)
	}
	initClient()

	nonce, err := getNonce(tf.From)
	if err != nil {
		log.Entry.Fatal(err)
	}

	signer := types.NewLondonSigner(big.NewInt(conf.Net.ChainID))

	amount, ok := big.NewInt(0).SetString(tf.Amount, 10)
	if !ok {
		log.Entry.Fatalf("invalid amount: %s", tf.Amount)
	}

	dataBs, err := buildTokenTransferData(tf.To, amount)
	if err != nil {
		log.Entry.Fatal(err)
	}

	to := common.HexToAddress(tf.Contract)
	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		GasTipCap: big.NewInt(0).SetUint64(tf.MaxPriorityFeePerGas),
		GasFeeCap: big.NewInt(0).SetUint64(tf.MaxFeePerGas),
		Gas:       tf.GasLimit,
		To:        &to,
		Value:     big.NewInt(0),
		Data:      dataBs,
	})

	tx, err = types.SignTx(tx, signer, key.ToECDSA())
	if err != nil {
		log.Entry.Fatal(err)
	}

	bs, err := tx.MarshalJSON()
	if err != nil {
		log.Entry.Fatal(err)
	}
	fmt.Println("------------ tx json ------------")
	fmt.Println(string(bs))

	data, err := tx.MarshalBinary()
	if err != nil {
		log.Entry.Fatal(err)
	}

	raw := hexutil.Encode(data)
	fmt.Println("\n------------ tx raw ------------")
	fmt.Println(raw)
	var res string
	err = node.SyncCall(&res, "eth_sendRawTransaction", raw)
	if err != nil {
		log.Entry.Fatal(err)
	}

	fmt.Printf("\n%s/tx/%s\n\n", conf.Net.Show, res)
}

func getNonce(address string) (uint64, error) {
	var res string
	err := node.SyncCall(&res, "eth_getTransactionCount", "0x"+address, "latest")
	if err != nil {
		return 0, err
	}
	nonce := big.NewInt(0)
	raw := cc.FormatHex(res)
	if len(raw) > 0 {
		_, ok := nonce.SetString(raw, 16)
		if !ok {
			return 0, errors.Errorf("invalid value: %s", res)
		}
	}
	return nonce.Uint64(), nil
}

// 0xa9059cbb
// 000000000000000000000000e94f9ce88501e57353ab2c31b1f3b47429c92b21
// 000000000000000000000000000000000000000000000000000000e8d4a51000
func buildTokenTransferData(to string, amount *big.Int) ([]byte, error) {
	to = strings.TrimPrefix(strings.ToLower(to), "0x")
	for len(to) < 64 {
		to = "0" + to
	}

	value := amount.Text(16)
	for len(value) < 64 {
		value = "0" + value
	}

	raw := "a9059cbb" + to + value

	bs, err := hex.DecodeString(raw)
	return bs, errors.WithStack(err)
}
