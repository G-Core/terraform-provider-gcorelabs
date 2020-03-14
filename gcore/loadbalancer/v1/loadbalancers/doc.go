/*
Package loadbalancer contains functionality for working GCLoud loadbalancers API resources

Example to List LoadBalancer

	listOpts := loadbalancers.ListOpts{
	}

	allPages, err := loadbalancers.List(loadbalancerClient).AllPages()
	if err != nil {
		panic(err)
	}

	allLoadBalancers, err := loadbalancers.ExtractLoadBalancers(allPages)
	if err != nil {
		panic(err)
	}

	for _, loadbalancer := range allLoadBalancers {
		fmt.Printf("%+v", loadbalancer)
	}

Example to Create a LoadBalancer

	createOpts := loadbalancers.CreateOpts{
	}

	loadbalancers, err := loadbalancers.Create(loadbalancerClient, createOpts).Extract()
	if err != nil {
		panic(err)
	}

Example to Delete a LoadBalancer

	loadbalancerID := "484cda0e-106f-4f4b-bb3f-d413710bbe78"
	err := loadbalancers.Delete(loadbalancerClient, loadbalancerID).ExtractErr()
	if err != nil {
		panic(err)
	}
*/
package loadbalancers
