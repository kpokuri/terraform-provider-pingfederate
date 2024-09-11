// Code generated by ping-terraform-plugin-framework-generator

package oauthclientregistrationpolicies

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/id"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/importprivatestate"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/providererror"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &oauthClientRegistrationPolicyResource{}
	_ resource.ResourceWithConfigure   = &oauthClientRegistrationPolicyResource{}
	_ resource.ResourceWithImportState = &oauthClientRegistrationPolicyResource{}

	customId = "policy_id"
)

func OauthClientRegistrationPolicyResource() resource.Resource {
	return &oauthClientRegistrationPolicyResource{}
}

type oauthClientRegistrationPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *oauthClientRegistrationPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_client_registration_policy"
}

func (r *oauthClientRegistrationPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type oauthClientRegistrationPolicyResourceModel struct {
	Configuration       types.Object `tfsdk:"configuration"`
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	ParentRef           types.Object `tfsdk:"parent_ref"`
	PluginDescriptorRef types.Object `tfsdk:"plugin_descriptor_ref"`
	PolicyId            types.String `tfsdk:"policy_id"`
}

func (r *oauthClientRegistrationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to create and manage an OAuth Client Registration Policy.",
		Attributes: map[string]schema.Attribute{
			"configuration": pluginconfiguration.ToSchema(),
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The plugin instance name. The name can be modified once the instance is created.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"parent_ref": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
				Optional:    true,
				Description: "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances.",
			},
			"plugin_descriptor_ref": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
				Required:    true,
				Description: "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created.",
			},
			"policy_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the plugin instance. The ID cannot be modified once the instance is created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
	id.ToSchema(&resp.Schema)
}

func (r *oauthClientRegistrationPolicyResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *oauthClientRegistrationPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if plan == nil {
		return
	}
	var state *oauthClientRegistrationPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if state == nil {
		return
	}
	var respDiags diag.Diagnostics
	plan.Configuration, respDiags = pluginconfiguration.MarkComputedAttrsUnknownOnChange(plan.Configuration, state.Configuration)
	resp.Diagnostics.Append(respDiags...)
	resp.Plan.Set(ctx, plan)
}

func (model *oauthClientRegistrationPolicyResourceModel) buildClientStruct() (*client.ClientRegistrationPolicy, diag.Diagnostics) {
	result := &client.ClientRegistrationPolicy{}
	var respDiags diag.Diagnostics
	var err error
	// configuration
	configurationValue, err := pluginconfiguration.ClientStruct(model.Configuration)
	if err != nil {
		respDiags.AddError(providererror.InternalProviderError, "Error building client struct for configuration: "+err.Error())
	} else {
		result.Configuration = *configurationValue
	}

	// name
	result.Name = model.Name.ValueString()
	// parent_ref
	if !model.ParentRef.IsNull() {
		parentRefValue := &client.ResourceLink{}
		parentRefAttrs := model.ParentRef.Attributes()
		parentRefValue.Id = parentRefAttrs["id"].(types.String).ValueString()
		result.ParentRef = parentRefValue
	}

	// plugin_descriptor_ref
	pluginDescriptorRefValue := client.ResourceLink{}
	pluginDescriptorRefAttrs := model.PluginDescriptorRef.Attributes()
	pluginDescriptorRefValue.Id = pluginDescriptorRefAttrs["id"].(types.String).ValueString()
	result.PluginDescriptorRef = pluginDescriptorRefValue

	// policy_id
	result.Id = model.PolicyId.ValueString()
	return result, respDiags
}

func (state *oauthClientRegistrationPolicyResourceModel) readClientResponse(response *client.ClientRegistrationPolicy, isImportRead bool) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// id
	state.Id = types.StringValue(response.Id)
	// configuration
	configurationValue, diags := pluginconfiguration.ToState(state.Configuration, &response.Configuration, isImportRead)
	respDiags.Append(diags...)

	state.Configuration = configurationValue
	// name
	state.Name = types.StringValue(response.Name)
	// parent_ref
	parentRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	var parentRefValue types.Object
	if response.ParentRef == nil {
		parentRefValue = types.ObjectNull(parentRefAttrTypes)
	} else {
		parentRefValue, diags = types.ObjectValue(parentRefAttrTypes, map[string]attr.Value{
			"id": types.StringValue(response.ParentRef.Id),
		})
		respDiags.Append(diags...)
	}

	state.ParentRef = parentRefValue
	// plugin_descriptor_ref
	pluginDescriptorRefAttrTypes := map[string]attr.Type{
		"id": types.StringType,
	}
	pluginDescriptorRefValue, diags := types.ObjectValue(pluginDescriptorRefAttrTypes, map[string]attr.Value{
		"id": types.StringValue(response.PluginDescriptorRef.Id),
	})
	respDiags.Append(diags...)

	state.PluginDescriptorRef = pluginDescriptorRefValue
	// policy_id
	state.PolicyId = types.StringValue(response.Id)
	return respDiags
}

func (r *oauthClientRegistrationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data oauthClientRegistrationPolicyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiCreateRequest := r.apiClient.OauthClientRegistrationPoliciesAPI.CreateDynamicClientRegistrationPolicy(config.AuthContext(ctx, r.providerConfig))
	apiCreateRequest = apiCreateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.OauthClientRegistrationPoliciesAPI.CreateDynamicClientRegistrationPolicyExecute(apiCreateRequest)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while creating the oauthClientRegistrationPolicy", err, httpResp, &customId)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, false)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *oauthClientRegistrationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	isImportRead, diags := importprivatestate.IsImportRead(ctx, req, resp)
	resp.Diagnostics.Append(diags...)

	var data oauthClientRegistrationPolicyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.OauthClientRegistrationPoliciesAPI.GetDynamicClientRegistrationPolicy(config.AuthContext(ctx, r.providerConfig), data.PolicyId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.AddResourceNotFoundWarning(ctx, &resp.Diagnostics, "oauthClientRegistrationPolicy", httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while reading the oauthClientRegistrationPolicy", err, httpResp, &customId)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, isImportRead)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *oauthClientRegistrationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data oauthClientRegistrationPolicyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.OauthClientRegistrationPoliciesAPI.UpdateDynamicClientRegistrationPolicy(config.AuthContext(ctx, r.providerConfig), data.PolicyId.ValueString())
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.OauthClientRegistrationPoliciesAPI.UpdateDynamicClientRegistrationPolicyExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while updating the oauthClientRegistrationPolicy", err, httpResp, &customId)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData, false)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *oauthClientRegistrationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data oauthClientRegistrationPolicyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	httpResp, err := r.apiClient.OauthClientRegistrationPoliciesAPI.DeleteDynamicClientRegistrationPolicy(config.AuthContext(ctx, r.providerConfig), data.PolicyId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpErrorCustomId(ctx, &resp.Diagnostics, "An error occurred while deleting the oauthClientRegistrationPolicy", err, httpResp, &customId)
	}
}

func (r *oauthClientRegistrationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to policy_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("policy_id"), req, resp)
	importprivatestate.MarkPrivateStateForImport(ctx, resp)
}