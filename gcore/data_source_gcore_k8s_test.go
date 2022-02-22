package gcore

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v1/clusters"
	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v1/pools"
	"github.com/G-Core/gcorelabscloud-go/gcore/keypair/v2/keypairs"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/availablenetworks"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/router/v1/routers"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	testClusterName      = "test-cluster"
	testClusterVersion   = "1.20.6"
	testClusterPoolName  = "test-pool"
	testPoolFlavor       = "g1-standard-1-2"
	testNodeCount        = 1
	testDockerVolumeSize = 10
	testDockerVolumeType = volumes.Standard
	testMinNodeCount     = 1
	testMaxNodeCount     = 1

	kpName = "testkp"
)

func TestAccK8sDataSource(t *testing.T) {
	fullName := "data.gcore_k8s.acctest"

	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	k8sClient, err := CreateTestClient(cfg.Provider, K8sPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	netClient, err := CreateTestClient(cfg.Provider, networksPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	subnetClient, err := CreateTestClient(cfg.Provider, subnetPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	kpClient, err := CreateTestClient(cfg.Provider, keypairsPoint, versionPointV2)
	if err != nil {
		t.Fatal(err)
	}

	netOpts := networks.CreateOpts{
		Name:         networkTestName,
		CreateRouter: true,
	}
	networkID, err := createTestNetwork(netClient, netOpts)
	if err != nil {
		t.Fatal(err)
	}
	defer networks.Delete(netClient, networkID)

	gw := net.ParseIP("")
	subnetOpts := subnets.CreateOpts{
		Name:                   subnetTestName,
		NetworkID:              networkID,
		ConnectToNetworkRouter: true,
		EnableDHCP:             true,
		GatewayIP:              &gw,
	}

	var gccidr gcorecloud.CIDR
	_, netIPNet, err := net.ParseCIDR(cidr)
	if err != nil {
		t.Fatal(err)
	}
	gccidr.IP = netIPNet.IP
	gccidr.Mask = netIPNet.Mask
	subnetOpts.CIDR = gccidr

	subnetID, err := CreateTestSubnet(subnetClient, subnetOpts)
	if err != nil {
		t.Fatal(err)
	}
	defer subnets.Delete(subnetClient, subnetID)

	// update our new network router so that the k8s nodes will have access to the Nexus
	// registry to download images
	if err := patchRouterForK8S(cfg.Provider, networkID); err != nil {
		t.Fatal(err)
	}

	pid, err := strconv.Atoi(os.Getenv("TEST_PROJECT_ID"))
	if err != nil {
		t.Fatal(err)
	}

	kpOpts := keypairs.CreateOpts{
		Name:      kpName,
		PublicKey: pkTest,
		ProjectID: pid,
	}
	keyPair, err := keypairs.Create(kpClient, kpOpts).Extract()
	defer keypairs.Delete(kpClient, keyPair.ID)

	k8sOpts := clusters.CreateOpts{
		Name:               testClusterName,
		FixedNetwork:       networkID,
		FixedSubnet:        subnetID,
		AutoHealingEnabled: true,
		KeyPair:            keyPair.ID,
		Version:            testClusterVersion,
		Pools: []pools.CreateOpts{{
			Name:             testClusterPoolName,
			FlavorID:         testPoolFlavor,
			NodeCount:        testNodeCount,
			DockerVolumeSize: testDockerVolumeSize,
			DockerVolumeType: testDockerVolumeType,
			MinNodeCount:     testMinNodeCount,
			MaxNodeCount:     testMaxNodeCount,
		}},
	}
	clusterID, err := CreateTestCluster(k8sClient, k8sOpts)
	if err != nil {
		t.Fatal(err)
	}
	defer clusters.Delete(k8sClient, clusterID)

	ipTemplate := fmt.Sprintf(`
			data "gcore_k8s" "acctest" {
			  %s
              %s
              cluster_id = "%s"
			}
		`, projectInfo(), regionInfo(), clusterID)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ipTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "cluster_id", clusterID),
				),
			},
		},
	})
}

func patchRouterForK8S(provider *gcorecloud.ProviderClient, networkID string) error {
	routersClient, err := CreateTestClient(provider, RouterPoint, versionPointV1)
	if err != nil {
		return err
	}

	aNetClient, err := CreateTestClient(provider, sharedNetworksPoint, versionPointV1)
	if err != nil {
		return err
	}

	availableNetworks, err := availablenetworks.ListAll(aNetClient)
	if err != nil {
		return err
	}
	var extNet availablenetworks.Network
	for _, an := range availableNetworks {
		if an.External {
			extNet = an
			break
		}
	}

	rs, err := routers.ListAll(routersClient, nil)
	if err != nil {
		return err
	}

	var router routers.Router
	for _, r := range rs {
		if strings.Contains(r.Name, networkID) {
			router = r
			break
		}
	}

	extSubnet := extNet.Subnets[0]
	routerOpts := routers.UpdateOpts{Routes: extSubnet.HostRoutes}
	_, err = routers.Update(routersClient, router.ID, routerOpts).Extract()
	if err != nil {
		return err
	}
	return nil
}

func CreateTestCluster(client *gcorecloud.ServiceClient, opts clusters.CreateOpts) (string, error) {
	res, err := clusters.Create(client, opts).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	clusterID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, K8sCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		clusterID, err := clusters.ExtractClusterIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve cluster ID from task info: %w", err)
		}
		return clusterID, nil
	},
	)
	if err != nil {
		return "", err
	}

	return clusterID.(string), nil
}
