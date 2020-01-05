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
	testAccRancher2ProjectRoleTemplateBindingType = "rancher2_project_role_template_binding"
)

var (
	testAccRancher2ProjectRoleTemplateBindingProject        string
	testAccRancher2ProjectRoleTemplateBindingConfig         string
	testAccRancher2ProjectRoleTemplateBindingUpdateConfig   string
	testAccRancher2ProjectRoleTemplateBindingRecreateConfig string
)

func init() {
	testAccRancher2ProjectRoleTemplateBindingProject = `
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  description = "Terraform project role template binding acceptance test"
}
`
	testAccRancher2ProjectRoleTemplateBindingConfig = testAccRancher2ProjectRoleTemplateBindingProject + `
resource "rancher2_project_role_template_binding" "foo" {
  name = "foo"
  project_id = "${rancher2_project.foo.id}"
  role_template_id = "project-member"
}
`

	testAccRancher2ProjectRoleTemplateBindingUpdateConfig = testAccRancher2ProjectRoleTemplateBindingProject + `
resource "rancher2_project_role_template_binding" "foo" {
  name = "foo"
  project_id = "${rancher2_project.foo.id}"
  role_template_id = "project-owner"
  user_id = "u-q2wg7"
}
 `

	testAccRancher2ProjectRoleTemplateBindingRecreateConfig = testAccRancher2ProjectRoleTemplateBindingProject + `
resource "rancher2_project_role_template_binding" "foo" {
  name = "foo"
  project_id = "${rancher2_project.foo.id}"
  role_template_id = "project-member"
}
 `
}

func TestAccRancher2ProjectRoleTemplateBinding_basic(t *testing.T) {
	var projectRole *managementClient.ProjectRoleTemplateBinding

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ProjectRoleTemplateBindingDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2ProjectRoleTemplateBindingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectRoleTemplateBindingExists(testAccRancher2ProjectRoleTemplateBindingType+".foo", projectRole),
					resource.TestCheckResourceAttr(testAccRancher2ProjectRoleTemplateBindingType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectRoleTemplateBindingType+".foo", "role_template_id", "project-member"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ProjectRoleTemplateBindingUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectRoleTemplateBindingExists(testAccRancher2ProjectRoleTemplateBindingType+".foo", projectRole),
					resource.TestCheckResourceAttr(testAccRancher2ProjectRoleTemplateBindingType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectRoleTemplateBindingType+".foo", "role_template_id", "project-owner"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectRoleTemplateBindingType+".foo", "user_id", "u-q2wg7"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2ProjectRoleTemplateBindingRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectRoleTemplateBindingExists(testAccRancher2ProjectRoleTemplateBindingType+".foo", projectRole),
					resource.TestCheckResourceAttr(testAccRancher2ProjectRoleTemplateBindingType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2ProjectRoleTemplateBindingType+".foo", "role_template_id", "project-member"),
				),
			},
		},
	})
}

func TestAccRancher2ProjectRoleTemplateBinding_disappears(t *testing.T) {
	var projectRole *managementClient.ProjectRoleTemplateBinding

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2ProjectRoleTemplateBindingDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2ProjectRoleTemplateBindingConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2ProjectRoleTemplateBindingExists(testAccRancher2ProjectRoleTemplateBindingType+".foo", projectRole),
					testAccRancher2ProjectRoleTemplateBindingDisappears(projectRole),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2ProjectRoleTemplateBindingDisappears(pro *managementClient.ProjectRoleTemplateBinding) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2ProjectRoleTemplateBindingType {
				continue
			}
			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			pro, err = client.ProjectRoleTemplateBinding.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.ProjectRoleTemplateBinding.Delete(pro)
			if err != nil {
				return fmt.Errorf("Error removing Project Role Template Binding: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"active"},
				Target:     []string{"removed"},
				Refresh:    projectRoleTemplateBindingStateRefreshFunc(client, pro.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for Project Role Template Binding (%s) to be removed: %s", pro.ID, waitErr)
			}
		}
		return nil

	}
}

func testAccCheckRancher2ProjectRoleTemplateBindingExists(n string, pro *managementClient.ProjectRoleTemplateBinding) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Project Role Template Binding ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundPro, err := client.ProjectRoleTemplateBinding.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Project Role Template Binding not found")
			}
			return err
		}

		pro = foundPro

		return nil
	}
}

func testAccCheckRancher2ProjectRoleTemplateBindingDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2ProjectRoleTemplateBindingType {
			continue
		}
		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.ProjectRoleTemplateBinding.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Project Role Template Binding still exists")
	}
	return nil
}
