package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/sda/terraform-provider-sda/internal/clients"

	"github.com/sda/terraform-provider-sda/internal/provider/device"
	"github.com/sda/terraform-provider-sda/internal/provider/document"
	"github.com/sda/terraform-provider-sda/internal/provider/gateway"
	"github.com/sda/terraform-provider-sda/internal/provider/link"
	"github.com/sda/terraform-provider-sda/internal/provider/resourcegroup"
	"github.com/sda/terraform-provider-sda/internal/provider/secret"
	"github.com/sda/terraform-provider-sda/internal/provider/tag"
	"github.com/sda/terraform-provider-sda/internal/provider/vault"
	"github.com/sda/terraform-provider-sda/internal/provider/license"
	"github.com/sda/terraform-provider-sda/internal/provider/project"
	"github.com/sda/terraform-provider-sda/internal/provider/role"
	"github.com/sda/terraform-provider-sda/internal/provider/user"
)

const (
	version = "0.1.0"
)

var (
	_ provider.Provider = &SDAProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SDAProvider{
			version: version,
		}
	}
}

type SDAProvider struct {
	version string
}

type SDAProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *SDAProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sda"
	resp.Version = p.version
}

func (p *SDAProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "SDA API Host.",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "SDA user account username: Provided via SDA_USERNAME environment variable.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "SDA user account password: Provided via SDA_PASSWORD environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *SDAProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring SDA client")

	// Retrieve provider data from configuration
	var config SDAProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	host := os.Getenv("SDA_HOST")
	username := os.Getenv("SDA_USERNAME")
	password := os.Getenv("SDA_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing SDA API Host",
			"The provider cannot create the SDA API client as there is a missing or empty value for the SDA API host. "+
				"Set the host value in the configuration or use the SDA_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing SDA API Username",
			"The provider cannot create the SDA API client as there is a missing or empty value for the SDA API username. "+
				"Set the username value in the configuration or use the SDA_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing SDA API Password",
			"The provider cannot create the SDA API client as there is a missing or empty value for the SDA API password. "+
				"Set the password value in the configuration or use the SDA_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "sda_host", host)
	ctx = tflog.SetField(ctx, "sda_username", username)
	ctx = tflog.SetField(ctx, "sda_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "sda_password")

	tflog.Debug(ctx, "Creating SDA client")

	// Create a new SDA REST client using the configuration values
	restclient, err := clients.NewRestClient(&host, &username, &password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create SDA API Client",
			"An unexpected error occurred when creating the SDA API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"SDA Client Error: "+err.Error(),
		)
		return
	}

	// Make the SDA client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = restclient
	resp.ResourceData = restclient

	tflog.Info(ctx, "Configured SDA client", map[string]any{"success": true})
}

func (p *SDAProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resourcegroup.NewResourceGroupResource,
		document.NewDocumentResource,
		gateway.NewGatewayResource,
		vault.NewVaultResource,
		secret.NewSecretResource,
		tag.NewTagResource,
		device.NewDeviceResource,
		link.NewLinkResource,
		license.NewLicenseResource,
		role.NewRoleResource,
		user.NewUserResource,
		project.NewProjectResource,
	}
}

func (p *SDAProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTenantDataSource,
	}
}

func (p *SDAProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}
