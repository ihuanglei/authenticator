package models

import (
	"github.com/ihuanglei/authenticator/pkg/consts"
	"xorm.io/builder"
)

func getResources(cond builder.Cond, page, limit int) (int64, []*resource, error) {
	if limit <= 0 {
		limit = consts.PageSize
	}
	if page <= 0 {
		page = 1
	}
	start := (page - 1) * limit
	var resources = make([]*resource, 0)
	count, err := _Engine.Where(cond).Limit(limit, start).FindAndCount(&resources)
	if err != nil {
		return 0, nil, err
	}
	return count, resources, nil
}
