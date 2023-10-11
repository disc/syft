package cpe

import (
	"fmt"
	"github.com/anchore/syft/internal"
	"github.com/anchore/syft/syft/pkg"
	"regexp"
	"strings"
)

var (
	vendorFromURLRegexp = regexp.MustCompile(`^https?://(www.)?(?P<vendor>.+)\.\w/?`)
)

func candidateVendorsForWordpressPlugin(p pkg.Package) fieldCandidateSet {
	metadata, ok := p.Metadata.(pkg.WordpressPluginMetadata)
	if !ok {
		return nil
	}

	vendors := newFieldCandidateSet()

	if metadata.AuthorURI != "" {
		matchMap := internal.MatchNamedCaptureGroups(vendorFromURLRegexp, metadata.AuthorURI)
		if vendor, ok := matchMap["vendor"]; ok && vendor != "" {
			vendors.addValue(vendor)
		}
	} else {
		// add plugin_name + _project as a vendor if no Author URI found
		vendors.addValue(fmt.Sprintf("%s_project", normalizeWordpressPluginName(p.Name)))
	}

	return vendors
}

func candidateProductsForWordpressPlugin(p pkg.Package) fieldCandidateSet {
	metadata, ok := p.Metadata.(pkg.WordpressPluginMetadata)
	if !ok {
		return nil
	}
	products := newFieldCandidateSet()

	products.addValue(normalizeWordpressPluginName(p.Name))
	products.addValue(normalizeWordpressPluginName(metadata.PluginName))

	return products
}

func normalizeWordpressPluginName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	for _, value := range []string{" "} {
		name = strings.ReplaceAll(name, value, "_")
	}
	return name
}
