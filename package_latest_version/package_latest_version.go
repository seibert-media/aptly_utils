package package_latest_version

import (
	"fmt"
	"sort"

	"github.com/bborbe/aptly_utils/package_name"
	aptly_password "github.com/bborbe/aptly_utils/password"
	aptly_repository "github.com/bborbe/aptly_utils/repository"
	aptly_url "github.com/bborbe/aptly_utils/url"
	aptly_user "github.com/bborbe/aptly_utils/user"
	aptly_version "github.com/bborbe/aptly_utils/version"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type PackageVersions func(
	url aptly_url.Url,
	user aptly_user.User,
	password aptly_password.Password,
	repository aptly_repository.Repository,
	packageName package_name.PackageName) ([]aptly_version.Version, error)

type PackageLatestVersion interface {
	PackageLatestVersion(
		url aptly_url.Url,
		user aptly_user.User,
		password aptly_password.Password,
		repository aptly_repository.Repository,
		packageName package_name.PackageName) (*aptly_version.Version, error)
}

type packageLatestVersion struct {
	packageVersions PackageVersions
}

func New(packageVersions PackageVersions) *packageLatestVersion {
	p := new(packageLatestVersion)
	p.packageVersions = packageVersions
	return p
}

func (p *packageLatestVersion) PackageLatestVersion(
	url aptly_url.Url,
	user aptly_user.User,
	password aptly_password.Password,
	repository aptly_repository.Repository,
	packageName package_name.PackageName) (*aptly_version.Version, error) {
	logger.Debugf("PackageLatestVersion")
	var err error
	var versions []aptly_version.Version
	if versions, err = p.packageVersions(url, user, password, repository, packageName); err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("package %s not found", packageName)
	}
	sort.Sort(aptly_version.VersionByName(versions))
	latestVersion := versions[len(versions)-1]
	logger.Debugf("found latest version %v for package %s", latestVersion, packageName)
	return &latestVersion, nil
}