//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"net"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	subnetTestName = "test-subnet"
	cidr           = "192.168.42.0/24"
)

func TestAccSubnetDataSource(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	clientNet, err := CreateTestClient(cfg.Provider, networksPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	clientSubnet, err := CreateTestClient(cfg.Provider, subnetPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	opts := networks.CreateOpts{
		Name: networkTestName,
	}

	networkID, err := createTestNetwork(clientNet, opts)
	if err != nil {
		t.Fatal(err)
	}

	defer deleteTestNetwork(clientNet, networkID)

	optsSubnet := subnets.CreateOpts{
		Name:      subnetTestName,
		NetworkID: networkID,
	}

	subnetID, err := CreateTestSubnet(clientSubnet, optsSubnet)
	if err != nil {
		t.Fatal(err)
	}

	fullName := "data.gcore_subnet.acctest"
	tpl := func(name string) string {
		return fmt.Sprintf(`
			data "gcore_subnet" "acctest" {
			  %s
              %s
              name = "%s"
			}
		`, projectInfo(), regionInfo(), name)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tpl(optsSubnet.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", optsSubnet.Name),
					resource.TestCheckResourceAttr(fullName, "id", subnetID),
					resource.TestCheckResourceAttr(fullName, "network_id", networkID),
				),
			},
		},
	})
}

func CreateTestSubnet(client *gcorecloud.ServiceClient, opts subnets.CreateOpts) (string, error) {
	var gccidr gcorecloud.CIDR
	_, netIPNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}
	gccidr.IP = netIPNet.IP
	gccidr.Mask = netIPNet.Mask
	opts.CIDR = gccidr

	res, err := subnets.Create(client, opts).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	subnetID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, SubnetCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		Subnet, err := subnets.ExtractSubnetIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve Subnet ID from task info: %w", err)
		}
		return Subnet, nil
	},
	)

	return subnetID.(string), err
}

func deleteTestSubnet(client *gcorecloud.ServiceClient, subnetID string) error {
	results, err := subnets.Delete(client, subnetID).Extract()
	if err != nil {
		return err
	}
	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, SubnetDeleting, func(task tasks.TaskID) (interface{}, error) {
		_, err := subnets.Get(client, subnetID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete subnet with ID: %s", subnetID)
		}
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			return nil, nil
		default:
			return nil, err
		}
	})
	return err
}
