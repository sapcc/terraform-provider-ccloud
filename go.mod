module github.com/sapcc/terraform-provider-ccloud

go 1.14

require (
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/runtime v0.19.5
	github.com/go-openapi/strfmt v0.19.3
	github.com/go-openapi/validate v0.19.3
	github.com/gophercloud/gophercloud v0.12.1-0.20200915182658-5690aed155c1
	github.com/gophercloud/utils v0.0.0-20200910022956-24a96b152a23
	github.com/hashicorp/terraform-plugin-sdk v1.15.0
	github.com/sapcc/gophercloud-sapcc v0.0.0-20200823173629-33e592571644
	github.com/sapcc/kubernikus v1.0.1-0.20200130145221-13142045b03f
	github.com/sapcc/limes v0.0.0-20200420152328-facea01fd1ab
	k8s.io/apimachinery v0.0.0-20190831074630-461753078381 // indirect
	k8s.io/client-go v11.0.0+incompatible
)
