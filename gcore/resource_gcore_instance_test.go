//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"net"
	"regexp"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/image/v1/images"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygroups"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func checkInstanceAttrs(resourceName string, opts *instances.CreateOpts) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if s.Empty() == true {
			return fmt.Errorf("State not updated")
		}

		checksStore := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceName, "name", opts.Names[0]),
			resource.TestCheckResourceAttr(resourceName, "flavor_id", opts.Flavor),
			resource.TestCheckResourceAttr(resourceName, "keypair_name", opts.Keypair),
			resource.TestCheckResourceAttr(resourceName, "password", opts.Password),
			resource.TestCheckResourceAttr(resourceName, "username", opts.Username),
		}

		// todo add check for interfaces/volumes/secgroups
		//for i, volume := range opts.Volumes {
		//	checksStore = append(checksStore,
		//		resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`volumes.%d.source`, i), volume.Source.String()),
		//		resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`volumes.%d.boot_index`, i), strconv.Itoa(volume.BootIndex)),
		//	)
		//}
		//
		//for i, iface := range opts.Interfaces {
		//	checksStore = append(checksStore,
		//		resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`interfaces.%d.type`, i), iface.Type.String()),
		//		resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`interfaces.%d.network_id`, i), iface.NetworkID),
		//		resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`interfaces.%d.subnet_id`, i), iface.SubnetID),
		//	)
		//}
		//
		//for i, secgroup := range opts.SecurityGroups {
		//	checksStore = append(checksStore,
		//		resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`security_groups.%d.id`, i), secgroup.ID),
		//	)
		//}

		for i, md := range opts.Metadata.Metadata {
			checksStore = append(checksStore,
				resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`metadata.%d.key`, i), md.Key),
				resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`metadata.%d.value`, i), md.Value),
			)
		}

		for i, cfg := range opts.Configuration.Metadata {
			checksStore = append(checksStore,
				resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`configuration.%d.key`, i), cfg.Key),
				resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`configuration.%d.value`, i), cfg.Value),
			)
		}

		return resource.ComposeTestCheckFunc(checksStore...)(s)
	}
}

