package codeListAPI

import (
	"encoding/json"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/dp-identity-api/api"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"strings"
	"testing"
)

var (
	newIdentity = identity.Model{
		Name:              "Peter Venkman",
		Email:             "venkman@whoyougunnacall.com",
		Deleted:           false,
		Migrated:          true,
		Password:          "There is no Dana only zuul!",
		TemporaryPassword: "",
		UserType:          "admin",
	}
)

func TestCreateIdentitySuccess(t *testing.T) {
	identityAPI := httpexpect.New(t, cfg.IdentityAPIURL)

	Convey("Given a valid identity", t, func() {
		Convey("When post request made to the API", func() {
			Convey("Then the identity is stored", func() {

				body, err := json.Marshal(newIdentity)
				So(err, ShouldBeNil)

				resp := identityAPI.POST("/identity", nil).
					WithHeader("Content-Type", "application/json").
					WithBytes(body).
					Expect().
					Status(http.StatusCreated)

				newID := resp.JSON().Object().Value("id").String().Raw()
				i, err := mongo.GetIdentity(cfg.MongoDB, collection, "id", newID)
				So(err, ShouldBeNil)

				So(i.ID, ShouldEqual, newID)
				So(i.Name, ShouldEqual, newIdentity.Name)
				So(i.Email, ShouldEqual, newIdentity.Email)

				tearDown(i.ID)
			})
		})
	})
}

func TestCreateIdentity_ValidationError(t *testing.T) {
	identityAPI := httpexpect.New(t, cfg.IdentityAPIURL)

	Convey("Given the request body is empty", t, func() {
		Convey("When post request is made to the API", func() {
			Convey("Then a Bad Request status is returned and no identity is stored", func() {
				resp := identityAPI.POST("/identity", nil).
					WithHeader("Content-Type", "application/json").
					Expect().
					Status(http.StatusBadRequest)

				So(getErrorBody(resp), ShouldEqual, api.ErrRequestBodyNil.Error())

				identities, err := mongo.GetIdentities(cfg.MongoDB, collection)
				So(err, ShouldBeNil)
				So(identities, ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given the identity.name is empty", t, func() {
		Convey("When post request is made to the API", func() {
			Convey("Then a Bad Request status is returned and no identity is stored", func() {
				b, err := json.Marshal(identity.Model{})
				So(err, ShouldBeNil)

				resp := identityAPI.POST("/identity", nil).
					WithHeader("Content-Type", "application/json").
					WithBytes(b).
					Expect().
					Status(http.StatusBadRequest)

				So(getErrorBody(resp), ShouldEqual, identity.ErrNameValidation.Error())

				identities, err := mongo.GetIdentities(cfg.MongoDB, collection)
				So(err, ShouldBeNil)
				So(identities, ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given the identity.email is empty", t, func() {
		Convey("When post request is made to the API", func() {
			Convey("Then a Bad Request status is returned and no identity is stored", func() {
				b, err := json.Marshal(identity.Model{Name: "Edmund Blackadder"})
				So(err, ShouldBeNil)

				resp := identityAPI.POST("/identity", nil).
					WithHeader("Content-Type", "application/json").
					WithBytes(b).
					Expect().
					Status(http.StatusBadRequest)

				So(getErrorBody(resp), ShouldEqual, identity.ErrEmailValidation.Error())

				identities, err := mongo.GetIdentities(cfg.MongoDB, collection)
				So(err, ShouldBeNil)
				So(identities, ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given the identity.password is empty", t, func() {
		Convey("When post request is made to the API", func() {
			Convey("Then a Bad Request status is returned and no identity is stored", func() {
				b, err := json.Marshal(identity.Model{Name: "Edmund Blackadder", Email: "captainB@thefrontline.com"})
				So(err, ShouldBeNil)

				resp := identityAPI.POST("/identity", nil).
					WithHeader("Content-Type", "application/json").
					WithBytes(b).
					Expect().
					Status(http.StatusBadRequest)

				So(getErrorBody(resp), ShouldEqual, identity.ErrPasswordValidation.Error())

				identities, err := mongo.GetIdentities(cfg.MongoDB, collection)
				So(err, ShouldBeNil)
				So(identities, ShouldHaveLength, 0)
			})
		})
	})
}

func getErrorBody(resp *httpexpect.Response) string {
	return strings.TrimSpace(resp.Body().Raw())
}

func tearDown(id string) {
	doc := &mongo.Doc{
		Database:   cfg.MongoDB,
		Collection: collection,
		Key:        "id",
		Value:      id,
	}

	err := mongo.Teardown(doc)
	if err != nil {
		log.ErrorC("failed to tear down identities docs", err, nil)
	}
}
