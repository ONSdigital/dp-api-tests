package importAPI

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-api-tests/config"
	"github.com/gavv/httpexpect"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPostJob_CreatesJob(t *testing.T) {

	importAPI := httpexpect.New(t, config.ImportAPIURL())

	jsonMap := make(map[string]interface{})
	json.Unmarshal([]byte(validJSON), &jsonMap)

	Convey("Given a valid json input to create a job", t, func() {

		Convey("The jobs endpoint returns 201 created", func() {

			r := importAPI.POST("/jobs").WithBytes([]byte(validJSON)).
				Expect().Status(http.StatusCreated).JSON().Object()

			fmt.Println(r)

			r.Value("recipe").Equal(jsonMap["recipe"])
			r.Value("recipe").Equal("b944be78-f56d-409b-9ebd-ab2b77ffe187")
			r.Value("state").Equal("created")

			r.Value("files").Array().Element(0).Object().Value("alias_name").Equal("v4")
			r.Value("files").Array().Element(0).Object().Value("url").Equal("https://s3-eu-west-1.amazonaws.com/dp-publish-content-test/OCIGrowth.csv")

			r.Value("links").Object().Value("instances").Array().Element(0).Object().Value("id").NotNull()
			r.Value("links").Object().Value("instances").Array().Element(0).Object().Value("href").NotNull()
		})

	})

}

func TestPostJob_InvalidInput(t *testing.T) {

	importAPI := httpexpect.New(t, config.ImportAPIURL())

	Convey("Given invalid json input to create a job", t, func() {

		Convey("The jobs endpoint returns 400 invalid json message ", func() {

			importAPI.POST("/jobs").WithBytes([]byte(invalidJSON)).
				Expect().Status(http.StatusBadRequest)
		})
	})
}
