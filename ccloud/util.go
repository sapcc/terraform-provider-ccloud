package ccloud

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func GetRegion(d *schema.ResourceData, config *Config) string {
	if v, ok := d.GetOk("region"); ok {
		return v.(string)
	}

	return config.Region
}

// List of headers that need to be redacted
var REDACT_HEADERS = []string{"x-auth-token", "x-auth-key", "x-service-token",
	"x-storage-token", "x-account-meta-temp-url-key", "x-account-meta-temp-url-key-2",
	"x-container-meta-temp-url-key", "x-container-meta-temp-url-key-2", "set-cookie",
	"x-subject-token"}

// RedactHeaders processes a headers object, returning a redacted list
func RedactHeaders(headers http.Header) (processedHeaders []string) {
	for name, header := range headers {
		for _, v := range header {
			if strSliceContains(REDACT_HEADERS, name) {
				processedHeaders = append(processedHeaders, fmt.Sprintf("%v: %v", name, "***"))
			} else {
				processedHeaders = append(processedHeaders, fmt.Sprintf("%v: %v", name, v))
			}
		}
	}
	return
}

// FormatHeaders processes a headers object plus a deliminator, returning a string
func FormatHeaders(headers http.Header, seperator string) string {
	redactedHeaders := RedactHeaders(headers)
	sort.Strings(redactedHeaders)

	return strings.Join(redactedHeaders, seperator)
}

// IsSliceContainsStr returns true if the string exists in given slice, ignore case.
func strSliceContains(sl []string, str string) bool {
	str = strings.ToLower(str)
	for _, s := range sl {
		if strings.ToLower(s) == str {
			return true
		}
	}
	return false
}

func expandToMapStringString(v map[string]interface{}) map[string]string {
	m := make(map[string]string)
	for key, val := range v {
		if strVal, ok := val.(string); ok {
			m[key] = strVal
		}
	}

	return m
}
