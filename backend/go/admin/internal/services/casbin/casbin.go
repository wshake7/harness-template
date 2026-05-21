package casbin

import (
	"strconv"
	"strings"

	"admin/internal/services/orm/query"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var E casbin.IEnforcer

var Adapter *gormadapter.Adapter

func New(db *gorm.DB) {
	var err error
	Adapter, err = gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(err)
	}
	sysCasbinModel := query.SysCasbinModel
	result, err := sysCasbinModel.Where(sysCasbinModel.IsEnabled.Is(true)).Select(sysCasbinModel.Content).First()
	if err != nil {
		panic(err)
	}
	m, err := model.NewModelFromString(normalizeModelContent(result.Content))
	if err != nil {
		panic(err)
	}
	E, err = casbin.NewSyncedCachedEnforcer(m, Adapter)
	if err != nil {
		panic(err)
	}
	E.EnableAutoSave(true)
	if err := normalizeLegacyPolicies(E); err != nil {
		panic(err)
	}
}

func normalizeModelContent(content string) string {
	content = strings.ReplaceAll(content, "p = sub_rule, obj_rule, act", "p = sub, obj, act")
	content = strings.ReplaceAll(content,
		"m = eval(p.sub_rule) && eval(p.obj_rule) && r.act == p.act",
		"m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && r.act == p.act",
	)
	return content
}

func normalizeLegacyPolicies(enforcer casbin.IEnforcer) error {
	policies, err := enforcer.GetPolicy()
	if err != nil {
		return err
	}
	for _, policy := range policies {
		if len(policy) < 3 {
			continue
		}
		sub, subOK := parseLegacySubjectRule(policy[0])
		obj, objOK := parseLegacyObjectRule(policy[1])
		if !subOK && !objOK {
			continue
		}
		next := append([]string(nil), policy...)
		if subOK {
			next[0] = sub
		}
		if objOK {
			next[1] = obj
		}
		if _, err := enforcer.RemovePolicy(policy); err != nil {
			return err
		}
		if _, err := enforcer.AddPolicy(next); err != nil {
			return err
		}
	}
	return nil
}

func parseLegacySubjectRule(value string) (string, bool) {
	value = strings.TrimSpace(value)
	const prefix = "r.sub == "
	if !strings.HasPrefix(value, prefix) {
		return "", false
	}
	sub, err := strconv.Unquote(strings.TrimSpace(strings.TrimPrefix(value, prefix)))
	if err != nil || sub == "" {
		return "", false
	}
	return sub, true
}

func parseLegacyObjectRule(value string) (string, bool) {
	value = strings.TrimSpace(value)
	const prefix = "keyMatch2(r.obj, "
	if !strings.HasPrefix(value, prefix) || !strings.HasSuffix(value, ")") {
		return "", false
	}
	obj, err := strconv.Unquote(strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(value, prefix), ")")))
	if err != nil || obj == "" {
		return "", false
	}
	return obj, true
}
