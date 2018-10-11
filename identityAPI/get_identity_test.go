package identityAPI

import (
	models "github.com/ONSdigital/dp-api-tests/identityAPIModels"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestGetIdentitySuccess(t *testing.T) {
	identityAPI := httpexpect.New(t, cfg.IdentityAPIURL)

	Convey("Given a valid auth token", t, func() {
		newID := DoCreateIdentitySuccess(newIdentity, identityAPI)

		r := models.NewTokenRequest{
			Email:    newIdentity.Email,
			Password: newIdentity.Password,
		}

		token := DoNewTokenRequestSuccess(r, identityAPI)

		Convey("When get request made to the API", func() {

			Convey("Then the expected identity is returned", func() {
				getIdentityResp := DoGetIdentity(token, identityAPI)

				respJSON := getIdentityResp.JSON().Object()

				So(respJSON.Value("id").String().Raw(), ShouldEqual, newID)
				So(respJSON.Value("email").String().Raw(), ShouldEqual, newIdentity.Email)
				So(respJSON.Value("user_type").String().Raw(), ShouldEqual, newIdentity.UserType)
				So(respJSON.Value("deleted").Boolean().Raw(), ShouldBeFalse)
				So(int64(respJSON.Value("token_ttl").Number().Raw()), ShouldEqual, time.Duration(time.Minute * 15).Nanoseconds())

				expected, err := mongo.GetIdentity("test", "identities", "id", newID)
				So(err, ShouldBeNil)

				created, err := time.Parse(time.RFC3339, respJSON.Value("created_date").String().Raw())
				So(expected.CreatedDate, ShouldEqual, created)
			})
		})

		Reset(func() {
			tearDown(newID)
		})
	})
}
