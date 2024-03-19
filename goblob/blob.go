package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/holiman/uint256"
	blobutil "github.com/offchainlabs/nitro/util/blobs"
)

// create eip4844 sidecar from blobs
func CreateSidecarAndVersionedHashes(blobs *[]kzg4844.Blob) (*types.BlobTxSidecar, []common.Hash, error) {
	commitments, versionedHashes, err := blobutil.ComputeCommitmentsAndHashes(*blobs)
	if err != nil {
		return nil, nil, err
	}

	proofs, err := blobutil.ComputeBlobProofs(*blobs, commitments)
	if err != nil {
		return nil, nil, err
	}

	return &types.BlobTxSidecar{
		Blobs:       *blobs,
		Commitments: commitments,
		Proofs:      proofs,
	}, versionedHashes, nil
}

// send blob tx with data in the blob
func CreateBlobTx(ethClient *ethclient.Client, privateKeyHex string, data []byte, receiverAddrHex string) (*types.BlobTx, error) {

	blobs, err := blobutil.EncodeBlobs(data)
	if err != nil {
		return nil, err
	}

	sidecar, versionedHashes, err := CreateSidecarAndVersionedHashes(&blobs)
	if err != nil {
		return nil, err
	}

	// chainid
	chainid, err := ethClient.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	// nonce
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := ethClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	// suggested gas price
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	tx := &types.BlobTx{
		ChainID:    uint256.MustFromBig(chainid),
		Nonce:      nonce,
		GasTipCap:  uint256.NewInt(1e9),           // 1gwei
		GasFeeCap:  uint256.MustFromBig(gasPrice), // gas price + 50gwei
		Gas:        21000,
		Value:      uint256.NewInt(0),
		Data:       nil,
		To:         common.HexToAddress(receiverAddrHex),
		BlobFeeCap: uint256.NewInt(3e10), // 30gwei
		BlobHashes: versionedHashes,
		Sidecar:    sidecar,
	}

	return tx, nil

}
