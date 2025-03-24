package bip32adapter

import (
	"Web3-Telegram-Wallet-Bot/internal/domain"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"strconv"
	"strings"
	"unicode"

	"go.opentelemetry.io/otel/trace"

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

type BIP32Adapter struct {
	tracer trace.Tracer
}

func New(tracer trace.Tracer) *BIP32Adapter {
	return &BIP32Adapter{tracer: tracer}
}

func (a *BIP32Adapter) GenerateHDWallet(ctx context.Context, userID int64) (*domain.HDWallet, string, error) {
	ctx, span := a.tracer.Start(ctx, "GenerateHDWallet")
	defer span.End()
	var mnemonic string
	var wlt *domain.HDWallet
	var err error
	for {
		mnemonic, err = a.generateMnemonic()
		if err != nil {
			err = errors.Wrap(err, "failed to generate mnemonic")
			span.RecordError(err)
			return nil, "", err
		}
		wlt, err = a.DeriveWalletFromMnemonic(ctx, mnemonic, userID)
		if err == nil {
			break
		}
		if !errors.Is(err, bip32.ErrInvalidPrivateKey) {
			err = errors.Wrap(err, "failed to derive wallet")
			span.RecordError(err)
			return nil, "", err
		}
	}
	return wlt, mnemonic, nil
}

func (a *BIP32Adapter) DeriveWalletFromMnemonic(
	ctx context.Context, mnemonic string, userID int64) (*domain.HDWallet, error) {
	_, span := a.tracer.Start(ctx, "DeriveWalletFromMnemonic")
	defer span.End()
	masterKey, err := a.deriveMasterKey(mnemonic)
	if err != nil {
		err = errors.Wrap(err, "failed to derive master key")
		span.RecordError(err)
		return nil, err
	}
	changeLevelKey, err := a.deriveChangeLevelKey(masterKey)
	if err != nil {
		err = errors.Wrap(err, "failed to derive change level key")
		span.RecordError(err)
		return nil, err
	}
	mkBytes, err := masterKey.Serialize()
	if err != nil {
		err = errors.Wrap(err, "failed to serialize master key")
		span.RecordError(err)
		return nil, err
	}
	clkBytes, err := changeLevelKey.Serialize()
	if err != nil {
		err = errors.Wrap(err, "failed to serialize change level key")
		span.RecordError(err)
		return nil, err
	}
	return &domain.HDWallet{UserID: userID, MasterKey: mkBytes,
		AddressManagementData: &domain.AddressManagementData{ChangeLevelKey: clkBytes}}, nil
}

func (a *BIP32Adapter) generateMnemonic() (string, error) {
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

func (a *BIP32Adapter) deriveMasterKey(mnemonic string) (*bip32.Key, error) {
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, errors.Wrap(err, "failed to derive master key from the seed")
	}
	return masterKey, nil
}

func (a *BIP32Adapter) deriveChangeLevelKey(masterKey *bip32.Key) (*bip32.Key, error) {
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

func (a *BIP32Adapter) GetAddress(
	ctx context.Context, changeLevelKeyBytes []byte, addressIndex uint32) (string, error) {
	_, span := a.tracer.Start(ctx, "GetAddress")
	defer span.End()
	changeLevelKey, err := bip32.Deserialize(changeLevelKeyBytes)
	if err != nil {
		err = errors.Wrap(err, "failed to deserialize change level key")
		span.RecordError(err)
		return "", err
	}
	// Derive Address-Level Private Key `m/44'/60'/0'/0/0
	addressLevelKey, err := changeLevelKey.NewChildKey(addressIndex)
	if err != nil {
		err = errors.Wrap(err, "failed to get new child key")
		span.RecordError(err)
		return "", err
	}
	privateKey, err := crypto.ToECDSA(addressLevelKey.Key)
	if err != nil {
		err = errors.Wrap(err, "failed to convert address to ECDSA")
		span.RecordError(err)
		return "", err
	}
	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		err = ErrPublicKeyConversion
		span.RecordError(err)
		return "", err
	}
	publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	address := crypto.Keccak256(publicKeyBytes)[12:]
	addressFormatted := a.toCheckSumAddress(hex.EncodeToString(address))
	return addressFormatted, nil
}

func (a *BIP32Adapter) toCheckSumAddress(address string) string {
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
