package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/sha3"
	"strconv"
	"strings"
	"unicode"
)

const (
	hardenedLevelStart = 2147483648 // 2**31 for hardened levels
)

func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate entropy")
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate mnemonic")
	}
	return mnemonic, nil
}

func GenerateMasterKey(mnemonic string) (*bip32.Key, error) {
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate master key from the seed")
	}
	return masterKey, nil
}

func GenerateChangeLevelKey(masterKey *bip32.Key) (*bip32.Key, error) {
	// Derive Change-Level Private Key`m/44'/60'/0'/0`
	changeLevelPath := []uint32{44 + hardenedLevelStart, 60 + hardenedLevelStart, 0 + hardenedLevelStart, 0}
	currentKey := masterKey
	var err error
	for _, level := range changeLevelPath {
		currentKey, err = currentKey.NewChildKey(level)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get new child key")
		}
	}
	return currentKey, nil
}

func GetAddress(changeLevelKeyHex string, addressIndex uint32) (string, error) {
	changeLevelKeyBytes, err := hex.DecodeString(changeLevelKeyHex)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode change level key")
	}
	changeLevelKey, err := bip32.Deserialize(changeLevelKeyBytes)
	if err != nil {
		return "", errors.Wrap(err, "failed to deserialize change level key")
	}
	// Derive Address-Level Private Key `m/44'/60'/0'/0/0
	addressLevelKey, err := changeLevelKey.NewChildKey(addressIndex)
	if err != nil {
		return "", errors.Wrap(err, "failed to get new child key")
	}
	privateKey, err := crypto.ToECDSA(addressLevelKey.Key)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert address to ECDSA")
	}
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	address := crypto.Keccak256(publicKeyBytes)[12:]
	addressFormatted := toCheckSumAddress(hex.EncodeToString(address))
	return addressFormatted, nil
}

func toCheckSumAddress(address string) string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(address))
	hashedAddress := hex.EncodeToString(hash.Sum(nil))
	checksumAddress := "0x"
	for i, symbol := range address {
		hashNibble, _ := strconv.ParseUint(string(hashedAddress[i]), 16, 4)
		if unicode.IsLetter(symbol) && hashNibble >= 8 {
			checksumAddress += strings.ToUpper(string(symbol))
		} else {
			checksumAddress += string(symbol)
		}
	}
	return checksumAddress
}
