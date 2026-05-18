package user_role_association

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure UserRoleAssociationResource implements CRUD interfaces
var _ resource.Resource = &UserRoleAssociationResource{}

// CREATE
func (r *UserRoleAssociationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserRoleAssociationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build URL
	url := fmt.Sprintf("%s/ident/v1/user_role_user_link/user/%s/user_role/%s", r.client.HostURL, plan.UserID.ValueString(), plan.UserRoleID.ValueString())

	// Create Payload body
	payload := map[string]string{}

	if !plan.ExpirationTimestamp.IsUnknown() && !plan.ExpirationTimestamp.IsNull() {
		payload["expiration_timestamp"] = plan.ExpirationTimestamp.ValueString()
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	reqHTTP, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating request: %s", err))
		return
	}
	reqHTTP.Header.Set("Content-Type", "application/json")

	resBody, err := r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		errSummary := "API Error"
		if strings.Contains(err.Error(), "status: 403") {
			errSummary = "API Error: Forbidden (Check Authentication and Permissions/Token Scope)"
		}

		resp.Diagnostics.AddError(errSummary, fmt.Sprintf("POST %s failed: %s\nResponse body:\n%s", reqHTTP.URL.String(), err, string(resBody)))
		return
	}

	if len(resBody) == 0 {
		resp.Diagnostics.AddError("API Error", "API returned an empty array for create user role to user link")
		return
	}

	// Try to unmarshal response if present
	var apiResp UserRoleAssociationAPIResponse
	if err = json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	state := UserRoleAssociationResourceModel{
		UserID:              types.StringValue(apiResp.UserID),
		UserRoleID:          types.StringValue(apiResp.UserRoleId),
		ObjectVersion:       types.Int64Value(apiResp.ObjectVersion),
		CreationUserID:      types.StringValue(apiResp.CreationUserID),
		UpdateUserID:        types.StringPointerValue(apiResp.UpdateUserID),
		CreationTimestamp:   types.StringValue(apiResp.CreationTimestamp),
		UpdateTimestamp:     types.StringPointerValue(apiResp.UpdateTimestamp),
		ExpirationTimestamp: types.StringPointerValue(apiResp.ExpirationTimestamp),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// READ
func (r *UserRoleAssociationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserRoleAssociationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/ident/v1/user_role_user_link/user/%s/user_role/%s", r.client.HostURL, state.UserID.ValueString(), state.UserRoleID.ValueString())

	reqHTTP, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating read request: %s", err))
		return
	}

	resBody, err := r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		if strings.Contains(err.Error(), "status: 404") {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error reading user-role link (user_id=%s, user_role_id=%s): %s\nResponse body:\n%s", state.UserID.ValueString(), state.UserRoleID.ValueString(), err, string(resBody)))
		return
	}

	if len(resBody) == 0 {
		resp.Diagnostics.AddError("API Error", "API returned an empty array for read user role to user link")
		return
	}

	var apiResp UserRoleAssociationAPIResponse
	if err = json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	state.ObjectVersion = types.Int64Value(apiResp.ObjectVersion)
	state.CreationUserID = types.StringValue(apiResp.CreationUserID)
	state.UpdateUserID = types.StringPointerValue(apiResp.UpdateUserID)
	state.CreationTimestamp = types.StringValue(apiResp.CreationTimestamp)
	state.UpdateTimestamp = types.StringPointerValue(apiResp.UpdateTimestamp)
	state.UserID = types.StringValue(apiResp.UserID)
	state.UserRoleID = types.StringValue(apiResp.UserRoleId)
	state.ExpirationTimestamp = types.StringPointerValue(apiResp.ExpirationTimestamp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// UPDATE
func (r *UserRoleAssociationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state UserRoleAssociationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.UserID.Equal(state.UserID) || !plan.UserRoleID.Equal(state.UserRoleID) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("User role associated to a user can't be updated"))
		return
	}

	payload := map[string]interface{}{
		"object_version": state.ObjectVersion.ValueInt64(),
	}

	if !plan.ExpirationTimestamp.Equal(state.ExpirationTimestamp) {
		if plan.ExpirationTimestamp.IsNull() {
			payload["expiration_timestamp"] = nil
		} else if !plan.ExpirationTimestamp.IsUnknown() {
			payload["expiration_timestamp"] = plan.ExpirationTimestamp.ValueString()
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	url := fmt.Sprintf("%s/ident/v1/user_role_user_link/user/%s/user_role/%s", r.client.HostURL, state.UserID.ValueString(), state.UserRoleID.ValueString())
	reqHTTP, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating update request: %s", err))
		return
	}
	reqHTTP.Header.Set("Content-Type", "application/json")

	resBody, err := r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		if strings.Contains(err.Error(), "status: 404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error updating user-role link: %s\nResponse body:\n%s", err, string(resBody)))
		return
	}

	if len(resBody) == 0 {
		resp.Diagnostics.AddError("API Error", "API returned an empty array for update user role to user link")
		return
	}

	var apiResp UserRoleAssociationAPIResponse
	if err = json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	state.ObjectVersion = types.Int64Value(apiResp.ObjectVersion)
	state.CreationUserID = types.StringValue(apiResp.CreationUserID)
	state.UpdateUserID = types.StringPointerValue(apiResp.UpdateUserID)
	state.CreationTimestamp = types.StringValue(apiResp.CreationTimestamp)
	state.UpdateTimestamp = types.StringPointerValue(apiResp.UpdateTimestamp)
	// Keep configured attributes aligned with the applied plan.
	state.UserID = plan.UserID
	state.UserRoleID = plan.UserRoleID
	state.ExpirationTimestamp = plan.ExpirationTimestamp

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// DELETE
func (r *UserRoleAssociationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserRoleAssociationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/ident/v1/user_role_user_link/user/%s/user_role/%s", r.client.HostURL, state.UserID.ValueString(), state.UserRoleID.ValueString())
	reqHTTP, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating delete request: %s", err))
		return
	}

	_, err = r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		if strings.Contains(err.Error(), "status: 404") {
			return
		}
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error deleting user-role link: %s", err))
	}
}
