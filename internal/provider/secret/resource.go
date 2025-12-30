package secret

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/sda/terraform-provider-sda/internal/clients"
)

var _ resource.Resource = &SecretResource{}
var _ resource.ResourceWithImportState = &SecretResource{}

func NewSecretResource() resource.Resource {
	return &SecretResource{}
}

type SecretResource struct {
	client *clients.Client
}

func (r *SecretResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

func (r *SecretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a secret resource in the SDA Assets Management Service.",
		Attributes: map[string]schema.Attribute{
			"secret_id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the secret.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vault_id": schema.StringAttribute{
				Optional:    true,
				Description: "ID of the vault this secret belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the secret.",
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: "Username of the secret.",
			},
			"secret_value": schema.StringAttribute{
				Required:    true,
				Description: "Secret value. Marked sensitive in Terraform UI.",
				Sensitive:   true,
			},
			"secret_type": schema.StringAttribute{
				Required:    true,
				Description: "Secret type. Marked sensitive in Terraform UI.",
			},
			"object_version": schema.Int64Attribute{
				Computed:    true,
				Description: "Version number used for optimistic locking.",
			},
			"creation_user_id": schema.StringAttribute{
				Computed:    true,
				Description: "User who created the secret.",
			},
			"update_user_id": schema.StringAttribute{
				Computed:    true,
				Description: "User who last updated the secret.",
			},
			"creation_timestamp": schema.StringAttribute{
				Computed:    true,
				Description: "Creation timestamp (ISO 8601).",
			},
			"update_timestamp": schema.StringAttribute{
				Computed:    true,
				Description: "Last-update timestamp (ISO 8601).",
			},
		},
	}
}

func (r *SecretResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*clients.Client)
}

func (r *SecretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("secret_id"), req, resp)
}
