package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// TestProvider instantiates the provider with default values and asserts that it is valid.
func TestProvider(t *testing.T) {
	testProv := New("test")()
	_, diags := testProv.(provider.Provider).Configure(context.Background(), provider.ConfigureRequest{
		Config: tfsdk.NewValue(map[string]tfsdk.Value{
			"solrcloud_endpoint": tfsdk.Value{
				Type:  types.StringType,
				Value: types.String{Value: "http://localhost:8983/solr"},
			},
		}),
	})

	assert.False(t, diags.HasError())
}

// TestProviderDataSources checks if the provider declares the expected data sources.
func TestProviderDataSources(t *testing.T) {
	testProv := New("test")()
	dataSources := testProv.(provider.Provider).DataSources(context.Background())

	// Replace "solrcloud_collections" with your actual data source name
	assert.Contains(t, dataSources, func() datasource.DataSource { return NewSolrCloudCollectionsDataSource() })
}

// TestProviderResources checks if the provider declares the expected resources.
func TestProviderResources(t *testing.T) {
	testProv := New("test")()
	resources := testProv.(provider.Provider).Resources(context.Background())

	// Add assertions for your resources here
	// assert.Contains(t, resources, func() resource.Resource { return NewYourResource() })
}
