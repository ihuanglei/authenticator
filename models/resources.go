package models

import (
	"github.com/ihuanglei/authenticator/pkg/convert"
	"github.com/ihuanglei/authenticator/pkg/mapper/st"
	"github.com/simplexwork/common"
	"xorm.io/builder"
)

// GetResources 获取资源数据
func GetResources(query st.ResourceQuery) (int64, []*st.ResourceDto, error) {
	page := query.Page
	limit := query.Limit
	cond := builder.NewCond()
	if common.Trim(query.Name) != "" {
		cond = cond.And(builder.Like{"name", query.Name + "%"})
	}
	count, resources, err := getResources(cond, page, limit)
	if err != nil {
		return 0, nil, err
	}
	var resourceDtos = make([]*st.ResourceDto, len(resources))
	err = convert.Map(&resources, &resourceDtos)
	if err != nil {
		return 0, nil, err
	}
	return count, resourceDtos, nil
}

// GetResourcesByIDs 根据编号获取资源数据
func GetResourcesByIDs(ids []common.ID) ([]*st.ResourceDto, error) {
	cond := builder.In("id", ids)
	_, resources, err := getResources(cond, 0, len(ids))
	if err != nil {
		return nil, err
	}
	var resourceDtos = make([]*st.ResourceDto, len(resources))
	err = convert.Map(&resources, &resourceDtos)
	if err != nil {
		return nil, err
	}
	return resourceDtos, nil
}
