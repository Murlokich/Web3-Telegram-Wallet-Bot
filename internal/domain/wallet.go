package domain

import "github.com/pkg/errors"

var (
	ErrInvalidAddressIndex = errors.New("invalid address index")
)

type HDWallet struct {
	UserID                int64
	MasterKey             []byte
	AddressManagementData *AddressManagementData
}

type AddressManagementData struct {
	ChangeLevelKey      []byte
	CurrentAddressIndex uint32
	LastAddressIndex    uint32
}

func (d *AddressManagementData) ValidateAddressIndex(index uint32) error {
	if index > d.LastAddressIndex {
		return ErrInvalidAddressIndex
	}
	return nil
}
