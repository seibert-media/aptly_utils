package repo_deleter

import (
	"fmt"

	aptly_model "github.com/bborbe/aptly_utils/model"
	aptly_requestbuilder_executor "github.com/bborbe/aptly_utils/requestbuilder_executor"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type UnPublishRepo func(api aptly_model.Api, repository aptly_model.Repository, distribution aptly_model.Distribution) error

type RepoDeleter interface {
	DeleteRepo(api aptly_model.Api, repository aptly_model.Repository, distribution aptly_model.Distribution) error
}

type repoDeleter struct {
	unPublishRepo              UnPublishRepo
	buildRequestAndExecute     aptly_requestbuilder_executor.RequestbuilderExecutor
	httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider
}

func New(buildRequestAndExecute aptly_requestbuilder_executor.RequestbuilderExecutor, httpRequestBuilderProvider http_requestbuilder.HttpRequestBuilderProvider, unPublishRepo UnPublishRepo) *repoDeleter {
	p := new(repoDeleter)
	p.buildRequestAndExecute = buildRequestAndExecute
	p.httpRequestBuilderProvider = httpRequestBuilderProvider
	p.unPublishRepo = unPublishRepo
	return p
}

func (c *repoDeleter) DeleteRepo(api aptly_model.Api, repository aptly_model.Repository, distribution aptly_model.Distribution) error {
	logger.Debugf("DeleteRepo - repo: %s distribution: %s", repository, distribution)
	err := c.unPublishRepo(api, repository, distribution)
	if err != nil {
		//return err
	}
	return c.deleteRepo(api, repository)
}

func (p *repoDeleter) deleteRepo(api aptly_model.Api, repository aptly_model.Repository) error {
	logger.Debugf("deleteRepo - repo: %s", repository)
	requestbuilder := p.httpRequestBuilderProvider.NewHttpRequestBuilder(fmt.Sprintf("%s/api/repos/%s", api.ApiUrl, repository))
	requestbuilder.AddBasicAuth(string(api.ApiUsername), string(api.ApiPassword))
	requestbuilder.SetMethod("DELETE")
	return p.buildRequestAndExecute.BuildRequestAndExecute(requestbuilder)
}
