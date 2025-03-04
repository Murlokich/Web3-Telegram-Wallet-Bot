package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"strconv"
	"strings"
	"unicode"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/sha3"
)

const (
	hardenedLevelStart = 2147483648 // 2**31 for hardened levels
	entropyBitSize     = 128        // for mnemonic of 12 words
)

var (
	ErrPublicKeyConversion = errors.New("cannot convert key to public key")
)

type HDWallet struct {
	Mnemonic       string
	MasterKey      []byte
	ChangeLevelKey []byte
}

func GenerateHDWallet() (*HDWallet, error) {
	var mnemonic string
	var wlt *HDWallet
	var err error
	for {
		mnemonic, err = generateMnemonic()
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate mnemonic")
		}
		wlt, err = DeriveWalletFromMnemonic(mnemonic)
		if err == nil {
			break
		}
		if !errors.Is(err, bip32.ErrInvalidPrivateKey) {
			return nil, errors.Wrap(err, "failed to derive wallet")
		}
	}
	return wlt, nil
}

func DeriveWalletFromMnemonic(mnemonic string) (*HDWallet, error) {
	masterKey, err := deriveMasterKey(mnemonic)
	if err != nil {
		return nil, errors.Wrap(err, "failed to derive master key")
	}
	changeLevelKey, err := deriveChangeLevelKey(masterKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to derive change level key")
	}
	mkBytes, err := masterKey.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "failed to serialize master key")
	}
	clkBytes, err := changeLevelKey.Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "failed to serialize change level key")
	}
	return &HDWallet{mnemonic, mkBytes, clkBytes}, nil
}

func generateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(entropyBitSize)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate entropy")
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate mnemonic")
	}
	return mnemonic, nil
}

func deriveMasterKey(mnemonic string) (*bip32.Key, error) {
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, errors.Wrap(err, "failed to derive master key from the seed")
	}
	return masterKey, nil
}

func deriveChangeLevelKey(masterKey *bip32.Key) (*bip32.Key, error) {
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

func GetAddress(changeLevelKeyBytes []byte, addressIndex uint32) (string, error) {
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
	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return "", ErrPublicKeyConversion
	}
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
