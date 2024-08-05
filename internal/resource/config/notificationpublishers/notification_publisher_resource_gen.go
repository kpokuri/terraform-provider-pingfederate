// Code generated by ping-terraform-plugin-framework-generator

package notificationpublishers

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
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/common/pluginconfiguration"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	_ resource.Resource                = &notificationPublisherResource{}
	_ resource.ResourceWithConfigure   = &notificationPublisherResource{}
	_ resource.ResourceWithImportState = &notificationPublisherResource{}
)

func NotificationPublisherResource() resource.Resource {
	return &notificationPublisherResource{}
}

type notificationPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

func (r *notificationPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_publisher"
}

func (r *notificationPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type notificationPublisherResourceModel struct {
	Configuration       types.Object `tfsdk:"configuration"`
	Name                types.String `tfsdk:"name"`
	ParentRef           types.Object `tfsdk:"parent_ref"`
	PluginDescriptorRef types.Object `tfsdk:"plugin_descriptor_ref"`
	PublisherId         types.String `tfsdk:"publisher_id"`
}

func (r *notificationPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to create and manage notification publisher plugin instances.",
		Attributes: map[string]schema.Attribute{
			"configuration": pluginconfiguration.ToSchema(),
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The plugin instance name. The name can be modified once the instance is created.<br>Note: Ignored when specifying a connection's adapter override.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"parent_ref": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource.",
					},
				},
				Optional:    true,
				Description: "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)",
			},
			"plugin_descriptor_ref": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:    true,
						Description: "The ID of the resource.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
				Required:    true,
				Description: "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
			},
			"publisher_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the plugin instance. The ID cannot be modified once the instance is created.<br>Note: Ignored when specifying a connection's adapter override.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

func (r *notificationPublisherResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *notificationPublisherResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if plan == nil {
		return
	}
	var state *notificationPublisherResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if state == nil {
		return
	}
	var respDiags diag.Diagnostics
	plan.Configuration, respDiags = pluginconfiguration.MarkComputedAttrsUnknownOnChange(plan.Configuration, state.Configuration)
	resp.Diagnostics.Append(respDiags...)
	resp.Plan.Set(ctx, plan)
}

func (model *notificationPublisherResourceModel) buildClientStruct() (*client.NotificationPublisher, diag.Diagnostics) {
	result := &client.NotificationPublisher{}
	var respDiags diag.Diagnostics
	var err error
	// configuration
	configurationValue, err := pluginconfiguration.ClientStruct(model.Configuration)
	if err != nil {
		respDiags.AddError("Error building client struct for configuration", err.Error())
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

	// publisher_id
	result.Id = model.PublisherId.ValueString()
	return result, respDiags
}

func (state *notificationPublisherResourceModel) readClientResponse(response *client.NotificationPublisher) diag.Diagnostics {
	var respDiags, diags diag.Diagnostics
	// configuration
	configurationValue, diags := pluginconfiguration.ToState(state.Configuration, &response.Configuration)
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
	// publisher_id
	state.PublisherId = types.StringValue(response.Id)
	return respDiags
}

func (r *notificationPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data notificationPublisherResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiCreateRequest := r.apiClient.NotificationPublishersAPI.CreateNotificationPublisher(config.AuthContext(ctx, r.providerConfig))
	apiCreateRequest = apiCreateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.NotificationPublishersAPI.CreateNotificationPublisherExecute(apiCreateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the notificationPublisher", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *notificationPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data notificationPublisherResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	responseData, httpResp, err := r.apiClient.NotificationPublishersAPI.GetNotificationPublisher(config.AuthContext(ctx, r.providerConfig), data.PublisherId.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while reading the notificationPublisher", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while reading the notificationPublisher", err, httpResp)
		}
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *notificationPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data notificationPublisherResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientData, diags := data.buildClientStruct()
	resp.Diagnostics.Append(diags...)
	apiUpdateRequest := r.apiClient.NotificationPublishersAPI.UpdateNotificationPublisher(config.AuthContext(ctx, r.providerConfig), data.PublisherId.ValueString())
	apiUpdateRequest = apiUpdateRequest.Body(*clientData)
	responseData, httpResp, err := r.apiClient.NotificationPublishersAPI.UpdateNotificationPublisherExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the notificationPublisher", err, httpResp)
		return
	}

	// Read response into the model
	resp.Diagnostics.Append(data.readClientResponse(responseData)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *notificationPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data notificationPublisherResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	httpResp, err := r.apiClient.NotificationPublishersAPI.DeleteNotificationPublisher(config.AuthContext(ctx, r.providerConfig), data.PublisherId.ValueString()).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the notificationPublisher", err, httpResp)
	}
}

func (r *notificationPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to publisher_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("publisher_id"), req, resp)
}
