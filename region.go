package xrmgo

import (
	"regexp"
)

var (
	region = map[string]string{
		"crm9.dynamics.com": "urn:crmgcc:dynamics.com",
		"crm7.dynamics.com": "urn:crmjpn:dynamics.com",
		"crm6.dynamics.com": "urn:crmoce:dynamics.com",
		"crm5.dynamics.com": "urn:crmapac:dynamics.com",
		"crm4.dynamics.com": "urn:crmemea:dynamics.com",
		"crm2.dynamics.com": "urn:crmsam:dynamics.com",
		"crm.dynamics.com":  "urn:crmna:dynamics.com",
	}
)

var (
	re = regexp.MustCompile("crm[\\d]?\\.dynamics\\.com")
)
