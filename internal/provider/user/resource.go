package user

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"

    "github.com/sda/terraform-provider-sda/internal/clients"
)

var _ resource.Resource = &UserResource{}
var _ resource.ResourceWithImportState = &UserResource{}

func NewUserResource() resource.Resource {
    return &UserResource{}
}

type UserResource struct {
    client *clients.Client
}

func (r *UserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Description: "Manages a user resource in the SDA Ident Service.",
        Attributes: map[string]schema.Attribute{
            "user_id": schema.StringAttribute{
                Computed:    true,
                Description: "Unique identifier for the user.",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "group_id": schema.StringAttribute{
                Optional:    true,
                Description: "Optional group id for the user.",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "first_name": schema.StringAttribute{
                Required:    true,
                Description: "First name of the user.",
            },
            "last_name": schema.StringAttribute{
                Required:    true,
                Description: "Last name of the user.",
            },
            "email": schema.StringAttribute{
                Required:    true,
                Description: "Email address of the user.",
            },
            "company_name": schema.StringAttribute{
                Optional:    true,
                Description: "Company name for the user.",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "phone_number": schema.StringAttribute{
                Optional:    true,
                Description: "Phone number for the user.",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "privacy_accepted": schema.BoolAttribute{
                Optional:    true,
                Description: "Whether the user has accepted privacy terms.",
            },
            "locale": schema.StringAttribute{
                Optional:    true,
                Description: "Locale of the user.",
            },
            "last_login_timestamp": schema.StringAttribute{
                Computed:    true,
                Description: "Last login timestamp (ISO 8601).",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "title": schema.StringAttribute{
                Optional:    true,
                Description: "Title of the user.",
            },
            "agree_to_contact": schema.BoolAttribute{
                Optional:    true,
                Description: "Whether the user agreed to be contacted.",
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
            "source": schema.StringAttribute{
                Computed:    true,
                Description: "Source of the user (SDA or SAML).",
                Default:     stringdefault.StaticString("SDA"),
            },
        },
    }
}

func (r *UserResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }
    r.client = req.ProviderData.(*clients.Client)
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("user_id"), req, resp)
}
