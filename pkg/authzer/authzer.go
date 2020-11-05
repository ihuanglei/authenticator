package authzer

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	"github.com/ihuanglei/authenticator/models"
)

const (
	_TableName = "at_rules"

	_ModelText = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && r.act == p.act || r.sub == "10000"
`
)

// NewAuthzer .
func NewAuthzer() *casbin.Enforcer {
	m, err := model.NewModelFromString(_ModelText)
	if err != nil {
		panic(err)
	}
	a, err := xormadapter.NewAdapterByEngineWithTableName(models.DefauleEngine(), _TableName)
	if err != nil {
		panic(err)
	}
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		panic(err)
	}
	return e
}
