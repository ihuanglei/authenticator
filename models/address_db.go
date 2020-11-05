package models

import (
	"github.com/ihuanglei/authenticator/pkg/consts"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/simplexwork/common"
	"xorm.io/builder"
)

func getAddresses(cond builder.Cond, page, limit int) (int64, []*userAddress, error) {
	if limit <= 0 {
		limit = consts.PageSize
	}
	if page <= 0 {
		page = 1
	}
	start := (page - 1) * limit
	var addresses = make([]*userAddress, 0)
	count, err := _Engine.Desc("create_time").Where(cond).Limit(limit, start).FindAndCount(&addresses)
	if err != nil {
		return 0, nil, err
	}
	return count, addresses, nil
}

// 根据地址获取
func getAddressByID(addressID common.ID) (*userAddress, error) {
	address := new(userAddress)
	has, err := _Engine.Where("address_id = ?", addressID).Get(address)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.ErrAddressNotFound
	}
	return address, nil
}

// 新增地址
func createAddress(address *userAddress) error {
	addressID, err := _IDWorker.Next()
	if err != nil {
		return err
	}
	address.AddressID = addressID
	address.CreateTime = common.Now()
	address.UpdateTime = address.CreateTime
	address.Status = consts.Normal
	session := _Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	_, err = session.Insert(address)
	if err != nil {
		return err
	}
	return session.Commit()
}

//  更新地址
func updateAddress(addressID common.ID, address *userAddress) error {
	address.UpdateTime = common.Now()
	session := _Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	_, err := session.Cols("name", "mobile", "province", "city", "county", "address", "zip", "update_time").Where("address_id = ?", addressID).Update(address)
	if err != nil {
		return err
	}
	return session.Commit()
}

// 删除地址
func delAddress(addressID common.ID) error {
	address := &userAddress{Status: consts.Delete, UpdateTime: common.Now()}
	session := _Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	_, err := session.Cols("status", "update_time").Where("address_id = ?", addressID).Update(address)
	if err != nil {
		return err
	}
	return session.Commit()
}
