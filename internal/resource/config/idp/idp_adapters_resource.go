package idp

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	internaljson "github.com/pingidentity/terraform-provider-pingfederate/internal/json"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &idpAdapterResource{}
	_ resource.ResourceWithConfigure   = &idpAdapterResource{}
	_ resource.ResourceWithImportState = &idpAdapterResource{}
)

// Define attribute types for object types
var (
	fieldsAttrTypes = map[string]attr.Type{
		"name":            types.StringType,
		"value":           types.StringType,
		"encrypted_value": types.StringType,
		"inherited":       types.BoolType,
	}
	tablesAttrTypes = map[string]attr.Type{
		"name": types.StringType,
		"rows": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"fields": types.ListType{
						ElemType: types.ObjectType{
							AttrTypes: fieldsAttrTypes,
						},
					},
					"default_row": types.BoolType,
				},
			},
		},
		"inherited": types.BoolType,
	}
	configurationAttrTypes = map[string]attr.Type{
		"tables": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: tablesAttrTypes,
			},
		},
		"fields": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: fieldsAttrTypes,
			},
		},
	}
	attributeContractAttrTypes = map[string]attr.Type{
		"core_attributes": types.SetType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name":      types.StringType,
					"pseudonym": types.BoolType,
					"masked":    types.BoolType,
				},
			},
		},
		"extended_attributes": types.SetType{
			ElemType: types.ObjectType{
				//TODO more duplication
				AttrTypes: map[string]attr.Type{
					"name":      types.StringType,
					"pseudonym": types.BoolType,
					"masked":    types.BoolType,
				},
			},
		},
		"unique_user_key_attribute": types.StringType,
		"mask_ognl_values":          types.BoolType,
		"inherited":                 types.BoolType,
	}
	attributeMappingAttrTypes = map[string]attr.Type{
		"attribute_contract_fulfillment": types.MapType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"source": types.ObjectType{
						//TODO remove this duplication?
						AttrTypes: map[string]attr.Type{
							"type": types.StringType,
							"id":   types.StringType,
						},
					},
					"value": types.StringType,
				},
			},
		},
		"issuance_criteria": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"conditional_criteria": types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"source": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"type": types.StringType,
									"id":   types.StringType,
								},
							},
							"attribute_name": types.StringType,
							"condition":      types.StringType,
							"value":          types.StringType,
							"error_result":   types.StringType,
						},
					},
				},
				"expression_criteria": types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"expression":   types.StringType,
							"error_result": types.StringType,
						},
					},
				},
			},
		},
		"inherited": types.BoolType,
	}
)

// IdpAdapterResource is a helper function to simplify the provider implementation.
func IdpAdapterResource() resource.Resource {
	return &idpAdapterResource{}
}

// idpAdapterResource is the resource implementation.
type idpAdapterResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type idpAdapterResourceModel struct {
	AuthnCtxClassRef    types.String `tfsdk:"authn_ctx_class_ref"`
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	PluginDescriptorRef types.Object `tfsdk:"plugin_descriptor_ref"`
	ParentRef           types.Object `tfsdk:"parent_ref"`
	Configuration       types.Object `tfsdk:"configuration"`
	AttributeMapping    types.Object `tfsdk:"attribute_mapping"`
	AttributeContract   types.Object `tfsdk:"attribute_contract"`
}

