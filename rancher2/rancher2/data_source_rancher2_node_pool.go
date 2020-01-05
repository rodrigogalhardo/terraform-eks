package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2NodePool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2NodePoolRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"node_template_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"delete_not_ready_after_secs": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"hostname_prefix": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_taints": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: taintFields(),
				},
			},
			"quantity": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"control_plane": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"etcd": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"worker": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"annotations": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2NodePoolRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	nodeTemplateID := d.Get("node_template_id").(string)

	filters := map[string]interface{}{
		"clusterId": clusterID,
		"name":      name,
	}
	if len(nodeTemplateID) > 0 {
		filters["nodeTemplateId"] = nodeTemplateID
	}
	listOpts := NewListOpts(filters)

	nodePools, err := client.NodePool.List(listOpts)
	if err != nil {
		return err
	}

	count := len(nodePools.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] node pool with name \"%s\" on cluster ID \"%s\" not found", name, clusterID)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d node pool with name \"%s\" on cluster ID \"%s\"", count, name, clusterID)
	}

	return flattenNodePool(d, &nodePools.Data[0])
}
