/*
Package network contains functionality for working GCLoud networks API resources

Example to List Network

	listOpts := networks.ListOpts{
	}

	allPages, err := networks.List(networkClient).AllPages()
	if err != nil {
		panic(err)
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		panic(err)
	}

	for _, network := range allNetworks {
		fmt.Printf("%+v", network)
	}

Example to Create a Network

	createOpts := networks.CreateOpts{
		Name:  "network_1",
		MTU:   1500,
	}

	networks, err := networks.Create(networkClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Network

	networkID := "484cda0e-106f-4f4b-bb3f-d413710bbe78"
	err := networks.Delete(networkClient, networkID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package networks
