package repo_lister

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	aptly_api "github.com/bborbe/aptly_utils/api"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

type ExecuteRequest func(req *http.Request) (resp *http.Response, err error)

type NewHttpRequestBuilder func(url string) http_requestbuilder.HttpRequestBuilder

type RepoLister interface {
	ListRepos(api aptly_api.Api) ([]map[string]string, error)
}

type repoVersion struct {
	executeRequest        ExecuteRequest
	newHttpRequestBuilder NewHttpRequestBuilder
}

var logger = log.DefaultLogger

func New(executeRequest ExecuteRequest, newHttpRequestBuilder NewHttpRequestBuilder) *repoVersion {
	p := new(repoVersion)
	p.newHttpRequestBuilder = newHttpRequestBuilder
	p.executeRequest = executeRequest
	return p
}

func (p *repoVersion) ListRepos(api aptly_api.Api) ([]map[string]string, error) {
	logger.Debugf("list repos")
	url := fmt.Sprintf("%s/api/repos", api.Url)
	requestbuilder := p.newHttpRequestBuilder(url)
	requestbuilder.AddBasicAuth(string(api.User), string(api.Password))
	requestbuilder.SetMethod("GET")
	requestbuilder.AddContentType("application/json")
	req, err := requestbuilder.Build()
	if err != nil {
		return nil, err
	}
	resp, err := p.executeRequest(req)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("request to %s failed with status %d", url, resp.StatusCode)
	}
	var repos []map[string]string
	err = json.Unmarshal(content, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}