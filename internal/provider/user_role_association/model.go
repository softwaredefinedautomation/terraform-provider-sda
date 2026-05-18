package user_role_association

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceModel maps resource schema attributes to Go types for CRUD operations.
type UserRoleAssociationResourceModel struct {
	UserID              types.String `tfsdk:"user_id"`
	UserRoleID          types.String `tfsdk:"user_role_id"`
	ExpirationTimestamp types.String `tfsdk:"expiration_timestamp"`
	ObjectVersion       types.Int64  `tfsdk:"object_version"`
	CreationUserID      types.String `tfsdk:"creation_user_id"`
	UpdateUserID        types.String `tfsdk:"update_user_id"`
	CreationTimestamp   types.String `tfsdk:"creation_timestamp"`
	UpdateTimestamp     types.String `tfsdk:"update_timestamp"`
}

// API response expected from create association endpoints.
type CreateUserRoleAssociationAPIResponse struct {
	ObjectVersion       int64   `json:"object_version"`
	CreationUserID      string  `json:"creation_user_id"`
	UpdateUserID        *string `json:"update_user_id"`
	DeleteUserId        *string `json:"delete_user_id"`
	CreationTimestamp   string  `json:"creation_timestamp"`
	UpdateTimestamp     *string `json:"update_timestamp"`
	UserRoleId          string  `json:"user_role_id"`
	UserID              string  `json:"user_id"`
	ExpirationTimestamp *string `json:"expiration_timestamp"`
}

type UserRoleAssociationAPIResponse = CreateUserRoleAssociationAPIResponse
