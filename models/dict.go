package models

import (
	"github.com/ihuanglei/authenticator/pkg/convert"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/ihuanglei/authenticator/pkg/mapper/st"
	"github.com/simplexwork/common"
)

// CreateDict 创建字典
func CreateDict(dictDto *st.DictDto) error {
	dict := new(dict)
	err := convert.Map(&dictDto, dict)
	if err != nil {
		return err
	}
	return createDict(dict)
}

// UpdateDict 更新字典
func UpdateDict(dictID common.ID, dictDto *st.DictDto) error {
	if _, err := getDictByID(dictID); err != nil {
		return err
	}
	dict := new(dict)
	err := convert.Map(&dictDto, dict)
	if err != nil {
		return err
	}
	return updateDict(dictID, dict)
}

// UpdateOneDict 更新字典
func UpdateOneDict(dictDto *st.DictDto) error {
	if common.IsEmpty(dictDto.TP) || common.IsEmpty(dictDto.Cate) {
		return errors.ErrArgument
	}
	dict := new(dict)
	err := convert.Map(&dictDto, dict)
	if err != nil {
		return err
	}
	return updateOne(dict)
}

// DelDict 删除字典
func DelDict(dictID common.ID) error {
	if _, err := getDictByID(dictID); err != nil {
		return err
	}
	return delDict(dictID)
}

// GetOneDict .
func GetOneDict(cate, tp string) (*st.DictDto, error) {
	if common.IsEmpty(cate) || common.IsEmpty(tp) {
		return nil, errors.ErrArgument
	}
	dict, err := getOneDict(cate, tp)
	if err != nil {
		return nil, err
	}
	dictDto := new(st.DictDto)
	if err := convert.Map(&dict, dictDto); err != nil {
		return nil, err
	}
	return dictDto, nil
}

// GetDictByName .
func GetDictByName(name ...string) ([]*st.DictDto, error) {
	dicts, err := getDictByName(name...)
	if err != nil {
		return nil, err
	}
	var dictDtos = make([]*st.DictDto, len(dicts))
	if err = convert.Map(&dicts, &dictDtos); err != nil {
		return nil, err
	}
	return dictDtos, nil
}

// GetDictByCate .
func GetDictByCate(cate string) ([]*st.DictDto, error) {
	dicts, err := getDictByCate(cate)
	if err != nil {
		return nil, err
	}
	var dictDtos = make([]*st.DictDto, len(dicts))
	err = convert.Map(&dicts, &dictDtos)
	if err != nil {
		return nil, err
	}
	return dictDtos, nil
}

// GetDictByID .
func GetDictByID(dictID common.ID) (*st.DictDto, error) {
	dict, err := getDictByID(dictID)
	if err != nil {
		return nil, err
	}
	dictDto := new(st.DictDto)
	if err := convert.Map(&dict, dictDto); err != nil {
		return nil, err
	}
	return dictDto, nil
}
