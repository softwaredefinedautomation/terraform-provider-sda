package resourcegroup

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResourceGroupResourceModel struct {
	GroupID           types.String `tfsdk:"group_id"`
	Name              types.String `tfsdk:"name"`
	GroupType         types.String `tfsdk:"group_type"`
	ParentGroupID     types.String `tfsdk:"parent_group_id"`
	IsSystemGroup     types.Bool   `tfsdk:"is_system_group"`
	ObjectVersion     types.Int64  `tfsdk:"object_version"`
	CreationUserID    types.String `tfsdk:"creation_user_id"`
	UpdateUserID      types.String `tfsdk:"update_user_id"`
	CreationTimestamp types.String `tfsdk:"creation_timestamp"`
	UpdateTimestamp   types.String `tfsdk:"update_timestamp"`
}

type ResourceGroupAPIResponse struct {
	ObjectVersion     int64   `json:"object_version"`
	CreationUserID    string  `json:"creation_user_id"`
	UpdateUserID      *string `json:"update_user_id"`
	CreationTimestamp string  `json:"creation_timestamp"`
	UpdateTimestamp   *string `json:"update_timestamp"`
	GroupID           string  `json:"group_id"`
	Name              string  `json:"name"`
	GroupType         string  `json:"group_type"`
	ParentGroupID     *string `json:"parent_group_id"`
	IsSystemGroup     bool    `json:"is_system_group"`
}