// GetSchema defines the schema for the resource.
func (r *idpAdapterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Idp Adapter",
		Attributes: map[string]schema.Attribute{
			"authn_ctx_class_ref": schema.StringAttribute{
				Description: "The fixed value that indicates how the user was authenticated.",
				Optional:    true,
			},
			//TODO don't add id in common schema
			"id": schema.StringAttribute{
				Description: "The ID of the plugin instance. The ID cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The plugin instance name. The name can be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
			},
			"plugin_descriptor_ref": schema.SingleNestedAttribute{
				Description: "Reference to the plugin descriptor for this instance. The plugin descriptor cannot be modified once the instance is created. Note: Ignored when specifying a connection's adapter override.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The ID of the resource.",
						Required:    true,
					},
					"location": schema.StringAttribute{
						Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
						Optional:    true,
					},
				},
			},
			"parent_ref": schema.SingleNestedAttribute{
				Description: "The reference to this plugin's parent instance. The parent reference is only accepted if the plugin type supports parent instances. Note: This parent reference is required if this plugin instance is used as an overriding plugin (e.g. connection adapter overrides)",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "The ID of the resource.",
						Required:    true,
					},
					"location": schema.StringAttribute{
						Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
						Optional:    true,
					},
				},
			},

			"configuration": schema.SingleNestedAttribute{
				Description: "Plugin instance configuration.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"tables": schema.ListNestedAttribute{
						Description: "List of configuration tables.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of the table.",
									Required:    true,
								},
								"rows": schema.ListNestedAttribute{
									Description: "List of table rows.",
									Optional:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"fields": schema.ListNestedAttribute{
												Description: "The configuration fields in the row.",
												Required:    true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"name": schema.StringAttribute{
															Description: "The name of the configuration field.",
															Required:    true,
														},
														"value": schema.StringAttribute{
															Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
															Optional:    true,
														},
														"encrypted_value": schema.StringAttribute{
															Description: "For encrypted or hashed fields, this attribute contains the encrypted representation of the field's value, if a value is defined. If you do not want to update the stored value, this attribute should be passed back unchanged.",
															Optional:    true,
														},
														"inherited": schema.BoolAttribute{
															Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
															Optional:    true,
														},
													},
												},
											},
											"default_row": schema.BoolAttribute{
												Description: "Whether this row is the default.",
												Optional:    true,
											},
										},
									},
								},
								"inherited": schema.BoolAttribute{
									Description: "Whether this table is inherited from its parent instance. If true, the rows become read-only. The default value is false.",
									Optional:    true,
								},
							},
						},
					},
					"fields": schema.ListNestedAttribute{
						Description: "List of configuration fields.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of the configuration field.",
									Required:    true,
								},
								"value": schema.StringAttribute{
									Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute. To update an encrypted or hashed field, specify the new value in this attribute.",
									Optional:    true,
								},
								"encrypted_value": schema.StringAttribute{
									Description: "For encrypted or hashed fields, this attribute contains the encrypted representation of the field's value, if a value is defined. If you do not want to update the stored value, this attribute should be passed back unchanged.",
									Optional:    true,
								},
								"inherited": schema.BoolAttribute{
									Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
									Optional:    true,
								},
							},
						},
					},
				},
			},

			"attribute_contract": schema.SingleNestedAttribute{
				Description: "The list of attributes that the IdP adapter provides.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"core_attributes": schema.SetNestedAttribute{
						Description: "A list of IdP adapter attributes that correspond to the attributes exposed by the IdP adapter type.",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
								},
								"pseudonym": schema.BoolAttribute{
									Description: "Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to false.",
									Optional:    true,
								},
								"masked": schema.BoolAttribute{
									Description: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
									Optional:    true,
								},
							},
						},
					},
					"extended_attributes": schema.SetNestedAttribute{
						Description: "A list of additional attributes that can be returned by the IdP adapter. The extended attributes are only used if the adapter supports them.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "The name of this attribute.",
									Required:    true,
								},
								"pseudonym": schema.BoolAttribute{
									Description: "Specifies whether this attribute is used to construct a pseudonym for the SP. Defaults to false.",
									Optional:    true,
								},
								"masked": schema.BoolAttribute{
									Description: "Specifies whether this attribute is masked in PingFederate logs. Defaults to false.",
									Optional:    true,
								},
							},
						},
					},
					"unique_user_key_attribute": schema.StringAttribute{
						Description: "The attribute to use for uniquely identify a user's authentication sessions.",
						Optional:    true,
					},
					"mask_ognl_values": schema.BoolAttribute{
						Description: "Whether or not all OGNL expressions used to fulfill an outgoing assertion contract should be masked in the logs. Defaults to false.",
						Optional:    true,
					},
					"inherited": schema.BoolAttribute{
						Description: "Whether this attribute contract is inherited from its parent instance. If true, the rest of the properties in this model become read-only. The default value is false.",
						Optional:    true,
					},
				},
			},

			"attribute_mapping": schema.SingleNestedAttribute{
				Description: "The attributes mapping from attribute sources to attribute targets.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					//TODO add attribute_sources
					/*"attribute_sources": schema.ListNestedAttribute{
						Description: "",
						NestedObject: schema.NestedAttributeObject{
							"ldap_attr_source": schema.SingleNestedAttribute{
								//TODO ldap specific schema here
							},
							"jdbc_attr_source": schema.SingleNestedAttribute{
								//TODO jdbc specific schema here
							},
							"custom_attr_source": schema.SingleNestedAttribute{
								//TODO custom specific schema here
							},
							Validators: []validator.Object{
								objectvalidator.AtLeastOneOf(
									path.MatchRoot("ldap_attr_source"),
									path.MatchRoot("jdbc_attr_source"),
									path.MatchRoot("custom_attr_source")),
							},
						},
					},*/
					"attribute_contract_fulfillment": schema.MapNestedAttribute{
						Description: "A list of mappings from attribute names to their fulfillment values.",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"source": schema.SingleNestedAttribute{
									Description: "The attribute value source.",
									Required:    true,
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Description: "The source type of this key.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.OneOf([]string{"TOKEN_EXCHANGE_PROCESSOR_POLICY", "ACCOUNT_LINK", "ADAPTER", "ASSERTION", "CONTEXT", "CUSTOM_DATA_STORE", "EXPRESSION", "JDBC_DATA_STORE", "LDAP_DATA_STORE", "PING_ONE_LDAP_GATEWAY_DATA_STORE", "MAPPED_ATTRIBUTES", "NO_MAPPING", "TEXT", "TOKEN", "REQUEST", "OAUTH_PERSISTENT_GRANT", "SUBJECT_TOKEN", "ACTOR_TOKEN", "PASSWORD_CREDENTIAL_VALIDATOR", "IDP_CONNECTION", "AUTHENTICATION_POLICY_CONTRACT", "CLAIMS", "LOCAL_IDENTITY_PROFILE", "EXTENDED_CLIENT_METADATA", "EXTENDED_PROPERTIES", "TRACKED_HTTP_PARAMS", "FRAGMENT", "INPUTS", "ATTRIBUTE_QUERY", "IDENTITY_STORE_USER", "IDENTITY_STORE_GROUP", "SCIM_USER", "SCIM_GROUP"}...),
											},
										},
										"id": schema.StringAttribute{
											Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
											Optional:    true,
										},
									},
								},
								"value": schema.StringAttribute{
									Description: "The value for this attribute.",
									Required:    true,
								},
							},
						},
					},
					"issuance_criteria": schema.SingleNestedAttribute{
						Description: "The issuance criteria that this transaction must meet before the corresponding attribute contract is fulfilled.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"conditional_criteria": schema.ListNestedAttribute{
								Description: "An issuance criterion that checks a source attribute against a particular condition and the expected value. If the condition is true then this issuance criterion passes, otherwise the criterion fails.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										//TODO any way to share these definitions
										"source": schema.SingleNestedAttribute{
											Description: "The attribute value source.",
											Required:    true,
											Attributes: map[string]schema.Attribute{
												"type": schema.StringAttribute{
													Description: "The source type of this key.",
													Required:    true,
													Validators: []validator.String{
														stringvalidator.OneOf([]string{"TOKEN_EXCHANGE_PROCESSOR_POLICY", "ACCOUNT_LINK", "ADAPTER", "ASSERTION", "CONTEXT", "CUSTOM_DATA_STORE", "EXPRESSION", "JDBC_DATA_STORE", "LDAP_DATA_STORE", "PING_ONE_LDAP_GATEWAY_DATA_STORE", "MAPPED_ATTRIBUTES", "NO_MAPPING", "TEXT", "TOKEN", "REQUEST", "OAUTH_PERSISTENT_GRANT", "SUBJECT_TOKEN", "ACTOR_TOKEN", "PASSWORD_CREDENTIAL_VALIDATOR", "IDP_CONNECTION", "AUTHENTICATION_POLICY_CONTRACT", "CLAIMS", "LOCAL_IDENTITY_PROFILE", "EXTENDED_CLIENT_METADATA", "EXTENDED_PROPERTIES", "TRACKED_HTTP_PARAMS", "FRAGMENT", "INPUTS", "ATTRIBUTE_QUERY", "IDENTITY_STORE_USER", "IDENTITY_STORE_GROUP", "SCIM_USER", "SCIM_GROUP"}...),
													},
												},
												"id": schema.StringAttribute{
													Description: "The attribute source ID that refers to the attribute source that this key references. In some resources, the ID is optional and will be ignored. In these cases the ID should be omitted. If the source type is not an attribute source then the ID can be omitted.",
													Optional:    true,
												},
											},
										},
										"attribute_name": schema.StringAttribute{
											Description: "The name of the attribute to use in this issuance criterion.",
											Required:    true,
										},
										"condition": schema.StringAttribute{
											Description: "The condition that will be applied to the source attribute's value and the expected value.",
											Required:    true,
											Validators: []validator.String{
												stringvalidator.OneOf([]string{"EQUALS", "EQUALS_CASE_INSENSITIVE", "EQUALS_DN", "NOT_EQUAL", "NOT_EQUAL_CASE_INSENSITIVE", "NOT_EQUAL_DN", "MULTIVALUE_CONTAINS", "MULTIVALUE_CONTAINS_CASE_INSENSITIVE", "MULTIVALUE_CONTAINS_DN", "MULTIVALUE_DOES_NOT_CONTAIN", "MULTIVALUE_DOES_NOT_CONTAIN_CASE_INSENSITIVE", "MULTIVALUE_DOES_NOT_CONTAIN_DN"}...),
											},
										},
										"value": schema.StringAttribute{
											Description: "The expected value of this issuance criterion.",
											Required:    true,
										},
										"error_result": schema.StringAttribute{
											Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
											Optional:    true,
										},
									},
								},
							},
							"expression_criteria": schema.ListNestedAttribute{
								Description: "An issuance criterion that uses a Boolean return value from an OGNL expression to determine whether or not it passes.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"expression": schema.StringAttribute{
											Required:    true,
											Description: "The OGNL expression to evaluate.",
										},
										"error_result": schema.StringAttribute{
											Optional:    true,
											Description: "The error result to return if this issuance criterion fails. This error result will show up in the PingFederate server logs.",
										},
									},
								},
							},
						},
					},
					"inherited": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether this attribute mapping is inherited from its parent instance. If true, the rest of the properties in this model become read-only. The default value is false.",
					},
				},
			},
		},
	}
}

