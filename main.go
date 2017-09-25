package main

import (
	"net/http"
	"log"
	"encoding/json"
	"fmt"
)

// RepoInfo contains the info the user will receive
type RepoInfo struct {
	Project 		string
	Owner			string
	Committer		string
	Commits			int
	Languages 		[]interface{}
}

// RepoStruct is used to decode json info about the repo
type RepoStruct struct {
	Title 			string			`json:"full_name"`
	Owner 			OwnerStruct		`json:"owner"`
}

// OwnerStruct is used to decode json info about the owner
type OwnerStruct struct {
	Login string `json:"login"`
}

// Contributor struct is used to decode json information about a contributor
type Contributor struct {
	Login string `json:"login"`
	Contributions int `json:"contributions"`
}

/* 	this function will request repo info of the repo specified by the url parameter
	from the github api and store it in the RepoInfo pointer called dest.*/
func getRepoInfo(url string, dest *RepoInfo) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 200 {
		var repoStruct RepoStruct
		json.NewDecoder(resp.Body).Decode(&repoStruct)
		dest.Project = repoStruct.Title
		dest.Owner = repoStruct.Owner.Login
	}
}

/* 	this function will request language info of the repo specified by the url parameter
	from the github api and store it in the RepoInfo pointer called dest.*/
func getLanguages(url string, dest *RepoInfo) {
	resp, err := http.Get(url + "/languages")
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 200 {
		langs := new(map[string]interface{})
		json.NewDecoder(resp.Body).Decode(&langs)

		for r := range *langs {
			dest.Languages = append(dest.Languages, r)
		}
	}
 }

/* 	this function will request contributor info of the repo specified by the url parameter
   	from the github api and will extract the top contributor and store his/her name as well
	as their commit count in the RepoInfo pointer called dest.*/
 func getTopContributor(url string, dest *RepoInfo) {
	 resp, err := http.Get(url  + "/contributors")
	 if err != nil {
		 log.Fatal(err)
	 }
	 if resp.StatusCode == 200 {
		 var contributors []Contributor
		 json.NewDecoder(resp.Body).Decode(&contributors)
		 dest.Committer = contributors[0].Login
		 dest.Commits = contributors[0].Contributions
	 }
 }

/* 	This is the main handler func. A url to the fetch info from the api of the user
	specified repo is created. Then the information is fetched and displayed as a respone
	to the user*/
func handlerFunc(w http.ResponseWriter, r *http.Request) {
	var info RepoInfo
	url := "https://api.github.com/repos" + r.URL.Path

	getRepoInfo(url, &info)
	getLanguages(url, &info)
	getTopContributor(url, &info)

	m, _ := json.MarshalIndent(info, "", "    ")
	fmt.Fprint(w, string(m))
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe("127.0.0.1:8081", nil)
}
