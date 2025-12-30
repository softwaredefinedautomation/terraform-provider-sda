package role

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"

    "github.com/sda/terraform-provider-sda/internal/clients"
)

var _ resource.Resource = &RoleResource{}
var _ resource.ResourceWithImportState = &RoleResource{}

func NewRoleResource() resource.Resource {
    return &RoleResource{}
}

type RoleResource struct {
    client *clients.Client
}

func (r *RoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_role"
}

func (r *RoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Manages a user role resource in the SDA Ident Service.",
        Attributes: map[string]schema.Attribute{
            "user_role_id": schema.StringAttribute{
                Computed:    true,
                Description: "Unique identifier for the user role.",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                Required:    true,
                Description: "Name of the user role.",
            },
            "group_id": schema.StringAttribute{
                Optional:    true,
                Description: "Optional group id for the role.",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "description": schema.StringAttribute{
                Optional:    true,
                Description: "Description of the role.",
            },
            "policies": schema.ListAttribute{
                ElementType: types.StringType,
                Optional:    true,
                Description: "List of policies for the role (string representation).",
            },
            "is_system_role": schema.BoolAttribute{
                Computed:    true,
                Description: "Whether this is a system role.",
            },
            "sso_group_mapping": schema.ListAttribute{
                ElementType: types.StringType,
                Optional:    true,
                Description: "SSO group mapping list.",
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
                Description: "Date and time when this object was first created (ISO 8601).",
            },
            "update_timestamp": schema.StringAttribute{
                Computed:    true,
                Description: "Date and time when this object was last modified (ISO 8601).",
            },
        },
    }
}

func (r *RoleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }
    r.client = req.ProviderData.(*clients.Client)
}

func (r *RoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("user_role_id"), req, resp)
}