func addOptionalIdpAdapterFields(ctx context.Context, addRequest *client.IdpAdapter, plan idpAdapterResourceModel) error {
	if internaltypes.IsDefined(plan.AuthnCtxClassRef) {
		addRequest.AuthnCtxClassRef = plan.AuthnCtxClassRef.ValueStringPointer()
	}

	if internaltypes.IsDefined(plan.ParentRef) {
		addRequest.ParentRef = &client.ResourceLink{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.ParentRef, false)), addRequest.ParentRef)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AttributeMapping) {
		addRequest.AttributeMapping = &client.IdpAdapterContractMapping{}
		// Don't include non-specified attributes, since PF will complain about issuance_criteria.expression_criteria
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeMapping, true)), addRequest.AttributeMapping)
		if err != nil {
			return err
		}
	}

	if internaltypes.IsDefined(plan.AttributeContract) {
		addRequest.AttributeContract = &client.IdpAdapterAttributeContract{}
		err := json.Unmarshal([]byte(internaljson.FromValue(plan.AttributeContract, false)), addRequest.AttributeContract)
		if err != nil {
			return err
		}
	}

	return nil
}

// Metadata returns the resource type name.
func (r *idpAdapterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_idp_adapters"
}

func (r *idpAdapterResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient

}

func readIdpAdapterResponse(ctx context.Context, r *client.IdpAdapter, state *idpAdapterResourceModel, diags *diag.Diagnostics) {
	state.AuthnCtxClassRef = internaltypes.StringTypeOrNil(r.AuthnCtxClassRef, false)
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Name)
	state.PluginDescriptorRef = internaltypes.ToStateResourceLink(ctx, &r.PluginDescriptorRef, diags)
	state.ParentRef = internaltypes.ToStateResourceLink(ctx, r.ParentRef, diags)

	var objectValueFromDiags diag.Diagnostics

	state.Configuration, objectValueFromDiags = types.ObjectValueFrom(ctx, configurationAttrTypes, r.Configuration)
	diags.Append(objectValueFromDiags...)

	if r.AttributeContract != nil {
		state.AttributeContract, objectValueFromDiags = types.ObjectValueFrom(ctx, attributeContractAttrTypes, r.AttributeContract)
		diags.Append(objectValueFromDiags...)
	}

	if r.AttributeMapping != nil {
		state.AttributeMapping, objectValueFromDiags = types.ObjectValueFrom(ctx, attributeMappingAttrTypes, r.AttributeMapping)
		diags.Append(objectValueFromDiags...)
	}
}

