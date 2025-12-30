package vault

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type VaultResourceModel struct {
	VaultID           types.String `tfsdk:"vault_id"`
	GroupID           types.String `tfsdk:"group_id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	ObjectVersion     types.Int64  `tfsdk:"object_version"`
	CreationUserID    types.String `tfsdk:"creation_user_id"`
	UpdateUserID      types.String `tfsdk:"update_user_id"`
	CreationTimestamp types.String `tfsdk:"creation_timestamp"`
	UpdateTimestamp   types.String `tfsdk:"update_timestamp"`
}

type VaultAPIResponse struct {
	ObjectVersion     int64   `json:"object_version"`
	CreationUserID    string  `json:"creation_user_id"`
	UpdateUserID      *string `json:"update_user_id"`
	CreationTimestamp string  `json:"creation_timestamp"`
	UpdateTimestamp   *string `json:"update_timestamp"`
	VaultID           string  `json:"vault_id"`
	GroupID           *string `json:"group_id"`
	Name              string  `json:"name"`
	Description       *string `json:"description"`
}
