package resourcegroup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings" // Added for error checking in Read and Create

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//-----------------------------------------------------------------
//         CREATE
//-----------------------------------------------------------------

func (r *ResourceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ResourceGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := map[string]interface{}{
		"name": plan.Name.ValueString(),
	}
	if !plan.GroupType.IsNull() {
		payload["group_type"] = plan.GroupType.ValueString()
	}
	if !plan.ParentGroupID.IsNull() {
		payload["parent_group_id"] = plan.ParentGroupID.ValueString()
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	url := fmt.Sprintf("%s/assets/v1/resource_group", r.client.HostURL)

	reqHTTP, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating request: %s", err))
		return
	}
	reqHTTP.Header.Set("Content-Type", "application/json")

	resBody, err := r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		// Provide more specific feedback for 403 Forbidden errors.
		errSummary := "API Error"
		if strings.Contains(err.Error(), "status: 403") {
			// This diagnostic helps the user understand the 403 issue they encountered.
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

	var apiResp ResourceGroupAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	if apiResp.GroupID == "" {
		resp.Diagnostics.AddError(
			"API Response Missing GroupID",
			fmt.Sprintf("The API did not return a group_id in response: %s", string(resBody)),
		)
		return
	}

	state := ResourceGroupResourceModel{
		ObjectVersion:     types.Int64Value(apiResp.ObjectVersion),
		CreationUserID:    types.StringValue(apiResp.CreationUserID),
		UpdateUserID:      types.StringPointerValue(apiResp.UpdateUserID),
		CreationTimestamp: types.StringValue(apiResp.CreationTimestamp),
		UpdateTimestamp:   types.StringPointerValue(apiResp.UpdateTimestamp),
		GroupID:           types.StringValue(apiResp.GroupID),
		Name:              types.StringValue(apiResp.Name),
		GroupType:         types.StringValue(apiResp.GroupType),
		ParentGroupID:     types.StringPointerValue(apiResp.ParentGroupID),
		IsSystemGroup:     types.BoolValue(apiResp.IsSystemGroup),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         READ
//-----------------------------------------------------------------
func (r *ResourceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ResourceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/assets/v1/resource_group/%s", r.client.HostURL, state.GroupID.ValueString())
	reqHTTP, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating read request: %s", err))
		return
	}

	resBody, err := r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		// Crucial Terraform pattern: If resource returns 404 Not Found, remove it from state.
		// We rely on the client's DoRequest error string including "status: 404".
		if strings.Contains(err.Error(), "status: 404") {
			resp.State.RemoveResource(ctx)
			return
		}

		// For any other error (e.g., 500, 403, network failure), return the error.
		resp.Diagnostics.AddError(
			"API Error",
			fmt.Sprintf(
				"Error reading resource group %s: %s\nResponse body:\n%s",
				state.GroupID.ValueString(),
				err,
				string(resBody),
			),
		)
		return
	}

	var apiResp ResourceGroupAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	//apiResp.GroupID = state.GroupID

	state = ResourceGroupResourceModel{
		ObjectVersion:     types.Int64Value(apiResp.ObjectVersion),
		CreationUserID:    types.StringValue(apiResp.CreationUserID),
		UpdateUserID:      types.StringPointerValue(apiResp.UpdateUserID),
		CreationTimestamp: types.StringValue(apiResp.CreationTimestamp),
		UpdateTimestamp:   types.StringPointerValue(apiResp.UpdateTimestamp),
		GroupID:           types.StringValue(apiResp.GroupID),
		Name:              types.StringValue(apiResp.Name),
		GroupType:         types.StringValue(apiResp.GroupType),
		ParentGroupID:     types.StringPointerValue(apiResp.ParentGroupID),
		IsSystemGroup:     types.BoolValue(apiResp.IsSystemGroup),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         UPDATE
//-----------------------------------------------------------------
func (r *ResourceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ResourceGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := map[string]interface{}{
		"name":           plan.Name.ValueString(),
		"object_version": state.ObjectVersion.ValueInt64(),
	}
	if !plan.GroupType.IsNull() {
		payload["group_type"] = plan.GroupType.ValueString()
	}
	if !plan.ParentGroupID.IsNull() {
		payload["parent_group_id"] = plan.ParentGroupID.ValueString()
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	url := fmt.Sprintf("%s/assets/v1/resource_group/%s", r.client.HostURL, state.GroupID.ValueString())

	reqHTTP, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating update request: %s", err))
		return
	}
	reqHTTP.Header.Set("Content-Type", "application/json")

	resBody, err := r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error updating resource group: %s", err))
		return
	}

	var apiResp ResourceGroupAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	state = ResourceGroupResourceModel{
		ObjectVersion:     types.Int64Value(apiResp.ObjectVersion),
		CreationUserID:    types.StringValue(apiResp.CreationUserID),
		UpdateUserID:      types.StringPointerValue(apiResp.UpdateUserID),
		CreationTimestamp: types.StringValue(apiResp.CreationTimestamp),
		UpdateTimestamp:   types.StringPointerValue(apiResp.UpdateTimestamp),
		GroupID:           types.StringValue(apiResp.GroupID),
		Name:              types.StringValue(apiResp.Name),
		GroupType:         types.StringValue(apiResp.GroupType),
		ParentGroupID:     types.StringPointerValue(apiResp.ParentGroupID),
		IsSystemGroup:     types.BoolValue(apiResp.IsSystemGroup),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         DELETE
//-----------------------------------------------------------------
func (r *ResourceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ResourceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/assets/v1/resource_group/%s", r.client.HostURL, state.GroupID.ValueString())
	reqHTTP, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating delete request: %s", err))
		return
	}

	_, err = r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		// A common pattern is to ignore 404 during deletion, as the resource is already gone.
		if strings.Contains(err.Error(), "status: 404") {
			return
		}

		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error deleting resource group: %s", err))
	}
}
