package identityAPI

import (
	"encoding/json"
	"github.com/ONSdigital/dp-api-tests/identityAPIModels"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
)

//DoCreateIdentity make a successful create identity request
func DoCreateIdentitySuccess(i identityAPIModels.API, api *httpexpect.Expect) string {
	resp := DoCreateIdentity(i, api)
	resp.Status(http.StatusCreated)
	return resp.JSON().Object().Value("id").String().Raw()
}

//DoCreateIdentity make a create identity request
func DoCreateIdentity(i identityAPIModels.API, identityAPI *httpexpect.Expect) *httpexpect.Response {
	body, err := json.Marshal(i)
	So(err, ShouldBeNil)

	return identityAPI.POST("/identity", nil).
		WithHeader("Content-Type", "application/json").
		WithBytes(body).
		Expect()
}

// DoNewTokenRequestSuccess make a auth/new token request
func DoNewTokenRequest(r identityAPIModels.NewTokenRequest, api *httpexpect.Expect) *httpexpect.Response {
	body, err := json.Marshal(r)
	So(err, ShouldBeNil)

	return api.POST("/token", nil).
		WithHeader("Content-Type", "application/json").
		WithBytes(body).
		Expect()
}

// DoNewTokenRequestSuccess make a successful auth/new token request
func DoNewTokenRequestSuccess(r identityAPIModels.NewTokenRequest, api *httpexpect.Expect) string {
	body, err := json.Marshal(r)
	So(err, ShouldBeNil)

	resp := api.POST("/token", nil).
		WithHeader("Content-Type", "application/json").
		WithBytes(body).
		Expect().
		Status(http.StatusOK)

	return resp.JSON().Object().Value("token").String().Raw()
}

// DoGetIdentity make a Get Identity request
func DoGetIdentity(token string, api *httpexpect.Expect) *httpexpect.Response {
	return api.GET("/identity", nil).
		WithHeader("Content-Type", "application/json").
		WithHeader("token", token).
		Expect().
		Status(http.StatusOK)
}
