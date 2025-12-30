package resourcegroup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/sda/terraform-provider-sda/internal/clients"
)

var _ resource.Resource = &ResourceGroupResource{}
var _ resource.ResourceWithImportState = &ResourceGroupResource{}

func NewResourceGroupResource() resource.Resource {
	return &ResourceGroupResource{}
}

type ResourceGroupResource struct {
	client *clients.Client
}

func (r *ResourceGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_group"
}

func (r *ResourceGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"group_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("ASSET"),
			},
			"parent_group_id": schema.StringAttribute{
				Optional: true,
			},
			"is_system_group": schema.BoolAttribute{
				Computed: true,
			},
			"object_version": schema.Int64Attribute{
				Computed: true,
			},
			"creation_user_id": schema.StringAttribute{
				Computed: true,
			},
			"update_user_id": schema.StringAttribute{
				Computed: true,
			},
			"creation_timestamp": schema.StringAttribute{
				Computed: true,
			},
			"update_timestamp": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *ResourceGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*clients.Client)
}

func (r *ResourceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("group_id"), req, resp)
}
