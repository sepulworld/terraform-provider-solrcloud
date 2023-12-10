package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure SolrCloudProvider satisfies various provider interfaces.
var _ provider.Provider = &SolrCloudProvider{}

// SolrCloudProvider defines the provider implementation.
type SolrCloudProvider struct {
	version string
}

// SolrCloudProviderModel describes the provider data model.
type solrCloudProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *SolrCloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "solrcloud"
	resp.Version = p.version
}

func (p *SolrCloudProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "The hostname of the SolrCloud API",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username for SolrCloud API authentication",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password for SolrCloud API authentication",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *SolrCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config solrCloudProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown SlorCloud API Host",
			"The provider cannot create the SolrCloud API client as there is an unknown configuration value for the SolrCloud API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SolrCLOUD_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown SolrCloud API Username",
			"The provider cannot create the SolrCloud API client as there is an unknown configuration value for the SolrCloud API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SOLCLOUD_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown SolrCloud API Password",
			"The provider cannot create the SolrCloud API client as there is an unknown configuration value for the SolrCloud API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SOLCLOUD_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("SOLRCLOUD_HOST")
	username := os.Getenv("SOLCLOUD_USERNAME")
	password := os.Getenv("SOLCLOUD_PASSWORD")

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
			"Missing SolrCloud API Host",
			"The provider cannot create the SolrCloud API client as there is a missing or empty value for the SolrCloud API host. "+
				"Set the host value in the configuration or use the SolrCLOUD_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing SolrCloud API Username",
			"The provider cannot create the SolrCloud API client as there is a missing or empty value for the SolrCloud API username. "+
				"Set the username value in the configuration or use the SOLCLOUD_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing SolrCloud API Password",
			"The provider cannot create the SolrCloud API client as there is a missing or empty value for the SolrCloud API password. "+
				"Set the password value in the configuration or use the SOLCLOUD_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "solrcloud_host", host)
	ctx = tflog.SetField(ctx, "solrcloud_username", username)
	ctx = tflog.SetField(ctx, "solrcloud_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "solrcloud_password")

	tflog.Debug(ctx, "Creating HashiCups client")

	client, err := NewClient(&host, &username, &password)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create SolrCloud API Client",
			"An unexpected error occurred when creating the SolrCloud API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"SolrCloud Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *SolrCloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCollectionResource,
	}
}

func (p *SolrCloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCollectionsDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SolrCloudProvider{
			version: version,
		}
	}
}
