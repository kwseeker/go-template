package geodata

import (
	"github.com/adrg/xdg"
	"github.com/golang/protobuf/proto"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common/errors"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common/platform"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common/platform/filesystem"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/infra/conf"
	"log"
	"path/filepath"
	"strings"
)

const (
	EnvAsset         = "v2ray.location.asset"
	DefaultGeoLoader = "standard"
)

var (
	loaders map[string]*GeoLoaderImpl
)

type Loader interface {
	LoadSite(filename, list string) ([]*conf.Domain, error)
	LoadIP(filename, country string) ([]*conf.CIDR, error)
}

type GeoLoader interface {
	Loader

	LoadGeoSite(list string) ([]*conf.Domain, error)
	LoadGeoSiteWithAttr(file string, siteWithAttr string) ([]*conf.Domain, error)
	LoadGeoIP(country string) ([]*conf.CIDR, error)
}

func getAssertLocation(file string) string {
	assetPath := platform.GetEnvValue(EnvAsset)
	defPath := filepath.Join(assetPath, file)
	relPath := filepath.Join("v2ray", file)
	fullPath, err := xdg.SearchDataFile(relPath)
	if err != nil {
		return defPath
	}
	return fullPath
}

func ReadAsset(file string) ([]byte, error) {
	return filesystem.ReadFile(getAssertLocation(file))
}

func loadSite(filename, list string) ([]*conf.Domain, error) {
	geositeBytes, err := ReadAsset(filename)
	if err != nil {
		return nil, errors.NewError("failed to open file: ", filename)
	}
	var geositeList conf.GeoSiteList
	if err := proto.Unmarshal(geositeBytes, &geositeList); err != nil {
		return nil, err
	}

	for _, site := range geositeList.Entry {
		if strings.EqualFold(site.CountryCode, list) {
			return site.Domain, nil
		}
	}

	return nil, errors.NewError("list not found in ", filename, ": ", list)
}

func loadIP(filename, country string) ([]*conf.CIDR, error) {
	geoipBytes, err := ReadAsset(filename)
	if err != nil {
		return nil, errors.NewError("failed to open file: ", filename)
	}
	var geoipList conf.GeoIPList
	if err := proto.Unmarshal(geoipBytes, &geoipList); err != nil {
		return nil, err
	}

	for _, geoip := range geoipList.Entry {
		if strings.EqualFold(geoip.CountryCode, country) {
			return geoip.Cidr, nil
		}
	}

	return nil, errors.NewError("country not found in ", filename, ": ", country)
}

type standardLoader struct{}

func (d standardLoader) LoadSite(filename, list string) ([]*conf.Domain, error) {
	return loadSite(filename, list)
}

func (d standardLoader) LoadIP(filename, country string) ([]*conf.CIDR, error) {
	return loadIP(filename, country)
}

type GeoLoaderImpl struct {
	Loader
}

func (l *GeoLoaderImpl) LoadGeoSite(list string) ([]*conf.Domain, error) {
	return l.LoadGeoSiteWithAttr("geosite.dat", list)
}

func (l *GeoLoaderImpl) LoadGeoSiteWithAttr(file string, siteWithAttr string) ([]*conf.Domain, error) {
	parts := strings.Split(siteWithAttr, "@")
	if len(parts) == 0 {
		return nil, errors.NewError("empty rule")
	}
	list := strings.TrimSpace(parts[0])
	attrVal := parts[1:]

	if len(list) == 0 {
		return nil, errors.NewError("empty list name in rule: ", siteWithAttr)
	}

	domains, err := l.LoadSite(file, list)
	if err != nil {
		return nil, err
	}

	attrs := parseAttrs(attrVal)
	if attrs.IsEmpty() {
		if strings.Contains(siteWithAttr, "@") {
			log.Println("empty attribute list: ", siteWithAttr)
		}
		return domains, nil
	}

	filteredDomains := make([]*conf.Domain, 0, len(domains))
	hasAttrMatched := false
	for _, domain := range domains {
		if attrs.Match(domain) {
			hasAttrMatched = true
			filteredDomains = append(filteredDomains, domain)
		}
	}
	if !hasAttrMatched {
		log.Println("attribute match no rule: geosite:", siteWithAttr)
	}

	return filteredDomains, nil
}

func (l *GeoLoaderImpl) LoadGeoIP(country string) ([]*conf.CIDR, error) {
	return l.LoadIP("geoip.dat", country)
}

func GetGeoLoader() GeoLoader {
	if loaders == nil {
		loaders = make(map[string]*GeoLoaderImpl)
	}
	loader := loaders[DefaultGeoLoader]
	if loader == nil {
		loader = &GeoLoaderImpl{standardLoader{}}
		loaders[DefaultGeoLoader] = loader
	}
	return loader
}
