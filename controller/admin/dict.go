package admin

import (
	"github.com/ihuanglei/authenticator/models"
	"github.com/ihuanglei/authenticator/pkg/context"
	"github.com/ihuanglei/authenticator/pkg/convert"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/ihuanglei/authenticator/pkg/mapper/st"
)

// CreateDict 创建字典
// @tags 管理 - 字典管理
// @Summary 创建字典
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param name query string false "名称"
// @Param cate query string false "类型"
// @Param value query string false "值"
// @Param tp query string false "业务类型"
// @Router /admin/dict/create [post]
// @Security AdminKeyAuth
func CreateDict(form st.DictForm, ctx *context.Context) {
	dictDto := new(st.DictDto)
	err := convert.Map(&form, dictDto)
	if err != nil {
		ctx.Error(err)
		return
	}
	if err := models.CreateDict(dictDto); err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONEmpty()
}

// UpdateDict 修改字典
// @tags 管理 - 字典管理
// @Summary 修改字典
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param id path string true "字典编号"
// @Param name query string false "名称"
// @Param cate query string false "类型"
// @Param value query string false "值"
// @Param tp query string false "业务类型"
// @Router /admin/dict/{id}/update [post]
// @Security AdminKeyAuth
func UpdateDict(form st.DictForm, ctx *context.Context) {
	dictID := ctx.ParamsID("dictID")
	if dictID <= 0 {
		ctx.BadRequestByError(errors.ErrArgument)
		return
	}
	dictDto := new(st.DictDto)
	err := convert.Map(&form, dictDto)
	if err != nil {
		ctx.Error(err)
		return
	}
	if err := models.UpdateDict(dictID, dictDto); err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONEmpty()
}

// UpdateOneDict 唯一
func UpdateOneDict(form st.DictForm, ctx *context.Context) {
	dictDto := new(st.DictDto)
	err := convert.Map(&form, dictDto)
	if err != nil {
		ctx.Error(err)
		return
	}
	if err := models.UpdateOneDict(dictDto); err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONEmpty()
}

// GetOneDict .
func GetOneDict(ctx *context.Context) {
	cate := ctx.QueryTrim("cate")
	tp := ctx.QueryTrim("tp")
	dictDto, err := models.GetOneDict(cate, tp)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSON(dictDto)
}

// DelDict 删除字典
func DelDict(ctx *context.Context) {
	dictID := ctx.ParamsID("dictID")
	if dictID <= 0 {
		ctx.BadRequestByError(errors.ErrArgument)
		return
	}
	if err := models.DelDict(dictID); err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONEmpty()
}

// GetDictByCate 根据类型查询字典
func GetDictByCate(ctx *context.Context) {
	cate := ctx.QueryTrim("cate")
	dicts, err := models.GetDictByCate(cate)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSON(dicts)
}