func (r *idpAdapterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan idpAdapterResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var pluginDescriptorRef client.ResourceLink
	err := json.Unmarshal([]byte(internaljson.FromValue(plan.PluginDescriptorRef, false)), &pluginDescriptorRef)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read plugin_descriptor_ref from plan", err.Error())
		return
	}

	var configuration client.PluginConfiguration
	err = json.Unmarshal([]byte(internaljson.FromValue(plan.Configuration, false)), &configuration)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read configuration from plan", err.Error())
		return
	}

	createIdpAdapter := client.NewIdpAdapter(plan.Id.ValueString(), plan.Name.ValueString(), pluginDescriptorRef, configuration)
	err = addOptionalIdpAdapterFields(ctx, createIdpAdapter, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for IdpAdapter", err.Error())
		return
	}
	requestJson, err := createIdpAdapter.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}

	apiCreateIdpAdapter := r.apiClient.IdpAdaptersApi.CreateIdpAdapter(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiCreateIdpAdapter = apiCreateIdpAdapter.Body(*createIdpAdapter)
	idpAdapterResponse, httpResp, err := r.apiClient.IdpAdaptersApi.CreateIdpAdapterExecute(apiCreateIdpAdapter)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the IdpAdapter", err, httpResp)
		return
	}
	responseJson, err := idpAdapterResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state idpAdapterResourceModel

	readIdpAdapterResponse(ctx, idpAdapterResponse, &state, &resp.Diagnostics)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *idpAdapterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state idpAdapterResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiReadIdpAdapter, httpResp, err := r.apiClient.IdpAdaptersApi.GetIdpAdapter(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()

	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while looking for an IdpAdapter", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := apiReadIdpAdapter.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readIdpAdapterResponse(ctx, apiReadIdpAdapter, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *idpAdapterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan idpAdapterResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	updateIdpAdapter := r.apiClient.IdpAdaptersApi.UpdateIdpAdapter(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	var pluginDescriptorRef client.ResourceLink
	err := json.Unmarshal([]byte(internaljson.FromValue(plan.PluginDescriptorRef, false)), &pluginDescriptorRef)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read plugin_descriptor_ref from plan", err.Error())
		return
	}

	var configuration client.PluginConfiguration
	err = json.Unmarshal([]byte(internaljson.FromValue(plan.Configuration, false)), &configuration)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read configuration from plan", err.Error())
		return
	}

	createUpdateRequest := client.NewIdpAdapter(plan.Id.ValueString(), plan.Name.ValueString(), pluginDescriptorRef, configuration)

	err = addOptionalIdpAdapterFields(ctx, createUpdateRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for IdpAdapter", err.Error())
		return
	}
	requestJson, err := createUpdateRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Update request: "+string(requestJson))
	}
	updateIdpAdapter = updateIdpAdapter.Body(*createUpdateRequest)
	updateIdpAdapterResponse, httpResp, err := r.apiClient.IdpAdaptersApi.UpdateIdpAdapterExecute(updateIdpAdapter)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating IdpAdapter", err, httpResp)
		return
	}
	// Log response JSON
	responseJson, err := updateIdpAdapterResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}
	// Read the response
	var state idpAdapterResourceModel
	readIdpAdapterResponse(ctx, updateIdpAdapterResponse, &state, &resp.Diagnostics)

	// Update computed values
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

}

// This config object is edit-only, so Terraform can't delete it. This method will just remove it from state
func (r *idpAdapterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *idpAdapterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
