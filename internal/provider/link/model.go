package link

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LinkResourceModel struct {
	SourceID          types.String `tfsdk:"source_id"`
	SourceType        types.String `tfsdk:"source_type"`
	DestinationID     types.String `tfsdk:"destination_id"`
	DestinationType   types.String `tfsdk:"destination_type"`
	MetaData          types.String `tfsdk:"meta_data"`
	ObjectVersion     types.Int64  `tfsdk:"object_version"`
	CreationUserID    types.String `tfsdk:"creation_user_id"`
	UpdateUserID      types.String `tfsdk:"update_user_id"`
	CreationTimestamp types.String `tfsdk:"creation_timestamp"`
	UpdateTimestamp   types.String `tfsdk:"update_timestamp"`
}

type AssetLinkMetaData struct {
	// DeviceToGatewayLinkMetaData
	Primary *bool `json:"primary,omitempty"`

	// DeviceToDeviceLinkMetaData
	Interface *string `json:"interface,omitempty"`
	Protocol  *string `json:"protocol,omitempty"`

	// ProjectToProjectLinkMetaData
	SourceVersionID      *string `json:"source_version_id,omitempty"`
	DestinationVersionID *string `json:"destination_version_id,omitempty"`

	// ProjectToDeviceLinkMetaData
	DeviceName                       *string `json:"device_name,omitempty"`
	DeviceType                       *string `json:"device_type,omitempty"`
	DeviceIPAddress                  *string `json:"device_ip_address,omitempty"`
	DeviceSubnet                     *string `json:"device_subnet,omitempty"`
	TargetProjectVersionID           *string `json:"target_project_version_id,omitempty"`
	TargetProjectSyncStatus          *string `json:"target_project_sync_status,omitempty"`
	TargetProjectSyncStatusTimestamp *string `json:"target_project_sync_status_timestamp,omitempty"`
	PreviousProjectSyncStatus        *string `json:"previous_project_sync_status,omitempty"`
	ProjectSyncType                  *string `json:"project_sync_type,omitempty"`
	ProjectSyncErrorMessage          *string `json:"project_sync_error_message,omitempty"`
	ProjectSyncJobID                 *string `json:"project_sync_job_id,omitempty"`

	// DocumentToAssetLinkMetaData / TagToAssetLinkMetaData
	AssetVersionID *string `json:"asset_version_id,omitempty"`
}

type LinkAPIResponse struct {
	ObjectVersion     int64              `json:"object_version"`
	CreationUserID    string             `json:"creation_user_id"`
	UpdateUserID      *string            `json:"update_user_id"`
	CreationTimestamp string             `json:"creation_timestamp"`
	UpdateTimestamp   *string            `json:"update_timestamp"`
	SourceID          string             `json:"source_id"`
	SourceType        string             `json:"source_type"`
	DestinationID     string             `json:"destination_id"`
	DestinationType   string             `json:"destination_type"`
	MetaData          *AssetLinkMetaData `json:"meta_data,omitempty"`
}
