package document

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DocumentResourceModel struct {
	DocumentID        types.String `tfsdk:"document_id"`
	GroupID           types.String `tfsdk:"group_id"`
	Name              types.String `tfsdk:"name"`
	DocumentType      types.String `tfsdk:"document_type"`
	LastVersionNumber types.Int64  `tfsdk:"last_version_number"`
	FilePath          types.String `tfsdk:"file_path"`
	FileName          types.String `tfsdk:"file_name"`
	CommitMessage     types.String `tfsdk:"commit_message"`
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

type CreateDocumentAPIResponse struct {
	UploadURLs        []S3MultipartUploadUrl `json:"upload_urls"`
	UploadID          string                 `json:"upload_id"`
	DocumentID        string                 `json:"document_id"`
	VersionID         string                 `json:"version_id"`
	ObjectVersion     int64                  `json:"object_version"`
	CreationUserID    string                 `json:"creation_user_id"`
	UpdateUserID      *string                `json:"update_user_id"`
	CreationTimestamp string                 `json:"creation_timestamp"`
	UpdateTimestamp   *string                `json:"update_timestamp"`
	GroupID           *string                `json:"group_id"`
	Name              string                 `json:"name"`
	LastVersionNumber int64                  `json:"last_version_number"`
	DocumentType      string                 `json:"document_type"`
}

type DocumentAPIResponse struct {
	ObjectVersion     int64   `json:"object_version"`
	CreationUserID    string  `json:"creation_user_id"`
	UpdateUserID      *string `json:"update_user_id"`
	CreationTimestamp string  `json:"creation_timestamp"`
	UpdateTimestamp   *string `json:"update_timestamp"`
	DocumentID        string  `json:"document_id"`
	GroupID           *string `json:"group_id"`
	Name              string  `json:"name"`
	LastVersionNumber int64   `json:"last_version_number"`
	DocumentType      string  `json:"document_type"`
}

type S3MultipartCompleteInfo struct {
	PartNumber int    `json:"part_number"`
	ETag       string `json:"etag"`
}

type CompleteMultipartUploadRequest struct {
	Parts    []S3MultipartCompleteInfo `json:"parts"`
	FileName string                    `json:"file_name"`
}
