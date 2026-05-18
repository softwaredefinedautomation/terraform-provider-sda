package user_role_association

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/sda/terraform-provider-sda/internal/clients"
)

var _ resource.Resource = &UserRoleAssociationResource{}
var _ resource.ResourceWithImportState = &UserRoleAssociationResource{}

func NewUserRoleAssociationResource() resource.Resource {
	return &UserRoleAssociationResource{}
}

type UserRoleAssociationResource struct {
	client *clients.Client
}

func (r *UserRoleAssociationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_role_association"
}

func (r *UserRoleAssociationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Link a user with a user role in the SDA Ident Service.",
		Attributes: map[string]schema.Attribute{
			"user_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier for the user.",
			},
			"user_role_id": schema.StringAttribute{
				Required:    true,
				Description: "Unique identifier for the user role.",
			},
			"expiration_timestamp": schema.StringAttribute{
				Optional:    true,
				Description: "Optional expiration date and time for this link (ISO 8601 format).",
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

func (r *UserRoleAssociationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*clients.Client)
}

func (r *UserRoleAssociationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("user_role_id"), req, resp)
}
