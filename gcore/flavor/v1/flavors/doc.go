/*
Package flavors contains functionality for working GCLoud flavors API resources

Example to List Flavors

	listOpts := flavors.ListOpts{
	}

	allPages, err := flavors.List(flavorClient).AllPages()
	if err != nil {
		panic(err)
	}

	allFlavors, err := flavors.ExtractFlavors(allPages)
	if err != nil {
		panic(err)
	}

	for _, flavor := range allFlavors {
		fmt.Printf("%+v", flavor)
	}

Example to Create a Flavor

	createOpts := flavors.CreateOpts{
	}

	flavors, err := flavors.Create(flavorClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Flavor

	flavorID := "484cda0e-106f-4f4b-bb3f-d413710bbe78"
	err := flavors.Delete(flavorClient, flavorID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package flavors
