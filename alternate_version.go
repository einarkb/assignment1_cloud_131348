package main

import (
	"net/http"
	"log"
	"encoding/json"
	"io"
	"fmt"
)

type Payload struct {
	Project 				string
	Owner					string
	TopContributor			string
	TopContributorCount 	interface{}
	Languages 				[]string
}

func getResponse(url string) io.ReadCloser{
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Body
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	var payLoad Payload

	mResp := new(map[string]interface{})
	json.NewDecoder(getResponse("https://api.github.com/repos" + r.URL.Path)).Decode(mResp)
	payLoad.Project = (*mResp)["full_name"].(string)
	payLoad.Owner = (((*mResp)["owner"].(map[string]interface{}))["login"]).(string)

	var con []interface{}
	json.NewDecoder(getResponse((*mResp)["contributors_url"].(string))).Decode(&con)
	payLoad.TopContributor = (con[0].(map[string]interface{})["login"]).(string)
	payLoad.TopContributorCount = con[0].(map[string]interface{})["contributions"]

	langResp := new(map[string]interface{})
	json.NewDecoder(getResponse((*mResp)["languages_url"].(string))).Decode(langResp)
	for r := range *langResp {
		payLoad.Languages = append(payLoad.Languages, r)
	}

	m, _ := json.MarshalIndent(payLoad, "", "    ")
	fmt.Fprint(w, string(m))
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
