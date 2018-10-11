package identityAPI

import (
	"github.com/ONSdigital/dp-api-tests/identityAPIModels"
	"github.com/ONSdigital/dp-api-tests/testDataSetup/mongo"
	"github.com/ONSdigital/dp-identity-api/schema"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"testing"
)

var (
	newIdentity = identityAPIModels.API{
		Name:              "Peter Venkman",
		Email:             "venkman@whoyougunnacall.com",
		Deleted:           false,
		Migrated:          true,
		Password:          "There is no Dana only zuul!",
		TemporaryPassword: false,
		UserType:          "admin",
	}
)

func TestCreateIdentitySuccess(t *testing.T) {
	identityAPI := httpexpect.New(t, cfg.IdentityAPIURL)

	Convey("Given a valid identity", t, func() {
		Convey("When post request made to the API", func() {
			Convey("Then the identity is stored", func() {
				resp := DoCreateIdentity(newIdentity, identityAPI)
				resp.Status(http.StatusCreated)

				newID := resp.JSON().Object().Value("id").String().Raw()
				i, err := mongo.GetIdentity(cfg.MongoDB, collection, "id", newID)
				So(err, ShouldBeNil)

				So(i.ID, ShouldEqual, newID)
				So(i.Name, ShouldEqual, newIdentity.Name)
				So(i.Email, ShouldEqual, newIdentity.Email)

				Convey("and the password is encrypted", func() {
					pwdErr := bcrypt.CompareHashAndPassword(i.HashedPassword(), []byte(newIdentity.Password))
					So(pwdErr, ShouldBeNil)
				})

				tearDown(i.ID)
			})
		})
	})
}

func TestCreateIdentity_EmailAlreadyAssociated(t *testing.T) {
	identityAPI := httpexpect.New(t, cfg.IdentityAPIURL)

	Convey("Given an email is already associated with an identity", t, func() {
		id := DoCreateIdentitySuccess(newIdentity, identityAPI)

		Convey("When a post request made to the API with the same email", func() {
			resp := DoCreateIdentity(newIdentity, identityAPI)

			Convey("Then a 409 status is returned", func() {
				resp.Status(http.StatusConflict)
				So(getErrorBody(resp), ShouldEqual, "active identity already exists with email")
			})
		})

		Reset(func() {
			tearDown(id)
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

				So(getErrorBody(resp), ShouldEqual, "error expected request body but was empty")

				identities, err := mongo.GetIdentities(cfg.MongoDB, collection)
				So(err, ShouldBeNil)
				So(identities, ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given the identity.name is empty", t, func() {
		Convey("When post request is made to the API", func() {
			Convey("Then a Bad Request status is returned and no identity is stored", func() {
				resp := DoCreateIdentity(identityAPIModels.API{}, identityAPI)
				resp.Status(http.StatusBadRequest)

				So(getErrorBody(resp), ShouldEqual, schema.ErrNameValidation.Error())

				identities, err := mongo.GetIdentities(cfg.MongoDB, collection)
				So(err, ShouldBeNil)
				So(identities, ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given the identity.email is empty", t, func() {
		Convey("When post request is made to the API", func() {
			Convey("Then a Bad Request status is returned and no identity is stored", func() {
				resp := DoCreateIdentity(identityAPIModels.API{Name: "Edmund Blackadder"}, identityAPI)
				resp.Status(http.StatusBadRequest)

				So(getErrorBody(resp), ShouldEqual, schema.ErrEmailValidation.Error())

				identities, err := mongo.GetIdentities(cfg.MongoDB, collection)
				So(err, ShouldBeNil)
				So(identities, ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given the identity.password is empty", t, func() {
		Convey("When post request is made to the API", func() {
			Convey("Then a Bad Request status is returned and no identity is stored", func() {
				resp := DoCreateIdentity(identityAPIModels.API{Name: "Edmund Blackadder", Email: "captainB@thefrontline.com"}, identityAPI)
				resp.Status(http.StatusBadRequest)

				So(getErrorBody(resp), ShouldEqual, schema.ErrPasswordValidation.Error())

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
