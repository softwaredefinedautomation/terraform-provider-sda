package license

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/sda/terraform-provider-sda/internal/clients"
)

var _ resource.Resource = &LicenseResource{}
var _ resource.ResourceWithImportState = &LicenseResource{}

func NewLicenseResource() resource.Resource {
	return &LicenseResource{}
}

type LicenseResource struct {
	client *clients.Client
}

func (r *LicenseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_license"
}

func (r *LicenseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a license resource in the SDA Assets Management Service. Handles file upload using multipart upload.",
		Attributes: map[string]schema.Attribute{
			"license_id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the license.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Resource group ID to which this license belongs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vendor_id": schema.StringAttribute{
				Required:    true,
				Description: "Vendor ID of the license.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"serial_id": schema.StringAttribute{
				Required:    true,
				Description: "Serial ID of the license.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"product": schema.StringAttribute{
				Required:    true,
				Description: "Product name for the license.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Type of license (COOPERATE, FLOATING, SINGLE, UPGRADE, TRIAL).",
				Default:     stringdefault.StaticString("FLOATING"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Status of the license (REQUESTED, ACTIVE, UPLOADED, EXPIRED, INVALID).",
				Default:     stringdefault.StaticString("REQUESTED"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"quantity": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Quantity of licenses.",
				Default:     int64default.StaticInt64(1),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Name of the license.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ide_config_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "IDE configuration ID associated with the license.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"expiration_timestamp": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Expiration timestamp for the license (ISO 8601 format).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"family": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "License family.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"company_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Company name for the license.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"product_key": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Product key for the license.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"container_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Container ID for the license.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"firm_code": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Firm code for the license.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"license_server": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "License server address.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"file_path": schema.StringAttribute{
				Optional:    true,
				Description: "Local file path of the license file to upload. Optional if no license file is needed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the uploaded license file.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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

func (r *LicenseResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*clients.Client)
}

func (r *LicenseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("license_id"), req, resp)
}