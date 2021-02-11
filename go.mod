module github.com/terraform-providers/terraform-provider-gcorelabs

go 1.14

replace github.com/G-Core/gcorelabscloud-go => ../gcorelabscloud-go //for local testing in case of both repos are changed

require (
	github.com/G-Core/gcorelabscloud-go v0.3.36
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.3.0
	github.com/mitchellh/mapstructure v1.4.0
)
