package user

import (
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type UserResourceModel struct {
    UserID            types.String `tfsdk:"user_id"`
    GroupID           types.String `tfsdk:"group_id"`
    FirstName         types.String `tfsdk:"first_name"`
    LastName          types.String `tfsdk:"last_name"`
    Email             types.String `tfsdk:"email"`
    CompanyName       types.String `tfsdk:"company_name"`
    PhoneNumber       types.String `tfsdk:"phone_number"`
    PrivacyAccepted   types.Bool   `tfsdk:"privacy_accepted"`
    Locale            types.String `tfsdk:"locale"`
    LastLoginTimestamp types.String `tfsdk:"last_login_timestamp"`
    Title             types.String `tfsdk:"title"`
    AgreeToContact    types.Bool   `tfsdk:"agree_to_contact"`
    ObjectVersion     types.Int64  `tfsdk:"object_version"`
    CreationUserID    types.String `tfsdk:"creation_user_id"`
    UpdateUserID      types.String `tfsdk:"update_user_id"`
    CreationTimestamp types.String `tfsdk:"creation_timestamp"`
    UpdateTimestamp   types.String `tfsdk:"update_timestamp"`
    Source            types.String `tfsdk:"source"`
}

// CreateUserAPIResponse describes the API response for create user
type CreateUserAPIResponse struct {
    ObjectVersion     int64   `json:"object_version"`
    CreationUserID    string  `json:"creation_user_id"`
    UpdateUserID      *string `json:"update_user_id"`
    CreationTimestamp string  `json:"creation_timestamp"`
    UpdateTimestamp   *string `json:"update_timestamp"`
    UserID            string  `json:"user_id"`
    GroupID           *string `json:"group_id"`
    FirstName         string  `json:"first_name"`
    LastName          string  `json:"last_name"`
    Email             string  `json:"email"`
    CompanyName       *string `json:"company_name"`
    PhoneNumber       *string `json:"phone_number"`
    PrivacyAccepted   *bool   `json:"privacy_accepted"`
    Locale            *string `json:"locale"`
    LastLoginTimestamp *string `json:"last_login_timestamp"`
    Title             *string `json:"title"`
    AgreeToContact    *bool   `json:"agree_to_contact"`
    Source            string  `json:"source"`
}

// UserAPIResponse used for read/update responses
type UserAPIResponse = CreateUserAPIResponse
