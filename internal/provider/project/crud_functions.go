package project

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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const multipartChunkSize = 5 * 1024 * 1024 // 5MB per part

//-----------------------------------------------------------------
//         CREATE
//-----------------------------------------------------------------

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read file to get size and name
	filePath := plan.FilePath.ValueString()
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		resp.Diagnostics.AddError("File Error", fmt.Sprintf("Error accessing file: %s", err))
		return
	}

	fileName := filepath.Base(filePath)
	fileSize := fileInfo.Size()

	// Calculate number of parts needed
	numParts := int((fileSize + multipartChunkSize - 1) / multipartChunkSize)

	// Pre-calculate MD5 hashes for each part
	var partMD5s []string
	var chunks [][]byte

	// Read and hash all parts
	{
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
	}

	// Create project with upload URLs
	payload := map[string]interface{}{
		"name":          plan.Name.ValueString(),
		"vendor_id":     plan.VendorID.ValueString(),
		"ide_config_id": plan.IdeConfigID.ValueString(),
		"file_name":     fileName,
		"parts":         numParts,
		"file_size":     fileSize,
		"part_md5s":     partMD5s,
	}

	if !plan.GroupID.IsUnknown() && !plan.GroupID.IsNull() {
		payload["group_id"] = plan.GroupID.ValueString()
	}
	if !plan.ProjectType.IsUnknown() && !plan.ProjectType.IsNull() {
		payload["project_type"] = plan.ProjectType.ValueString()
	}
	if !plan.Description.IsUnknown() && !plan.Description.IsNull() {
		payload["description"] = plan.Description.ValueString()
	}
	if !plan.SecretID.IsUnknown() && !plan.SecretID.IsNull() {
		payload["secret_id"] = plan.SecretID.ValueString()
	}
	if !plan.AttachedLicenses.IsUnknown() && !plan.AttachedLicenses.IsNull() {
		var licenses []string
		resp.Diagnostics.Append(plan.AttachedLicenses.ElementsAs(ctx, &licenses, false)...)
		if !resp.Diagnostics.HasError() {
			payload["attached_licenses"] = licenses
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	url := fmt.Sprintf("%s/assets/v1/project", r.client.HostURL)
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

	var createResp CreateProjectAPIResponse
	if err := json.Unmarshal(resBody, &createResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	if createResp.ProjectID == "" {
		resp.Diagnostics.AddError(
			"API Response Missing ProjectID",
			fmt.Sprintf("The API did not return a project_id in response: %s", string(resBody)),
		)
		return
	}

	// Upload file parts
	completeParts := []S3MultipartCompleteInfo{}
	for i, uploadURL := range createResp.UploadURLs {
		chunk := chunks[i]

		// Upload chunk
		uploadReq, err := http.NewRequest(http.MethodPut, uploadURL.UploadURL, bytes.NewReader(chunk))
		if err != nil {
			resp.Diagnostics.AddError("Upload Error", fmt.Sprintf("Error creating upload request: %s", err))
			return
		}

		// MD5 is already included in the pre-signed URL, no need to set it in headers

		uploadResp, err := http.DefaultClient.Do(uploadReq)
		if err != nil {
			resp.Diagnostics.AddError("Upload Error", fmt.Sprintf("Error uploading part %d: %s", uploadURL.PartNumber, err))
			return
		}
		defer uploadResp.Body.Close()

		if uploadResp.StatusCode != http.StatusOK {
			uploadBody, _ := io.ReadAll(uploadResp.Body)
			resp.Diagnostics.AddError(
				"Upload Error",
				fmt.Sprintf("Error uploading part %d: status %d, body: %s", uploadURL.PartNumber, uploadResp.StatusCode, string(uploadBody)),
			)
			return
		}

		// Get ETag from response
		etag := uploadResp.Header.Get("ETag")
		etag = strings.Trim(etag, "\"")

		completeParts = append(completeParts, S3MultipartCompleteInfo{
			PartNumber: uploadURL.PartNumber,
			ETag:       etag,
		})
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

	completeURL := fmt.Sprintf("%s/assets/v1/project/%s/version/%s/complete_upload/%s",
		r.client.HostURL,
		createResp.ProjectID,
		createResp.VersionID,
		createResp.UploadID,
	)

	completeReq, err := http.NewRequest(http.MethodPost, completeURL, bytes.NewReader(completeBody))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating complete request: %s", err))
		return
	}
	completeReq.Header.Set("Content-Type", "application/json")

	_, err = r.client.DoRequest(completeReq, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error completing upload: %s", err))
		return
	}

	// Build attached_licenses list
	var attachedLicensesList basetypes.ListValue
	if createResp.AttachedLicenses != nil && len(createResp.AttachedLicenses) > 0 {
		licensesElements := make([]types.String, len(createResp.AttachedLicenses))
		for i, lic := range createResp.AttachedLicenses {
			licensesElements[i] = types.StringValue(lic)
		}
		var diags = resp.Diagnostics
		attachedLicensesList, _ = types.ListValueFrom(ctx, types.StringType, licensesElements)
		resp.Diagnostics = diags
	} else {
		attachedLicensesList = types.ListNull(types.StringType)
	}

	// Build state
	state := ProjectResourceModel{
		ObjectVersion:     types.Int64Value(createResp.ObjectVersion),
		CreationUserID:    types.StringValue(createResp.CreationUserID),
		UpdateUserID:      types.StringPointerValue(createResp.UpdateUserID),
		CreationTimestamp: types.StringValue(createResp.CreationTimestamp),
		UpdateTimestamp:   types.StringPointerValue(createResp.UpdateTimestamp),
		ProjectID:         types.StringValue(createResp.ProjectID),
		GroupID:           types.StringPointerValue(createResp.GroupID),
		Name:              types.StringValue(createResp.Name),
		VendorID:          types.StringValue(createResp.VendorID),
		IdeConfigID:       types.StringValue(createResp.IdeConfigID),
		ProjectType:       types.StringValue(createResp.ProjectType),
		LastVersionNumber: types.Int64Value(createResp.LastVersionNumber),
		Description:       types.StringPointerValue(createResp.Description),
		SecretID:          types.StringPointerValue(createResp.SecretID),
		AttachedLicenses:  attachedLicensesList,
		FilePath:          plan.FilePath,
		FileName:          types.StringValue(fileName),
		VersionID:         types.StringValue(createResp.VersionID),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         READ
//-----------------------------------------------------------------
func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/assets/v1/project/%s", r.client.HostURL, state.ProjectID.ValueString())
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
				"Error reading project %s: %s\nResponse body:\n%s",
				state.ProjectID.ValueString(),
				err,
				string(resBody),
			),
		)
		return
	}

	var apiResp ProjectAPIResponse
	if err := json.Unmarshal(resBody, &apiResp); err != nil {
		resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
		return
	}

	// Build attached_licenses list
	var attachedLicensesList basetypes.ListValue
	if apiResp.AttachedLicenses != nil && len(apiResp.AttachedLicenses) > 0 {
		licensesElements := make([]types.String, len(apiResp.AttachedLicenses))
		for i, lic := range apiResp.AttachedLicenses {
			licensesElements[i] = types.StringValue(lic)
		}
		var diags = resp.Diagnostics
		attachedLicensesList, _ = types.ListValueFrom(ctx, types.StringType, licensesElements)
		resp.Diagnostics = diags
	} else {
		attachedLicensesList = types.ListNull(types.StringType)
	}

	state.ObjectVersion = types.Int64Value(apiResp.ObjectVersion)
	state.CreationUserID = types.StringValue(apiResp.CreationUserID)
	state.UpdateUserID = types.StringPointerValue(apiResp.UpdateUserID)
	state.CreationTimestamp = types.StringValue(apiResp.CreationTimestamp)
	state.UpdateTimestamp = types.StringPointerValue(apiResp.UpdateTimestamp)
	state.ProjectID = types.StringValue(apiResp.ProjectID)
	state.GroupID = types.StringPointerValue(apiResp.GroupID)
	state.Name = types.StringValue(apiResp.Name)
	state.VendorID = types.StringValue(apiResp.VendorID)
	state.IdeConfigID = types.StringValue(apiResp.IdeConfigID)
	state.ProjectType = types.StringValue(apiResp.ProjectType)
	state.LastVersionNumber = types.Int64Value(apiResp.LastVersionNumber)
	state.Description = types.StringPointerValue(apiResp.Description)
	state.SecretID = types.StringPointerValue(apiResp.SecretID)
	state.AttachedLicenses = attachedLicensesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         UPDATE
//-----------------------------------------------------------------
func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ProjectResourceModel
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

	// Include attached_licenses if changed
	if !plan.AttachedLicenses.Equal(state.AttachedLicenses) {
		if plan.AttachedLicenses.IsNull() {
			payload["attached_licenses"] = []string{}
		} else {
			var licenses []string
			resp.Diagnostics.Append(plan.AttachedLicenses.ElementsAs(ctx, &licenses, false)...)
			if !resp.Diagnostics.HasError() {
				payload["attached_licenses"] = licenses
			}
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
		return
	}

	url := fmt.Sprintf("%s/assets/v1/project/%s", r.client.HostURL, state.ProjectID.ValueString())

	reqHTTP, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating update request: %s", err))
		return
	}
	reqHTTP.Header.Set("Content-Type", "application/json")

	resBody, err := r.client.DoRequest(reqHTTP, nil)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error updating project: %s", err))
		return
	}

	var apiResp ProjectAPIResponse
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

	// Build attached_licenses list
	var attachedLicensesList basetypes.ListValue
	if plan.AttachedLicenses.IsNull() {
		attachedLicensesList = types.ListNull(types.StringType)
	} else if apiResp.AttachedLicenses != nil && len(apiResp.AttachedLicenses) > 0 {
		licensesElements := make([]types.String, len(apiResp.AttachedLicenses))
		for i, lic := range apiResp.AttachedLicenses {
			licensesElements[i] = types.StringValue(lic)
		}
		var diags = resp.Diagnostics
		attachedLicensesList, _ = types.ListValueFrom(ctx, types.StringType, licensesElements)
		resp.Diagnostics = diags
	} else {
		attachedLicensesList = types.ListNull(types.StringType)
	}

	state.ObjectVersion = types.Int64Value(apiResp.ObjectVersion)
	state.CreationUserID = types.StringValue(apiResp.CreationUserID)
	state.UpdateUserID = types.StringPointerValue(apiResp.UpdateUserID)
	state.CreationTimestamp = types.StringValue(apiResp.CreationTimestamp)
	state.UpdateTimestamp = types.StringPointerValue(apiResp.UpdateTimestamp)
	state.ProjectID = types.StringValue(apiResp.ProjectID)
	state.GroupID = types.StringPointerValue(apiResp.GroupID)
	state.Name = types.StringValue(apiResp.Name)
	state.VendorID = types.StringValue(apiResp.VendorID)
	state.IdeConfigID = types.StringValue(apiResp.IdeConfigID)
	state.ProjectType = types.StringValue(apiResp.ProjectType)
	state.LastVersionNumber = types.Int64Value(apiResp.LastVersionNumber)
	state.Description = types.StringPointerValue(apiResp.Description)
	state.SecretID = types.StringPointerValue(apiResp.SecretID)
	state.AttachedLicenses = attachedLicensesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

//-----------------------------------------------------------------
//         DELETE
//-----------------------------------------------------------------
func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/assets/v1/project/%s", r.client.HostURL, state.ProjectID.ValueString())
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

		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error deleting project: %s", err))
	}
}