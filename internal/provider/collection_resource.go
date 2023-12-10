package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &collectionResource{}
	_ resource.ResourceWithConfigure = &collectionResource{}
)

// NewCollectionResource is a helper function to simplify the provider implementation.
func NewCollectionResource() resource.Resource {
	return &collectionResource{}
}

// orderResource is the resource implementation.
type collectionResource struct {
	client Client
}

// CollectionResourceModel is the model for the solrcloud_collection resource.
type CollectionResourceModel struct {
	Name              types.String   `tfsdk:"name"`
	RouterName        types.String   `tfsdk:"router_name"`
	NumShards         types.Int64    `tfsdk:"num_shards"`
	ReplicationFactor types.Int64    `tfsdk:"replication_factor"`
	Shards            []types.String `tfsdk:"shards"`
}

// Configure adds the provider configured client to the resource.
func (r *collectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = *client
}

// Schema defines the schema for the resource.
func (r *collectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the collection to be created.",
			},
			"router_name": schema.StringAttribute{
				Optional:    true,
				Description: "The router name to be used. Possible values are implicit or compositeId.",
			},
			"num_shards": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of shards to be created as part of the collection.",
			},
			"replication_factor": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of replicas to be created for each shard.",
			},
			"shards": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "The shard names to use when creating this collection.",
			},
		},
	}
}

func (r *collectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection"
}

// Create creates the resource and sets the initial Terraform state.
func (r *collectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CollectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert plan.Shards from []types.String to []string
	var shards []string
	for _, shard := range plan.Shards {
		shards = append(shards, shard.String())
	}

	err := r.client.CreateCollection(plan.Name.String(), int(plan.NumShards.ValueInt64()), plan.RouterName.String(), int(plan.ReplicationFactor.ValueInt64()), shards)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating collection",
			"Could not create collection, unexpected error: "+err.Error(),
		)
		return
	}

	// Todo set the state

}

// Read refreshes the Terraform state with the latest data.
func (r *collectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *collectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *collectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
