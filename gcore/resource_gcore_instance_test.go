package gcore

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func checkInstanceAttrs(resourceName string, opts *instances.CreateOpts) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if s.Empty() == true {
			return fmt.Errorf("State not updated")
		}

		checksStore := []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceName, "name.0", opts.Names[0]),
			resource.TestCheckResourceAttr(resourceName, "flavor_id", opts.Flavor),
		}

		for i, volume := range opts.Volumes {
			checksStore = append(checksStore,
				resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`volumes.%d.name`, i), volume.Name),
				resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`volumes.%d.source`, i), volume.Source.String()),
				resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`volumes.%d.boot_index`, i), strconv.Itoa(volume.BootIndex)),
				resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`volumes.%d.size`, i), strconv.Itoa(volume.Size)),
				resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`volumes.%d.type_name`, i), volume.TypeName.String()),
				resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`volumes.%d.image_id`, i), volume.ImageID),
				resource.TestCheckResourceAttr(resourceName, fmt.Sprintf(`volumes.%d.attachment_tag`, i), volume.AttachmentTag),
			)
		}

		return resource.ComposeTestCheckFunc(checksStore...)(s)
	}
}

func TestAccInstance(t *testing.T) {
	volumes := []instances.CreateVolumeOpts{
		{
			Name:          "boot volume",
			Source:        "image",
			BootIndex:     0,
			Size:          5,
			TypeName:      "ssd_hiiops",
			ImageID:       GCORE_IMAGE,
			AttachmentTag: "some tag",
		},
		{
			Name:      "empty volume",
			Source:    "new-volume",
			BootIndex: 1,
			Size:      5,
			TypeName:  "ssd_hiiops",
		},
	}
	interfaces := []instances.InterfaceOpts{
		{
			Type:      "subnet",
			NetworkID: GCORE_PRIV_NET,
			SubnetID:  GCORE_PRIV_SUBNET,
		},
	}
	secgroups := []gcorecloud.ItemID{
		{
			ID: GCORE_SECGROUP,
		},
	}
	metadata := instances.MetadataSetOpts{}
	metadata.Metadata = []instances.MetadataOpts{
		{
			Key:   "somekey",
			Value: "somevalue",
		},
	}

	createFixt := instances.CreateOpts{
		Names:          []string{"create_instance"},
		NameTemplates:  []string{},
		Flavor:         "g1-standard-2-4",
		Password:       "password",
		Username:       "user",
		Keypair:        "mykey",
		Volumes:        volumes,
		Interfaces:     interfaces,
		SecurityGroups: secgroups,
		Metadata:       &metadata,
		Configuration:  &metadata,
	}

	type Params struct {
		Name           []string
		Flavor         string
		Volumes        []map[string]string
		Interfaces     []map[string]string
		SecurityGroups []map[string]string
		MetaData       []map[string]string
		Configuration  []map[string]string
	}

	create := Params{
		Name:   []string{"create_instance"},
		Flavor: "g1-standard-2-4",
		Volumes: []map[string]string{
			{"name": "boot volume", "source": "image", "boot_index": "0", "size": "5",
				"type_name": "ssd_hiiops", "image_id": GCORE_IMAGE, "attachment_tag": "some tag"},
			{"name": "empty volume", "source": "new-volume", "boot_index": "1", "size": "5",
				"type_name": "ssd_hiiops", "image_id": "", "attachment_tag": ""},
		},
		Interfaces: []map[string]string{
			{"type": "subnet", "network_id": GCORE_PRIV_NET, "subnet_id": GCORE_PRIV_SUBNET},
		},
		SecurityGroups: []map[string]string{{"id": GCORE_SECGROUP, "name": "default"}},
		MetaData:       []map[string]string{{"key": "somekey", "value": "somevalue"}},
		Configuration:  []map[string]string{{"key": "somekey", "value": "somevalue"}},
	}

	instanceTemplate := func(params *Params) string {
		template := `
		locals {`

		template += fmt.Sprintf(`
			names = ["%s"]
			volumes = [`, params.Name[0])

		for i, _ := range params.Volumes {
			template += fmt.Sprintf(`
			{
				source = "%s"
				type_name = "%s"
				size = %s
				name = "%s"
				boot_index = %s
				`, params.Volumes[i]["source"], params.Volumes[i]["type_name"], params.Volumes[i]["size"],
				params.Volumes[i]["name"], params.Volumes[i]["boot_index"])
			if params.Volumes[i]["image_id"] != "" && params.Volumes[i]["attachment_tag"] != "" {
				template += fmt.Sprintf(`
				image_id = "%s"
				attachment_tag = "%s"
			},`, params.Volumes[i]["image_id"], params.Volumes[i]["attachment_tag"])
			} else {
				template += fmt.Sprintf(`
				image_id = %s
				attachment_tag = %s
			},`, "null", "null")
			}
		}

		template += fmt.Sprintf(`]
			interfaces = [`)
		for i, _ := range params.Interfaces {
			if params.Interfaces[i]["network_id"] != "" {
				template += fmt.Sprintf(`
			{
				type = "%s"
				network_id = "%s"
				subnet_id = "%s"
				floating_ip = [
				{
					source = null
					existing_floating_id = null
				},
			]
			},`, params.Interfaces[i]["type"], params.Interfaces[i]["network_id"], params.Interfaces[i]["subnet_id"])
			} else {
				template += fmt.Sprintf(`
			{
				type = "%s"
				subnet_id = "%s"
				network_id = null
			},`, params.Interfaces[i]["type"], params.Interfaces[i]["subnet_id"])
			}
		}
		template += fmt.Sprintf(`]
			security_groups = [`)
		for i, _ := range params.SecurityGroups {
			template += fmt.Sprintf(`
			{
				id = "%s"
				name = "%s"
			},`, params.SecurityGroups[i]["id"], params.SecurityGroups[i]["name"])
		}
		template += fmt.Sprintf(`]
			metadata = [`)
		for i, _ := range params.MetaData {
			template += fmt.Sprintf(`
			{
				key = "%s"
				value = "%s"
			},`, params.MetaData[i]["key"], params.MetaData[i]["value"])
		}
		template += fmt.Sprintf(`]
			configuration = [`)
		for i, _ := range params.Configuration {
			template += fmt.Sprintf(`
			{
				key = "%s"
				value = "%s"
			},`, params.Configuration[i]["key"], params.Configuration[i]["value"])
		}
		template += fmt.Sprintf(`]
        }
        resource "gcore_instance" "acctest" {
           flavor_id =  "%s"
           name = local.names

           dynamic volumes {
           iterator = vol
           for_each = local.volumes
           content {
           		boot_index = vol.value.boot_index
    			source = vol.value.source
                type_name = vol.value.type_name
                size = vol.value.size
                name = vol.value.name
                image_id = vol.value.image_id
                attachment_tag = vol.value.attachment_tag
    			}
  			}

        	dynamic interfaces {
			iterator = ifaces
			for_each = local.interfaces
			content {
				type = ifaces.value.type
                network_id = ifaces.value.network_id
				subnet_id = ifaces.value.subnet_id
                dynamic floating_ip {
                iterator = fip
                for_each = ifaces.value.floating_ip
                content {
                	source = fip.value.source
                    existing_floating_id = fip.value.existing_floating_id
      			    }
                }
                }
  			}

			dynamic security_groups {
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

            %s
			%s

		`, params.Flavor, regionInfo(), projectInfo())
		log.Println(template)
		return template + "\n}"
	}

	fullName := "gcore_instance.acctest"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckInstance(t) },
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
		},
	})
}

func testAccInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, InstancePoint)
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
