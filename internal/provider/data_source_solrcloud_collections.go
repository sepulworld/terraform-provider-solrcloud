package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &collectionsDataSource{}
	_ datasource.DataSourceWithConfigure = &collectionsDataSource{}
)

func NewCollectionsDataSource() datasource.DataSource {
	return &collectionsDataSource{}
}

type collectionsDataSource struct {
	client Client
}

// Configure adds the provider configured client to the data source.
func (d *collectionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *solrcloud.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = *client
}

func (d *collectionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collections"
}

// Schema defines the schema for the data source.
func (d *collectionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"collections": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

type collectionsDataSourceModel struct {
	Collections []types.String `tfsdk:"collections"`
}

func (d *collectionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state collectionsDataSourceModel

	collections, err := d.client.GetCollections()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch collections",
			fmt.Sprintf("Unable to fetch collections: %s", err),
		)
		return
	}

	for _, collectionName := range collections.Collections {
		state.Collections = append(state.Collections, types.StringValue(collectionName))
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
