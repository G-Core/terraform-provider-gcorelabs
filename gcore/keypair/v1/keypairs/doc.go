/*
Package keypair contains functionality for working GCLoud keypairs API resources

Example to List KeyPair

	listOpts := keypairs.ListOpts{
	}

	allPages, err := keypairs.List(client).AllPages()
	if err != nil {
		panic(err)
	}

	allKeyPairs, err := keypairs.ExtractKeyPairs(allPages)
	if err != nil {
		panic(err)
	}

	for _, keypair := range allKeyPairs {
		fmt.Printf("%+v", keypair)
	}

Example to Create a KeyPair

	createOpts := keypairs.CreateOpts{
		Name": "alice",
		PublicKey: "",
	}

	keypairs, err := keypairs.Create(client, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a KeyPair

	keypairID := "alica"
	err := keypairs.Delete(client, keypairID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package keypairs
