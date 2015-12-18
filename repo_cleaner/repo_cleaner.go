package repo_deleter

import (
	"github.com/bborbe/aptly_utils/version"
	"github.com/bborbe/log"
)

type DeletePackagesByKey func(url string, user string, password string, repo string, keys []string) error
type ListPackages func(url string, user string, password string, repo string) ([]map[string]string, error)

type RepoCleaner interface {
	CleanRepo(url string, user string, password string, repo string) error
}

type repoCleaner struct {
	deletePackagesByKey DeletePackagesByKey
	listPackages        ListPackages
}

var logger = log.DefaultLogger

func New(deletePackagesByKey DeletePackagesByKey, listPackages ListPackages) *repoCleaner {
	r := new(repoCleaner)
	r.deletePackagesByKey = deletePackagesByKey
	r.listPackages = listPackages
	return r
}

func (r *repoCleaner) CleanRepo(url string, user string, password string, repo string) error {
	logger.Debugf("clean repo: %s", repo)
	keys, err := r.findKeysToDelete(url, user, password, repo)
	if err != nil {
		return err
	}
	return r.deletePackagesByKey(url, user, password, repo, keys)
}

func (r *repoCleaner) findKeysToDelete(url string, user string, password string, repo string) ([]string, error) {
	logger.Debugf("find keys to delete repo: %s", repo)
	packages, err := r.listPackages(url, user, password, repo)
	if err != nil {
		return nil, err
	}
	return packagesToKeys(packages), nil
}

func packagesToKeys(packages []map[string]string) []string {
	latestVersions := make(map[string]map[string]string)
	var keys []string
	for _, currentPackage := range packages {
		name := currentPackage["Package"]
		if latestPackage, ok := latestVersions[name]; ok {
			var packageToDelete map[string]string
			logger.Tracef("compare %s < %s", currentPackage["Version"], latestPackage["Version"])
			if version.Less(version.Version(currentPackage["Version"]), version.Version(latestPackage["Version"])) {
				packageToDelete = currentPackage
			} else {
				logger.Tracef("set latest version %s %s", currentPackage["Package"], currentPackage["Version"])
				latestVersions[name] = currentPackage
				packageToDelete = latestPackage
			}
			keys = append(keys, packageToDelete["Key"])
			logger.Debugf("mark package %s %s to delete", packageToDelete["Package"], packageToDelete["Version"])
		} else {
			latestVersions[name] = currentPackage
		}
	}
	for _, p := range latestVersions {
		logger.Debugf("keep package %s %s to delete", p["Package"], p["Version"])
	}
	return keys
}
