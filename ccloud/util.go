package ccloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
)

// CheckDeleted checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeleted(d *schema.ResourceData, err error, msg string) error {
	if _, ok := err.(gophercloud.ErrDefault404); ok {
		d.SetId("")
		return nil
	}

	return fmt.Errorf("%s %s: %s", msg, d.Id(), err)
}

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
		if strVal, ok := val.(bool); ok {
			m[key] = fmt.Sprintf("%t", strVal)
		}
	}

	return m
}

func expandToStringSlice(v []interface{}) []string {
	s := make([]string, len(v))
	for i, val := range v {
		if strVal, ok := val.(string); ok {
			s[i] = strVal
		}
	}

	return s
}

func normalizeJsonString(v interface{}) string {
	json, _ := structure.NormalizeJsonString(v)
	return json
}

func validateURL(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := url.ParseRequestURI(value)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q URL is not valid: %s", k, err))
	}
	return
}

func validateJsonObject(v interface{}, k string) ([]string, []error) {
	if v == nil || v.(string) == "" {
		return nil, []error{fmt.Errorf("%q value must not be empty", k)}
	}

	var j map[string]interface{}
	s := v.(string)

	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return nil, []error{fmt.Errorf("%q must be a JSON object: %s", k, err)}
	}

	return nil, nil
}

func validateJsonArray(v interface{}, k string) ([]string, []error) {
	if v == nil || v.(string) == "" {
		return nil, []error{fmt.Errorf("%q value must not be empty", k)}
	}

	var j []interface{}
	s := v.(string)

	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return nil, []error{fmt.Errorf("%q must be a JSON array: %s", k, err)}
	}

	return nil, nil
}

func diffSuppressJsonObject(k, old, new string, d *schema.ResourceData) bool {
	if strSliceContains([]string{"{}", "null", ""}, old) &&
		strSliceContains([]string{"{}", "null", ""}, new) {
		return true
	}
	return false
}

func diffSuppressJsonArray(k, old, new string, d *schema.ResourceData) bool {
	if strSliceContains([]string{"[]", "null", ""}, old) &&
		strSliceContains([]string{"[]", "null", ""}, new) {
		return true
	}
	return false
}
