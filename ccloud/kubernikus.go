package ccloud

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/sapcc/kubernikus/pkg/api/client/operations"

	"github.com/gophercloud/gophercloud"
	osClient "github.com/gophercloud/utils/client"
)

var (
	httpMethods = []string{
		"GET",
		"POST",
		"PATCH",
		"DELETE",
		"PUT",
		"HEAD",
		"OPTIONS",
		"CONNECT",
		"TRACE",
	}
	maskHeader = strings.ToLower("X-Auth-Token:")
)

type kubernikus struct {
	operations.ClientService
	provider *gophercloud.ProviderClient
}

type kubernikusLogger struct{}

func (kubernikusLogger) Printf(format string, args ...interface{}) {
	if len(format) == 0 || format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
}

func (kubernikusLogger) Debugf(format string, args ...interface{}) {
	for _, arg := range args {
		if v, ok := arg.(string); ok {
			str := deleteEmpty(strings.Split(v, "\n"))
			cycle := "Response"
			if len(str) > 0 {
				for _, method := range httpMethods {
					if strings.HasPrefix(str[0], method) {
						cycle = "Request"
						break
					}
				}
			}
			printed := false

			for i, s := range str {
				if i == 0 && cycle == "Request" {
					v := strings.SplitN(s, " ", 3)
					if len(v) > 1 {
						log.Printf("[DEBUG] Kubernikus %s URL: %s %s", cycle, v[0], v[1])
					}
				} else if i == 0 && cycle == "Response" {
					v := strings.SplitN(s, " ", 2)
					if len(v) > 1 {
						log.Printf("[DEBUG] Kubernikus %s Code: %s", cycle, v[1])
					}
				} else if i == len(str)-1 {
					debugInfo, err := formatJSON([]byte(s))
					if err != nil {
						printHeaders(cycle, &printed)
						log.Print(s)
					} else {
						log.Printf("[DEBUG] Kubernikus %s Body: %s", cycle, debugInfo)
					}
				} else if strings.HasPrefix(strings.ToLower(s), maskHeader) {
					printHeaders(cycle, &printed)
					v := strings.SplitN(s, ":", 2)
					log.Printf("%s: ***", v[0])
				} else {
					printHeaders(cycle, &printed)
					log.Print(s)
				}
			}
		}
	}
}

func newKubernikusV1(c *Config, eo gophercloud.EndpointOpts) (*kubernikus, error) {
	var err error
	var endpoint string
	var kurl *url.URL

	if v, ok := c.EndpointOverrides["kubernikus"]; ok {
		if e, ok := v.(string); ok && e != "" {
			endpoint = e
		}
	}

	if endpoint == "" && !reflect.DeepEqual(eo, gophercloud.EndpointOpts{}) {
		eo.ApplyDefaults("kubernikus")
		endpoint, err = c.OsClient.EndpointLocator(eo)
		if err != nil {
			return nil, err
		}
	}

	if kurl, err = url.Parse(endpoint); err != nil {
		return nil, fmt.Errorf("Parsing the Kubernikus URL failed: %s", err)
	}

	transport := httptransport.New(kurl.Host, kurl.EscapedPath(), []string{kurl.Scheme})

	if v, ok := c.OsClient.HTTPClient.Transport.(*osClient.RoundTripper); ok && v.Logger != nil {
		// enable JSON debug for Kubernikus
		transport.SetLogger(kubernikusLogger{})
		transport.Debug = true
	}

	operations := operations.New(transport, strfmt.Default)

	return &kubernikus{operations, c.OsClient}, nil
}

func (k *kubernikus) authFunc() runtime.ClientAuthInfoWriterFunc {
	return runtime.ClientAuthInfoWriterFunc(
		func(req runtime.ClientRequest, reg strfmt.Registry) error {
			err := req.SetHeaderParam("X-Auth-Token", k.provider.Token())
			if err != nil {
				log.Printf("[DEBUG] Kubernikus auth func cannot set X-Auth-Token header value: %v", err)
			}
			err = req.SetHeaderParam("User-Agent", k.provider.UserAgent.Join())
			if err != nil {
				log.Printf("[DEBUG] Kubernikus auth func cannot set User-Agent header value: %v", err)
			}
			return nil
		})
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if strings.TrimSpace(str) != "" {
			r = append(r, str)
		}
	}
	return r
}

func printHeaders(cycle string, printed *bool) {
	if !*printed {
		log.Printf("[DEBUG] Kubernikus %s Headers:", cycle)
		*printed = true
	}
}

// formatJSON is a function to pretty-format a JSON body.
// It will also mask known fields which contain sensitive information.
func formatJSON(raw []byte) (string, error) {
	var rawData interface{}

	err := json.Unmarshal(raw, &rawData)
	if err != nil {
		return string(raw), fmt.Errorf("unable to parse OpenStack JSON: %s", err)
	}

	data, ok := rawData.(map[string]interface{})
	if !ok {
		pretty, err := json.MarshalIndent(rawData, "", "  ")
		if err != nil {
			return string(raw), fmt.Errorf("unable to re-marshal OpenStack JSON: %s", err)
		}

		return string(pretty), nil
	}

	// Strip kubeconfig
	if _, ok := data["kubeconfig"].(string); ok {
		data["kubeconfig"] = "***"
	}

	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return string(raw), fmt.Errorf("unable to re-marshal OpenStack JSON: %s", err)
	}

	return string(pretty), nil
}
