/*
Package volume contains functionality for working GCLoud volumes API resources

Example to List Volume

	listOpts := volumes.ListOpts{
	}

	allPages, err := volumes.List(volumeClient).AllPages()
	if err != nil {
		panic(err)
	}

	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		panic(err)
	}

	for _, volume := range allVolumes {
		fmt.Printf("%+v", volume)
	}

Example to Create a Volume

	createOpts := volumes.CreateOpts{
	}

	volumes, err := volumes.Create(volumeClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a Volume

	volumeID := "484cda0e-106f-4f4b-bb3f-d413710bbe78"
	err := volumes.Delete(volumeClient, volumeID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package volumes
