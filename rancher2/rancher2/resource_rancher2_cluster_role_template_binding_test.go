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
	testAccRancher2ClusterRoleTemplateBindingType = "rancher2_cluster_role_template_binding"
)

var (
	testAccRancher2ClusterRoleTemplateBindingConfig         string
	testAccRancher2ClusterRoleTemplateBindingUpdateConfig   string
	testAccRancher2ClusterRoleTemplateBindingRecreateConfig string
)

func init() {
	testAccRancher2ClusterRoleTemplateBindingConfig = `
resource "rancher2_cluster_role_template_binding" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  role_template_id = "cluster-admin"
}
`

	testAccRancher2ClusterRoleTemplateBindingUpdateConfig = `
resource "rancher2_cluster_role_template_binding" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  role_template_id = "projects-create"
}
 `

	testAccRancher2ClusterRoleTemplateBindingRecreateConfig = `
resource "rancher2_cluster_role_template_binding" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  role_template_id = "cluster-admin"
}
 `
}

func TestAccRancher2ClusterRoleTemplateBinding_basic(t *testing.T) {
	var clusterRole *managementClient.ClusterRoleTemplateBinding

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterRoleTemplateBindingDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2ClusterRoleTemplateBindingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterRoleTemplateBindingExists(testAccRancher2ClusterRoleTemplateBindingType+".foo", clusterRole),
					resource.TestCheckResourceAttr(testAccRancher2ClusterRoleTemplateBindingType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterRoleTemplateBindingType+".foo", "cluster_id", testAccRancher2ClusterID),
					resource.TestCheckResourceAttr(testAccRancher2ClusterRoleTemplateBindingType+".foo", "role_template_id", "cluster-admin"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ClusterRoleTemplateBindingUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterRoleTemplateBindingExists(testAccRancher2ClusterRoleTemplateBindingType+".foo", clusterRole),
					resource.TestCheckResourceAttr(testAccRancher2ClusterRoleTemplateBindingType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterRoleTemplateBindingType+".foo", "cluster_id", testAccRancher2ClusterID),
					resource.TestCheckResourceAttr(testAccRancher2ClusterRoleTemplateBindingType+".foo", "role_template_id", "projects-create"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ClusterRoleTemplateBindingRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterRoleTemplateBindingExists(testAccRancher2ClusterRoleTemplateBindingType+".foo", clusterRole),
					resource.TestCheckResourceAttr(testAccRancher2ClusterRoleTemplateBindingType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ClusterRoleTemplateBindingType+".foo", "cluster_id", testAccRancher2ClusterID),
					resource.TestCheckResourceAttr(testAccRancher2ClusterRoleTemplateBindingType+".foo", "role_template_id", "cluster-admin"),
				),
			},
		},
	})
}

func TestAccRancher2ClusterRoleTemplateBinding_disappears(t *testing.T) {
	var clusterRole *managementClient.ClusterRoleTemplateBinding

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ClusterRoleTemplateBindingDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2ClusterRoleTemplateBindingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ClusterRoleTemplateBindingExists(testAccRancher2ClusterRoleTemplateBindingType+".foo", clusterRole),
					testAccRancher2ClusterRoleTemplateBindingDisappears(clusterRole),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2ClusterRoleTemplateBindingDisappears(pro *managementClient.ClusterRoleTemplateBinding) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2ClusterRoleTemplateBindingType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			pro, err = client.ClusterRoleTemplateBinding.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.ClusterRoleTemplateBinding.Delete(pro)
			if err != nil {
				return fmt.Errorf("Error removing Cluster Role Template Binding: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active"},
				Target:     []string{"removed"},
				Refresh:    clusterRoleTemplateBindingStateRefreshFunc(client, pro.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for Cluster Role Template Binding (%s) to be removed: %s", pro.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2ClusterRoleTemplateBindingExists(n string, pro *managementClient.ClusterRoleTemplateBinding) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cluster Role Template Binding ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundPro, err := client.ClusterRoleTemplateBinding.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Cluster Role Template Binding not found")
			}
			return err
		}

		pro = foundPro

		return nil
	}
}

func testAccCheckRancher2ClusterRoleTemplateBindingDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2ClusterRoleTemplateBindingType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.ClusterRoleTemplateBinding.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Cluster Role Template Binding still exists")
	}
	return nil
}
