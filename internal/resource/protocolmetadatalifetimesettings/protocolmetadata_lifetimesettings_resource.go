package protocolMetadataLifetimeSettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	config "github.com/pingidentity/terraform-provider-pingfederate/internal/resource"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &protocolMetadataLifetimeSettingsResource{}
	_ resource.ResourceWithConfigure   = &protocolMetadataLifetimeSettingsResource{}
	_ resource.ResourceWithImportState = &protocolMetadataLifetimeSettingsResource{}
)

// ProtocolMetadataLifetimeSettingsResource is a helper function to simplify the provider implementation.
func ProtocolMetadataLifetimeSettingsResource() resource.Resource {
	return &protocolMetadataLifetimeSettingsResource{}
}

// protocolMetadataLifetimeSettingsResource is the resource implementation.
type protocolMetadataLifetimeSettingsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type protocolMetadataLifetimeSettingsResourceModel struct {
	Id            types.String `tfsdk:"id"`
	CacheDuration types.Int64  `tfsdk:"cache_duration"`
	ReloadDelay   types.Int64  `tfsdk:"reload_delay"`
}

// GetSchema defines the schema for the resource.
func (r *protocolMetadataLifetimeSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	protocolMetadataLifetimeSettingsResourceSchema(ctx, req, resp, false)
}

func protocolMetadataLifetimeSettingsResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a ProtocolMetadataLifetimeSettings.",
		Attributes: map[string]schema.Attribute{
			"cache_duration": schema.Int64Attribute{
				Description: "This field adjusts the validity of your metadata in minutes. The default value is 1440 (1 day).",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown()},
			},
			"reload_delay": schema.Int64Attribute{
				Description: "This field adjusts the frequency of automatic reloading of SAML metadata in minutes. The default value is 1440 (1 day).",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}

	config.AddCommonSchema(&schema, false)
	resp.Schema = schema
}
func addOptionalProtocolMetadataLifetimeSettingsFields(ctx context.Context, addRequest *client.MetadataLifetimeSettings, plan protocolMetadataLifetimeSettingsResourceModel) error {

	if internaltypes.IsDefined(plan.CacheDuration) {
		addRequest.CacheDuration = plan.CacheDuration.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.ReloadDelay) {
		addRequest.ReloadDelay = plan.ReloadDelay.ValueInt64Pointer()
	}
	return nil

}

// Metadata returns the resource type name.
func (r *protocolMetadataLifetimeSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_protocolmetadata_lifetimesettings"
}

func (r *protocolMetadataLifetimeSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readProtocolMetadataLifetimeSettingsResponse(ctx context.Context, r *client.MetadataLifetimeSettings, state *protocolMetadataLifetimeSettingsResourceModel, expectedValues *protocolMetadataLifetimeSettingsResourceModel) {
	state.Id = types.StringValue("id")
	state.CacheDuration = types.Int64Value(*r.CacheDuration)
	state.ReloadDelay = types.Int64Value(*r.ReloadDelay)
}

func (r *protocolMetadataLifetimeSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan protocolMetadataLifetimeSettingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateProtocolMetadataLifetimeSettings := r.apiClient.ProtocolMetadataApi.UpdateLifetimeSettings(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	createUpdateRequest := client.NewMetadataLifetimeSettings()
	err := addOptionalProtocolMetadataLifetimeSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for ProtocolMetadataLifetimeSettings", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	updateProtocolMetadataLifetimeSettings = updateProtocolMetadataLifetimeSettings.Body(*createUpdateRequest)
	protocolMetadataLifetimeSettingsResponse, httpResp, err := r.apiClient.ProtocolMetadataApi.UpdateLifetimeSettingsExecute(updateProtocolMetadataLifetimeSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the ProtocolMetadataLifetimeSettings", err, httpResp)
		return
	}
	responseJson, err := protocolMetadataLifetimeSettingsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state protocolMetadataLifetimeSettingsResourceModel

	readProtocolMetadataLifetimeSettingsResponse(ctx, protocolMetadataLifetimeSettingsResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *protocolMetadataLifetimeSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readProtocolMetadataLifetimeSettings(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readProtocolMetadataLifetimeSettings(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state protocolMetadataLifetimeSettingsResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadProtocolMetadataLifetimeSettings, httpResp, err := apiClient.ProtocolMetadataApi.GetLifetimeSettings(config.ProviderBasicAuthContext(ctx, providerConfig)).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a ProtocolMetadataLifetimeSettings", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadProtocolMetadataLifetimeSettings.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readProtocolMetadataLifetimeSettingsResponse(ctx, apiReadProtocolMetadataLifetimeSettings, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *protocolMetadataLifetimeSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateProtocolMetadataLifetimeSettings(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateProtocolMetadataLifetimeSettings(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var plan protocolMetadataLifetimeSettingsResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateProtocolMetadataLifetimeSettings := apiClient.ProtocolMetadataApi.UpdateLifetimeSettings(config.ProviderBasicAuthContext(ctx, providerConfig))
	createUpdateRequest := client.NewMetadataLifetimeSettings()
	err := addOptionalProtocolMetadataLifetimeSettingsFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for ProtocolMetadataLifetimeSettings", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	updateProtocolMetadataLifetimeSettings = updateProtocolMetadataLifetimeSettings.Body(*createUpdateRequest)
	protocolMetadataLifetimeSettingsResponse, httpResp, err := apiClient.ProtocolMetadataApi.UpdateLifetimeSettingsExecute(updateProtocolMetadataLifetimeSettings)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the ProtocolMetadataLifetimeSettings", err, httpResp)
		return
	}
	responseJson, err := protocolMetadataLifetimeSettingsResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state protocolMetadataLifetimeSettingsResourceModel

	readProtocolMetadataLifetimeSettingsResponse(ctx, protocolMetadataLifetimeSettingsResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// This config object is edit-only, so Terraform can't delete it.
func (r *protocolMetadataLifetimeSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *protocolMetadataLifetimeSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocation(ctx, req, resp)
}
func importLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}