package link

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/sda/terraform-provider-sda/internal/clients"
)

var _ resource.Resource = &LinkResource{}
var _ resource.ResourceWithImportState = &LinkResource{}

func NewLinkResource() resource.Resource {
	return &LinkResource{}
}

type LinkResource struct {
	client *clients.Client
}

func (r *LinkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_link"
}

func (r *LinkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an asset link resource in the SDA Assets Management Service. Links connect two assets together with optional metadata.",
		Attributes: map[string]schema.Attribute{
			"source_id": schema.StringAttribute{
				Required:    true,
				Description: "The source ID of the asset link. This is the ID of the asset (resource) from which the link is created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_type": schema.StringAttribute{
				Required:    true,
				Description: "The source type of the asset link (DEVICE, DOCUMENT, GATEWAY, LICENSE, PROJECT, TAG, VAULT, SECRET, etc.).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"destination_id": schema.StringAttribute{
				Required:    true,
				Description: "The destination ID of the asset link. This is the ID of the asset (resource) to which the link is created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"destination_type": schema.StringAttribute{
				Required:    true,
				Description: "The destination type of the asset link (DEVICE, DOCUMENT, GATEWAY, LICENSE, PROJECT, TAG, VAULT, SECRET, etc.).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"meta_data": schema.StringAttribute{
				Optional:    true,
				Description: "Metadata for the asset link in JSON format. The structure depends on the asset types being linked.",
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

func (r *LinkResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*clients.Client)
}

func (r *LinkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: source_type/source_id/destination_type/destination_id
	// We'll parse this in a custom way
	idParts := parseImportID(req.ID)
	if len(idParts) != 4 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: source_type/source_id/destination_type/destination_id, got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("source_type"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("source_id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("destination_type"), idParts[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("destination_id"), idParts[3])...)
}

// Helper function to parse import ID
func parseImportID(id string) []string {
	parts := []string{}
	current := ""
	for _, char := range id {
		if char == '/' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
