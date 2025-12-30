package license

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LicenseResourceModel struct {
	LicenseID         types.String `tfsdk:"license_id"`
	GroupID           types.String `tfsdk:"group_id"`
	VendorID          types.String `tfsdk:"vendor_id"`
	SerialID          types.String `tfsdk:"serial_id"`
	Product           types.String `tfsdk:"product"`
	Type              types.String `tfsdk:"type"`
	Status            types.String `tfsdk:"status"`
	Quantity          types.Int64  `tfsdk:"quantity"`
	Name              types.String `tfsdk:"name"`
	IdeConfigID       types.String `tfsdk:"ide_config_id"`
	ExpirationTime    types.String `tfsdk:"expiration_timestamp"`
	Family            types.String `tfsdk:"family"`
	CompanyName       types.String `tfsdk:"company_name"`
	ProductKey        types.String `tfsdk:"product_key"`
	ContainerID       types.String `tfsdk:"container_id"`
	FirmCode          types.String `tfsdk:"firm_code"`
	LicenseServer     types.String `tfsdk:"license_server"`
	FilePath          types.String `tfsdk:"file_path"`
	FileName          types.String `tfsdk:"file_name"`
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

type CreateLicenseAPIResponse struct {
	UploadURLs        []S3MultipartUploadUrl `json:"upload_urls"`
	UploadID          string                 `json:"upload_id"`
	LicenseID         string                 `json:"license_id"`
	ObjectVersion     int64                  `json:"object_version"`
	CreationUserID    string                 `json:"creation_user_id"`
	UpdateUserID      *string                `json:"update_user_id"`
	CreationTimestamp string                 `json:"creation_timestamp"`
	UpdateTimestamp   *string                `json:"update_timestamp"`
	GroupID           *string                `json:"group_id"`
	VendorID          string                 `json:"vendor_id"`
	SerialID          string                 `json:"serial_id"`
	Product           string                 `json:"product"`
	Type              string                 `json:"type"`
	Status            string                 `json:"status"`
	Quantity          int64                  `json:"quantity"`
	Name              *string                `json:"name"`
	IdeConfigID       *string                `json:"ide_config_id"`
	ExpirationTime    *string                `json:"expiration_timestamp"`
	Family            *string                `json:"family"`
	CompanyName       *string                `json:"company_name"`
	ProductKey        *string                `json:"product_key"`
	ContainerID       *string                `json:"container_id"`
	FirmCode          *string                `json:"firm_code"`
	LicenseServer     *string                `json:"license_server"`
}

type LicenseAPIResponse struct {
	ObjectVersion     int64   `json:"object_version"`
	CreationUserID    string  `json:"creation_user_id"`
	UpdateUserID      *string `json:"update_user_id"`
	CreationTimestamp string  `json:"creation_timestamp"`
	UpdateTimestamp   *string `json:"update_timestamp"`
	LicenseID         string  `json:"license_id"`
	GroupID           *string `json:"group_id"`
	VendorID          string  `json:"vendor_id"`
	SerialID          string  `json:"serial_id"`
	Product           string  `json:"product"`
	Type              string  `json:"type"`
	Status            string  `json:"status"`
	Quantity          int64   `json:"quantity"`
	Name              *string `json:"name"`
	IdeConfigID       *string `json:"ide_config_id"`
	ExpirationTime    *string `json:"expiration_timestamp"`
	Family            *string `json:"family"`
	CompanyName       *string `json:"company_name"`
	ProductKey        *string `json:"product_key"`
	ContainerID       *string `json:"container_id"`
	FirmCode          *string `json:"firm_code"`
	LicenseServer     *string `json:"license_server"`
}

type S3MultipartCompleteInfo struct {
	PartNumber int    `json:"part_number"`
	ETag       string `json:"etag"`
}

type CompleteMultipartUploadRequest struct {
	Parts    []S3MultipartCompleteInfo `json:"parts"`
	FileName string                    `json:"file_name"`
}