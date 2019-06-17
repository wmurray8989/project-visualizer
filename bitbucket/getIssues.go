package bitbucket

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/wmurray8989/project-visualizer/config"
)

type issueResultStruct struct {
	StartAt    int    `json:"startAt"`
	MaxResults int    `json:"maxResults"`
	Total      int    `json:"total"`
	Issues     []Card `json:"issues"`
}

// GetIssues queries bitbucket for a list of cards matching the given criteria
func GetIssues(conf *config.Config) issueResultStruct {
	client := http.Client{
		Timeout: time.Duration(time.Second * 10),
	}

	// Build jql
	jql := ""
	if len(conf.Epic) > 0 {
		jql = jql + "\"epic link\"=\"" + conf.Epic + "\""
	}
	if len(conf.Assignee) > 0 {
		if len(jql) > 0 {
			jql = jql + " AND "
		}
		jql = jql + "\"assignee\"=\"" + conf.Assignee + "\""
	}
	if len(conf.Status) > 0 {
		if len(jql) > 0 {
			jql = jql + " AND "
		}
		jql = jql + "\"status\"=\"" + conf.Status + "\""
	}
	if len(jql) == 0 {
		// Require that some filters be set
		return issueResultStruct{}
	}

	req, err := http.NewRequest("GET", "https://crosschx.atlassian.net/rest/api/3/search", nil)
	if err != nil {
		log.Print(err)
		return issueResultStruct{}
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(conf.Username, conf.Password)
	params := req.URL.Query()
	params.Add("jql", jql)
	req.URL.RawQuery = params.Encode()

	fmt.Println(req.URL.RequestURI())

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return issueResultStruct{}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return issueResultStruct{}
	}

	var results = issueResultStruct{}
	json.Unmarshal(body, &results)

	return results
}
