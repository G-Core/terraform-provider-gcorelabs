/*
Package clustertemplates contains functionality for working GCore magnum cluster templates API resources

Example to List Cluster Templates

	listOpts := clustertemplates.ListOpts{
	}

	allPages, err := clustertemplates.List(clustertemplateClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allNetworks, err := clustertemplates.ExtractClusterTemplates(allPages)
	if err != nil {
		panic(err)
	}

	for _, clustertemplates := range allNetworks {
		fmt.Printf("%+v", clustertemplates)
	}

Example to Create a Cluster Template

	iTrue := true
	createOpts := clustertemplates.CreateOpts{
		Name:         "clustertemplate_1",
		ExternalNetworkID   "id",
		ImageID             "image_id"
		KeyPairID           "keypair_id",
		Name                "name",
		DockerVolumeSize    10,
	}

	clustertemplates, err := clustertemplates.Create(clustertemplateClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Network

	clustertemplateID := "484cda0e-106f-4f4b-bb3f-d413710bbe78"
	err := clustertemplates.Delete(clustertemplateClient, clustertemplateID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package clustertemplates
