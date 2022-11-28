package commons

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/apparentlymart/go-versions/versions"
	"github.com/apparentlymart/go-versions/versions/constraints"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/configs"
	httputils "github.com/tenable/terrascan/pkg/utils/http"
	"go.uber.org/zap"
)

const (
	apiVersion             = "v1"
	terraformVersionHeader = "X-Terraform-Version"
)

// VersionConstraints ...
type VersionConstraints = constraints.IntersectionSpec

// Requirements ...
type Requirements map[addrs.Provider]VersionConstraints

// ProviderVersions ...
type ProviderVersions struct {
	Versions []struct {
		Version   string   `json:"version"`
		Protocols []string `json:"protocols"`
	} `json:"versions"`
	Warnings []string `json:"warnings"`
}

// ProviderVersionList fetches all the versions of terraform providers
func ProviderVersionList(ctx context.Context, addr addrs.Provider, terraformVersion string) (versions.List, []string, error) {
	zap.S().Debugf("fetching list of providers metadata, hostname: %q, type: %q, namespace: %q, ", addr.Hostname.String(), addr.Namespace, addr.Type)

	if addr.Hostname.String() == "" {
		return nil, nil, fmt.Errorf("error preparing the providers list endpoint, error: hostname can't be empty")
	}

	endpointURL, err := url.Parse(path.Join(apiVersion, "providers", addr.Namespace, addr.Type, "versions"))
	if err != nil {
		return nil, nil, fmt.Errorf("error preparing the providers list endpoint, error: %s", err.Error())
	}

	endpointURL.Host = addr.Hostname.String()
	endpointURL.Scheme = "https"

	headers := http.Header{
		terraformVersionHeader: []string{terraformVersion},
	}

	resp, err := httputils.SendRequest("GET", endpointURL.String(), "", nil, headers)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// Great!
	case http.StatusNotFound:
		return nil, nil, fmt.Errorf("provider registry %s does not have a provider named %s", addr.Hostname.ForDisplay(), addr.Type)
	default:
		return nil, nil, fmt.Errorf("could not query provider registry for %s: %s", addr.String(), resp.Status)
	}

	var body ProviderVersions

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&body); err != nil {
		return nil, nil, fmt.Errorf("could not query provider registry for %s: %s", addr.String(), err.Error())
	}

	if len(body.Versions) == 0 {
		return nil, body.Warnings, nil
	}

	versionList := make(versions.List, 0, len(body.Versions))
	for _, v := range body.Versions {
		ver, err := versions.ParseVersion(v.Version)
		if err != nil {
			zap.S().Warnf("registry response includes invalid version string %q: %s", v.Version, err)
			continue
		}
		versionList = append(versionList, ver)
	}
	// versionList.Newest()

	return versionList, body.Warnings, nil
}

var versionCache = make(map[string]string)

// LatestProviderVersion returns the latest published version for the asked provider.
// It returns "0.0.0" in case its not available
func LatestProviderVersion(addr addrs.Provider, terraformVersion string) string {

	// check if the cache has the version info
	if v, found := versionCache[fmt.Sprintf("%s-%s", addr.Type, terraformVersion)]; found {
		return v
	}
	versionList, _, err := ProviderVersionList(context.TODO(), addr, terraformVersion)
	if err != nil {
		zap.S().Errorf("failed to fetch latest version for terraform provider, error: %s", err.Error())
	}
	// update cache
	versionCache[fmt.Sprintf("%s-%s", addr.Type, terraformVersion)] = versionList.Newest().String()
	return versionList.Newest().String()
}

// ParseVersionConstraints parses a "Ruby-like" version constraint string
// into a VersionConstraints value.
func ParseVersionConstraints(str string) (VersionConstraints, error) {
	return constraints.ParseRubyStyleMulti(str)
}

// GetModuleProviderVersion gets the provider version form the 'required_providers' block for module.
// if the 'required_providers' is not defined, it returns empty string
func GetModuleProviderVersion(m *configs.Module) string {
	version := ""
	if m == nil || m.ProviderRequirements == nil {
		return version
	}
	for _, requiredProvider := range m.ProviderRequirements.RequiredProviders {
		version = requiredProvider.Requirement.Required[0].String()
	}

	// trim version string
	s := strings.Split(version, " ")
	if len(s) > 1 {
		version = s[1]
	}

	return version
}

// GetFileProviderVersion gets the provider version form the 'required_providers' block for the file.
// if the 'required_providers' is not defined, it returns empty string
func GetFileProviderVersion(f *configs.File) string {
	version := ""
	if f == nil || f.RequiredProviders == nil {
		return version
	}
	for _, requiredProvider := range f.RequiredProviders[0].RequiredProviders {
		version = requiredProvider.Requirement.Required[0].String()
	}

	// trim version string
	s := strings.Split(version, " ")
	if len(s) > 1 {
		version = s[1]
	}

	return version
}
