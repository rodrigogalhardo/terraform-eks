package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenClusterAlertGroup(d *schema.ResourceData, in *managementClient.ClusterAlertGroup) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)
	d.Set("cluster_id", in.ClusterID)

	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	d.Set("group_interval_seconds", int(in.GroupIntervalSeconds))
	d.Set("group_wait_seconds", int(in.GroupWaitSeconds))
	d.Set("recipients", flattenRecipients(in.Recipients))
	d.Set("repeat_interval_seconds", int(in.RepeatIntervalSeconds))

	err := d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}
	return nil
}

// Expanders

func expandClusterAlertGroup(in *schema.ResourceData) *managementClient.ClusterAlertGroup {
	obj := &managementClient.ClusterAlertGroup{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.ClusterID = in.Get("cluster_id").(string)

	if v, ok := in.Get("description").(string); ok && len(v) > 0 {
		obj.Description = v
	}

	obj.GroupIntervalSeconds = int64(in.Get("group_interval_seconds").(int))
	obj.GroupWaitSeconds = int64(in.Get("group_wait_seconds").(int))

	if v, ok := in.Get("recipients").([]interface{}); ok && len(v) > 0 {
		obj.Recipients = expandRecipients(v)
	}

	obj.RepeatIntervalSeconds = int64(in.Get("repeat_interval_seconds").(int))

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}
