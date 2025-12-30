package secret

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SecretResourceModel maps the Terraform resource state
type SecretResourceModel struct {
	SecretID          types.String `tfsdk:"secret_id"`
	VaultID           types.String `tfsdk:"vault_id"`
	Name              types.String `tfsdk:"name"`
	Username          types.String `tfsdk:"username"`
	Value             types.String `tfsdk:"secret_value"`
	Type              types.String `tfsdk:"secret_type"`
	ObjectVersion     types.Int64  `tfsdk:"object_version"`
	CreationUserID    types.String `tfsdk:"creation_user_id"`
	UpdateUserID      types.String `tfsdk:"update_user_id"`
	CreationTimestamp types.String `tfsdk:"creation_timestamp"`
	UpdateTimestamp   types.String `tfsdk:"update_timestamp"`
}

// SecretAPIResponse models the JSON response from the secrets API
type SecretAPIResponse struct {
	ObjectVersion     int64   `json:"object_version"`
	CreationUserID    string  `json:"creation_user_id"`
	UpdateUserID      *string `json:"update_user_id"`
	CreationTimestamp string  `json:"creation_timestamp"`
	UpdateTimestamp   *string `json:"update_timestamp"`
	SecretID          string  `json:"secret_id"`
	VaultID           *string `json:"vault_id"`
	Name              string  `json:"name"`
	Username          string  `json:"username"`
	Value             *string `json:"secret_value"`
	Type              string  `json:"secret_type"`
}
