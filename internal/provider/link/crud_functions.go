package link

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

func (r *LinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan LinkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := map[string]interface{}{}

	// Include metadata if provided
	if !plan.MetaData.IsUnknown() && !plan.MetaData.IsNull() {
		var metaData AssetLinkMetaData
		if err := json.Unmarshal([]byte(plan.MetaData.ValueString()), &metaData); err == nil {
			payload["meta_data"] = metaData
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	url := fmt.Sprintf("%s/assets/v1/link/%s/%s/%s/%s",
		r.client.HostURL,
		plan.SourceType.ValueString(),
		plan.SourceID.ValueString(),
		plan.DestinationType.ValueString(),
		plan.DestinationID.ValueString(),
	)

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

	var apiResp LinkAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	state := LinkResourceModel{
		ObjectVersion:     types.Int64Value(apiResp.ObjectVersion),
		CreationUserID:    types.StringValue(apiResp.CreationUserID),
		UpdateUserID:      types.StringPointerValue(apiResp.UpdateUserID),
		CreationTimestamp: types.StringValue(apiResp.CreationTimestamp),
		UpdateTimestamp:   types.StringPointerValue(apiResp.UpdateTimestamp),
		SourceID:          types.StringValue(apiResp.SourceID),
		SourceType:        types.StringValue(apiResp.SourceType),
		DestinationID:     types.StringValue(apiResp.DestinationID),
		DestinationType:   types.StringValue(apiResp.DestinationType),
	}

	// Build metadata JSON string
	if apiResp.MetaData != nil {
		metaDataJSON, err := json.Marshal(apiResp.MetaData)
		if err == nil {
			state.MetaData = types.StringValue(string(metaDataJSON))
		} else {
			state.MetaData = types.StringNull()
		}
	} else {
		state.MetaData = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         READ
//-----------------------------------------------------------------
func (r *LinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state LinkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/assets/v1/link/%s/%s/%s/%s",
		r.client.HostURL,
		state.SourceType.ValueString(),
		state.SourceID.ValueString(),
		state.DestinationType.ValueString(),
		state.DestinationID.ValueString(),
	)

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
				"Error reading link %s/%s -> %s/%s: %s\nResponse body:\n%s",
				state.SourceType.ValueString(),
				state.SourceID.ValueString(),
				state.DestinationType.ValueString(),
				state.DestinationID.ValueString(),
				err,
				string(resBody),
			),
		)
		return
	}

	var apiResp LinkAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	state = LinkResourceModel{
		ObjectVersion:     types.Int64Value(apiResp.ObjectVersion),
		CreationUserID:    types.StringValue(apiResp.CreationUserID),
		UpdateUserID:      types.StringPointerValue(apiResp.UpdateUserID),
		CreationTimestamp: types.StringValue(apiResp.CreationTimestamp),
		UpdateTimestamp:   types.StringPointerValue(apiResp.UpdateTimestamp),
		SourceID:          types.StringValue(apiResp.SourceID),
		SourceType:        types.StringValue(apiResp.SourceType),
		DestinationID:     types.StringValue(apiResp.DestinationID),
		DestinationType:   types.StringValue(apiResp.DestinationType),
	}

	// Build metadata JSON string
	if apiResp.MetaData != nil {
		metaDataJSON, err := json.Marshal(apiResp.MetaData)
		if err == nil {
			state.MetaData = types.StringValue(string(metaDataJSON))
		} else {
			state.MetaData = types.StringNull()
		}
	} else {
		state.MetaData = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         UPDATE
//-----------------------------------------------------------------
func (r *LinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state LinkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := map[string]interface{}{
		"object_version": state.ObjectVersion.ValueInt64(),
	}

	// Include metadata if changed
	if !plan.MetaData.Equal(state.MetaData) {
		if plan.MetaData.IsNull() {
			payload["meta_data"] = nil
		} else {
			var metaData AssetLinkMetaData
			if err := json.Unmarshal([]byte(plan.MetaData.ValueString()), &metaData); err == nil {
				payload["meta_data"] = metaData
			}
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	url := fmt.Sprintf("%s/assets/v1/link/%s/%s/%s/%s",
		r.client.HostURL,
		state.SourceType.ValueString(),
		state.SourceID.ValueString(),
		state.DestinationType.ValueString(),
		state.DestinationID.ValueString(),
	)

	reqHTTP, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating update request: %s", err))
		return
	}
	reqHTTP.Header.Set("Content-Type", "application/json")

	resBody, err := r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error updating link: %s", err))
		return
	}

	var apiResp LinkAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	// Preserve null in state when user removed metadata
	if plan.MetaData.IsNull() {
		apiResp.MetaData = nil
	}

	state = LinkResourceModel{
		ObjectVersion:     types.Int64Value(apiResp.ObjectVersion),
		CreationUserID:    types.StringValue(apiResp.CreationUserID),
		UpdateUserID:      types.StringPointerValue(apiResp.UpdateUserID),
		CreationTimestamp: types.StringValue(apiResp.CreationTimestamp),
		UpdateTimestamp:   types.StringPointerValue(apiResp.UpdateTimestamp),
		SourceID:          types.StringValue(apiResp.SourceID),
		SourceType:        types.StringValue(apiResp.SourceType),
		DestinationID:     types.StringValue(apiResp.DestinationID),
		DestinationType:   types.StringValue(apiResp.DestinationType),
	}

	// Build metadata JSON string
	if apiResp.MetaData != nil {
		metaDataJSON, err := json.Marshal(apiResp.MetaData)
		if err == nil {
			state.MetaData = types.StringValue(string(metaDataJSON))
		} else {
			state.MetaData = types.StringNull()
		}
	} else {
		state.MetaData = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         DELETE
//-----------------------------------------------------------------
func (r *LinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state LinkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/assets/v1/link/%s/%s/%s/%s",
		r.client.HostURL,
		state.SourceType.ValueString(),
		state.SourceID.ValueString(),
		state.DestinationType.ValueString(),
		state.DestinationID.ValueString(),
	)

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

		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error deleting link: %s", err))
	}
}
