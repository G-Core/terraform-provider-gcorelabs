package gcore

import (
	"fmt"
	"net"
	"strconv"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/router/v1/routers"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func checkRouterAttrs(resourceName string, opts *routers.CreateOpts) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if s.Empty() == true {
			return fmt.Errorf("State not updated")
		}

		checksStore := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceName, "name", opts.Name),
		}

		mapopts, _ := opts.ToRouterCreateMap()
		_, ok := mapopts["external_gateway_info"]
		if ok {
			checksStore = append(checksStore,
				resource.TestCheckResourceAttr(resourceName, "external_gateway_info.0.type", opts.ExternalGatewayInfo.Type.String()),
				resource.TestCheckResourceAttr(resourceName, "external_gateway_info.0.enable_snat", strconv.FormatBool(*opts.ExternalGatewayInfo.EnableSNat)),
			)
		}

		if len(opts.Interfaces) > 0 {
			for i, iface := range opts.Interfaces {
				checksStore = append(checksStore,
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`interfaces.%d.type`, i), iface.Type.String()),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`interfaces.%d.subnet_id`, i), iface.SubnetID),
				)
			}
		}

		for i, r := range opts.Routes {
			checksStore = append(checksStore,
				resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`routes.%d.destination`, i), r.Destination.String()),
				resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`routes.%d.nexthop`, i), r.NextHop.String()),
			)
		}

		return resource.ComposeTestCheckFunc(checksStore...)(s)
	}
}

func TestAccRouter(t *testing.T) {
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
		Name:         networkTestName,
		CreateRouter: false,
	}

	networkID, err := createTestNetwork(clientNet, opts)
	if err != nil {
		t.Fatal(err)
	}

	defer networks.Delete(clientNet, networkID)

	gw := net.ParseIP("")
	optsSubnet := subnets.CreateOpts{
		Name:                   subnetTestName,
		NetworkID:              networkID,
		ConnectToNetworkRouter: false,
		GatewayIP:              &gw,
	}

	var gccidr gcorecloud.CIDR
	_, netIPNet, err := net.ParseCIDR(cidr)
	if err != nil {
		t.Fatal(err)
	}
	gccidr.IP = netIPNet.IP
	gccidr.Mask = netIPNet.Mask
	optsSubnet.CIDR = gccidr

	subnetID, err := CreateTestSubnet(clientSubnet, optsSubnet)
	if err != nil {
		t.Fatal(err)
	}

	var dst1 gcorecloud.CIDR
	snat1 := true

	_, netIPNet, _ = net.ParseCIDR(cidr)
	dst1.IP = netIPNet.IP
	dst1.Mask = netIPNet.Mask

	createFixt := routers.CreateOpts{
		Name: "create_router",
		ExternalGatewayInfo: routers.GatewayInfo{
			Type:       "default",
			EnableSNat: &snat1,
		},
		Interfaces: []routers.Interface{
			{
				Type:     "subnet",
				SubnetID: subnetID,
			},
		},
		Routes: []subnets.HostRoute{
			{
				Destination: dst1,
				NextHop:     net.ParseIP("192.168.42.2"),
			},
		},
	}

	type Params struct {
		Name           string
		ExtGatewayInfo []map[string]string
		Interfaces     []map[string]string
		Routes         []map[string]string
	}

	create := Params{
		Name:           "create_router",
		ExtGatewayInfo: []map[string]string{{"type": "default", "enable_snat": "true", "network_id": ""}},
		Interfaces:     []map[string]string{{"type": "subnet", "subnet_id": subnetID}},
		Routes:         []map[string]string{{"destination": "192.168.42.0/24", "nexthop": "192.168.42.2"}},
	}

	RouterTemplate := func(params *Params) string {
		template := `
		locals {
            external_gateway_info = [`
		for i := range params.ExtGatewayInfo {
			template += fmt.Sprintf(`
			{
				type = "%s"
				enable_snat = %s
				network_id = "%s"
			},`, params.ExtGatewayInfo[i]["type"], params.ExtGatewayInfo[i]["enable_snat"], params.ExtGatewayInfo[i]["network_id"])
		}

		template += fmt.Sprintf(`]
			interfaces = [`)
		for i := range params.Interfaces {
			template += fmt.Sprintf(`
			{
				type = "%s"
				subnet_id = "%s"
			},`, params.Interfaces[i]["type"], params.Interfaces[i]["subnet_id"])
		}

		template += fmt.Sprintf(`]
			routes = [`)
		for i := range params.Routes {
			template += fmt.Sprintf(`
			{
				destination = "%s"
				nexthop = "%s"
			},`, params.Routes[i]["destination"], params.Routes[i]["nexthop"])
		}

		template += fmt.Sprintf(`]
        }
        resource "gcore_router" "acctest" {
			name = "%s"

			dynamic external_gateway_info {
			iterator = egi
			for_each = local.external_gateway_info
			content {
				type = egi.value.type
				enable_snat = egi.value.enable_snat
				network_id = egi.value.network_id
				}
			}
          
        	dynamic interfaces {
			iterator = ifaces
			for_each = local.interfaces
			content {
				type = ifaces.value.type
				subnet_id = ifaces.value.subnet_id
				}
  			}

			dynamic routes {
			iterator = rs
			for_each = local.routes
			content {	
				destination = rs.value.destination
				nexthop = rs.value.nexthop
				}
			}

            %[2]s
			%[3]s

		`, params.Name, regionInfo(), projectInfo())
		return template + "\n}"
	}

	fullName := "gcore_router.acctest"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: RouterTemplate(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					checkRouterAttrs(fullName, &createFixt),
				),
			},
		},
	})
}

func testAccRouterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, RouterPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_router" {
			continue
		}

		_, err := routers.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("router still exists")
		}
	}

	return nil
}
