package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func newTestEnforcer(log *logrus.Logger) *casbin.Enforcer {
	enforcer, err := casbin.NewEnforcer("../internal/casbin/model.conf", "../internal/casbin/policy.csv")
	if err != nil {
		log.Fatalf("Error load casbin configuration : %+v", err)
	}

	return enforcer
}

func newRequest(method string, url string, requestBody string) *http.Request {
	request := httptest.NewRequest(method, url, strings.NewReader(requestBody))
	request.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return request
}

func newRequestWithToken(method string, url string, requestBody string, token string) *http.Request {
	request := newRequest(method, url, requestBody)
	bearer := fmt.Sprintf("Bearer %s", token)
	request.Header.Add("Authorization", bearer)

	return request
}
