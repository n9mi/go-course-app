package test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/n9mi/go-course-app/internal/helper"
	"github.com/n9mi/go-course-app/internal/model"
	"github.com/stretchr/testify/require"
)

var (
	registerUrl = "http://127.0.0.1:5000/api/auth/register"
	loginUrl    = "http://127.0.0.1:5000/api/auth/login"
)

func TestRegister(t *testing.T) {
	testItems := map[string]TestSchema{
		"Auth_Register_User_OK": {
			"request_name":     "John Doe",
			"request_email":    "johndoe@mail.com",
			"request_password": "password",
			"expected_code":    http.StatusOK,
		},
		"Auth_Register_User_BAD_REQUEST_name": {
			"request_name":     "",
			"request_email":    "johndoe@mail.com",
			"request_password": "password",
			"expected_code":    http.StatusBadRequest,
		},
		"Auth_Register_User_BAD_REQUEST_email": {
			"request_name":     "John Doe",
			"request_email":    "thisisnotemail",
			"request_password": "password",
			"expected_code":    http.StatusBadRequest,
		},
		"Auth_Register_User_BAD_REQUEST_password": {
			"request_name":     "John Doe",
			"request_email":    "johndoe@mail.com",
			"request_password": "pwd",
			"expected_code":    http.StatusBadRequest,
		},
	}

	for testName, testProps := range testItems {
		t.Run(testName, func(t *testing.T) {
			requestBody :=
				fmt.Sprintf(`{"name": "%s", "email": "%s", "password": "%s"}`,
					testProps["request_name"], testProps["request_email"], testProps["request_password"])
			request := newRequest(fiber.MethodPost, registerUrl, requestBody)

			response, err := app.Test(request)
			require.Nil(t, err)

			responseBody, err := io.ReadAll(response.Body)
			require.Nil(t, err)

			testReponse := new(TestResponse[any])
			err = json.Unmarshal(responseBody, testReponse)
			require.Nil(t, err)

			require.Equal(t, response.StatusCode, testProps["expected_code"])
		})
	}
}

func TestLogin(t *testing.T) {
	testItems := map[string]TestSchema{
		"Auth_Login_Admin_OK": {
			"request_email":    validAdminData.Email,
			"request_password": "password",
			"expected_code":    http.StatusOK,
		},
		"Auth_Login_Admin_UNAUTHORIZED": {
			"request_email":    validAdminData.Email,
			"request_password": "wrongpassword",
			"expected_code":    http.StatusUnauthorized,
		},
		"Auth_Login_User_OK": {
			"request_email":    validUserData.Email,
			"request_password": "password",
			"expected_code":    http.StatusOK,
		},
		"Auth_Login_User_UNAUTHORIZED": {
			"request_email":    validUserData.Email,
			"request_password": "wrongpassword",
			"expected_code":    http.StatusUnauthorized,
		},
		"Auth_Login_Any_UNAUTHORIZED_email": {
			"request_email":    "randomemail@mail.com",
			"request_password": "randompassword",
			"expected_code":    http.StatusUnauthorized,
		},
		"Auth_Login_Any_BAD_REQUEST_email_1": {
			"request_email":    "notanemail",
			"request_password": "notapassword",
			"expected_code":    http.StatusBadRequest,
		},
		"Auth_Login_Any_BAD_REQUEST_email_2": {
			"request_email":    "",
			"request_password": "anypassword",
			"expected_code":    http.StatusBadRequest,
		},
		"Auth_Login_Any_BAD_REQUEST_password": {
			"request_email":    "randomemail@mail.com",
			"request_password": "",
			"expected_code":    http.StatusBadRequest,
		},
	}

	for testName, testProps := range testItems {
		t.Run(testName, func(t *testing.T) {
			requestBody :=
				fmt.Sprintf(`{"email": "%s", "password": "%s"}`,
					testProps["request_email"], testProps["request_password"])
			fmt.Println(requestBody)
			request := newRequest(fiber.MethodPost, loginUrl, requestBody)

			response, err := app.Test(request)
			require.Nil(t, err)

			responseBody, err := io.ReadAll(response.Body)
			require.Nil(t, err)

			testResponse := new(TestResponse[model.TokenResponse])
			err = json.Unmarshal(responseBody, testResponse)
			require.Nil(t, err)

			require.Equal(t, testProps["expected_code"], response.StatusCode)

			if response.StatusCode == fiber.StatusOK {
				if testProps["request_email"] == validAdminData.Email {
					fmt.Println("TEST RESPONSE", testResponse)
					validAdminData.Token = testResponse.Data.AccessToken
					require.NotEmpty(t, validAdminData.Token)

					authData, err := helper.VerifyAccessToken(context.Background(), viperConfig, redisClient, log,
						validAdminData.Token)
					require.Nil(t, err)

					validAdminData.ID = authData.ID
					validAdminData.Roles = authData.Roles
				} else if testProps["request_email"] == validUserData.Email {
					fmt.Println("TEST RESPONSE", testResponse)
					validUserData.Token = testResponse.Data.AccessToken
					require.NotEmpty(t, validUserData.Token)

					authData, err := helper.VerifyAccessToken(context.Background(), viperConfig, redisClient, log,
						validUserData.Token)
					require.Nil(t, err)

					validUserData.ID = authData.ID
					validUserData.Roles = authData.Roles
				}
			}
		})
	}
}
