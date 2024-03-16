package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/holiman/uint256"
	"github.com/labstack/echo/v4"
)

type MethodType = string

const (
	Mint     MethodType = "mint"
	Transfer MethodType = "transfer"
)

type Inscript struct {
	Method string `json:"method"`
	Ticker string `json:"ticker"`
	Img    string `json:"img"`
	Amount int    `json:"amount"`
	To     string `json:"to"`
}

var rpc = sepolia_rpc()

func mintHandler(c echo.Context) error {
	fmt.Println("go mint")
	// Read the request body
	var data Inscript
	// requestBody, err := io.ReadAll(c.Request().Body)
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Error reading request body: %s", err.Error()))
	}

	dataByte, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}

	// log.Println(string(dataByte))
	blobs, commitments, proofs, versionedHashes, err := EncodeBlobs(dataByte)
	if err != nil {
		log.Fatalf("failed to compute commitments: %v", err)
	}

	// send blob tx

	var privateKeyHex = privateKey()
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to convert private key: %v", err)
	}

	// send tx
	client, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Println("send from address: %v", fromAddress)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	chainID, err := client.NetworkID(context.Background())

	log.Println("chainid:%v", chainID.Int64())

	if err != nil {
		log.Fatal(err)
	}

	val, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Error getting suggested gas price: %v", err)
	}
	var nok bool
	gasPrice, nok := uint256.FromBig(val)
	if nok {
		log.Fatalf("gas price is too high! got %v", val.String())
	}

	priorityGasPrice256 := gasPrice

	tx := types.NewTx(&types.BlobTx{
		ChainID:    uint256.NewInt(11155111),
		Nonce:      uint64(nonce),
		GasTipCap:  priorityGasPrice256,
		GasFeeCap:  gasPrice.Add(uint256.NewInt(5e10), gasPrice),// gas price + 50gwei
		Gas:        231072,
		Value:      uint256.NewInt(0),
		Data:       nil,
		To:         common.Address{0x03, 0x04, 0x05},
		BlobFeeCap: uint256.NewInt(3e10), // 30gwei
		BlobHashes: versionedHashes,
		Sidecar:    &types.BlobTxSidecar{Blobs: blobs, Commitments: commitments, Proofs: proofs},
	})



	// dynaTx := types.NewTx(&types.DynamicFeeTx{
	// 	ChainID:    chainID,
	// 	Nonce:      uint64(nonce),
	// 	GasTipCap:  priorityGasPrice256.ToBig(),
	// 	GasFeeCap:  gasPrice.Add(uint256.NewInt(5000000000), gasPrice).ToBig(),
	// 	Gas:        231072,
	// 	Value: (big.NewInt(10000000000000000)),
	// 	To:         &common.Address{0x03, 0x04, 0x05},
	// })
	// signedTx, _ := types.SignTx(tx, types.NewCancunSigner(chainId), key)
	// err = client.SendTransaction(context.Background(), signedTx)

	signedTx, err := types.SignTx(tx, types.NewCancunSigner(chainID), privateKey)
	// signedTx, err := types.SignTx(dynaTx, types.NewCancunSigner(chainID), privateKey)

	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(signedTx.Hash())

	// check tx statuc
	// receipt, err := client.tran(context.Background(), signedTx.Hash())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Sprintf("receipt: %+v", receipt)

	return c.JSON(http.StatusOK, "go backend get mint data success fully")
}

// func mintHandler(c echo.Context) error {
// 	fmt.Println("go mint")
// 	// Read the request body
// 	var data Inscript
// 	// requestBody, err := io.ReadAll(c.Request().Body)
// 	err := c.Bind(&data)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Error reading request body: %s", err.Error()))
// 	}

// 	dataByte,err := json.Marshal(data)
// 	if err != nil {
// 		log.Fatalf("JSON marshaling failed: %s", err)
// 	}

// 	var b = Blob{}
// 	b.FromData(dataByte)

// 	// send blob tx

// 	var privateKeyHex = "ef4ea702da4db75fdaf15344ad1871865f22f80eaa2b99fcd03a4e2aecb788ce"
// 	privateKey, err := crypto.HexToECDSA(privateKeyHex)
// 	if err != nil {
// 		log.Fatalf("Failed to convert private key: %v", err)
// 	}

// 	// send tx
// 	client, err := ethclient.Dial(rpc)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	publicKey := privateKey.Public()
// 	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
// 	if !ok {
// 		log.Fatal("error casting public key to ECDSA")
// 	}

// 	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
// 	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	blobTx := createEmptyBlobTx(nonce,[]kzg4844.Blob{kzg4844.Blob(b)},privateKey, true)
// 	chainID, err := client.NetworkID(context.Background())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	signedTx, err := types.SignTx(blobTx, types.LatestSignerForChainID(chainID), privateKey)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	err = client.SendTransaction(context.Background(), signedTx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(signedTx.Hash())

// 	// check tx statuc
// 	// receipt, err := client.tran(context.Background(), signedTx.Hash())
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// fmt.Sprintf("receipt: %+v", receipt)

// 	return c.JSON(http.StatusOK, fmt.Sprintf("go backend get mint data success fully"))
// }

func transferHandler(c echo.Context) error {
	fmt.Println("go transfer")
	// Read the request body
	var data Inscript
	// requestBody, err := io.ReadAll(c.Request().Body)
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("Error reading request body: %s", err.Error()))
	}

	dataByte, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}

	// log.Println(string(dataByte))
	blobs, commitments, proofs, versionedHashes, err := EncodeBlobs(dataByte)
	if err != nil {
		log.Fatalf("failed to compute commitments: %v", err)
	}

	// send blob tx

	var privateKeyHex = privateKey()
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to convert private key: %v", err)
	}

	// send tx
	client, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Println("send from address:", fromAddress)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	chainID, err := client.NetworkID(context.Background())

	log.Println("chainid:", chainID.Int64())

	if err != nil {
		log.Fatal(err)
	}

	val, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Error getting suggested gas price: %v", err)
	}
	var nok bool
	gasPrice, nok := uint256.FromBig(val)
	if nok {
		log.Fatalf("gas price is too high! got %v", val.String())
	}

	priorityGasPrice256 := gasPrice

	tx := types.NewTx(&types.BlobTx{
		ChainID:    uint256.NewInt(11155111),
		Nonce:      uint64(nonce),
		GasTipCap:  priorityGasPrice256,
		GasFeeCap:  gasPrice,
		Gas:        231072,
		To:         common.Address{0x03, 0x04, 0x05},
		BlobFeeCap: uint256.NewInt(1000),
		BlobHashes: versionedHashes,
		Sidecar:    &types.BlobTxSidecar{Blobs: blobs, Commitments: commitments, Proofs: proofs},
	})
	// signedTx, _ := types.SignTx(tx, types.NewCancunSigner(chainId), key)
	// err = client.SendTransaction(context.Background(), signedTx)

	signedTx, err := types.SignTx(tx, types.NewCancunSigner(chainID), privateKey)

	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(signedTx.Hash())

	// check tx statuc
	// receipt, err := client.tran(context.Background(), signedTx.Hash())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Sprintf("receipt: %+v", receipt)

	return c.JSON(http.StatusOK, fmt.Sprintf("go backend get mint data success fully"))
}
