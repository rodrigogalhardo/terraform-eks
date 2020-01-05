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
	testAccRancher2MultiClusterAppType = "rancher2_multi_cluster_app"
)

var (
	testAccRancher2MultiClusterAppProject        string
	testAccRancher2MultiClusterAppNamespace      string
	testAccRancher2MultiClusterAppConfig         string
	testAccRancher2MultiClusterAppUpdateConfig   string
	testAccRancher2MultiClusterAppRecreateConfig string
)

func init() {
	testAccRancher2MultiClusterAppProject = `
resource "rancher2_project" "foo" {
  name = "foo"
  cluster_id = "` + testAccRancher2ClusterID + `"
  description = "Terraform app acceptance test"
  resource_quota {
    project_limit {
      limits_cpu = "2000m"
      limits_memory = "2000Mi"
      requests_storage = "2Gi"
    }
    namespace_default_limit {
      limits_cpu = "500m"
      limits_memory = "500Mi"
      requests_storage = "1Gi"
    }
  }
}
`

	testAccRancher2MultiClusterAppConfig = testAccRancher2MultiClusterAppProject + `
resource "rancher2_multi_cluster_app" "foo" {
  catalog_name = "library"
  name = "foo"
  targets {
    project_id = "${rancher2_project.foo.id}"
  }
  template_name = "docker-registry"
  template_version = "1.8.1"
  answers {
    values = {
      "ingress_host" = "test.xip.io"
    }
  }
  roles = ["project-member"]
}
`

	testAccRancher2MultiClusterAppUpdateConfig = testAccRancher2MultiClusterAppProject + `
resource "rancher2_multi_cluster_app" "foo" {
  catalog_name = "library"
  name = "foo"
  targets {
    project_id = "${rancher2_project.foo.id}"
  }
  template_name = "docker-registry"
  template_version = "1.8.1"
  answers {
    values = {
      "ingress_host" = "test2.xip.io"
    }
  }
  roles = ["cluster-admin"]
}
`

	testAccRancher2MultiClusterAppRecreateConfig = testAccRancher2MultiClusterAppProject + `
resource "rancher2_multi_cluster_app" "foo" {
  catalog_name = "library"
  name = "foo"
  targets {
    project_id = "${rancher2_project.foo.id}"
  }
  template_name = "docker-registry"
  template_version = "1.8.1"
  answers {
    values = {
      "ingress_host" = "test.xip.io"
    }
  }
  roles = ["project-member"]
}
`
}

func TestAccRancher2MultiClusterApp_basic(t *testing.T) {
	var app *managementClient.MultiClusterApp

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2MultiClusterAppDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2MultiClusterAppConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2MultiClusterAppExists(testAccRancher2MultiClusterAppType+".foo", app),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "template_version_id", "cattle-global-data:library-docker-registry-1.8.1"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "answers.0.values.ingress_host", "test.xip.io"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "roles.0", "project-member"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2MultiClusterAppUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2MultiClusterAppExists(testAccRancher2MultiClusterAppType+".foo", app),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "template_version_id", "cattle-global-data:library-docker-registry-1.8.1"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "answers.0.values.ingress_host", "test2.xip.io"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "roles.0", "cluster-admin"),
				),
			},
			resource.TestStep{
				Config: testAccRancher2MultiClusterAppRecreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2MultiClusterAppExists(testAccRancher2MultiClusterAppType+".foo", app),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "name", "foo"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "template_version_id", "cattle-global-data:library-docker-registry-1.8.1"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "answers.0.values.ingress_host", "test.xip.io"),
					resource.TestCheckResourceAttr(testAccRancher2MultiClusterAppType+".foo", "roles.0", "project-member"),
				),
			},
		},
	})
}

func TestAccRancher2MultiClusterApp_disappears(t *testing.T) {
	var app *managementClient.MultiClusterApp

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancher2MultiClusterAppDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancher2MultiClusterAppConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancher2MultiClusterAppExists(testAccRancher2MultiClusterAppType+".foo", app),
					testAccRancher2MultiClusterAppDisappears(app),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRancher2MultiClusterAppDisappears(mca *managementClient.MultiClusterApp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != testAccRancher2MultiClusterAppType {
				continue
			}

			client, err := testAccProvider.Meta().(*Config).ManagementClient()
			if err != nil {
				return err
			}

			mca, err := client.MultiClusterApp.ByID(rs.Primary.ID)
			if err != nil {
				if IsNotFound(err) {
					return nil
				}
				return err
			}

			err = client.MultiClusterApp.Delete(mca)
			if err != nil {
				return fmt.Errorf("Error removing multi cluster app: %s", err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"removing"},
				Target:     []string{"removed"},
				Refresh:    multiClusterAppStateRefreshFunc(client, rs.Primary.ID),
				Timeout:    10 * time.Minute,
				Delay:      1 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, waitErr := stateConf.WaitForState()
			if waitErr != nil {
				return fmt.Errorf(
					"[ERROR] waiting for multi cluster app (%s) to be removed: %s", rs.Primary.ID, waitErr)
			}
		}
		return nil
	}
}

func testAccCheckRancher2MultiClusterAppExists(n string, mca *managementClient.MultiClusterApp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No multi cluster app ID is set")
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		foundMultiClusterApp, err := client.MultiClusterApp.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return fmt.Errorf("Multi cluster app not found")
			}
			return err
		}

		mca = foundMultiClusterApp

		return nil
	}
}

func testAccCheckRancher2MultiClusterAppDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccRancher2MultiClusterAppType {
			continue
		}

		client, err := testAccProvider.Meta().(*Config).ManagementClient()
		if err != nil {
			return err
		}

		_, err = client.MultiClusterApp.ByID(rs.Primary.ID)
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("Multi cluster app still exists")
	}
	return nil
}
