package tag

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TagResourceModel struct {
	Name              types.String `tfsdk:"name"`
	Color             types.String `tfsdk:"color"`
	Icon              types.String `tfsdk:"icon"`
	ObjectVersion     types.Int64  `tfsdk:"object_version"`
	CreationUserID    types.String `tfsdk:"creation_user_id"`
	UpdateUserID      types.String `tfsdk:"update_user_id"`
	CreationTimestamp types.String `tfsdk:"creation_timestamp"`
	UpdateTimestamp   types.String `tfsdk:"update_timestamp"`
}

type TagAPIResponse struct {
	ObjectVersion     int64   `json:"object_version"`
	CreationUserID    string  `json:"creation_user_id"`
	UpdateUserID      *string `json:"update_user_id"`
	CreationTimestamp string  `json:"creation_timestamp"`
	UpdateTimestamp   *string `json:"update_timestamp"`
	Name              string  `json:"name"`
	Color             *string `json:"color"`
	Icon              *string `json:"icon"`
}
