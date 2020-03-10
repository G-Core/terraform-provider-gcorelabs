/*
Package nodegroups contains functionality for working GCore magnum cluster nodegroups API resources

Example to List Cluster Nodegroups

	listOpts := nodegroups.ListOpts{
	}

	allPages, err := nodegroups.List(client, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allNodegroups, err := nodegroups.ExtractClusterNodegroups(allPages)
	if err != nil {
		panic(err)
	}

	for _, nodegroups := range allNodegroups {
		fmt.Printf("%+v", nodegroups)
	}

Example to Create a Cluster NodeGroup

	iTrue := true
	createOpts := nodegroups.CreateOpts{
	}

	nodegroups, err := nodegroups.Create(client, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a NodeGroup

	nodegroupID := "484cda0e-106f-4f4b-bb3f-d413710bbe78"
	err := nodegroups.Delete(client, nodegroupID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package nodegroups
