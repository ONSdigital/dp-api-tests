package codeListAPI

import (
	"encoding/json"
	models "github.com/ONSdigital/dp-api-tests/identityAPIModels"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"strings"
	"testing"
)

func TestNewToken_Success(t *testing.T) {
	identityAPI := httpexpect.New(t, cfg.IdentityAPIURL)

	Convey("Given an identity exists", t, func() {
		identityID := createIdentity(identityAPI)

		Convey("When a newTokenRequest made with the correct password", func() {
			authRequest := models.NewTokenRequest{
				Email:       newIdentity.Email,
				Password: newIdentity.Password,
			}

			body, err := json.Marshal(authRequest)
			So(err, ShouldBeNil)

			resp := identityAPI.POST("/token", nil).
				WithHeader("Content-Type", "application/json").
				WithBytes(body).
				Expect()

			Convey("Then a 200 status and auth token are returned", func() {
				resp.Status(http.StatusOK)

				token := resp.JSON().Object().Value("token").String().Raw()
				So(token, ShouldNotBeEmpty)
			})
		})

		Reset(func() {
			tearDown(identityID)
		})
	})
}

func TestNewToken_PasswordIncorrect(t *testing.T) {
	identityAPI := httpexpect.New(t, cfg.IdentityAPIURL)

	Convey("Given an identity exists", t, func() {
		identityID := createIdentity(identityAPI)

		Convey("When a newTokenRequest is made with an incorrect password", func() {
			authRequest := models.NewTokenRequest{
				Email:       newIdentity.Email,
				Password: "this password is incorrect",
			}

			body, err := json.Marshal(authRequest)
			So(err, ShouldBeNil)

			resp := identityAPI.POST("/token", nil).
				WithHeader("Content-Type", "application/json").
				WithBytes(body).
				Expect()

			Convey("Then a 403 status is returned", func() {
				resp.Status(http.StatusForbidden)
				So("authentication unsuccessful", ShouldEqual, strings.TrimSpace(resp.Body().Raw()))
			})
		})

		Reset(func() {
			tearDown(identityID)
		})
	})
}

func TestNewToken_UserNotFound(t *testing.T) {
	identityAPI := httpexpect.New(t, cfg.IdentityAPIURL)

	Convey("Given no user exists with the provided identity ID", t, func() {

		Convey("When a newTokenRest is made", func() {
			authRequest := models.NewTokenRequest{
				Email:       newIdentity.Email,
				Password: newIdentity.Password,
			}

			body, err := json.Marshal(authRequest)
			So(err, ShouldBeNil)

			resp := identityAPI.POST("/token", nil).
				WithHeader("Content-Type", "application/json").
				WithBytes(body).
				Expect()

			Convey("Then a 404 status is returned", func() {
				resp.Status(http.StatusNotFound)
				So("authentication unsuccessful user not found", ShouldEqual, strings.TrimSpace(resp.Body().Raw()))
			})
		})
	})
}

func createIdentity(identityAPI *httpexpect.Expect) string {
	body, err := json.Marshal(newIdentity)
	So(err, ShouldBeNil)

	resp := identityAPI.POST("/identity", nil).
		WithHeader("Content-Type", "application/json").
		WithBytes(body).
		Expect().
		Status(http.StatusCreated)

	return resp.JSON().Object().Value("id").String().Raw()
}
