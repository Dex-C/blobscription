package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

func createEmptyBlobTx(nonce uint64, blobs []kzg4844.Blob, key *ecdsa.PrivateKey, withSidecar bool) *types.Transaction {
	blobtx := createEmptyBlobTxInner(nonce, blobs, withSidecar)
	signer := types.NewCancunSigner(uint256.NewInt(11155111).ToBig())
	return types.MustSignNewTx(key, signer, blobtx)
}

func createEmptyBlobTxInner(nonce uint64, blobs []kzg4844.Blob, withSidecar bool) *types.BlobTx {
	commitment, _ := kzg4844.BlobToCommitment(blobs[0])
	proof, _ := kzg4844.ComputeBlobProof(blobs[0], commitment)
	sidecar := &types.BlobTxSidecar{
		Blobs:       blobs,
		Commitments: []kzg4844.Commitment{commitment},
		Proofs:      []kzg4844.Proof{proof},
	}
	blobtx := &types.BlobTx{
		ChainID:    uint256.NewInt(11155111),
		Nonce:      nonce,
		BlobFeeCap: uint256.NewInt(5000000000),
		GasTipCap:  uint256.NewInt(2000000000),   // a.k.a. maxPriorityFeePerGas
		GasFeeCap:  uint256.NewInt(100000000000), // a.k.a. maxFeePerGas
		Gas:        131072,
		To:         common.Address{0x03, 0x04, 0x05},
		BlobHashes: sidecar.BlobHashes(),
	}
	if withSidecar {
		blobtx.Sidecar = sidecar
	}
	return blobtx
}

func EncodeBlobs(data []byte) ([]kzg4844.Blob, []kzg4844.Commitment, []kzg4844.Proof, []common.Hash, error) {
	var b Blob
	b.FromData(data)
	var (
		blobs           = []kzg4844.Blob{*b.KZGBlob()}
		commits         []kzg4844.Commitment
		proofs          []kzg4844.Proof
		versionedHashes []common.Hash
	)
	for _, blob := range blobs {
		commit, err := kzg4844.BlobToCommitment(blob)
		if err != nil {
			log.Fatalf("commit error:%v", err.Error())
			return nil, nil, nil, nil, err
		}
		commits = append(commits, commit)

		proof, err := kzg4844.ComputeBlobProof(blob, commit)
		if err != nil {
			log.Fatalf(err.Error())
			log.Fatalf("proof error:%v", err.Error())
			return nil, nil, nil, nil, err
		}
		proofs = append(proofs, proof)

		versionedHashes = append(versionedHashes, kZGToVersionedHash(commit))
	}
	return blobs, commits, proofs, versionedHashes, nil
}

func encodeBlobs(data []byte) []kzg4844.Blob {

	blobs := []kzg4844.Blob{{}}
	blobIndex := 0
	fieldIndex := -1
	for i := 0; i < len(data); i += 31 {
		fieldIndex++
		if fieldIndex == params.BlobTxFieldElementsPerBlob {
			blobs = append(blobs, kzg4844.Blob{})
			blobIndex++
			fieldIndex = 0
		}
		max := i + 31
		if max > len(data) {
			max = len(data)
		}
		copy(blobs[blobIndex][fieldIndex*32:], data[i:max])
	}
	return blobs
}

// kZGToVersionedHash implements kzg_to_versioned_hash from EIP-4844
func kZGToVersionedHash(kzg kzg4844.Commitment) common.Hash {
	h := sha256.Sum256(kzg[:])
	var blobCommitmentVersionKZG uint8 = 0x01
	h[0] = blobCommitmentVersionKZG

	return h
}
