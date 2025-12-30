package vault

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/sda/terraform-provider-sda/internal/clients"
)

var _ resource.Resource = &VaultResource{}
var _ resource.ResourceWithImportState = &VaultResource{}

func NewVaultResource() resource.Resource {
	return &VaultResource{}
}

type VaultResource struct {
	client *clients.Client
}

func (r *VaultResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vault"
}

func (r *VaultResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a vault resource in the SDA Assets Management Service.",
		Attributes: map[string]schema.Attribute{
			"vault_id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the vault.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Description: "Resource group ID to which this vault belongs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the vault.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the vault.",
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

func (r *VaultResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*clients.Client)
}

func (r *VaultResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("vault_id"), req, resp)
}
