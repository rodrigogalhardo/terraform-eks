package rancher2

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	testAccRancher2NodeDriverType   = "rancher2_node_driver"
	testAccRancher2NodeDriverConfig = `
resource "rancher2_node_driver" "foo" {
    active = false
    builtin = false
    checksum = "0x0"
    description = "Foo description"
    external_id = "foo_external"
    name = "foo"
    ui_url = "local://ui"
    url = "local://"
	whitelist_domains = ["*.foo.com"]
	annotations = {
		foo = "bar"
	}
	labels = {
		foo = "baz"
	}
}
`
	testAccRancher2NodeDriverUpdateConfig = `
resource "rancher2_node_driver" "foo" {
    active = false
    builtin = false
    checksum = "0x1"
    description= "Foo description - updated"
    external_id = "external"
    name = "foo"
    ui_url = "local://ui/updated"
    url = "local://updated"
    whitelist_domains = ["*.foo.com", "updated.foo.com"]
	annotations = {
		foo = "updated"
		bar = "added"
	}
	labels = {
		foo = "updated"
		bar = "added"
	}
}
 `
	testAccRancher2NodeDriverRecreateConfig = `
resource "rancher2_node_driver" "foo" {
    active = false
    builtin = false
    checksum = "0x0"
    description = "Foo description"
    external_id = "foo_external"
    name = "foo"
    ui_url = "local://ui"
    url = "local://"
	whitelist_domains = ["*.foo.com"]
	annotations = {
		foo = "bar"
	}
	labels = {
		foo = "baz"
	}
}
`
)

func TestAccRancher2NodeDriver_basic(t *testing.T) {
	var nodeDriver *managementClient.NodeDriver
	name := testAccRancher2NodeDriverType + ".foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeDriverDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeDriverConfig,
				// Some annotation and labels are computed, as such the
				// subsequent plan would not be empty
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeDriverExists(name, nodeDriver),
					resource.TestCheckResourceAttr(name, "active", "false"),
					resource.TestCheckResourceAttr(name, "builtin", "false"),
					resource.TestCheckResourceAttr(name, "checksum", "0x0"),
					resource.TestCheckResourceAttr(name, "description", "Foo description"),
					resource.TestCheckResourceAttr(name, "external_id", "foo_external"),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "ui_url", "local://ui"),
					resource.TestCheckResourceAttr(name, "url", "local://"),
					resource.TestCheckResourceAttr(name, "whitelist_domains.0", "*.foo.com"),
					resource.TestCheckResourceAttr(name, "annotations.foo", "bar"),
					resource.TestCheckResourceAttr(name, "labels.foo", "baz"),
					resource.TestCheckResourceAttr(name, "labels.cattle.io/creator", "norman"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeDriverUpdateConfig,
				// Some annotation and labels are computed, as such the
				// subsequent plan would not be empty
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeDriverExists(name, nodeDriver),
					resource.TestCheckResourceAttr(name, "active", "false"),
					resource.TestCheckResourceAttr(name, "builtin", "false"),
					resource.TestCheckResourceAttr(name, "checksum", "0x1"),
					resource.TestCheckResourceAttr(name, "description", "Foo description - updated"),
					resource.TestCheckResourceAttr(name, "external_id", "external"),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "ui_url", "local://ui/updated"),
					resource.TestCheckResourceAttr(name, "url", "local://updated"),
					resource.TestCheckResourceAttr(name, "whitelist_domains.0", "*.foo.com"),
					resource.TestCheckResourceAttr(name, "whitelist_domains.1", "updated.foo.com"),
					resource.TestCheckResourceAttr(name, "annotations.foo", "updated"),
					resource.TestCheckResourceAttr(name, "annotations.bar", "added"),
					resource.TestCheckResourceAttr(name, "annotations.lifecycle.cattle.io/create.node-driver-controller", "true"),
					resource.TestCheckResourceAttr(name, "labels.foo", "updated"),
					resource.TestCheckResourceAttr(name, "labels.bar", "added"),
					resource.TestCheckResourceAttr(name, "labels.cattle.io/creator", "norman"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2NodeDriverRecreateConfig,
				// Some annotation and labels are computed, as such the
				// subsequent plan would not be empty
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeDriverExists(name, nodeDriver),
					resource.TestCheckResourceAttr(name, "active", "false"),
					resource.TestCheckResourceAttr(name, "builtin", "false"),
					resource.TestCheckResourceAttr(name, "checksum", "0x0"),
					resource.TestCheckResourceAttr(name, "description", "Foo description"),
					resource.TestCheckResourceAttr(name, "external_id", "foo_external"),
					resource.TestCheckResourceAttr(name, "name", "foo"),
					resource.TestCheckResourceAttr(name, "ui_url", "local://ui"),
					resource.TestCheckResourceAttr(name, "url", "local://"),
					resource.TestCheckResourceAttr(name, "whitelist_domains.0", "*.foo.com"),
					resource.TestCheckResourceAttr(name, "annotations.foo", "bar"),
					resource.TestCheckResourceAttr(name, "labels.foo", "baz"),
					resource.TestCheckResourceAttr(name, "labels.cattle.io/creator", "norman"),
				),
			},
		},
	})
}

func TestAccRancher2NodeDriver_disappears(t *testing.T) {
	var nodeDriver *managementClient.NodeDriver

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2NodeDriverDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2NodeDriverConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2NodeDriverExists(testAccRancher2NodeDriverType+".foo", nodeDriver),
					testAccRancher2NodeDriverDisappears(nodeDriver),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2NodeDriverDisappears(nodeDriver *managementClient.NodeDriver) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2NodeDriverType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			nodeDriver, err := client.NodeDriver.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.NodeDriver.Delete(nodeDriver)
			if err != nil {
				return fmt.Errorf("Error removing Node Driver: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"removing"},
				Target:     []string{"removed"},
				Refresh:    nodeDriverStateRefreshFunc(client, nodeDriver.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for node driver (%s) to be removed: %s", nodeDriver.ID, waitErr)
			}
		}
		return nil
	}
}

func testAccCheckRancher2NodeDriverExists(n string, nodeDriver *managementClient.NodeDriver) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Node Driver ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundNodeDriver, err := client.NodeDriver.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Node Driver not found")
			}
			return err
		}

		nodeDriver = foundNodeDriver

		return nil
	}
}

func testAccCheckRancher2NodeDriverDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2NodeDriverType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.NodeDriver.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Node Driver still exists")
	}
	return nil
}
