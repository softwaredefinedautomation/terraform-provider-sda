package document

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/sda/terraform-provider-sda/internal/clients"
)

var _ resource.Resource = &DocumentResource{}
var _ resource.ResourceWithImportState = &DocumentResource{}

func NewDocumentResource() resource.Resource {
	return &DocumentResource{}
}

type DocumentResource struct {
	client *clients.Client
}

func (r *DocumentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_document"
}

func (r *DocumentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a document resource in the SDA Assets Management Service. Handles file upload using multipart upload.",
		Attributes: map[string]schema.Attribute{
			"document_id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the document.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Optional:    true,
				Description: "Resource group ID to which this document belongs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the document.",
			},
			"document_type": schema.StringAttribute{
				Required:    true,
				Description: "Type of document (PDF, MD, CSV, DOCX, TXT, XML, HTML, JSON, OTHERS).",
			},
			"last_version_number": schema.Int64Attribute{
				Computed:    true,
				Description: "Last version number of the document.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"file_path": schema.StringAttribute{
				Required:    true,
				Description: "Local file path of the document to upload.",
			},
			"file_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the uploaded file.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"commit_message": schema.StringAttribute{
				Optional:    true,
				Description: "Commit message for the document version.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"version_id": schema.StringAttribute{
				Computed:    true,
				Description: "Version ID of the uploaded document.",
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

func (r *DocumentResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*clients.Client)
}

func (r *DocumentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("document_id"), req, resp)
}
