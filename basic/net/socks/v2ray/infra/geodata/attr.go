package geodata

import (
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/infra/conf"
	"strings"
)

type AttributeList struct {
	matcher []AttributeMatcher
}

func (al *AttributeList) Match(domain *conf.Domain) bool {
	for _, matcher := range al.matcher {
		if !matcher.Match(domain) {
			return false
		}
	}
	return true
}

func (al *AttributeList) IsEmpty() bool {
	return len(al.matcher) == 0
}

func parseAttrs(attrs []string) *AttributeList {
	al := new(AttributeList)
	for _, attr := range attrs {
		trimmedAttr := strings.ToLower(strings.TrimSpace(attr))
		if len(trimmedAttr) == 0 {
			continue
		}
		al.matcher = append(al.matcher, BooleanMatcher(trimmedAttr))
	}
	return al
}

type AttributeMatcher interface {
	Match(*conf.Domain) bool
}

type BooleanMatcher string

func (m BooleanMatcher) Match(domain *conf.Domain) bool {
	for _, attr := range domain.Attribute {
		if strings.EqualFold(attr.GetKey(), string(m)) {
			return true
		}
	}
	return false
}
