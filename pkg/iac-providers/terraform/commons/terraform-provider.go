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
	"github.com/hashicorp/terraform/addrs"
	hclConfigs "github.com/hashicorp/terraform/configs"
	httputils "github.com/tenable/terrascan/pkg/utils/http"
	"go.uber.org/zap"
)

const (
	apiVersion             = "v1"
	terraformVersionHeader = "X-Terraform-Version"
)

var versionCache = make(map[string]string)

// ProviderVersions ...
type ProviderVersions struct {
	Versions []struct {
		Version   string   `json:"version"`
		Protocols []string `json:"protocols"`
	} `json:"versions"`
	Warnings []string `json:"warnings"`
}

// providerVersionList fetches all the versions of terraform providers
func providerVersionList(ctx context.Context, addr addrs.Provider, terraformVersion string) (versions.List, []string, error) {
	zap.S().Debugf("fetching list of providers metadata, hostname: %q, type: %q, namespace: %q, ", addr.Hostname.String(), addr.Namespace, addr.Type)

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
	return versionList, body.Warnings, nil
}

// latestProviderVersion returns the latest published version for the asked provider.
// It returns "0.0.0" in case its not available
func latestProviderVersion(addr addrs.Provider, terraformVersion string) string {

	// check if the cache has the version info
	if v, found := versionCache[fmt.Sprintf("%s-%s", addr.Type, terraformVersion)]; found {
		return v
	}
	versionList, _, err := providerVersionList(context.TODO(), addr, terraformVersion)
	if err != nil {
		zap.S().Warnf("failed to fetch latest version for terraform provider, error: %s", err.Error())
		return versionList.Newest().String()
	}
	// update cache
	versionCache[fmt.Sprintf("%s-%s", addr.Type, terraformVersion)] = versionList.Newest().String()
	return versionList.Newest().String()
}

// GetProviderVersion identifies the version constraints for file resources.
func GetProviderVersion(f *hclConfigs.File, addr addrs.Provider, terraformVersion string) string {
	version := ""

	for _, rps := range f.RequiredProviders {
		if rp, exist := rps.RequiredProviders[addr.Type]; exist {
			version = trimVersionConstraints(rp.Requirement.Required.String())
			break
		}
	}

	// older version of terraform (terraform version < 1.x) may have version in provider block
	if len(version) == 0 {
		for _, pc := range f.ProviderConfigs {
			if pc.Name == addr.Type {
				version = trimVersionConstraints(pc.Version.Required.String())
				break
			}
		}
	}

	// fetch latest version
	if len(version) == 0 {
		version = latestProviderVersion(addr, terraformVersion)
	}
	v, err := versions.ParseVersion(version)
	if err != nil {
		zap.S().Warnf("failed to parse provider version: %s", err.Error())
		return ""
	}
	return v.String()
}

// GetModuleProviderVersion identifies the version constraints for module resources.
func GetModuleProviderVersion(module *hclConfigs.Module, addr addrs.Provider, terraformVersion string) string {
	version := ""

	if rp, exist := module.ProviderRequirements.RequiredProviders[addr.Type]; exist {
		version = trimVersionConstraints(rp.Requirement.Required.String())
	}

	// older version of terraform (terraform version < 1.x) may have version in provider block
	if len(version) == 0 {
		if pc, exist := module.ProviderConfigs[addr.Type]; exist {
			version = trimVersionConstraints(pc.Version.Required.String())
		}
	}

	// fetch latest version
	if len(version) == 0 {
		version = latestProviderVersion(addr, terraformVersion)
	}
	v, err := versions.ParseVersion(version)
	if err != nil {
		zap.S().Warnf("failed to parse provider version: %s", err.Error())
		return ""
	}
	return v.String()
}

// ResolveProvider resolves provider addr
func ResolveProvider(resource *hclConfigs.Resource, requiredProviders []*hclConfigs.RequiredProviders) addrs.Provider {
	implied, err := addrs.ParseProviderPart(resource.Addr().ImpliedProvider())
	if err != nil {
		zap.S().Warnf("failed to parse provider namespace or type: %s", err.Error())
		return addrs.NewDefaultProvider("aws")
	}
	for _, rp := range requiredProviders {
		if provider, exists := rp.RequiredProviders[implied]; exists {
			return provider.Type
		}
	}
	return addrs.ImpliedProviderForUnqualifiedType(implied)
}

// trimVersionConstraints trim version constraints from string.
// e.g. "~> 3.0.2" will become "3.0.2"
func trimVersionConstraints(v string) string {
	s := strings.Split(v, " ")
	if len(s) > 1 {
		v = s[1]
	}
	return v
}
