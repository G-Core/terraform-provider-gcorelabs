/*
Package cluster contains functionality for working GCore magnum clusters API resources

Example to List Cluster

	listOpts := clusters.ListOpts{
	}

	allPages, err := clusters.List(clusterClient, listOpts).AllPages()
	if err != nil {
		panic(err)
	}

	allNetworks, err := clusters.ExtractClusters(allPages)
	if err != nil {
		panic(err)
	}

	for _, clustertemplates := range allNetworks {
		fmt.Printf("%+v", clustertemplates)
	}

Example to Create a Cluster

	createOpts := clusters.CreateOpts{
		Name:         		 "cluster_1",
		ExternalNetworkId:   "id",
		ImageId:             "image_id"
		KeyPairID:           "keypair_id",
		Name:                "name",
		DockerVolumeSize:    10,
	}

	clusters, err := clusters.Create(clusterClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Network

	clusterID := "484cda0e-106f-4f4b-bb3f-d413710bbe78"
	err := clusters.Delete(clusterClient, clusterID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package clusters
