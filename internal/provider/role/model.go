package role

import (
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type RoleResourceModel struct {
    UserRoleID       types.String `tfsdk:"user_role_id"`
    Name             types.String `tfsdk:"name"`
    GroupID          types.String `tfsdk:"group_id"`
    Description      types.String `tfsdk:"description"`
    Policies         types.List   `tfsdk:"policies"`
    IsSystemRole     types.Bool   `tfsdk:"is_system_role"`
    SsoGroupMapping  types.List   `tfsdk:"sso_group_mapping"`
    ObjectVersion    types.Int64  `tfsdk:"object_version"`
    CreationUserID   types.String `tfsdk:"creation_user_id"`
    UpdateUserID     types.String `tfsdk:"update_user_id"`
    CreationTimestamp types.String `tfsdk:"creation_timestamp"`
    UpdateTimestamp  types.String `tfsdk:"update_timestamp"`
}

type Policy struct {
    PolicyID   string   `json:"policy_id"`
    Name       string   `json:"name"`
    Action     []string `json:"action"`
    Resource   []string `json:"resource"`
    Description *string `json:"description"`
}

type CreateRoleAPIResponse struct {
    ObjectVersion    int64    `json:"object_version"`
    CreationUserID   string   `json:"creation_user_id"`
    UpdateUserID     *string  `json:"update_user_id"`
    CreationTimestamp string  `json:"creation_timestamp"`
    UpdateTimestamp  *string  `json:"update_timestamp"`
    UserRoleID       string   `json:"user_role_id"`
    Name             string   `json:"name"`
    GroupID          *string  `json:"group_id"`
    Description      *string  `json:"description"`
    Policies         []Policy `json:"policies"`
    IsSystemRole     bool     `json:"is_system_role"`
    SsoGroupMapping  []string `json:"sso_group_mapping"`
}

type RoleAPIResponse = CreateRoleAPIResponse
