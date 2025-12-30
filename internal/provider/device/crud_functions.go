package device

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

//-----------------------------------------------------------------
//         CREATE
//-----------------------------------------------------------------

func (r *DeviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DeviceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := map[string]interface{}{}

	// Required fields
	payload["name"] = plan.Name.ValueString()
	payload["vendor_id"] = plan.VendorID.ValueString()
	payload["ide_config_id"] = plan.IdeConfigID.ValueString()

	// Connection configuration (required)
	var connConfig ConnectionConfiguration
	resp.Diagnostics.Append(plan.ConnectionConfig.As(ctx, &connConfig, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return
	}
	payload["connection_configuration"] = connConfig

	// Optional fields
	if !plan.GroupID.IsUnknown() && !plan.GroupID.IsNull() {
		payload["group_id"] = plan.GroupID.ValueString()
	}
	if !plan.MetaData.IsUnknown() && !plan.MetaData.IsNull() {
		var metaData map[string]interface{}
		if err := json.Unmarshal([]byte(plan.MetaData.ValueString()), &metaData); err == nil {
			payload["meta_data"] = metaData
		}
	}
	if !plan.DeviceType.IsUnknown() && !plan.DeviceType.IsNull() {
		payload["device_type"] = plan.DeviceType.ValueString()
	}
	if !plan.Description.IsUnknown() && !plan.Description.IsNull() {
		payload["description"] = plan.Description.ValueString()
	}
	if !plan.SecretID.IsUnknown() && !plan.SecretID.IsNull() {
		payload["secret_id"] = plan.SecretID.ValueString()
	}
	if !plan.FtpConfig.IsUnknown() && !plan.FtpConfig.IsNull() {
		var ftpConfig FtpConfiguration
		resp.Diagnostics.Append(plan.FtpConfig.As(ctx, &ftpConfig, basetypes.ObjectAsOptions{})...)
		if !resp.Diagnostics.HasError() {
			payload["ftp_configuration"] = ftpConfig
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	url := fmt.Sprintf("%s/assets/v1/device", r.client.HostURL)

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

	var apiResp DeviceAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	if apiResp.DeviceID == "" {
		resp.Diagnostics.AddError(
			"API Response Missing DeviceID",
			fmt.Sprintf("The API did not return a device_id in response: %s", string(resBody)),
		)
		return
	}

	state := buildDeviceState(ctx, &apiResp, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         READ
//-----------------------------------------------------------------
func (r *DeviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DeviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/assets/v1/device/%s", r.client.HostURL, state.DeviceID.ValueString())
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
				"Error reading device %s: %s\nResponse body:\n%s",
				state.DeviceID.ValueString(),
				err,
				string(resBody),
			),
		)
		return
	}

	var apiResp DeviceAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	state = buildDeviceState(ctx, &apiResp, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         UPDATE
//-----------------------------------------------------------------
func (r *DeviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state DeviceResourceModel
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
		payload["name"] = plan.Name.ValueString()
	}

	// Include group_id if changed
	if !plan.GroupID.Equal(state.GroupID) {
		if plan.GroupID.IsNull() {
			payload["group_id"] = ""
		} else {
			payload["group_id"] = plan.GroupID.ValueString()
		}
	}

	// Include connection_configuration if changed
	if !plan.ConnectionConfig.Equal(state.ConnectionConfig) {
		var connConfig ConnectionConfiguration
		resp.Diagnostics.Append(plan.ConnectionConfig.As(ctx, &connConfig, basetypes.ObjectAsOptions{})...)
		if !resp.Diagnostics.HasError() {
			partialConn := PartialConnectionConfiguration{
				IPAddress:        &connConfig.IPAddress,
				Port:             &connConfig.Port,
				SubnetMask:       connConfig.SubnetMask,
				GatewayIPAddress: connConfig.GatewayIPAddress,
			}
			payload["connection_configuration"] = partialConn
		}
	}

	// Include meta_data if changed
	if !plan.MetaData.Equal(state.MetaData) {
		if plan.MetaData.IsNull() {
			payload["meta_data"] = nil
		} else {
			var metaData map[string]interface{}
			if err := json.Unmarshal([]byte(plan.MetaData.ValueString()), &metaData); err == nil {
				payload["meta_data"] = metaData
			}
		}
	}

	// Include description if changed
	if !plan.Description.Equal(state.Description) {
		if plan.Description.IsNull() {
			payload["description"] = ""
		} else {
			payload["description"] = plan.Description.ValueString()
		}
	}

	// Include secret_id if changed
	if !plan.SecretID.Equal(state.SecretID) {
		if plan.SecretID.IsNull() {
			payload["secret_id"] = ""
		} else {
			payload["secret_id"] = plan.SecretID.ValueString()
		}
	}

	// Include ftp_configuration if changed
	if !plan.FtpConfig.Equal(state.FtpConfig) {
		if plan.FtpConfig.IsNull() {
			payload["ftp_configuration"] = nil
		} else {
			var ftpConfig FtpConfiguration
			resp.Diagnostics.Append(plan.FtpConfig.As(ctx, &ftpConfig, basetypes.ObjectAsOptions{})...)
			if !resp.Diagnostics.HasError() {
				payload["ftp_configuration"] = ftpConfig
			}
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	url := fmt.Sprintf("%s/assets/v1/device/%s", r.client.HostURL, state.DeviceID.ValueString())

	reqHTTP, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating update request: %s", err))
		return
	}
	reqHTTP.Header.Set("Content-Type", "application/json")

	resBody, err := r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error updating device: %s", err))
		return
	}

	var apiResp DeviceAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	// Preserve null values in state when user removed them
	if plan.GroupID.IsNull() {
		apiResp.GroupID = nil
	}
	if plan.Description.IsNull() {
		apiResp.Description = nil
	}
	if plan.SecretID.IsNull() {
		apiResp.SecretID = nil
	}
	if plan.FtpConfig.IsNull() {
		apiResp.FtpConfig = nil
	}

	state = buildDeviceState(ctx, &apiResp, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         DELETE
//-----------------------------------------------------------------
func (r *DeviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DeviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/assets/v1/device/%s", r.client.HostURL, state.DeviceID.ValueString())
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

		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error deleting device: %s", err))
	}
}

//-----------------------------------------------------------------
//         HELPER FUNCTIONS
//-----------------------------------------------------------------

func buildDeviceState(ctx context.Context, apiResp *DeviceAPIResponse, diags *diag.Diagnostics) DeviceResourceModel {
	state := DeviceResourceModel{
		ObjectVersion:     types.Int64Value(apiResp.ObjectVersion),
		CreationUserID:    types.StringValue(apiResp.CreationUserID),
		UpdateUserID:      types.StringPointerValue(apiResp.UpdateUserID),
		CreationTimestamp: types.StringValue(apiResp.CreationTimestamp),
		UpdateTimestamp:   types.StringPointerValue(apiResp.UpdateTimestamp),
		DeviceID:          types.StringValue(apiResp.DeviceID),
		GroupID:           types.StringPointerValue(apiResp.GroupID),
		Name:              types.StringValue(apiResp.Name),
		VendorID:          types.StringValue(apiResp.VendorID),
		IdeConfigID:       types.StringValue(apiResp.IdeConfigID),
		DeviceType:        types.StringValue(apiResp.DeviceType),
		Description:       types.StringPointerValue(apiResp.Description),
		SecretID:          types.StringPointerValue(apiResp.SecretID),
	}

	// Build connection configuration object
	connConfigAttrs := map[string]attr.Value{
		"ip_address":         types.StringValue(apiResp.ConnectionConfig.IPAddress),
		"port":               types.Int64Value(apiResp.ConnectionConfig.Port),
		"subnet_mask":        types.StringPointerValue(apiResp.ConnectionConfig.SubnetMask),
		"gateway_ip_address": types.StringPointerValue(apiResp.ConnectionConfig.GatewayIPAddress),
	}
	connConfigObj, diag := types.ObjectValue(ConnectionConfigurationObjectType().AttrTypes, connConfigAttrs)
	diags.Append(diag...)
	state.ConnectionConfig = connConfigObj

	// Build metadata JSON string
	if apiResp.MetaData != nil && len(apiResp.MetaData) > 0 {
		metaDataJSON, err := json.Marshal(apiResp.MetaData)
		if err == nil {
			state.MetaData = types.StringValue(string(metaDataJSON))
		} else {
			state.MetaData = types.StringNull()
		}
	} else {
		state.MetaData = types.StringNull()
	}

	// Build FTP configuration object if present
	if apiResp.FtpConfig != nil {
		ftpConfigAttrs := map[string]attr.Value{
			"ip_address":     types.StringValue(apiResp.FtpConfig.IPAddress),
			"port":           types.Int64Value(apiResp.FtpConfig.Port),
			"protocol":       types.StringPointerValue(apiResp.FtpConfig.Protocol),
			"secret_id":      types.StringPointerValue(apiResp.FtpConfig.SecretID),
			"root_directory": types.StringPointerValue(apiResp.FtpConfig.RootDirectory),
		}
		ftpConfigObj, diag := types.ObjectValue(FtpConfigurationObjectType().AttrTypes, ftpConfigAttrs)
		diags.Append(diag...)
		state.FtpConfig = ftpConfigObj
	} else {
		state.FtpConfig = types.ObjectNull(FtpConfigurationObjectType().AttrTypes)
	}

	return state
}
