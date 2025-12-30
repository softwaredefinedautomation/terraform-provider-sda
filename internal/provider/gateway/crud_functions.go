package gateway

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

//-----------------------------------------------------------------
//         CREATE
//-----------------------------------------------------------------

func (r *GatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan GatewayResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := map[string]interface{}{}
	// Name is required, so always include it
	payload["name"] = plan.Name.ValueString()

	// Include optional fields only if known and non-null
	if !plan.GroupID.IsUnknown() && !plan.GroupID.IsNull() {
		payload["group_id"] = plan.GroupID.ValueString()
	}
	if !plan.Description.IsUnknown() && !plan.Description.IsNull() {
		payload["description"] = plan.Description.ValueString()
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	url := fmt.Sprintf("%s/assets/v1/gateway", r.client.HostURL)

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

	var apiResp GatewayAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	if apiResp.GatewayID == "" {
		resp.Diagnostics.AddError(
			"API Response Missing GatewayID",
			fmt.Sprintf("The API did not return a gateway_id in response: %s", string(resBody)),
		)
		return
	}

	state := GatewayResourceModel{
		ObjectVersion:     types.Int64Value(apiResp.ObjectVersion),
		CreationUserID:    types.StringValue(apiResp.CreationUserID),
		UpdateUserID:      types.StringPointerValue(apiResp.UpdateUserID),
		CreationTimestamp: types.StringValue(apiResp.CreationTimestamp),
		UpdateTimestamp:   types.StringPointerValue(apiResp.UpdateTimestamp),
		GatewayID:         types.StringValue(apiResp.GatewayID),
		GroupID:           types.StringPointerValue(apiResp.GroupID),
		Name:              types.StringValue(apiResp.Name),
		Description:       types.StringPointerValue(apiResp.Description),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         READ
//-----------------------------------------------------------------
func (r *GatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state GatewayResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/assets/v1/gateway/%s", r.client.HostURL, state.GatewayID.ValueString())
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
				"Error reading gateway %s: %s\nResponse body:\n%s",
				state.GatewayID.ValueString(),
				err,
				string(resBody),
			),
		)
		return
	}

	var apiResp GatewayAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	state = GatewayResourceModel{
		ObjectVersion:     types.Int64Value(apiResp.ObjectVersion),
		CreationUserID:    types.StringValue(apiResp.CreationUserID),
		UpdateUserID:      types.StringPointerValue(apiResp.UpdateUserID),
		CreationTimestamp: types.StringValue(apiResp.CreationTimestamp),
		UpdateTimestamp:   types.StringPointerValue(apiResp.UpdateTimestamp),
		GatewayID:         types.StringValue(apiResp.GatewayID),
		GroupID:           types.StringPointerValue(apiResp.GroupID),
		Name:              types.StringValue(apiResp.Name),
		Description:       types.StringPointerValue(apiResp.Description),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         UPDATE
//-----------------------------------------------------------------
func (r *GatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state GatewayResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := map[string]interface{}{
		"object_version": state.ObjectVersion.ValueInt64(),
	}

	// Include name if changed
	if !plan.Name.Equal(state.Name) {
		if plan.Name.IsNull() {
			payload["name"] = nil
		} else {
			payload["name"] = plan.Name.ValueString()
		}
	}

	// Include group_id if changed
	if !plan.GroupID.Equal(state.GroupID) {
		if plan.GroupID.IsNull() {
			payload["group_id"] = ""
		} else {
			payload["group_id"] = plan.GroupID.ValueString()
		}
	}

	// if plan.GroupID.IsNull() {
	// 	// Explicitly clear the group_id when user removed it
	// 	payload["group_id"] = nil
	// } else {
	// 	payload["group_id"] = plan.GroupID.ValueString()
	// }

	// Include description if changed
	if !plan.Description.Equal(state.Description) {
		if plan.Description.IsNull() {
			payload["description"] = nil
		} else {
			payload["description"] = plan.Description.ValueString()
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	url := fmt.Sprintf("%s/assets/v1/gateway/%s", r.client.HostURL, state.GatewayID.ValueString())

	reqHTTP, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating update request: %s", err))
		return
	}
	reqHTTP.Header.Set("Content-Type", "application/json")

	resBody, err := r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error updating gateway: %s", err))
		return
	}

	var apiResp GatewayAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	if plan.GroupID.IsNull() {
		// User removed group_id, preserve null in state
		apiResp.GroupID = nil
	}

	state = GatewayResourceModel{
		ObjectVersion:     types.Int64Value(apiResp.ObjectVersion),
		CreationUserID:    types.StringValue(apiResp.CreationUserID),
		UpdateUserID:      types.StringPointerValue(apiResp.UpdateUserID),
		CreationTimestamp: types.StringValue(apiResp.CreationTimestamp),
		UpdateTimestamp:   types.StringPointerValue(apiResp.UpdateTimestamp),
		GatewayID:         types.StringValue(apiResp.GatewayID),
		//GroupID:           types.StringPointerValue(apiResp.GroupID),
		Name:        types.StringValue(apiResp.Name),
		Description: types.StringPointerValue(apiResp.Description),
	}

	// Safely map group_id pointer -> Terraform value
	if apiResp.GroupID == nil || *apiResp.GroupID == "" {
		state.GroupID = types.StringNull()
	} else {
		state.GroupID = types.StringValue(*apiResp.GroupID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         DELETE
//-----------------------------------------------------------------
func (r *GatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state GatewayResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/assets/v1/gateway/%s", r.client.HostURL, state.GatewayID.ValueString())
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

		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error deleting gateway: %s", err))
	}
}
