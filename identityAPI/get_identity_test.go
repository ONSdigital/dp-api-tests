package identityAPI

import (
	"github.com/ONSdigital/dp-api-tests/identityAPIModels"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

var (
	defaultIdentity = &identityAPIModels.API{
		Name:              "John Paul Jones",
		Email:             "blackdog@ons.gov.uk",
		Password:          "foo",
		UserType:          "bar",
		TemporaryPassword: false,
		Migrated:          false,
		Deleted:           false,
	}
)

func TestGetIdentitySuccess(t *testing.T) {
	identityAPI := httpexpect.New(t, cfg.IdentityAPIURL)

	Convey("Given a request contains a token", t, func() {
		Convey("When a GET request is made to the API", func() {
			Convey("Return a default identity", func() {

				resp := identityAPI.GET("/identity", nil).
					WithHeader("token", "1234").
					Expect().
					Status(http.StatusOK).JSON().Object()

				So(resp.Value("name").String().Raw(), ShouldEqual, defaultIdentity.Name)
				So(resp.Value("email").String().Raw(), ShouldEqual, defaultIdentity.Email)
				So(resp.Value("password").String().Raw(), ShouldEqual, defaultIdentity.Password)
				So(resp.Value("user_type").String().Raw(), ShouldEqual, defaultIdentity.UserType)
				So(resp.Value("temporary_password").Boolean().Raw(), ShouldEqual, defaultIdentity.TemporaryPassword)
				So(resp.Value("migrated").Boolean().Raw(), ShouldEqual, defaultIdentity.Migrated)
				So(resp.Value("deleted").Boolean().Raw(), ShouldEqual, defaultIdentity.Deleted)
			})
		})
	})
}

func TestGetIdentityNoTokenError(t *testing.T) {
	identityAPI := httpexpect.New(t, cfg.IdentityAPIURL)

	Convey("Given a request that does not contain a token", t, func() {
		Convey("When a GET request is made to the API", func() {
			Convey("Return a 401 Status Unauthorised response.", func() {

				identityAPI.GET("/identity", nil).Expect().Status(http.StatusUnauthorized)

			})
		})
	})
}
