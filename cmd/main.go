package main

import (
	"context"
	"flag"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

func main() {
	// 读取运行参数
	var rpcURI = flag.String("swap", "", "rpc uri of eth")
	var sourceKey = flag.String("key", "", "eth address private key")
	var oldTxID = flag.String("txid", "", "txid to speed up")
	var gas = flag.Int64("gas", 10, "gas price value in gwei")
	var newGasLimit = flag.Uint64("limit", 0, "gas limit of tx\ndefault is 0, when limit is 0, it will keep the old gas limit in the original tx")
	var h = flag.Bool("h", false, "help message")
	flag.Parse()
	if *h {
		flag.Usage()
		return
	}
	if len(*rpcURI) == 0 || len(*sourceKey) == 0 || len(*oldTxID) == 0 {
		flag.Usage()
		return
	}
	client, err := ethclient.Dial(*rpcURI)
	if err != nil {
		log.Fatalf("eth client dial error: [%T] %s", err, err.Error())
	}
	rpcTx, isPending, err := client.TransactionByHash(
		context.Background(),
		common.HexToHash(*oldTxID),
	)
	if err != nil {
		log.Fatalf("rpc get tx error: %s", err.Error())
	}
	if rpcTx == nil {
		log.Fatalf("rpc can't get tx: %s", *oldTxID)
	}
	if !isPending {
		log.Fatal("rpc get tx is not pending")
	}
	oldGasPrice := rpcTx.GasPrice()
	gasPrice := big.NewInt(*gas)
	gasPrice = gasPrice.Mul(gasPrice, big.NewInt(1000000000))
	if gasPrice.Cmp(oldGasPrice) < 0 {
		log.Fatalf("gap price %s<%s", gasPrice.String(), oldGasPrice.String())
	}
	var inputBs []byte
	if len(rpcTx.Data()) > 0 {
		inputBs = rpcTx.Data()
	}
	gasLimit := rpcTx.Gas()
	if gasLimit < *newGasLimit {
		gasLimit = *newGasLimit
	}
	nonce := rpcTx.Nonce()
	value := rpcTx.Value()
	netVer, err := client.NetworkID(context.Background())
	if err != nil {
		log.Panicf("error net ver: %s", err.Error())
	}
	log.Printf("netVer: %d", netVer)

	ethTx := types.NewTransaction(
		nonce,
		*rpcTx.To(),
		value,
		gasLimit,
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
	signedTx, err := types.SignTx(ethTx, types.NewEIP155Signer(netVer), privateKey)
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
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Panicf("send tx err: %s", err.Error())
	}
	log.Printf("send result: %s", txHash)
}
