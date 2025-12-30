package gateway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/sda/terraform-provider-sda/internal/clients"
)

var _ resource.Resource = &GatewayResource{}
var _ resource.ResourceWithImportState = &GatewayResource{}

func NewGatewayResource() resource.Resource {
	return &GatewayResource{}
}

type GatewayResource struct {
	client *clients.Client
}

func (r *GatewayResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway"
}

func (r *GatewayResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a gateway resource in the SDA Assets Management Service.",
		Attributes: map[string]schema.Attribute{
			"gateway_id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the gateway.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				//Computed:    true,
				Optional:    true,
				Description: "Resource group ID to which this gateway belongs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(), // Copy state to plan if not configured
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the gateway.",
			},
			"description": schema.StringAttribute{
				//Computed:    true,
				Optional:    true,
				Description: "Description of the gateway.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(), // Copy state to plan if not configured
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

func (r *GatewayResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*clients.Client)
}

func (r *GatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("gateway_id"), req, resp)
}
