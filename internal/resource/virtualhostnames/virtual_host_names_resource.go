package virtualHostNames

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	config "github.com/pingidentity/terraform-provider-pingfederate/internal/resource"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &virtualHostNamesResource{}
	_ resource.ResourceWithConfigure   = &virtualHostNamesResource{}
	_ resource.ResourceWithImportState = &virtualHostNamesResource{}
)

// VirtualHostNamesResource is a helper function to simplify the provider implementation.
func VirtualHostNamesResource() resource.Resource {
	return &virtualHostNamesResource{}
}

// virtualHostNamesResource is the resource implementation.
type virtualHostNamesResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type virtualHostNamesResourceModel struct {
	Id               types.String `tfsdk:"id"`
	VirtualHostNames types.Set    `tfsdk:"virtual_host_names"`
}

// GetSchema defines the schema for the resource.
func (r *virtualHostNamesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	virtualHostNamesResourceSchema(ctx, req, resp, false)
}

func virtualHostNamesResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a VirtualHostNames.",
		Attributes: map[string]schema.Attribute{
			"virtual_host_names": schema.SetAttribute{
				Description: "List of virtual host names.",
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown()},
			},
		},
	}

	// Set attributes in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{""})
	}
	config.AddCommonSchema(&schema, false)
	resp.Schema = schema
}
func addOptionalVirtualHostNamesFields(ctx context.Context, addRequest *client.VirtualHostNameSettings, plan virtualHostNamesResourceModel) error {

	if internaltypes.IsDefined(plan.VirtualHostNames) {
		var slice []string
		plan.VirtualHostNames.ElementsAs(ctx, &slice, false)
		addRequest.VirtualHostNames = slice
	}
	return nil

}

// Metadata returns the resource type name.
func (r *virtualHostNamesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_host_names"
}

func (r *virtualHostNamesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readVirtualHostNamesResponse(ctx context.Context, r *client.VirtualHostNameSettings, state *virtualHostNamesResourceModel, expectedValues *virtualHostNamesResourceModel) {
	state.Id = types.StringValue("id")
	state.VirtualHostNames = internaltypes.GetStringSet(r.VirtualHostNames)
}

func (r *virtualHostNamesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan virtualHostNamesResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createVirtualHostNames := client.NewVirtualHostNameSettings()
	err := addOptionalVirtualHostNamesFields(ctx, createVirtualHostNames, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for VirtualHostNames", err.Error())
		return
	}
	requestJson, err := createVirtualHostNames.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateVirtualHostNames := r.apiClient.VirtualHostNamesApi.UpdateVirtualHostNamesSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateVirtualHostNames = apiCreateVirtualHostNames.Body(*createVirtualHostNames)
	virtualHostNamesResponse, httpResp, err := r.apiClient.VirtualHostNamesApi.UpdateVirtualHostNamesSettingsExecute(apiCreateVirtualHostNames)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the VirtualHostNames", err, httpResp)
		return
	}
	responseJson, err := virtualHostNamesResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state virtualHostNamesResourceModel

	readVirtualHostNamesResponse(ctx, virtualHostNamesResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *virtualHostNamesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readVirtualHostNames(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readVirtualHostNames(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state virtualHostNamesResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadVirtualHostNames, httpResp, err := apiClient.VirtualHostNamesApi.GetVirtualHostNamesSettings(config.ProviderBasicAuthContext(ctx, providerConfig)).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a VirtualHostNames", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadVirtualHostNames.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readVirtualHostNamesResponse(ctx, apiReadVirtualHostNames, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *virtualHostNamesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateVirtualHostNames(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateVirtualHostNames(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan virtualHostNamesResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state virtualHostNamesResourceModel
	req.State.Get(ctx, &state)
	updateVirtualHostNames := apiClient.VirtualHostNamesApi.UpdateVirtualHostNamesSettings(config.ProviderBasicAuthContext(ctx, providerConfig))
	createUpdateRequest := client.NewVirtualHostNameSettings()
	err := addOptionalVirtualHostNamesFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for VirtualHostNames", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateVirtualHostNames = updateVirtualHostNames.Body(*createUpdateRequest)
	updateVirtualHostNamesResponse, httpResp, err := apiClient.VirtualHostNamesApi.UpdateVirtualHostNamesSettingsExecute(updateVirtualHostNames)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating VirtualHostNames", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateVirtualHostNamesResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readVirtualHostNamesResponse(ctx, updateVirtualHostNamesResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// This config object is edit-only, so Terraform can't delete it.
func (r *virtualHostNamesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *virtualHostNamesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
