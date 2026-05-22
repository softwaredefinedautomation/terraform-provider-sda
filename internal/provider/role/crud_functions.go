package role

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

// CREATE
func (r *RoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RoleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := map[string]interface{}{
		"name": plan.Name.ValueString(),
	}

	if !plan.GroupID.IsUnknown() && !plan.GroupID.IsNull() {
		payload["group_id"] = plan.GroupID.ValueString()
	}
	if !plan.Description.IsUnknown() && !plan.Description.IsNull() {
		payload["description"] = plan.Description.ValueString()
	}

	// policies (list of strings representing policy ids or JSON)
	if !plan.Policies.IsUnknown() && !plan.Policies.IsNull() {
		var policyStrings []string
		if diags := plan.Policies.ElementsAs(ctx, &policyStrings, false); diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		// Unmarshal each JSON string back into a map or object
		var policyObjects []map[string]interface{}
		for _, policyStr := range policyStrings {
			var policyMap map[string]interface{}
			if err := json.Unmarshal([]byte(policyStr), &policyMap); err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error unmarshalling policy JSON string: %s. Content: %s", err, policyStr))
				return
			}
			policyObjects = append(policyObjects, policyMap)
		}
		payload["policies"] = policyObjects
	}

	if !plan.SsoGroupMapping.IsUnknown() && !plan.SsoGroupMapping.IsNull() {
		var sso []string
		if diags := plan.SsoGroupMapping.ElementsAs(ctx, &sso, false); diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		} else {
			payload["sso_group_mapping"] = sso
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	url := fmt.Sprintf("%s/ident/v1/user_role", r.client.HostURL)
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
		resp.Diagnostics.AddError("API Error", "API returned an empty response body after creating the role")
		return
	}

	// We must unmarshal the response body
	var apiResp RoleAPIResponse
	if err = json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding create API response: %s. Body: %s", err, string(resBody)))
		return
	}

	state := RoleResourceModel{
		ObjectVersion:     types.Int64Value(apiResp.ObjectVersion),
		CreationUserID:    types.StringValue(apiResp.CreationUserID),
		UpdateUserID:      types.StringPointerValue(apiResp.UpdateUserID),
		CreationTimestamp: types.StringValue(apiResp.CreationTimestamp),
		UpdateTimestamp:   types.StringPointerValue(apiResp.UpdateTimestamp),
		UserRoleID:        types.StringValue(apiResp.UserRoleID),
		Name:              types.StringValue(apiResp.Name),
		GroupID:           types.StringPointerValue(apiResp.GroupID),
		Description:       types.StringPointerValue(apiResp.Description),
		IsSystemRole:      types.BoolValue(apiResp.IsSystemRole),
	}

	// Preserve provided lists in state
	if !plan.Policies.IsUnknown() {
		state.Policies = plan.Policies
	}
	if !plan.SsoGroupMapping.IsUnknown() {
		state.SsoGroupMapping = plan.SsoGroupMapping
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// READ
func (r *RoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RoleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/ident/v1/user_role/%s", r.client.HostURL, state.UserRoleID.ValueString())
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
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error reading role %s: %s\nResponse body:\n%s", state.UserRoleID.ValueString(), err, string(resBody)))
		return
	}

	if len(resBody) == 0 {
		resp.Diagnostics.AddError("API Error", "API returned an empty array for the role")
		return
	}

	// We must unmarshal the response body
	var apiResp RoleAPIResponse
	if err = json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding read API response: %s. Body: %s", err, string(resBody)))
		return
	}

	state.ObjectVersion = types.Int64Value(apiResp.ObjectVersion)
	state.CreationUserID = types.StringValue(apiResp.CreationUserID)
	state.UpdateUserID = types.StringPointerValue(apiResp.UpdateUserID)
	state.CreationTimestamp = types.StringValue(apiResp.CreationTimestamp)
	state.UpdateTimestamp = types.StringPointerValue(apiResp.UpdateTimestamp)
	state.UserRoleID = types.StringValue(apiResp.UserRoleID)
	state.Name = types.StringValue(apiResp.Name)
	state.GroupID = types.StringPointerValue(apiResp.GroupID)
	state.Description = types.StringPointerValue(apiResp.Description)
	state.IsSystemRole = types.BoolValue(apiResp.IsSystemRole)

	// Read policies from API response
	if apiResp.Policies != nil {
		var policyStrings []string
		for _, p := range apiResp.Policies {
			policyJSON, err := json.Marshal(p)
			if err != nil {
				resp.Diagnostics.AddError("Decode Error",
					fmt.Sprintf("Error marshalling policy: %s", err))
				return
			}

			// Strip server-generated fields that aren't in user's config
			var policyMap map[string]interface{}
			if err := json.Unmarshal(policyJSON, &policyMap); err != nil {
				resp.Diagnostics.AddError("Decode Error",
					fmt.Sprintf("Error unmarshalling policy: %s", err))
				return
			}
			delete(policyMap, "policy_id")

			cleanJSON, err := json.Marshal(policyMap)
			if err != nil {
				resp.Diagnostics.AddError("Decode Error",
					fmt.Sprintf("Error marshalling cleaned policy: %s", err))
				return
			}
										
			policyStrings = append(policyStrings, string(cleanJSON))
		}
		policyList, diags := types.ListValueFrom(ctx, types.StringType, policyStrings)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		state.Policies = policyList
	} else {
		state.Policies = types.ListNull(types.StringType)
	}

	// Read SSO group mapping from API response
	if apiResp.SsoGroupMapping != nil {
		ssoList, diags := types.ListValueFrom(ctx, types.StringType, apiResp.SsoGroupMapping)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		state.SsoGroupMapping = ssoList
	} else {
		state.SsoGroupMapping = types.ListNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// UPDATE
func (r *RoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state RoleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := map[string]interface{}{
		"object_version": state.ObjectVersion.ValueInt64(),
	}

	if !plan.Name.Equal(state.Name) {
		payload["name"] = plan.Name.ValueString()
	}
	if !plan.GroupID.Equal(state.GroupID) {
		if plan.GroupID.IsNull() {
			payload["group_id"] = ""
		} else {
			payload["group_id"] = plan.GroupID.ValueString()
		}
	}
	if !plan.Description.Equal(state.Description) {
		if plan.Description.IsNull() {
			payload["description"] = ""
		} else {
			payload["description"] = plan.Description.ValueString()
		}
	}

	if !plan.Policies.Equal(state.Policies) {
		if plan.Policies.IsNull() {
			payload["policies"] = []map[string]interface{}{} // Send empty object list
		} else {
			var policyStrings []string
			if diags := plan.Policies.ElementsAs(ctx, &policyStrings, false); diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
			// Unmarshal each JSON string back into a map or object
			var policyObjects []map[string]interface{}
			for _, policyStr := range policyStrings {
				var policyMap map[string]interface{}
				if err := json.Unmarshal([]byte(policyStr), &policyMap); err != nil {
					resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error unmarshalling policy JSON string: %s. Content: %s", err, policyStr))
					return
				}
				policyObjects = append(policyObjects, policyMap)
			}
			payload["policies"] = policyObjects
		}
	}

	if !plan.SsoGroupMapping.Equal(state.SsoGroupMapping) {
		if plan.SsoGroupMapping.IsNull() {
			payload["sso_group_mapping"] = []string{}
		} else {
			var sso []string
			if diags := plan.SsoGroupMapping.ElementsAs(ctx, &sso, false); diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
			payload["sso_group_mapping"] = sso
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	url := fmt.Sprintf("%s/ident/v1/user_role/%s", r.client.HostURL, state.UserRoleID.ValueString())
	reqHTTP, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating update request: %s", err))
		return
	}
	reqHTTP.Header.Set("Content-Type", "application/json")

	resBody, err := r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error updating role: %s\nResponse body:\n%s", err, string(resBody)))
		return
	}

	if len(resBody) == 0 {
		resp.Diagnostics.AddError("API Error", "API returned an empty array for the role")
		return
	}

	// We must unmarshal the response body
	var apiResp RoleAPIResponse
	if err = json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding update API response: %s. Body: %s", err, string(resBody)))
		return
	}

	state.ObjectVersion = types.Int64Value(apiResp.ObjectVersion)
	state.CreationUserID = types.StringValue(apiResp.CreationUserID)
	state.UpdateUserID = types.StringPointerValue(apiResp.UpdateUserID)
	state.CreationTimestamp = types.StringValue(apiResp.CreationTimestamp)
	state.UpdateTimestamp = types.StringPointerValue(apiResp.UpdateTimestamp)
	state.UserRoleID = types.StringValue(apiResp.UserRoleID)
	state.Name = types.StringValue(apiResp.Name)
	state.GroupID = types.StringPointerValue(apiResp.GroupID)
	state.Description = types.StringPointerValue(apiResp.Description)
	state.IsSystemRole = types.BoolValue(apiResp.IsSystemRole)

	// Preserve plan policies and sso_group_mapping in state to avoid inconsistent results
	if !plan.Policies.IsNull() && !plan.Policies.IsUnknown() {
		state.Policies = plan.Policies
	}

	if !plan.SsoGroupMapping.IsNull() && !plan.SsoGroupMapping.IsUnknown() {
		state.SsoGroupMapping = plan.SsoGroupMapping
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// DELETE
func (r *RoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RoleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/ident/v1/user_role/%s", r.client.HostURL, state.UserRoleID.ValueString())
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
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error deleting role: %s", err))
	}
}
