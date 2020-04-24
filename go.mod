module github.com/sapcc/terraform-provider-ccloud

go 1.12

require (
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/runtime v0.19.5
	github.com/go-openapi/strfmt v0.19.3
	github.com/go-openapi/validate v0.19.3
	github.com/gophercloud/gophercloud v0.10.1-0.20200424014253-c3bfe50899e5
	github.com/gophercloud/utils v0.0.0-20200423144003-7c72efc7435d
	github.com/hashicorp/terraform-plugin-sdk v1.10.0
	github.com/sapcc/gophercloud-sapcc v0.0.0-20200423152112-c948c38c38fb
	github.com/sapcc/kubernikus v1.0.1-0.20200130145221-13142045b03f
	github.com/sapcc/limes v0.0.0-20200420152328-facea01fd1ab
	k8s.io/apimachinery v0.0.0-20190831074630-461753078381 // indirect
	k8s.io/client-go v11.0.0+incompatible
)
