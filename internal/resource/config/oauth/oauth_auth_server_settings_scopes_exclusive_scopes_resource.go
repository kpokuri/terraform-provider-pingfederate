package oauth

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthAuthServerSettingsScopesExclusiveScopesResource{}
	_ resource.ResourceWithConfigure   = &oauthAuthServerSettingsScopesExclusiveScopesResource{}
	_ resource.ResourceWithImportState = &oauthAuthServerSettingsScopesExclusiveScopesResource{}
)

// OauthAuthServerSettingsScopesExclusiveScopesResource is a helper function to simplify the provider implementation.
func OauthAuthServerSettingsScopesExclusiveScopesResource() resource.Resource {
	return &oauthAuthServerSettingsScopesExclusiveScopesResource{}
}

// oauthAuthServerSettingsScopesExclusiveScopesResource is the resource implementation.
type oauthAuthServerSettingsScopesExclusiveScopesResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type oauthAuthServerSettingsScopesExclusiveScopesResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Dynamic     types.Bool   `tfsdk:"dynamic"`
}

// GetSchema defines the schema for the resource.
func (r *oauthAuthServerSettingsScopesExclusiveScopesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	oauthAuthServerSettingsScopesExclusiveScopesResourceSchema(ctx, req, resp, false)
}

func oauthAuthServerSettingsScopesExclusiveScopesResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a OauthAuthServerSettingsScopesExclusiveScopes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Computed attribute tied to the name property of this resource.",
				Computed:    true,
				Optional:    false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Description: "The name of the scope.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the scope that appears when the user is prompted for authorization.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dynamic": schema.BoolAttribute{
				Description: "True if the scope is dynamic. (Defaults to false)",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}

	// Set attributes in string list
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"name", "description"})
	}
	resp.Schema = schema
}

func (r *oauthAuthServerSettingsScopesExclusiveScopesResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var model oauthAuthServerSettingsScopesExclusiveScopesResourceModel
	req.Plan.Get(ctx, &model)
	if model.Dynamic.ValueBool() && (model.Name.ValueString() != "" || !model.Name.IsNull()) {
		{
			containsAsteriskPrefix := strings.Index(model.Name.ValueString(), "*")
			if containsAsteriskPrefix != 0 {
				resp.Diagnostics.AddError("Dynamic property is set to true with Name property incorrectly specified!", "The Name property must be prefixed with an \"*\". For example, \"*example\"")
			}
		}
	}
}

func addOptionalOauthAuthServerSettingsScopesExclusiveScopesFields(ctx context.Context, addRequest *client.ScopeEntry, plan oauthAuthServerSettingsScopesExclusiveScopesResourceModel) error {

	if internaltypes.IsDefined(plan.Name) {
		addRequest.Name = plan.Name.ValueString()
	}
	if internaltypes.IsDefined(plan.Description) {
		addRequest.Description = plan.Description.ValueString()
	}
	if internaltypes.IsDefined(plan.Dynamic) {
		addRequest.Dynamic = plan.Dynamic.ValueBoolPointer()
	}
	return nil

}

// Metadata returns the resource type name.
func (r *oauthAuthServerSettingsScopesExclusiveScopesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_auth_server_settings_scopes_exclusive_scopes"
}

