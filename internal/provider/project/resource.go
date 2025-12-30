package project

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/sda/terraform-provider-sda/internal/clients"
)

var _ resource.Resource = &ProjectResource{}
var _ resource.ResourceWithImportState = &ProjectResource{}

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

type ProjectResource struct {
	client *clients.Client
}

func (r *ProjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a project resource in the SDA Assets Management Service. Handles file upload using multipart upload.",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the project.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Description: "Resource group ID to which this project belongs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the project.",
			},
			"vendor_id": schema.StringAttribute{
				Required:    true,
				Description: "Vendor ID of the project.",
			},
			"ide_config_id": schema.StringAttribute{
				Required:    true,
				Description: "IDE configuration ID for the project.",
			},
			"project_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Type of project (STANDARD, LIBRARY, GENERIC).",
				Default:     stringdefault.StaticString("STANDARD"),
			},
			"last_version_number": schema.Int64Attribute{
				Computed:    true,
				Description: "Last version number of the project.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the project.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secret_id": schema.StringAttribute{
				Optional:    true,
				Description: "Secret ID for project credentials.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"attached_licenses": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of license IDs attached to this project.",
			},
			"file_path": schema.StringAttribute{
				Required:    true,
				Description: "Local file path of the project file to upload.",
			},
			"file_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the uploaded file.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"version_id": schema.StringAttribute{
				Computed:    true,
				Description: "Version ID of the uploaded project.",
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

func (r *ProjectResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*clients.Client)
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("project_id"), req, resp)
}