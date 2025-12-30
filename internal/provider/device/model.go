package device

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceResourceModel struct {
	DeviceID          types.String `tfsdk:"device_id"`
	GroupID           types.String `tfsdk:"group_id"`
	Name              types.String `tfsdk:"name"`
	VendorID          types.String `tfsdk:"vendor_id"`
	IdeConfigID       types.String `tfsdk:"ide_config_id"`
	ConnectionConfig  types.Object `tfsdk:"connection_configuration"`
	MetaData          types.String `tfsdk:"meta_data"`
	DeviceType        types.String `tfsdk:"device_type"`
	Description       types.String `tfsdk:"description"`
	SecretID          types.String `tfsdk:"secret_id"`
	FtpConfig         types.Object `tfsdk:"ftp_configuration"`
	ObjectVersion     types.Int64  `tfsdk:"object_version"`
	CreationUserID    types.String `tfsdk:"creation_user_id"`
	UpdateUserID      types.String `tfsdk:"update_user_id"`
	CreationTimestamp types.String `tfsdk:"creation_timestamp"`
	UpdateTimestamp   types.String `tfsdk:"update_timestamp"`
}

type ConnectionConfiguration struct {
	IPAddress        string  `tfsdk:"ip_address" json:"ip_address"`
	Port             int64   `tfsdk:"port" json:"port"`
	SubnetMask       *string `tfsdk:"subnet_mask" json:"subnet_mask,omitempty"`
	GatewayIPAddress *string `tfsdk:"gateway_ip_address" json:"gateway_ip_address,omitempty"`
}

type PartialConnectionConfiguration struct {
	IPAddress        *string `json:"ip_address,omitempty"`
	Port             *int64  `json:"port,omitempty"`
	SubnetMask       *string `json:"subnet_mask,omitempty"`
	GatewayIPAddress *string `json:"gateway_ip_address,omitempty"`
}

type FtpConfiguration struct {
	IPAddress     string  `tfsdk:"ip_address" json:"ip_address"`
	Port          int64   `tfsdk:"port" json:"port"`
	Protocol      *string `tfsdk:"protocol" json:"protocol,omitempty"`
	SecretID      *string `tfsdk:"secret_id" json:"secret_id,omitempty"`
	RootDirectory *string `tfsdk:"root_directory" json:"root_directory,omitempty"`
}

type DeviceAPIResponse struct {
	ObjectVersion     int64                   `json:"object_version"`
	CreationUserID    string                  `json:"creation_user_id"`
	UpdateUserID      *string                 `json:"update_user_id"`
	CreationTimestamp string                  `json:"creation_timestamp"`
	UpdateTimestamp   *string                 `json:"update_timestamp"`
	DeviceID          string                  `json:"device_id"`
	GroupID           *string                 `json:"group_id"`
	Name              string                  `json:"name"`
	VendorID          string                  `json:"vendor_id"`
	IdeConfigID       string                  `json:"ide_config_id"`
	ConnectionConfig  ConnectionConfiguration `json:"connection_configuration"`
	MetaData          map[string]interface{}  `json:"meta_data,omitempty"`
	DeviceType        string                  `json:"device_type"`
	Description       *string                 `json:"description"`
	SecretID          *string                 `json:"secret_id"`
	FtpConfig         *FtpConfiguration       `json:"ftp_configuration,omitempty"`
}
