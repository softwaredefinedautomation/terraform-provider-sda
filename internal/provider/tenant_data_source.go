package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/sda/terraform-provider-sda/internal/clients"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &tenantDataSource{}
	_ datasource.DataSourceWithConfigure = &tenantDataSource{}
)

func NewTenantDataSource() datasource.DataSource {
	return &tenantDataSource{}
}

type tenantDataSource struct {
	client *clients.Client
}

// **1. Framework Data Model:** Used for Terraform state/schema mapping.
// MUST use `tfsdk` tags and framework types (types.String).
type tenantModel struct {
	TenantId          types.String `tfsdk:"tenant_id"`
	OwnerId           types.String `tfsdk:"owner_id"`
	CompanyName       types.String `tfsdk:"company_name"`
	CreationTimestamp types.String `tfsdk:"creation_timestamp"`
	ImportDemoData    types.Bool   `tfsdk:"import_demo_data"`
}

// **2. API Response Model:** Used ONLY for receiving raw JSON from the API.
// MUST use standard Go types (string) and `json` tags.
type apiTenantModel struct {
	TenantId          string `json:"tenant_id"`
	OwnerId           string `json:"owner_id"`
	CompanyName       string `json:"company_name"`
	CreationTimestamp string `json:"creation_timestamp"`
	ImportDemoData    bool   `json:"import_demo_data"`
}

// tenantDataSourceModel maps the top-level "tenant" attribute.
type tenantDataSourceModel struct {
	Tenant tenantModel `tfsdk:"tenant"`
}

// Configure adds the provider configured client to the data source.
func (d *tenantDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*clients.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *clients.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *tenantDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tenant"
}

// Schema defines the schema for the data source.
func (d *tenantDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"tenant": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"tenant_id":          schema.StringAttribute{Computed: true},
					"owner_id":           schema.StringAttribute{Computed: true},
					"company_name":       schema.StringAttribute{Computed: true},
					"creation_timestamp": schema.StringAttribute{Computed: true},
					"import_demo_data":   schema.BoolAttribute{Computed: true},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *tenantDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tenantDataSourceModel

	tenant, err := d.GetTenant()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SDA Tenant",
			err.Error(),
		)
		return
	}

	// Map response body to model
	tenantState := tenantModel{
		TenantId:          types.StringValue(tenant.TenantId),
		OwnerId:           types.StringValue(tenant.OwnerId),
		CompanyName:       types.StringValue(tenant.CompanyName),
		CreationTimestamp: types.StringValue(tenant.CreationTimestamp),
		ImportDemoData:    types.BoolValue(tenant.ImportDemoData),
	}

	state.Tenant = tenantState

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *tenantDataSource) GetTenant() (apiTenantModel, error) { // Note the return type: apiTenantModel
	// Use standard http request construction
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/ident/v1/tenant", d.client.HostURL), nil)
	if err != nil {
		return apiTenantModel{}, err
	}

	// d.client.DoRequest executes the request and returns the raw body
	body, err := d.client.DoRequest(req, nil)
	if err != nil {
		return apiTenantModel{}, err
	}

	// Unmarshal into the plain Go struct
	apiTenant := apiTenantModel{}
	err = json.Unmarshal(body, &apiTenant)
	if err != nil {
		return apiTenantModel{}, err
	}

	return apiTenant, nil
}
