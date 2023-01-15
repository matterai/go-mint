package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	nft "example.com/nft"
)

func main() {
	client, err := ethclient.Dial(fmt.Sprintf("https://polygon-mumbai.g.alchemy.com/v2/%s", os.Getenv("ALCHEMY_KEY")))
	if err != nil {
		log.Fatal(err)
	}

	privateKeyECDSA, err := crypto.HexToECDSA(os.Getenv("WALLET_PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasLimit := uint64(3000000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println(nonce)
	log.Println(gasLimit)
	log.Println(gasPrice)
	log.Println(fromAddress)
	log.Println(chainID)

	contract_address := common.HexToAddress("0xA609EFd836D7f8Ed8bA4A0b319029562bAb902e0")

	// ctr.MintNFT(toAddress, "ipfs://QmTBwqTFJTCH6ZMo64YfYg5pXwt6FvCdfT3Mmc5LyLWfpS")
	instance, err := nft.NewNft(contract_address, client)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, chainID)

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice

	toAddress := common.HexToAddress("0x5c8F7E1B4Ebe1f48BB96E5a3c943CC87260Eb7eB")
	tx, err := instance.MintNFT(auth, toAddress, "ipfs://QmTBwqTFJTCH6ZMo64YfYg5pXwt6FvCdfT3Mmc5LyLWfpS")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", tx.Hash().Hex())
}
