package config

import (
	"github.com/casbin/casbin/v2"
	"github.com/sirupsen/logrus"
)

func NewEnforcer(log *logrus.Logger) *casbin.Enforcer {
	enforcer, err := casbin.NewEnforcer("./internal/casbin/model.conf", "./internal/casbin/policy.csv")
	if err != nil {
		log.Fatalf("Error load casbin configuration : %+v", err)
	}

	return enforcer
}
