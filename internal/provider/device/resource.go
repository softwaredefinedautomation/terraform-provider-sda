package device

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/sda/terraform-provider-sda/internal/clients"
)

var _ resource.Resource = &DeviceResource{}
var _ resource.ResourceWithImportState = &DeviceResource{}

func NewDeviceResource() resource.Resource {
	return &DeviceResource{}
}

type DeviceResource struct {
	client *clients.Client
}

func (r *DeviceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

func (r *DeviceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a device resource in the SDA Assets Management Service.",
		Attributes: map[string]schema.Attribute{
			"device_id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the device.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Description: "Resource group ID to which this device belongs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the device.",
			},
			"vendor_id": schema.StringAttribute{
				Required:    true,
				Description: "Vendor ID of the device.",
			},
			"ide_config_id": schema.StringAttribute{
				Required:    true,
				Description: "IDE configuration ID for the device.",
			},
			"connection_configuration": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Connection configuration for the device.",
				Attributes: map[string]schema.Attribute{
					"ip_address": schema.StringAttribute{
						Required:    true,
						Description: "IP address of the device.",
					},
					"port": schema.Int64Attribute{
						Required:    true,
						Description: "Port number for the device connection.",
					},
					"subnet_mask": schema.StringAttribute{
						Optional:    true,
						Description: "Subnet mask for the device.",
					},
					"gateway_ip_address": schema.StringAttribute{
						Optional:    true,
						Description: "Gateway IP address for the device.",
					},
				},
			},
			"meta_data": schema.StringAttribute{
				Optional:    true,
				Description: "Metadata for the device in JSON format.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"device_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Type of device (PLC, IPC, HMI, AGV, ROBOT, DRIVE, OTHER).",
				Default:     stringdefault.StaticString("PLC"),
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the device.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secret_id": schema.StringAttribute{
				Optional:    true,
				Description: "Secret ID for device credentials.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ftp_configuration": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "FTP configuration for the device.",
				Attributes: map[string]schema.Attribute{
					"ip_address": schema.StringAttribute{
						Required:    true,
						Description: "IP address of the FTP server on the device.",
					},
					"port": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: "Port of the FTP server on the device.",
						Default:     int64default.StaticInt64(22),
					},
					"protocol": schema.StringAttribute{
						Optional:    true,
						Description: "Protocol used by the FTP server (FTP, SFTP).",
					},
					"secret_id": schema.StringAttribute{
						Optional:    true,
						Description: "Secret ID for FTP server credentials.",
					},
					"root_directory": schema.StringAttribute{
						Optional:    true,
						Description: "Root directory of the FTP server on the device.",
					},
				},
			},
			"object_version": schema.Int64Attribute{
				Computed:    true,
				Description: "Version number of the object, used for optimistic locking and change tracking.",
			},
			"creation_user_id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier of the user who created this object.",
			},
			"update_user_id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier of the user who last updated this object.",
			},
			"creation_timestamp": schema.StringAttribute{
				Computed:    true,
				Description: "Date and time when this object was first created (ISO 8601 format).",
			},
			"update_timestamp": schema.StringAttribute{
				Computed:    true,
				Description: "Date and time when this object was last modified (ISO 8601 format).",
			},
		},
	}
}

func (r *DeviceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*clients.Client)
}

func (r *DeviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("device_id"), req, resp)
}

// Helper function to get connection configuration object type
func ConnectionConfigurationObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"ip_address":         types.StringType,
			"port":               types.Int64Type,
			"subnet_mask":        types.StringType,
			"gateway_ip_address": types.StringType,
		},
	}
}

// Helper function to get FTP configuration object type
func FtpConfigurationObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"ip_address":     types.StringType,
			"port":           types.Int64Type,
			"protocol":       types.StringType,
			"secret_id":      types.StringType,
			"root_directory": types.StringType,
		},
	}
}
