package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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

	receiverAddrHex := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
	blobtx, err := CreateBlobTx(client, privateKeyHex, dataByte, receiverAddrHex)

	if err != nil {
		log.Fatal("error creating blob tx %v", err)
	}

	tx := types.NewTx(blobtx)
	signedTx, err := types.SignTx(tx, types.NewCancunSigner(blobtx.ChainID.ToBig()), privateKey)

	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(signedTx.Hash())

	return c.JSON(http.StatusOK, "go backend get mint data success fully")
}

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

	receiverAddrHex := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
	blobtx, err := CreateBlobTx(client, privateKeyHex, dataByte, receiverAddrHex)

	if err != nil {
		log.Fatal("error creating blob tx %v", err)
	}

	tx := types.NewTx(blobtx)
	signedTx, err := types.SignTx(tx, types.NewCancunSigner(blobtx.ChainID.ToBig()), privateKey)

	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(signedTx.Hash())

	return c.JSON(http.StatusOK, "go backend get mint data success fully")
}
