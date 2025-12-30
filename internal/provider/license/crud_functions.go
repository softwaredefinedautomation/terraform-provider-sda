package license

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const multipartChunkSize = 5 * 1024 * 1024 // 5MB per part

//-----------------------------------------------------------------
//         CREATE
//-----------------------------------------------------------------

func (r *LicenseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan LicenseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var fileName string
	var fileSize int64
	var numParts int
	var partMD5s []string
	var chunks [][]byte
	hasFile := !plan.FilePath.IsUnknown() && !plan.FilePath.IsNull()

	// Read file if provided
	if hasFile {
		filePath := plan.FilePath.ValueString()
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			resp.Diagnostics.AddError("File Error", fmt.Sprintf("Error accessing file: %s", err))
			return
		}

		fileName = filepath.Base(filePath)
		fileSize = fileInfo.Size()

		// Calculate number of parts needed
		numParts = int((fileSize + multipartChunkSize - 1) / multipartChunkSize)

		// Pre-calculate MD5 hashes for each part
		file, err := os.Open(filePath)
		if err != nil {
			resp.Diagnostics.AddError("File Error", fmt.Sprintf("Error opening file: %s", err))
			return
		}
		defer file.Close()

		partMD5s = make([]string, numParts)
		chunks = make([][]byte, numParts)
		buffer := make([]byte, multipartChunkSize)

		for i := 0; i < numParts; i++ {
			n, readErr := file.Read(buffer)
			if readErr != nil && readErr != io.EOF {
				resp.Diagnostics.AddError("File Error", fmt.Sprintf("Error reading file chunk: %s", readErr))
				return
			}

			// Make a copy of the chunk
			chunk := make([]byte, n)
			copy(chunk, buffer[:n])
			chunks[i] = chunk

			hash := md5.Sum(chunk)
			partMD5s[i] = base64.StdEncoding.EncodeToString(hash[:])
		}
	} else {
		// Default file name from API spec
		fileName = "SDA_License.vhdx"
		numParts = 1
	}

	// Create license with upload URLs
	payload := map[string]interface{}{
		"vendor_id": plan.VendorID.ValueString(),
		"serial_id": plan.SerialID.ValueString(),
		"product":   plan.Product.ValueString(),
		"parts":     numParts,
	}

	if hasFile {
		payload["file_name"] = fileName
		payload["part_md5s"] = partMD5s
	} else {
		payload["file_name"] = fileName
	}

	if !plan.GroupID.IsUnknown() && !plan.GroupID.IsNull() {
		payload["group_id"] = plan.GroupID.ValueString()
	}
	if !plan.Type.IsUnknown() && !plan.Type.IsNull() {
		payload["type"] = plan.Type.ValueString()
	}
	if !plan.Status.IsUnknown() && !plan.Status.IsNull() {
		payload["status"] = plan.Status.ValueString()
	}
	if !plan.Quantity.IsUnknown() && !plan.Quantity.IsNull() {
		payload["quantity"] = plan.Quantity.ValueInt64()
	}
	if !plan.Name.IsUnknown() && !plan.Name.IsNull() {
		payload["name"] = plan.Name.ValueString()
	}
	if !plan.IdeConfigID.IsUnknown() && !plan.IdeConfigID.IsNull() {
		payload["ide_config_id"] = plan.IdeConfigID.ValueString()
	}
	if !plan.ExpirationTime.IsUnknown() && !plan.ExpirationTime.IsNull() {
		payload["expiration_timestamp"] = plan.ExpirationTime.ValueString()
	}
	if !plan.Family.IsUnknown() && !plan.Family.IsNull() {
		payload["family"] = plan.Family.ValueString()
	}
	if !plan.CompanyName.IsUnknown() && !plan.CompanyName.IsNull() {
		payload["company_name"] = plan.CompanyName.ValueString()
	}
	if !plan.ProductKey.IsUnknown() && !plan.ProductKey.IsNull() {
		payload["product_key"] = plan.ProductKey.ValueString()
	}
	if !plan.ContainerID.IsUnknown() && !plan.ContainerID.IsNull() {
		payload["container_id"] = plan.ContainerID.ValueString()
	}
	if !plan.FirmCode.IsUnknown() && !plan.FirmCode.IsNull() {
		payload["firm_code"] = plan.FirmCode.ValueString()
	}
	if !plan.LicenseServer.IsUnknown() && !plan.LicenseServer.IsNull() {
		payload["license_server"] = plan.LicenseServer.ValueString()
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	url := fmt.Sprintf("%s/assets/v1/license", r.client.HostURL)
	reqHTTP, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating request: %s", err))
		return
	}
	reqHTTP.Header.Set("Content-Type", "application/json")

	resBody, err := r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error creating license: %s\nResponse body:\n%s", err, string(resBody)))
		return
	}

	var apiResp CreateLicenseAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	// Upload parts if file provided
	var completeParts []S3MultipartCompleteInfo
	if hasFile {
		if len(apiResp.UploadURLs) != numParts {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("Mismatch in number of upload URLs: expected %d, got %d", numParts, len(apiResp.UploadURLs)))
			return
		}

		completeParts = make([]S3MultipartCompleteInfo, numParts)
		for i, uploadURL := range apiResp.UploadURLs {
			partReq, err := http.NewRequest(http.MethodPut, uploadURL.UploadURL, bytes.NewReader(chunks[i]))
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating part upload request: %s", err))
				return
			}

			partResp, err := http.DefaultClient.Do(partReq)
			if err != nil {
				resp.Diagnostics.AddError("Upload Error", fmt.Sprintf("Error uploading part %d: %s", i+1, err))
				return
			}
			defer partResp.Body.Close()

			if partResp.StatusCode != http.StatusOK {
				bodyBytes, _ := io.ReadAll(partResp.Body)
				resp.Diagnostics.AddError("Upload Error", fmt.Sprintf("Failed to upload part %d: status %d, body: %s", i+1, partResp.StatusCode, string(bodyBytes)))
				return
			}

			etag := strings.Trim(partResp.Header.Get("ETag"), `"`)
			completeParts[i] = S3MultipartCompleteInfo{
				PartNumber: uploadURL.PartNumber,
				ETag:       etag,
			}
		}

		// Complete multipart upload
		completePayload := CompleteMultipartUploadRequest{
			Parts:    completeParts,
			FileName: fileName,
		}
		completeBody, err := json.Marshal(completePayload)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling complete payload: %s", err))
			return
		}

		completeURL := fmt.Sprintf("%s/assets/v1/license/%s/file?upload_id=%s", r.client.HostURL, apiResp.LicenseID, apiResp.UploadID)
		completeReq, err := http.NewRequest(http.MethodPost, completeURL, bytes.NewReader(completeBody))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating complete request: %s", err))
			return
		}
		completeReq.Header.Set("Content-Type", "application/json")

		_, err = r.client.DoRequest(completeReq, nil)
		if err != nil {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error completing multipart upload: %s", err))
			return
		}
	}

	// Build state from API response
	state := LicenseResourceModel{
		ObjectVersion:     types.Int64Value(apiResp.ObjectVersion),
		CreationUserID:    types.StringValue(apiResp.CreationUserID),
		UpdateUserID:      types.StringPointerValue(apiResp.UpdateUserID),
		CreationTimestamp: types.StringValue(apiResp.CreationTimestamp),
		UpdateTimestamp:   types.StringPointerValue(apiResp.UpdateTimestamp),
		LicenseID:         types.StringValue(apiResp.LicenseID),
		GroupID:           types.StringPointerValue(apiResp.GroupID),
		VendorID:          types.StringValue(apiResp.VendorID),
		SerialID:          types.StringValue(apiResp.SerialID),
		Product:           types.StringValue(apiResp.Product),
		Type:              types.StringValue(apiResp.Type),
		Status:            types.StringValue(apiResp.Status),
		Quantity:          types.Int64Value(apiResp.Quantity),
		Name:              types.StringPointerValue(apiResp.Name),
		IdeConfigID:       types.StringPointerValue(apiResp.IdeConfigID),
		ExpirationTime:    types.StringPointerValue(apiResp.ExpirationTime),
		Family:            types.StringPointerValue(apiResp.Family),
		CompanyName:       types.StringPointerValue(apiResp.CompanyName),
		ProductKey:        types.StringPointerValue(apiResp.ProductKey),
		ContainerID:       types.StringPointerValue(apiResp.ContainerID),
		FirmCode:          types.StringPointerValue(apiResp.FirmCode),
		LicenseServer:     types.StringPointerValue(apiResp.LicenseServer),
		FileName:          types.StringValue(fileName),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         READ
//-----------------------------------------------------------------
func (r *LicenseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state LicenseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/assets/v1/license/%s", r.client.HostURL, state.LicenseID.ValueString())
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
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error reading license %s: %s\nResponse body:\n%s", state.LicenseID.ValueString(), err, string(resBody)))
		return
	}

	var apiResp LicenseAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	state.ObjectVersion = types.Int64Value(apiResp.ObjectVersion)
	state.CreationUserID = types.StringValue(apiResp.CreationUserID)
	state.UpdateUserID = types.StringPointerValue(apiResp.UpdateUserID)
	state.CreationTimestamp = types.StringValue(apiResp.CreationTimestamp)
	state.UpdateTimestamp = types.StringPointerValue(apiResp.UpdateTimestamp)
	state.LicenseID = types.StringValue(apiResp.LicenseID)
	state.GroupID = types.StringPointerValue(apiResp.GroupID)
	state.VendorID = types.StringValue(apiResp.VendorID)
	state.SerialID = types.StringValue(apiResp.SerialID)
	state.Product = types.StringValue(apiResp.Product)
	state.Type = types.StringValue(apiResp.Type)
	state.Status = types.StringValue(apiResp.Status)
	state.Quantity = types.Int64Value(apiResp.Quantity)
	state.Name = types.StringPointerValue(apiResp.Name)
	state.IdeConfigID = types.StringPointerValue(apiResp.IdeConfigID)
	state.ExpirationTime = types.StringPointerValue(apiResp.ExpirationTime)
	state.Family = types.StringPointerValue(apiResp.Family)
	state.CompanyName = types.StringPointerValue(apiResp.CompanyName)
	state.ProductKey = types.StringPointerValue(apiResp.ProductKey)
	state.ContainerID = types.StringPointerValue(apiResp.ContainerID)
	state.FirmCode = types.StringPointerValue(apiResp.FirmCode)
	state.LicenseServer = types.StringPointerValue(apiResp.LicenseServer)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         UPDATE
//-----------------------------------------------------------------
func (r *LicenseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state LicenseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := map[string]interface{}{
		"object_version": state.ObjectVersion.ValueInt64(),
	}

	// Group ID
	if !plan.GroupID.Equal(state.GroupID) {
		if plan.GroupID.IsNull() {
			payload["group_id"] = nil
		} else {
			payload["group_id"] = plan.GroupID.ValueString()
		}
	}

	// Name
	if !plan.Name.Equal(state.Name) {
		if plan.Name.IsNull() {
			payload["name"] = nil
		} else {
			payload["name"] = plan.Name.ValueString()
		}
	}

	// IDE Config ID
	if !plan.IdeConfigID.Equal(state.IdeConfigID) {
		if plan.IdeConfigID.IsNull() {
			payload["ide_config_id"] = nil
		} else {
			payload["ide_config_id"] = plan.IdeConfigID.ValueString()
		}
	}

	// Expiration Timestamp
	if !plan.ExpirationTime.Equal(state.ExpirationTime) {
		if plan.ExpirationTime.IsNull() {
			payload["expiration_timestamp"] = nil
		} else {
			payload["expiration_timestamp"] = plan.ExpirationTime.ValueString()
		}
	}

	// Family
	if !plan.Family.Equal(state.Family) {
		if plan.Family.IsNull() {
			payload["family"] = nil
		} else {
			payload["family"] = plan.Family.ValueString()
		}
	}

	// Company Name
	if !plan.CompanyName.Equal(state.CompanyName) {
		if plan.CompanyName.IsNull() {
			payload["company_name"] = nil
		} else {
			payload["company_name"] = plan.CompanyName.ValueString()
		}
	}

	// Product Key
	if !plan.ProductKey.Equal(state.ProductKey) {
		if plan.ProductKey.IsNull() {
			payload["product_key"] = nil
		} else {
			payload["product_key"] = plan.ProductKey.ValueString()
		}
	}

	// Container ID
	if !plan.ContainerID.Equal(state.ContainerID) {
		if plan.ContainerID.IsNull() {
			payload["container_id"] = nil
		} else {
			payload["container_id"] = plan.ContainerID.ValueString()
		}
	}

	// Firm Code
	if !plan.FirmCode.Equal(state.FirmCode) {
		if plan.FirmCode.IsNull() {
			payload["firm_code"] = nil
		} else {
			payload["firm_code"] = plan.FirmCode.ValueString()
		}
	}

	// License Server
	if !plan.LicenseServer.Equal(state.LicenseServer) {
		if plan.LicenseServer.IsNull() {
			payload["license_server"] = nil
		} else {
			payload["license_server"] = plan.LicenseServer.ValueString()
		}
	}

	// Status
	if !plan.Status.Equal(state.Status) {
		payload["status"] = plan.Status.ValueString()
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	url := fmt.Sprintf("%s/assets/v1/license/%s", r.client.HostURL, state.LicenseID.ValueString())
	reqHTTP, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating update request: %s", err))
		return
	}
	reqHTTP.Header.Set("Content-Type", "application/json")

	resBody, err := r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error updating license: %s\nResponse body:\n%s", err, string(resBody)))
		return
	}

	var apiResp LicenseAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	state.ObjectVersion = types.Int64Value(apiResp.ObjectVersion)
	state.CreationUserID = types.StringValue(apiResp.CreationUserID)
	state.UpdateUserID = types.StringPointerValue(apiResp.UpdateUserID)
	state.CreationTimestamp = types.StringValue(apiResp.CreationTimestamp)
	state.UpdateTimestamp = types.StringPointerValue(apiResp.UpdateTimestamp)
	state.LicenseID = types.StringValue(apiResp.LicenseID)
	state.GroupID = types.StringPointerValue(apiResp.GroupID)
	state.VendorID = types.StringValue(apiResp.VendorID)
	state.SerialID = types.StringValue(apiResp.SerialID)
	state.Product = types.StringValue(apiResp.Product)
	state.Type = types.StringValue(apiResp.Type)
	state.Status = types.StringValue(apiResp.Status)
	state.Quantity = types.Int64Value(apiResp.Quantity)
	state.Name = types.StringPointerValue(apiResp.Name)
	state.IdeConfigID = types.StringPointerValue(apiResp.IdeConfigID)
	state.ExpirationTime = types.StringPointerValue(apiResp.ExpirationTime)
	state.Family = types.StringPointerValue(apiResp.Family)
	state.CompanyName = types.StringPointerValue(apiResp.CompanyName)
	state.ProductKey = types.StringPointerValue(apiResp.ProductKey)
	state.ContainerID = types.StringPointerValue(apiResp.ContainerID)
	state.FirmCode = types.StringPointerValue(apiResp.FirmCode)
	state.LicenseServer = types.StringPointerValue(apiResp.LicenseServer)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         DELETE
//-----------------------------------------------------------------
func (r *LicenseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state LicenseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/assets/v1/license/%s", r.client.HostURL, state.LicenseID.ValueString())
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

		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error deleting license: %s", err))
	}
}