package models

import (
	"time"

	"github.com/ihuanglei/authenticator/pkg/consts"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/simplexwork/common"
	"xorm.io/builder"
)

// 删除字典
func delDict(dictID common.ID) error {
	dict := &dict{UpdateTime: common.Now(), Status: consts.Delete}
	session := _Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	if _, err := session.Where("dict_id = ?", dictID).Update(dict); err != nil {
		return err
	}
	return session.Commit()
}

// 更新字典
func createDict(dict *dict) error {
	dictID, err := _IDWorker.Next()
	if err != nil {
		return err
	}
	dict.DictID = dictID
	dict.Status = consts.Normal
	dict.CreateTime = common.Now()
	dict.UpdateTime = common.DateTime(time.Date(1, 1, 1, 0, 0, 0, 0, time.Local))

	session := _Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	if _, err := session.Insert(dict); err != nil {
		return err
	}
	return session.Commit()
}

// 更新字典
func updateDict(dictID common.ID, dict *dict) error {
	dict.UpdateTime = common.Now()
	session := _Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	if _, err := session.Cols("cate", "name", "value", "tp", "update_time").Where("dict_id = ?", dictID).Update(dict); err != nil {
		return err
	}
	return session.Commit()
}

//  根据ID获取内容
func getDictByID(dictID common.ID) (*dict, error) {
	var dict dict
	has, err := _Engine.Where("dict_id = ? AND status = ?", dictID, consts.Normal).Get(&dict)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.ErrDictNotFound
	}
	return &dict, nil
}

//  根据名称获取内容
func getDictByName(name ...string) ([]*dict, error) {
	return getDict(builder.Eq{"status": consts.Normal}.And(builder.In("name", name)))
}

//  根据类型获取内容
func getDictByCate(cate string) ([]*dict, error) {
	return getDict(builder.Eq{"cate": cate, "status": consts.Normal})
}

// 获取字段
func getDict(cond builder.Cond) ([]*dict, error) {
	var dicts = make([]*dict, 0)
	err := _Engine.Where(cond).Find(&dicts)
	if err != nil {
		return nil, err
	}
	return dicts, nil
}

func getOneDict(cate, tp string) (*dict, error) {
	dicts, err := getDict(builder.Eq{"cate": cate, "tp": tp, "status": consts.Normal})
	if err != nil {
		return nil, err
	}
	if len(dicts) == 0 {
		return nil, errors.ErrDictNotFound
	}
	return dicts[0], nil
}

// 根据cate和name,指定唯一数据进行更新或新增
func updateOne(dict *dict) error {
	d, err := getOneDict(dict.Cate, dict.TP)
	if e, ok := err.(errors.Error); ok && e.Code() == errors.ErrDictNotFound.Code() {
		return createDict(dict)
	} else if err != nil {
		return err
	}
	return updateDict(d.DictID, dict)
}
