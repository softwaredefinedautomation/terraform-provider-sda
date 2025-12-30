package project

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProjectResourceModel struct {
	ProjectID         types.String `tfsdk:"project_id"`
	GroupID           types.String `tfsdk:"group_id"`
	Name              types.String `tfsdk:"name"`
	VendorID          types.String `tfsdk:"vendor_id"`
	IdeConfigID       types.String `tfsdk:"ide_config_id"`
	ProjectType       types.String `tfsdk:"project_type"`
	LastVersionNumber types.Int64  `tfsdk:"last_version_number"`
	Description       types.String `tfsdk:"description"`
	SecretID          types.String `tfsdk:"secret_id"`
	AttachedLicenses  types.List   `tfsdk:"attached_licenses"`
	FilePath          types.String `tfsdk:"file_path"`
	FileName          types.String `tfsdk:"file_name"`
	VersionID         types.String `tfsdk:"version_id"`
	ObjectVersion     types.Int64  `tfsdk:"object_version"`
	CreationUserID    types.String `tfsdk:"creation_user_id"`
	UpdateUserID      types.String `tfsdk:"update_user_id"`
	CreationTimestamp types.String `tfsdk:"creation_timestamp"`
	UpdateTimestamp   types.String `tfsdk:"update_timestamp"`
}

type S3MultipartUploadUrl struct {
	PartNumber int    `json:"part_number"`
	UploadURL  string `json:"upload_url"`
}

type CreateProjectAPIResponse struct {
	UploadURLs        []S3MultipartUploadUrl `json:"upload_urls"`
	UploadID          string                 `json:"upload_id"`
	ProjectID         string                 `json:"project_id"`
	VersionID         string                 `json:"version_id"`
	ObjectVersion     int64                  `json:"object_version"`
	CreationUserID    string                 `json:"creation_user_id"`
	UpdateUserID      *string                `json:"update_user_id"`
	CreationTimestamp string                 `json:"creation_timestamp"`
	UpdateTimestamp   *string                `json:"update_timestamp"`
	GroupID           *string                `json:"group_id"`
	Name              string                 `json:"name"`
	VendorID          string                 `json:"vendor_id"`
	IdeConfigID       string                 `json:"ide_config_id"`
	ProjectType       string                 `json:"project_type"`
	LastVersionNumber int64                  `json:"last_version_number"`
	Description       *string                `json:"description"`
	SecretID          *string                `json:"secret_id"`
	AttachedLicenses  []string               `json:"attached_licenses,omitempty"`
}

type ProjectAPIResponse struct {
	ObjectVersion     int64    `json:"object_version"`
	CreationUserID    string   `json:"creation_user_id"`
	UpdateUserID      *string  `json:"update_user_id"`
	CreationTimestamp string   `json:"creation_timestamp"`
	UpdateTimestamp   *string  `json:"update_timestamp"`
	ProjectID         string   `json:"project_id"`
	GroupID           *string  `json:"group_id"`
	Name              string   `json:"name"`
	VendorID          string   `json:"vendor_id"`
	IdeConfigID       string   `json:"ide_config_id"`
	ProjectType       string   `json:"project_type"`
	LastVersionNumber int64    `json:"last_version_number"`
	Description       *string  `json:"description"`
	SecretID          *string  `json:"secret_id"`
	AttachedLicenses  []string `json:"attached_licenses,omitempty"`
}

type S3MultipartCompleteInfo struct {
	PartNumber int    `json:"part_number"`
	ETag       string `json:"etag"`
}

type CompleteMultipartUploadRequest struct {
	Parts    []S3MultipartCompleteInfo `json:"parts"`
	FileName string                    `json:"file_name"`
}