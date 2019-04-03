package ccloud

import (
	"strings"

	"github.com/sapcc/limes"
)

var (
	SERVICES = map[string]map[string]limes.Unit{
		"database": {
			"cfm_share_capacity": limes.UnitBytes,
		},
		"compute": {
			"cores":     limes.UnitNone,
			"instances": limes.UnitNone,
			"ram":       limes.UnitMebibytes,
		},
		"volumev2": {
			"capacity":  limes.UnitGibibytes,
			"snapshots": limes.UnitNone,
			"volumes":   limes.UnitNone,
		},
		"network": {
			"floating_ips":         limes.UnitNone,
			"networks":             limes.UnitNone,
			"ports":                limes.UnitNone,
			"rbac_policies":        limes.UnitNone,
			"routers":              limes.UnitNone,
			"security_group_rules": limes.UnitNone,
			"security_groups":      limes.UnitNone,
			"subnet_pools":         limes.UnitNone,
			"subnets":              limes.UnitNone,
			"healthmonitors":       limes.UnitNone,
			"l7policies":           limes.UnitNone,
			"listeners":            limes.UnitNone,
			"loadbalancers":        limes.UnitNone,
			"pools":                limes.UnitNone,
			"pool_members":         limes.UnitNone,
		},
		"dns": {
			"zones":      limes.UnitNone,
			"recordsets": limes.UnitNone,
		},
		"sharev2": {
			"share_networks":    limes.UnitNone,
			"share_capacity":    limes.UnitGibibytes,
			"shares":            limes.UnitNone,
			"snapshot_capacity": limes.UnitGibibytes,
			"share_snapshots":   limes.UnitNone,
		},
		"object-store": {
			"capacity": limes.UnitBytes,
		},
	}
)

func toString(r interface{}) string {
	switch v := r.(type) {
	case *limes.ProjectResourceReport:
		return limes.ValueWithUnit{v.Quota, v.Unit}.String()
	case *limes.DomainResourceReport:
		return limes.ValueWithUnit{v.DomainQuota, v.Unit}.String()
	}
	return ""
}

func sanitize(s string) string {
	return strings.Replace(s, "-", "", -1)
}
