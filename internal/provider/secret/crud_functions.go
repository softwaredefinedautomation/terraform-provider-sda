package secret

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

// Create - create a secret
func (r *SecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SecretResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := map[string]interface{}{}
	// required fields
	payload["name"] = plan.Name.ValueString()

    // required fields
    payload["username"] = plan.Username.ValueString()

	// required fields
	payload["secret_type"] = plan.Type.ValueString()

	// Value is required (sensitive)
	payload["secret_value"] = plan.Value.ValueString()

	// optional vault_id
	if !plan.VaultID.IsUnknown() && !plan.VaultID.IsNull() {
		payload["vault_id"] = plan.VaultID.ValueString()
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	url := fmt.Sprintf("%s/assets/v1/secret", r.client.HostURL)

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

		resp.Diagnostics.AddError(
			errSummary,
			fmt.Sprintf(
				"POST %s failed: %s\nResponse body:\n%s",
				reqHTTP.URL.String(),
				err,
				string(resBody),
			),
		)
		return
	}

	var apiResp SecretAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	if apiResp.SecretID == "" {
		resp.Diagnostics.AddError(
			"API Response Missing SecretID",
			fmt.Sprintf("The API did not return a secret_id in response: %s", string(resBody)),
		)
		return
	}

	state := SecretResourceModel{
		ObjectVersion:     types.Int64Value(apiResp.ObjectVersion),
		CreationUserID:    types.StringValue(apiResp.CreationUserID),
		UpdateUserID:      types.StringPointerValue(apiResp.UpdateUserID),
		CreationTimestamp: types.StringValue(apiResp.CreationTimestamp),
		UpdateTimestamp:   types.StringPointerValue(apiResp.UpdateTimestamp),
		SecretID:          types.StringValue(apiResp.SecretID),
		Name:              types.StringValue(apiResp.Name),
		Username:          types.StringValue(apiResp.Username),
		// Type:			   types.StringValue(apiResp.Type),
		// Value:             types.StringPointerValue(apiResp.Value),

		Value: plan.Value,
		Type:  plan.Type,
	}

	// Map optional vault_id pointer -> Terraform value
	if apiResp.VaultID == nil || *apiResp.VaultID == "" {
		state.VaultID = types.StringNull()
	} else {
		state.VaultID = types.StringValue(*apiResp.VaultID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read - read secret
func (r *SecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SecretResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/assets/v1/secret/%s", r.client.HostURL, state.SecretID.ValueString())
	reqHTTP, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating read request: %s", err))
		return
	}

	resBody, err := r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		// If resource returns 404 Not Found, remove it from state.
		if strings.Contains(err.Error(), "status: 404") {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"API Error",
			fmt.Sprintf(
				"Error reading secret %s: %s\nResponse body:\n%s",
				state.SecretID.ValueString(),
				err,
				string(resBody),
			),
		)
		return
	}

	var apiResp SecretAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	state = SecretResourceModel{
		ObjectVersion:     types.Int64Value(apiResp.ObjectVersion),
		CreationUserID:    types.StringValue(apiResp.CreationUserID),
		UpdateUserID:      types.StringPointerValue(apiResp.UpdateUserID),
		CreationTimestamp: types.StringValue(apiResp.CreationTimestamp),
		UpdateTimestamp:   types.StringPointerValue(apiResp.UpdateTimestamp),
		SecretID:          types.StringValue(apiResp.SecretID),
		Name:              types.StringValue(apiResp.Name),
		Username:          types.StringValue(apiResp.Username),
		Type:              types.StringValue(apiResp.Type),
		Value:             state.Value,
	}

	if apiResp.VaultID == nil || *apiResp.VaultID == "" {
		state.VaultID = types.StringNull()
	} else {
		state.VaultID = types.StringValue(*apiResp.VaultID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update - update secret
func (r *SecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state SecretResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := map[string]interface{}{
		"object_version": state.ObjectVersion.ValueInt64(),
	}

	// If name changed
	if !plan.Name.Equal(state.Name) {
		if plan.Name.IsNull() {
			payload["name"] = nil
		} else {
			payload["name"] = plan.Name.ValueString()
		}
	}

	// If username changed
	if !plan.Username.Equal(state.Username) {
		if plan.Username.IsNull() {
			payload["username"] = nil
		} else {
			payload["username"] = plan.Username.ValueString()
		}
	}

	// If value changed
	if !plan.Value.Equal(state.Value) {
		if plan.Value.IsNull() {
			payload["secret_value"] = nil
		} else {
			payload["secret_value"] = plan.Value.ValueString()
		}
	}

	// If value changed
	if !plan.Type.Equal(state.Type) {
		if plan.Type.IsNull() {
			payload["secret_type"] = nil
		} else {
			payload["secret_type"] = plan.Type.ValueString()
		}
	}

	// vault_id change
	if !plan.VaultID.Equal(state.VaultID) {
		if plan.VaultID.IsNull() {
			payload["vault_id"] = ""
		} else {
			payload["vault_id"] = plan.VaultID.ValueString()
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	url := fmt.Sprintf("%s/assets/v1/secret/%s", r.client.HostURL, state.SecretID.ValueString())

	reqHTTP, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating update request: %s", err))
		return
	}
	reqHTTP.Header.Set("Content-Type", "application/json")

	resBody, err := r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error updating secret: %s", err))
		return
	}

	var apiResp SecretAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	// preserve nulls if user explicitly set them in plan
	if plan.VaultID.IsNull() {
		apiResp.VaultID = nil
	}
	if plan.Value.IsNull() {
		apiResp.Value = nil
	}

	state = SecretResourceModel{
		ObjectVersion:     types.Int64Value(apiResp.ObjectVersion),
		CreationUserID:    types.StringValue(apiResp.CreationUserID),
		UpdateUserID:      types.StringPointerValue(apiResp.UpdateUserID),
		CreationTimestamp: types.StringValue(apiResp.CreationTimestamp),
		UpdateTimestamp:   types.StringPointerValue(apiResp.UpdateTimestamp),
		SecretID:          types.StringValue(apiResp.SecretID),
		Name:              types.StringValue(apiResp.Name),
		Username:          types.StringValue(apiResp.Username),
		Type:              types.StringValue(apiResp.Type),
		Value:             state.Value,
	}

	if apiResp.VaultID == nil || *apiResp.VaultID == "" {
		state.VaultID = types.StringNull()
	} else {
		state.VaultID = types.StringValue(*apiResp.VaultID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Delete - delete secret
func (r *SecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SecretResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/assets/v1/secret/%s", r.client.HostURL, state.SecretID.ValueString())
	reqHTTP, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating delete request: %s", err))
		return
	}

	_, err = r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		// Ignore 404 during deletion, as the resource is already gone.
		if strings.Contains(err.Error(), "status: 404") {
			return
		}

		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error deleting secret: %s", err))
	}
}
