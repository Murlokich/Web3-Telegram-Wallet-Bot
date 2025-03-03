package wallet //nolint:testpackage // I want to test unexported/private functions

// All test data from these tests was generated by the website: https://iancoleman.io/bip39/#english

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tyler-smith/go-bip32"
)

const (
	cycles                = 3
	wordSetLength         = 12
	privateKeyBytesLength = 82
)

func TestGenerateMnemonic(t *testing.T) {
	t.Parallel()
	for range cycles {
		mnemonic, err := generateMnemonic()
		require.NoError(t, err)
		words := strings.Split(mnemonic, " ")
		assert.Len(t, words, wordSetLength)
	}
}

func TestGenerateMasterKey(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name      string
		mnemonic  string
		masterKey string
	}{
		{
			name:      "correct1",
			mnemonic:  "custom enemy fuel drum fever involve final rule pipe border tuna nasty",
			masterKey: "xprv9s21ZrQH143K2afG55BsGK5scL5JNNKM4gSF1VEyB7c44G7B9RWSYTjTNC8Fj58P4YEG4gD6XFd2ig4kTn6gQdVirhjUwpAsG4BHS7RXUj3", //nolint:lll
		},
		{
			name:      "correct2",
			mnemonic:  "hedgehog hill orange glove occur ridge team before puzzle settle alpha divert",
			masterKey: "xprv9s21ZrQH143K2nv8wgBYtuGLdvpYaxZzcRPPsCEK8humTApEDT2HDGFS9mCfzkWQQzNSQk8dZ982GXLxqVU7x4zuV64LwgM8rmjcawtMT5M", //nolint:lll
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mk, err := generateMasterKey(tc.mnemonic)
			require.NoError(t, err)
			assert.Equal(t, tc.masterKey, mk.String())
		})
	}
}

func TestGenerateChangeLevelKey(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name           string
		masterKey      string
		changeLevelKey string
	}{
		{
			name:           "correct1",
			masterKey:      "xprv9s21ZrQH143K3BuCRtyBdPuXShHCjYAnJ9RiGvFA698MWr4zjTriXRXBQT5rk2BAhpshqr1KZDaoZHDS2nuAyytAx3xHty3EhUzLfXwMPq2", //nolint:lll
			changeLevelKey: "xprvA1fyHUhuNm38cXe8hR6mpqwsx2Vji79Ry9KXiW5rswCzyvkWiaPw4ruTufH2trX2odxNYBNoXEZpRA33gVd1SwUcgsaMSfnuiYXDGrYPCTH", //nolint:lll
		},
		{
			name:           "correct2",
			masterKey:      "xprv9s21ZrQH143K2a6r45RLbzvQSth9t4PYvrtSzrE7DKweTs1buEfAyaUjeLu5XBudQ3vusX3t6cvW73SJEj41fxgNtbUVVVruKRKQnv3nW9i", //nolint:lll
			changeLevelKey: "xprvA2BZAtHcXFLXojowrzxfbP4AXxgPyUa4KEdCpShfmsNJ8mFNnjH8PKGRTewkjHMriDcghLyKdtNjsFhJnfHUuu4mhc95M6bDXAkUUz4NKXa", //nolint:lll
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mk, err := bip32.B58Deserialize(tc.masterKey)
			require.NoError(t, err)
			changeLevelKey, err := generateChangeLevelKey(mk)
			require.NoError(t, err)
			assert.Equal(t, tc.changeLevelKey, changeLevelKey.String())
		})
	}
}

func TestGetAddress(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name           string
		changeLevelKey string
		addresses      []string
	}{
		{
			name:           "correct1",
			changeLevelKey: "xprvA1fyHUhuNm38cXe8hR6mpqwsx2Vji79Ry9KXiW5rswCzyvkWiaPw4ruTufH2trX2odxNYBNoXEZpRA33gVd1SwUcgsaMSfnuiYXDGrYPCTH", //nolint:lll
			addresses: []string{
				"0xe3e27A2506BA4215985Aa38E0332b4fBC9BB6D27",
				"0xDeF649F7Bc734E0A6b73e1E6e79e779E1a570793",
				"0xDe9bFb3e6B607697fDbb99B6C7dFd0FF4211d400",
			},
		},
		{
			name:           "correct2",
			changeLevelKey: "xprvA2BZAtHcXFLXojowrzxfbP4AXxgPyUa4KEdCpShfmsNJ8mFNnjH8PKGRTewkjHMriDcghLyKdtNjsFhJnfHUuu4mhc95M6bDXAkUUz4NKXa", //nolint:lll
			addresses: []string{
				"0xf8c8278356BE8DF8F3967355aF2587591bb6B32D",
				"0x9AE2ECD3c8FD4C9c459C8c574Ce3bfB0f99d1C45",
				"0xecb672967CBd130063425B52F96b633d8C65Ff81",
				"0x8dB8dBE916e4ce2e68A500bEF2808888FB021ec4",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			clk, err := bip32.B58Deserialize(tc.changeLevelKey)
			require.NoError(t, err)
			clkBytes, err := clk.Serialize()
			require.NoError(t, err)
			var computedAddress string
			for i := range tc.addresses {
				computedAddress, err = GetAddress(clkBytes, uint32(i))
				require.NoError(t, err)
				assert.Equal(t, tc.addresses[i], computedAddress)
			}
		})
	}
}

func TestGenerateAddress(t *testing.T) {
	t.Parallel()
	for range cycles {
		wallet, err := GenerateHDWallet()
		require.NoError(t, err)
		words := strings.Split(wallet.Mnemonic, " ")
		assert.Len(t, words, wordSetLength)
		assert.Len(t, wallet.MasterKey, privateKeyBytesLength)
		assert.Len(t, wallet.ChangeLevelKey, privateKeyBytesLength)
	}
}