func (r *oauthAuthServerSettingsScopesExclusiveScopesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readOauthAuthServerSettingsScopesExclusiveScopesResponse(ctx context.Context, r *client.ScopeEntry, state *oauthAuthServerSettingsScopesExclusiveScopesResourceModel, expectedValues *oauthAuthServerSettingsScopesExclusiveScopesResourceModel) {
	state.Id = types.StringValue(r.Name)
	state.Name = types.StringValue(r.Name)
	state.Description = types.StringValue(r.Description)
	state.Dynamic = types.BoolPointerValue(r.Dynamic)
}

func (r *oauthAuthServerSettingsScopesExclusiveScopesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauthAuthServerSettingsScopesExclusiveScopesResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOauthAuthServerSettingsScopesExclusiveScopes := client.NewScopeEntry(plan.Name.ValueString(), plan.Description.ValueString())
	err := addOptionalOauthAuthServerSettingsScopesExclusiveScopesFields(ctx, createOauthAuthServerSettingsScopesExclusiveScopes, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OauthAuthServerSettingsScopesExclusiveScopes", err.Error())
		return
	}
	requestJson, err := createOauthAuthServerSettingsScopesExclusiveScopes.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateOauthAuthServerSettingsScopesExclusiveScopes := r.apiClient.OauthAuthServerSettingsApi.AddExclusiveScope(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateOauthAuthServerSettingsScopesExclusiveScopes = apiCreateOauthAuthServerSettingsScopesExclusiveScopes.Body(*createOauthAuthServerSettingsScopesExclusiveScopes)
	oauthAuthServerSettingsScopesExclusiveScopesResponse, httpResp, err := r.apiClient.OauthAuthServerSettingsApi.AddExclusiveScopeExecute(apiCreateOauthAuthServerSettingsScopesExclusiveScopes)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the OauthAuthServerSettingsScopesExclusiveScopes", err, httpResp)
		return
	}
	responseJson, err := oauthAuthServerSettingsScopesExclusiveScopesResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state oauthAuthServerSettingsScopesExclusiveScopesResourceModel

	readOauthAuthServerSettingsScopesExclusiveScopesResponse(ctx, oauthAuthServerSettingsScopesExclusiveScopesResponse, &state, &plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *oauthAuthServerSettingsScopesExclusiveScopesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readOauthAuthServerSettingsScopesExclusiveScopes(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readOauthAuthServerSettingsScopesExclusiveScopes(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var state oauthAuthServerSettingsScopesExclusiveScopesResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadOauthAuthServerSettingsScopesExclusiveScopes, httpResp, err := apiClient.OauthAuthServerSettingsApi.GetExclusiveScope(config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for a OauthAuthServerSettingsScopesExclusiveScopes", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadOauthAuthServerSettingsScopesExclusiveScopes.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readOauthAuthServerSettingsScopesExclusiveScopesResponse(ctx, apiReadOauthAuthServerSettingsScopesExclusiveScopes, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *oauthAuthServerSettingsScopesExclusiveScopesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateOauthAuthServerSettingsScopesExclusiveScopes(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateOauthAuthServerSettingsScopesExclusiveScopes(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan oauthAuthServerSettingsScopesExclusiveScopesResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state oauthAuthServerSettingsScopesExclusiveScopesResourceModel
	req.State.Get(ctx, &state)
	updateOauthAuthServerSettingsScopesExclusiveScopes := apiClient.OauthAuthServerSettingsApi.UpdateExclusiveScope(config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())
	createUpdateRequest := client.NewScopeEntry(plan.Name.ValueString(), plan.Description.ValueString())
	err := addOptionalOauthAuthServerSettingsScopesExclusiveScopesFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for OauthAuthServerSettingsScopesExclusiveScopes", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateOauthAuthServerSettingsScopesExclusiveScopes = updateOauthAuthServerSettingsScopesExclusiveScopes.Body(*createUpdateRequest)
	updateOauthAuthServerSettingsScopesExclusiveScopesResponse, httpResp, err := apiClient.OauthAuthServerSettingsApi.UpdateExclusiveScopeExecute(updateOauthAuthServerSettingsScopesExclusiveScopes)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating OauthAuthServerSettingsScopesExclusiveScopes", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateOauthAuthServerSettingsScopesExclusiveScopesResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	readOauthAuthServerSettingsScopesExclusiveScopesResponse(ctx, updateOauthAuthServerSettingsScopesExclusiveScopesResponse, &state, &plan)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// // Delete deletes the resource and removes the Terraform state on success.
func (r *oauthAuthServerSettingsScopesExclusiveScopesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	deleteOauthAuthServerSettingsScopesExclusiveScopes(ctx, req, resp, r.apiClient, r.providerConfig)
}
func deleteOauthAuthServerSettingsScopesExclusiveScopes(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from state
	var state oauthAuthServerSettingsScopesExclusiveScopesResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	httpResp, err := apiClient.OauthAuthServerSettingsApi.RemoveExclusiveScope(config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting a OauthAuthServerSettingsScopesExclusiveScopes", err, httpResp)
		return
	}

}

func (r *oauthAuthServerSettingsScopesExclusiveScopesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importOauthAuthServerSettingsScopesExclusiveScopesLocation(ctx, req, resp)
}
func importOauthAuthServerSettingsScopesExclusiveScopesLocation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}