func TestAccInstance(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	clientImage, err := CreateTestClient(cfg.Provider, imagesPoint, versionPointV1)
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

	clientSec, err := CreateTestClient(cfg.Provider, securityGroupPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	imgs, err := images.ListAll(clientImage, nil)
	if err != nil {
		t.Fatal(err)
	}

	var img images.Image
	for _, i := range imgs {
		if i.OsDistro == testOsDistro {
			img = i
			break
		}
	}
	if img.ID == "" {
		t.Fatalf("images with os_distro='%s' does not exist", testOsDistro)
	}

	opts := networks.CreateOpts{
		Name: networkTestName,
	}

	networkID, err := createTestNetwork(clientNet, opts)
	if err != nil {
		t.Fatal(err)
	}

	defer networks.Delete(clientNet, networkID)

	optsSubnet := subnets.CreateOpts{
		Name:      subnetTestName,
		NetworkID: networkID,
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

	volumes := []instances.CreateVolumeOpts{
		{
			Source:    "existing-volume",
			BootIndex: 0,
		},
		{
			Source:    "existing-volume",
			BootIndex: 1,
		},
	}
	interfaces := []instances.InterfaceOpts{
		{
			Type:      "subnet",
			NetworkID: networkID,
			SubnetID:  subnetID,
		},
	}
	update_interfaces := []instances.InterfaceOpts{
		{
			Type:     "subnet",
			SubnetID: subnetID,
		},
	}

	sgs, err := securitygroups.ListAll(clientSec)
	if err != nil {
		t.Fatal(err)
	}

	secgroups := []gcorecloud.ItemID{{ID: sgs[0].ID}}
	update_sg := []gcorecloud.ItemID{
		{
			ID: "someid",
		},
	}
	metadata := instances.MetadataSetOpts{}
	metadata.Metadata = []instances.MetadataOpts{
		{
			Key:   "somekey",
			Value: "somevalue",
		},
	}
	update_metadata := instances.MetadataSetOpts{}
	update_metadata.Metadata = []instances.MetadataOpts{
		{
			Key:   "newsomekey",
			Value: "newsomevalue",
		},
	}

	createFixt := instances.CreateOpts{
		Names:          []string{"create_instance"},
		NameTemplates:  []string{},
		Flavor:         "g1-standard-2-4",
		Password:       "password",
		Username:       "user",
		Keypair:        "acctest",
		Volumes:        volumes,
		Interfaces:     interfaces,
		SecurityGroups: secgroups,
		Metadata:       &metadata,
		Configuration:  &metadata,
	}

	update_interfaceFixt := createFixt
	update_interfaceFixt.Interfaces = update_interfaces

	update_secgroupsFixt := createFixt
	update_interfaceFixt.SecurityGroups = update_sg

	updateFixt := createFixt
	updateFixt.Flavor = "g1-standard-2-8"
	updateFixt.Metadata = &update_metadata
	updateFixt.Configuration = &update_metadata

	type Params struct {
		Name           []string
		Flavor         string
		Password       string
		Username       string
		Keypair        string
		Publickey      string
		Image          string
		Interfaces     []map[string]string
		SecurityGroups []map[string]string
		MetaData       []map[string]string
		Configuration  []map[string]string
	}

	create := Params{
		Name:     []string{"create_instance"},
		Flavor:   "g1-standard-2-4",
		Password: "password",
		Username: "user",
		Keypair:  "acctest",
		Publickey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC1bdbQYquD/swsZpFPXagY9KvhlNUTKYMdhRNtlGglAMgRxJS3Q0V74BNElJtP+UU/" +
			"AbZD4H2ZAwW3PLLD/maclnLlrA48xg/ez9IhppBop0WADZ/nB4EcvQfR/Db7nHDTZERW6EiiGhV6CkHVasK2sY/WNRXqPveeWUlwCqtSnU90l/" +
			"s9kQCoEfkM2auO6ppJkVrXbs26vcRclS8KL7Cff4HwdVpV7b+edT5seZdtrFUCbkEof9D9nGpahNvg8mYWf0ofx4ona4kaXm1NdPID+ljvE/" +
			"dbYUX8WZRmyLjMvVQS+VxDJtsiDQIVtwbC4w+recqwDvHhLWwoeczsbEsp ondi@ds",
		Image: img.ID,
		Interfaces: []map[string]string{
			{"type": "subnet", "network_id": networkID, "subnet_id": subnetID},
		},
		SecurityGroups: []map[string]string{{"id": sgs[0].ID, "name": sgs[0].Name}},
		MetaData:       []map[string]string{{"key": "somekey", "value": "somevalue"}},
		Configuration:  []map[string]string{{"key": "somekey", "value": "somevalue"}},
	}

	update_interface := create
	update_interface.Interfaces = []map[string]string{{"type": "subnet", "subnet_id": subnetID}}

	update_secgroups := create
	update_secgroups.SecurityGroups = []map[string]string{{"id": "someid", "name": "somegroup"}}

	update := create
	update.Flavor = "g1-standard-2-8"
	update.MetaData = []map[string]string{{"key": "newsomekey", "value": "newsomevalue"}}
	update.Configuration = []map[string]string{{"key": "newsomekey", "value": "newsomevalue"}}

	instanceTemplate := func(params *Params) string {
		template := `
		locals {`

		template += fmt.Sprintf(`
			names = "%s"
            volumes_ids = [gcore_volume.first_volume.id, gcore_volume.second_volume.id]`, params.Name[0])

		template += fmt.Sprint(`
			interfaces = [`)
		for i := range params.Interfaces {
			template += fmt.Sprintf(`
			{
				type = "%s"
				network_id = "%s"
				subnet_id = "%s"
                fip_source = null
                existing_fip_id = null
                port_id = null
                ip_address = null
				
			},`, params.Interfaces[i]["type"], params.Interfaces[i]["network_id"], params.Interfaces[i]["subnet_id"])

		}
		template += fmt.Sprint(`]
			security_groups = [`)
		for i := range params.SecurityGroups {
			template += fmt.Sprintf(`
			{
				id = "%s"
				name = "%s"
			},`, params.SecurityGroups[i]["id"], params.SecurityGroups[i]["name"])
		}
		template += fmt.Sprint(`]
			metadata = [`)
		for i := range params.MetaData {
			template += fmt.Sprintf(`
			{
				key = "%s"
				value = "%s"
			},`, params.MetaData[i]["key"], params.MetaData[i]["value"])
		}
		template += fmt.Sprint(`]
			configuration = [`)
		for i := range params.Configuration {
			template += fmt.Sprintf(`
			{
				key = "%s"
				value = "%s"
			},`, params.Configuration[i]["key"], params.Configuration[i]["value"])
		}
		template += fmt.Sprintf(`]
        }

        resource "gcore_volume" "first_volume" {
  			name = "boot volume"
  			type_name = "ssd_hiiops"
  			size = 5
  			image_id = "%[1]s"
  			%[7]s
			%[8]s
		}

		resource "gcore_volume" "second_volume" {
  			name = "second volume"
  			type_name = "ssd_hiiops"
  			size = 5
  			%[7]s
			%[8]s
		}

        resource "gcore_keypair" "kp" {
  			sshkey_name = "%[2]s"
            public_key = "%[3]s"
            %[8]s
		}

        resource "gcore_instance" "acctest" {
			flavor_id = "%[4]s"
           	name = local.names
           	keypair_name = gcore_keypair.kp.sshkey_name
           	password = "%[5]s"
           	username = "%[6]s"

			dynamic volume {
		  	iterator = vol
		  	for_each = local.volumes_ids
		  	content {
				boot_index = index(local.volumes_ids, vol.value)
				source = "existing-volume"
				volume_id = vol.value
				}
		  	}

			dynamic interface {
			iterator = ifaces
			for_each = local.interfaces
			content {
				type = ifaces.value.type
				network_id = ifaces.value.network_id
				subnet_id = ifaces.value.subnet_id
                fip_source = ifaces.value.fip_source
				existing_fip_id = ifaces.value.existing_fip_id
                port_id = ifaces.value.port_id
                ip_address = ifaces.value.ip_address
				}
			}

			dynamic security_group {
			iterator = sg
			for_each = local.security_groups
			content {	
				id = sg.value.id
				name = sg.value.name
				}
			}

			dynamic metadata {
			iterator = md
			for_each = local.metadata
			content {	
				key = md.value.key
				value = md.value.value
				}
			}

			dynamic configuration {
			iterator = cfg
			for_each = local.configuration
			content {	
				key = cfg.value.key
				value = cfg.value.value
				}
			}

            %[7]s
			%[8]s

		`, params.Image, params.Keypair, params.Publickey, params.Flavor, params.Password, params.Username, regionInfo(), projectInfo())
		return template + "\n}"
	}

	fullName := "gcore_instance.acctest"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: instanceTemplate(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					checkInstanceAttrs(fullName, &createFixt),
				),
			},
			{
				Config: instanceTemplate(&update_interface),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					checkInstanceAttrs(fullName, &update_interfaceFixt),
				),
			},
			{
				Config: instanceTemplate(&update_secgroups),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					checkInstanceAttrs(fullName, &update_secgroupsFixt),
				),
				ExpectError: regexp.MustCompile("not found"),
			},
		},
	})
}

func testAccInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, InstancePoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_instance" {
			continue
		}

		_, err := instances.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Instance still exists")
		}
	}

	return nil
}
