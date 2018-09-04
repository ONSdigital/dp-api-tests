package codeListAPI

import (
	"encoding/json"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/dp-identity-api/identity"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
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
		Convey("When post is called", func() {
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
