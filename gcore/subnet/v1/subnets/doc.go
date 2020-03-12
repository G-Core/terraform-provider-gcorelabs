/*
Package subnet contains functionality for working GCLoud subnets API resources

Example to List Subnet

	listOpts := subnets.ListOpts{
	}

	allPages, err := subnets.List(subnetClient).AllPages()
	if err != nil {
		panic(err)
	}

	allSubnets, err := subnets.ExtractSubnets(allPages)
	if err != nil {
		panic(err)
	}

	for _, subnet := range allSubnets {
		fmt.Printf("%+v", subnet)
	}

Example to Create a Subnet

	createOpts := subnets.CreateOpts{
	}

	subnets, err := subnets.Create(subnetClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Subnet

	subnetID := "484cda0e-106f-4f4b-bb3f-d413710bbe78"
	err := subnets.Delete(subnetClient, subnetID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package subnets
