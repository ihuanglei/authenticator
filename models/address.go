package models

import (
	"github.com/ihuanglei/authenticator/pkg/convert"
	"github.com/ihuanglei/authenticator/pkg/mapper/st"
	"github.com/simplexwork/common"
)

// CreateAddress 新增地址
func CreateAddress(userID common.ID, addressDto st.AddressDto) error {
	address := new(userAddress)
	err := convert.Map(addressDto, address)
	if err != nil {
		return err
	}
	address.UserID = userID
	return createAddress(address)
}

// UpdateAddress 更新地址
func UpdateAddress(userID common.ID, addressID common.ID, addressDto st.AddressDto) error {
	_, err := getUserByID(userID)
	if err != nil {
		return err
	}
	oAddress, err := getAddressByID(addressID)
	if err != nil {
		return err
	}
	if oAddress.UserID != userID {
		return nil
	}
	address := new(userAddress)
	err = convert.Map(addressDto, address)
	if err != nil {
		return err
	}
	return updateAddress(addressID, address)
}

// DelAddress 删除地址
func DelAddress(userID common.ID, addressID common.ID) error {
	_, err := getUserByID(userID)
	if err != nil {
		return err
	}
	oAddress, err := getAddressByID(addressID)
	if err != nil {
		return err
	}
	if oAddress.UserID != userID {
		return nil
	}
	return delAddress(addressID)
}
