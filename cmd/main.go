package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/parnurzeal/gorequest"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// StRpcRespError rpc 错误
type StRpcRespError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func (e *StRpcRespError) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}

// StRpcReq rpc请求
type StRpcReq struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// StRpcResp rpc返回
type StRpcResp struct {
	ID    string          `json:"id"`
	Error *StRpcRespError `json:"error"`
}

// StEthTransaction 交易
type StEthTransaction struct {
	From             string      `json:"from"`
	Gas              string      `json:"gas"`
	GasPrice         string      `json:"gasPrice"`
	Hash             string      `json:"hash"`
	Input            string      `json:"input"`
	Nonce            string      `json:"nonce"`
	R                string      `json:"r"`
	S                string      `json:"s"`
	To               string      `json:"to"`
	TransactionIndex interface{} `json:"transactionIndex"`
	Type             string      `json:"type"`
	V                string      `json:"v"`
	Value            string      `json:"value"`
}

// doReq 发送请求
func doReq(client *gorequest.SuperAgent, method string, arqs []interface{}, resp interface{}) error {
	_, body, errs := client.
		Send(StRpcReq{
			Jsonrpc: "1.0",
			ID:      "1",
			Method:  method,
			Params:  arqs,
		}).EndBytes()
	if errs != nil {
		return errs[0]
	}
	err := json.Unmarshal(body, resp)
	if err != nil {
		return err
	}
	return nil
}

// EthRpcNetVersion 获取block信息
// "1": Ethereum Mainnet
// "2": Morden Testnet (deprecated)
// "3": Ropsten Testnet
// "4": Rinkeby Testnet
// "42": Kovan Testnet
func EthRpcNetVersion(rpcURI string) (int64, error) {
	resp := struct {
		StRpcResp
		Result int64 `json:"result,string"`
	}{}
	err := doReq(
		ethClient(rpcURI),
		"net_version",
		nil,
		&resp,
	)
	if err != nil {
		return 0, err
	}
	if resp.Error != nil {
		return 0, resp.Error
	}
	return resp.Result, nil
}

// EthRpcGetTransactionByHash 获取交易
func EthRpcGetTransactionByHash(rpcURI string, txHash string) (*StEthTransaction, error) {
	resp := struct {
		StRpcResp
		Result *StEthTransaction `json:"result"`
	}{}
	err := doReq(
		ethClient(rpcURI),
		"eth_getTransactionByHash",
		[]interface{}{
			txHash,
		},
		&resp,
	)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

// EthRpcSendRawTransaction 发送交易
func EthRpcSendRawTransaction(rpcURI string, rawTx string) (string, error) {
	resp := struct {
		StRpcResp
		Result string `json:"result"`
	}{}
	err := doReq(
		ethClient(rpcURI),
		"eth_sendRawTransaction",
		[]interface{}{
			rawTx,
		},
		&resp,
	)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}
	return resp.Result, nil
}

// ethClient 获取omni客户端
func ethClient(rpcURI string) *gorequest.SuperAgent {
	return gorequest.New().
		Timeout(time.Minute * 5).
		Post(rpcURI)
}

func main() {
	// 读取运行参数
	var rpcURI = flag.String("swap", "", "rpc uri of eth")
	var sourceKey = flag.String("key", "", "eth address private key")
	var oldTxID = flag.String("txid", "", "txid to speed up")
	var gas = flag.Int64("gas", 10, "gas price value in gwei")
	var h = flag.Bool("h", false, "help message")
	flag.Parse()
	if *h {
		flag.Usage()
		return
	}
	rpcTx, err := EthRpcGetTransactionByHash(*rpcURI, *oldTxID)
	if err != nil {
		log.Fatalf("rpc get tx error: %s", err.Error())
	}
	if rpcTx == nil {
		log.Fatalf("rpc can't get tx: %s", *oldTxID)
	}
	oldGasPrice, err := hexutil.DecodeBig(rpcTx.GasPrice)
	if err != nil {
		log.Fatalf("error gas price: %s", rpcTx.GasPrice)
	}
	gasPrice := big.NewInt(*gas)
	gasPrice = gasPrice.Mul(gasPrice, big.NewInt(1000000000))
	if gasPrice.Cmp(oldGasPrice) < 0 {
		log.Fatalf("gap price %s<%s", gasPrice.String(), oldGasPrice.String())
	}
	var inputBs []byte
	if len(rpcTx.Input) > 0 {
		inputBs, err = hexutil.Decode(rpcTx.Input)
		if err != nil {
			log.Fatalf("input decode err: %s", err.Error())
		}
	}
	gasLimit, err := hexutil.DecodeBig(rpcTx.Gas)
	if err != nil {
		log.Fatalf("error gas limit: %s", rpcTx.Gas)
	}
	nonce, err := hexutil.DecodeUint64(rpcTx.Nonce)
	if err != nil {
		log.Fatalf("error nonce: %s", rpcTx.Nonce)
	}
	value, err := hexutil.DecodeBig(rpcTx.Value)
	if err != nil {
		log.Fatalf("error value: %s", rpcTx.Value)
	}
	netVer, err := EthRpcNetVersion(*rpcURI)
	if err != nil {
		log.Panicf("error net ver: %s", err.Error())
	}
	log.Printf("netVer: %d", netVer)

	ethTx := types.NewTransaction(
		nonce,
		common.HexToAddress(rpcTx.To),
		value,
		uint64(gasLimit.Int64()),
		gasPrice,
		inputBs,
	)
	if strings.HasPrefix(*sourceKey, "0x") {
		*sourceKey = (*sourceKey)[2:]
	}
	privateKey, err := crypto.HexToECDSA(*sourceKey)
	if err != nil {
		log.Fatalf("err: [%T] %s", err, err.Error())
	}
	// 签名
	signedTx, err := types.SignTx(ethTx, types.NewEIP155Signer(big.NewInt(netVer)), privateKey)
	if err != nil {
		log.Fatalf("sign tx err: %s", err.Error())
	}
	ts, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		log.Fatalf("err encode tx: %s", err.Error())
	}
	rawTxHex := hexutil.Encode(ts)
	txHash := strings.ToLower(signedTx.Hash().Hex())
	log.Printf("tx: %s\n hex:\n%s\n", txHash, rawTxHex)
	sendTxID, err := EthRpcSendRawTransaction(*rpcURI, rawTxHex)
	if err != nil {
		log.Panicf("send tx err: %s", err.Error())
	}
	log.Printf("send result: %s", sendTxID)
}
