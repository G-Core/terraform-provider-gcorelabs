package gcore

import (
	"fmt"
	"github.com/G-Core/gcorelabscloud-go/gcore/lifecyclepolicy/v1/lifecyclepolicy"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestAccLifecyclePolicy(t *testing.T) {
	// Templates
	resName := "acctest"
	fullLPName := lifecyclePolicyResource + "." + resName
	volumeId := "gcore_volume." + resName + ".id"
	cronScheduleConfig := func(cron lifecyclepolicy.CreateCronScheduleOpts) string {
		return fmt.Sprintf(`
	schedule {
		resource_name_template = "%s"
		max_quantity = %d
		cron {
			timezone = "%s"
			hour = "%s"
		}
	}`, cron.ResourceNameTemplate, cron.MaxQuantity, cron.Timezone, cron.Hour)
	}
	intervalScheduleConfig := func(interval lifecyclepolicy.CreateIntervalScheduleOpts) string {
		return fmt.Sprintf(`
	schedule {
		resource_name_template = "%s"
		max_quantity = %d
		retention_time {
			hours = %d
		}
		interval {
			weeks = %d
		}
	}`, interval.ResourceNameTemplate, interval.MaxQuantity, interval.RetentionTime.Hours, interval.Weeks)
	}
	malformedScheduleConfig := `
	schedule {
		max_quantity = 1
		interval {
			weeks = 1
		}
		cron {
			week = "1"
		}
	}`
	volumeConfig := fmt.Sprintf(`
resource "gcore_volume" "%s" {
	%s
	%s
	name = "test-volume"
	type_name = "standard"
	size = 1
}`, resName, projectInfo(), regionInfo())
	policyConfig := func(opts lifecyclepolicy.CreateOpts, schedules string) string {
		var volumes string
		for _, id := range opts.VolumeIds {
			volumes += fmt.Sprintf(`
	volume {
		id = %s
	}`, id)
		}
		return fmt.Sprintf(`
resource "%s" "%s" {
	%s
	%s
	name = "%s"
	status = "%s"
	%s
	%s
}`, lifecyclePolicyResource, resName, projectInfo(), regionInfo(), opts.Name, opts.Status, volumes, schedules)
	}

	// Options
	create := lifecyclepolicy.CreateOpts{
		Name:      "policy0",
		Status:    lifecyclepolicy.PolicyStatusPaused,
		VolumeIds: []string{},
	}
	update1 := lifecyclepolicy.CreateOpts{
		Name:      "policy1",
		Status:    lifecyclepolicy.PolicyStatusActive,
		VolumeIds: []string{volumeId},
	}
	update2 := lifecyclepolicy.CreateOpts{
		Name:      "policy2",
		Status:    lifecyclepolicy.PolicyStatusActive,
		VolumeIds: []string{},
	}
	cronSchedule := lifecyclepolicy.CreateCronScheduleOpts{
		CommonCreateScheduleOpts: lifecyclepolicy.CommonCreateScheduleOpts{
			ResourceNameTemplate: "template_0",
			MaxQuantity:          3,
		},
		Timezone: "Europe/London",
		Hour:     "2,8",
	}
	intervalSchedule := lifecyclepolicy.CreateIntervalScheduleOpts{
		CommonCreateScheduleOpts: lifecyclepolicy.CommonCreateScheduleOpts{
			ResourceNameTemplate: "template_1",
			MaxQuantity:          4,
			RetentionTime: &lifecyclepolicy.RetentionTimer{
				Hours: 100,
			},
		},
		Weeks: 1,
	}

	// Tests
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccLifecyclePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: volumeConfig + policyConfig(create, cronScheduleConfig(cronSchedule)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullLPName),
					resource.TestCheckResourceAttr(fullLPName, "name", create.Name),
					resource.TestCheckResourceAttr(fullLPName, "status", create.Status.String()),
					resource.TestCheckResourceAttr(fullLPName, "volume.#", "0"),
					resource.TestCheckResourceAttr(fullLPName, "action", lifecyclepolicy.PolicyActionVolumeSnapshot.String()),
					resource.TestCheckResourceAttrSet(fullLPName, "user_id"),
					resource.TestCheckResourceAttr(fullLPName, "schedule.#", "1"),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.max_quantity", strconv.Itoa(cronSchedule.MaxQuantity)),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.interval.#", "0"),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.cron.#", "1"),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.cron.0.timezone", cronSchedule.Timezone),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.cron.0.hour", cronSchedule.Hour),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.cron.0.minute", "0"),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.cron.0.month", "*"),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.resource_name_template", cronSchedule.ResourceNameTemplate),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.retention_time.#", "0"),
					resource.TestCheckResourceAttrSet(fullLPName, "schedule.0.id"),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.type", lifecyclepolicy.ScheduleTypeCron.String()),
				),
			},
			{
				Config: volumeConfig + policyConfig(update1, cronScheduleConfig(cronSchedule)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullLPName),
					resource.TestCheckResourceAttr(fullLPName, "name", update1.Name),
					resource.TestCheckResourceAttr(fullLPName, "status", update1.Status.String()),
					resource.TestCheckResourceAttr(fullLPName, "volume.#", "1"),
				),
			},
			{
				Config: volumeConfig + policyConfig(update2, cronScheduleConfig(cronSchedule)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullLPName),
					resource.TestCheckResourceAttr(fullLPName, "name", update2.Name),
					resource.TestCheckResourceAttr(fullLPName, "volume.#", "0"),
				),
			},
			{ // Delete policy, so we can test another schedule.
				// TODO: For some reason, it doesn't call Create otherwise, even though "schedule" is ForceNew
				Config: volumeConfig,
			},
			{
				Config: volumeConfig + policyConfig(create, intervalScheduleConfig(intervalSchedule)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullLPName),
					resource.TestCheckResourceAttr(fullLPName, "schedule.#", "1"),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.max_quantity", strconv.Itoa(intervalSchedule.MaxQuantity)),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.interval.#", "1"),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.interval.0.weeks", strconv.Itoa(intervalSchedule.Weeks)),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.interval.0.days", "0"),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.cron.#", "0"),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.resource_name_template", intervalSchedule.ResourceNameTemplate),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.retention_time.#", "1"),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.retention_time.0.hours", strconv.Itoa(intervalSchedule.RetentionTime.Hours)),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.retention_time.0.days", "0"),
					resource.TestCheckResourceAttrSet(fullLPName, "schedule.0.id"),
					resource.TestCheckResourceAttr(fullLPName, "schedule.0.type", lifecyclepolicy.ScheduleTypeInterval.String()),
				),
			},
			{ // Delete policy, so we can test another schedule.
				// TODO: For some reason, it doesn't call Create otherwise, even though "schedule" is ForceNew
				Config: volumeConfig,
			},
			{
				Config:      volumeConfig + policyConfig(create, malformedScheduleConfig),
				ExpectError: regexp.MustCompile("exactly one of interval and cron blocks should be provided"),
			},
		},
	})
}

func testAccLifecyclePolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	volumesClient, err := CreateTestClient(config.Provider, volumesPoint, versionPointV1)
	if err != nil {
		return err
	}
	lifecyclePolicyClient, err := CreateTestClient(config.Provider, lifecyclePolicyPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "gcore_volume" {
			_, err := networks.Get(volumesClient, rs.Primary.ID).Extract()
			if err == nil {
				return fmt.Errorf("volume still exists")
			}
			if !strings.Contains(err.Error(), "not found") {
				return err
			}
		} else if rs.Type == lifecyclePolicyResource {
			id, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("error converting lifecycle policy ID to integer: %s", err)
			}
			_, err = lifecyclepolicy.Get(lifecyclePolicyClient, id, lifecyclepolicy.GetOpts{}).Extract()
			if err == nil {
				return fmt.Errorf("policy still exists")
			}
			if !strings.Contains(err.Error(), "not exist") {
				return err
			}
		}
	}
	return nil
}
