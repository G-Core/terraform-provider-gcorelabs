module github.com/terraform-providers/terraform-provider-gcorelabs

go 1.16

//replace github.com/G-Core/gcorelabscloud-go => /home/ondi/go/src/github.com/G-Core/gcorelabscloud-go

require (
	github.com/G-Core/g-dns-sdk-go v0.1.2
	github.com/G-Core/gcorelabs-storage-sdk-go v0.0.9
	github.com/G-Core/gcorelabscdn-go v0.1.2
	github.com/G-Core/gcorelabscloud-go v0.4.35
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/terraform-exec v0.15.0 // indirect
	github.com/hashicorp/terraform-plugin-docs v0.5.1 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.7.0
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/zclconf/go-cty v1.10.0 // indirect
)
