package codeListAPI

import (
	"github.com/ONSdigital/dp-bolt/bolt"
	"fmt"
)

type codeListObject struct {
	codeList     string
	label        string
	edition      string
	id           string
	codeListLink string
	editionLink  string
	editionsLink string
	codesLink    string
}

var (
	gibsonGuitars2017 = codeListObject{
		codeList:     "gibson-guitars",
		label:        "gibson",
		edition:      "2017",
		id:           "2017",
		codeListLink: "(.*)\\/code-lists\\/gibson-guitars",
		editionsLink: "(.*)\\/code-lists\\/gibson-guitars/editions$",
		editionLink:  "(.*)\\/code-lists\\/gibson-guitars/editions/2017$",
		codesLink:    "(.*)\\/code-lists\\/gibson-guitars/editions/2017/codes$",
	}

	gibsonGuitars2018 = codeListObject{
		codeList:     "gibson-guitars",
		label:        "gibson",
		edition:      "2018",
		id:           "2018",
		codeListLink: "(.*)\\/code-lists\\/gibson-guitars",
		editionsLink: "(.*)\\/code-lists\\/gibson-guitars/editions$",
		editionLink:  "(.*)\\/code-lists\\/gibson-guitars/editions/2018$",
		codesLink:    "(.*)\\/code-lists\\/gibson-guitars/editions/2018/codes$",
	}

	fenderGuitars2017 = codeListObject{
		codeList:     "fender-guitars",
		label:        "fender",
		edition:      "2017",
		id:           "2017",
		codeListLink: "(.*)\\/code-lists\\/fender-guitars",
		editionsLink: "(.*)\\/code-lists\\/fender-guitars/editions$",
		editionLink:  "(.*)\\/code-lists\\/fender-guitars/editions/2017$",
		codesLink:    "(.*)\\/code-lists\\/fender-guitars/editions/2017/codes$",
	}

	fenderGuitars2018 = codeListObject{
		codeList:     "fender-guitars",
		label:        "fender",
		edition:      "2018",
		id:           "2018",
		codeListLink: "(.*)\\/code-lists\\/fender-guitars",
		editionsLink: "(.*)\\/code-lists\\/fender-guitars/editions$",
		editionLink:  "(.*)\\/code-lists\\/fender-guitars/editions/2018$",
		codesLink:    "(.*)\\/code-lists\\/fender-guitars/editions/2018/codes$",
	}

	tearItDown = bolt.Stmt{Query: "MATCH(n:`_api_test`) DETACH DELETE n", Params: nil}

	edition2017 = "2017"
	edition2018 = "2018"

	gibson2017 = []bolt.Stmt{
		{
			Query:  "CREATE (node:`_code_list`:`_code_list_gibson-guitars`:`_api_test` { label:'gibson', edition: {ed}})",
			Params: bolt.Params{"ed": edition2017},
		},
		{
			Query:  "MATCH (parent:`_code_list`:`_code_list_gibson-guitars` {edition: {ed}}) WITH parent CREATE (node:`_code`:`_api_test` { value: {v}})-[:usedBy { label: {l}}]->(parent)",
			Params: bolt.Params{"ed": edition2017, "v": "Les_Paul", "l": "Les Paul"},
		},
		{
			Query:  "MATCH (parent:`_code_list`:`_code_list_gibson-guitars` {edition: {ed}}) WITH parent CREATE (node:`_code`:`_api_test` { value: {v}})-[:usedBy { label: {l}}]->(parent)",
			Params: bolt.Params{"ed": edition2017, "v": "SG", "l": "SG"},
		},
		{
			Query:  "MATCH (parent:`_code_list`:`_code_list_gibson-guitars` {edition: {ed}}) WITH parent CREATE (node:`_code`:`_api_test` { value: {v}})-[:usedBy { label: {l}}]->(parent)",
			Params: bolt.Params{"ed": edition2017, "v": "Explorer", "l": "Explorer"},
		},
		{
			Query:  "MATCH (parent:`_code_list`:`_code_list_gibson-guitars` {edition: {ed}}) WITH parent CREATE (node:`_code`:`_api_test` { value: {v}})-[:usedBy { label: {l}}]->(parent)",
			Params: bolt.Params{"ed": edition2017, "v": "Flying_V", "l": "Flying V"},
		},
	}

	gibson2018 = []bolt.Stmt{
		{
			Query:  "CREATE (node:`_code_list`:`_code_list_gibson-guitars`:`_api_test` { label:'gibson', edition: {ed}})",
			Params: bolt.Params{"ed": edition2018}},
		{
			Query:  "MATCH (parent:`_code_list`:`_code_list_gibson-guitars` {edition: {ed}}) WITH parent CREATE (node:`_code`:`_api_test` { value: {v}})-[:usedBy { label: {l}}]->(parent)",
			Params: bolt.Params{"ed": edition2018, "v": "Les_Paul", "l": "Les Paul"},
		},
		{
			Query:  "MATCH (parent:`_code_list`:`_code_list_gibson-guitars` {edition: {ed}}) WITH parent CREATE (node:`_code`:`_api_test` { value: {v}})-[:usedBy { label: {l}}]->(parent)",
			Params: bolt.Params{"ed": edition2018, "v": "SG", "l": "SG"},
		},
		{
			Query:  "MATCH (parent:`_code_list`:`_code_list_gibson-guitars` {edition: {ed}}) WITH parent CREATE (node:`_code`:`_api_test` { value: {v}})-[:usedBy { label: {l}}]->(parent)",
			Params: bolt.Params{"ed": edition2018, "v": "Explorer", "l": "Explorer"},
		},
		{
			Query:  "MATCH (parent:`_code_list`:`_code_list_gibson-guitars` {edition: {ed}}) WITH parent CREATE (node:`_code`:`_api_test` { value: {v}})-[:usedBy { label: {l}}]->(parent)",
			Params: bolt.Params{"ed": edition2018, "v": "Flying_V", "l": "Flying V"},
		},
		{
			Query:  "MATCH (parent:`_code_list`:`_code_list_gibson-guitars` {edition: {ed}}) WITH parent CREATE (node:`_code`:`_api_test` { value: {v}})-[:usedBy { label: {l}}]->(parent)",
			Params: bolt.Params{"ed": edition2018, "v": "Firebird", "l": "Firebird"},
		},
	}

	fender2017 = []bolt.Stmt{
		{
			Query:  "CREATE (node:`_code_list`:`_code_list_fender-guitars`:`_api_test` { label:'fender', edition: {ed}})",
			Params: bolt.Params{"ed": edition2017}},
		{
			Query:  "MATCH (parent:`_code_list`:`_code_list_fender-guitars` {edition: {ed}}) WITH parent CREATE (node:`_code`:`_api_test` { value: {v}})-[:usedBy { label: {l}}]->(parent)",
			Params: bolt.Params{"ed": edition2017, "v": "Strat", "l": "Stratocaster"},
		},
		{
			Query:  "MATCH (parent:`_code_list`:`_code_list_fender-guitars` {edition: {ed}}) WITH parent CREATE (node:`_code`:`_api_test` { value: {v}})-[:usedBy { label: {l}}]->(parent)",
			Params: bolt.Params{"ed": edition2017, "v": "Tele", "l": "Telecaster"},
		},
	}

	fender2018 = []bolt.Stmt{
		{
			Query:  "CREATE (node:`_code_list`:`_code_list_fender-guitars`:`_api_test` { label:'fender', edition: {ed}})",
			Params: bolt.Params{"ed": edition2018}},
		{
			Query:  "MATCH (parent:`_code_list`:`_code_list_fender-guitars` {edition: {ed}}) WITH parent CREATE (node:`_code`:`_api_test` { value: {v}})-[:usedBy { label: {l}}]->(parent)",
			Params: bolt.Params{"ed": edition2018, "v": "Strat", "l": "Stratocaster"},
		},
		{
			Query:  "MATCH (parent:`_code_list`:`_code_list_fender-guitars` {edition: {ed}}) WITH parent CREATE (node:`_code`:`_api_test` { value: {v}})-[:usedBy { label: {l}}]->(parent)",
			Params: bolt.Params{"ed": edition2018, "v": "Tele", "l": "Telecaster"},
		},
		{
			Query:  "MATCH (parent:`_code_list`:`_code_list_fender-guitars` {edition: {ed}}) WITH parent CREATE (node:`_code`:`_api_test` { value: {v}})-[:usedBy { label: {l}}]->(parent)",
			Params: bolt.Params{"ed": edition2018, "v": "Jazzmaster", "l": "Jazzmaster"},
		},
	}
)

func (c codeListObject) CodeLink(code string) string {
	return fmt.Sprintf("(.*)\\/code-lists\\/%s/editions/%s/codes/%s$", c.codeList, c.edition, code)
}

func AllTestData() []bolt.Stmt {
	d := append(gibson2017, gibson2018...)
	d = append(d, fender2017...)
	return append(d, fender2018...)
}
