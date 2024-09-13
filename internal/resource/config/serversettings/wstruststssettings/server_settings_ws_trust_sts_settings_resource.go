// Code generated by ping-terraform-plugin-framework-generator

package serversettingswstruststssettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingfederate-go-client/v1210/configurationapi"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

var (
	emptyStringSetDefault, _  = types.SetValue(types.StringType, nil)
	resourceLinkSetDefault, _ = types.SetValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id": types.StringType,
		},
	}, nil)
	usersSetDefault, _ = types.SetValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"username": types.StringType,
			"password": types.StringType,
		},
	}, nil)
)

func (r *serverSettingsWsTrustStsSettingsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config serverSettingsWsTrustStsSettingsResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if config.BasicAuthnEnabled.ValueBool() {
		if config.Users.IsNull() || (internaltypes.IsDefined(config.Users) && len(config.Users.Elements()) == 0) {
			resp.Diagnostics.AddError("'basic_authn_enabled' can only be true if users are defined", "")
		}
	} else if internaltypes.IsDefined(config.Users) && len(config.Users.Elements()) > 0 {
		resp.Diagnostics.AddError("users can only be defined if 'basic_authn_enabled' is true", "")
	}

	if config.ClientCertAuthnEnabled.ValueBool() {
		if !config.RestrictByIssuerCert.ValueBool() && !config.RestrictBySubjectDn.ValueBool() {
			resp.Diagnostics.AddError("'client_cert_authn_enabled' can only be true if 'restrict_by_issuer_cert' or 'restrict_by_subject_dn' is true", "")
		}
	}

	if config.RestrictByIssuerCert.ValueBool() {
		if config.IssuerCerts.IsNull() || (internaltypes.IsDefined(config.IssuerCerts) && len(config.IssuerCerts.Elements()) == 0) {
			resp.Diagnostics.AddError("if 'restrict_by_issuer_cert' is true, issuer certs must be defined", "")
		}
		if !config.ClientCertAuthnEnabled.ValueBool() {
			resp.Diagnostics.AddError("'restrict_by_issuer_cert' can only be true if 'client_cert_authn_enabled' is true", "")
		}
	}

	if config.RestrictBySubjectDn.ValueBool() {
		if config.SubjectDns.IsNull() || (internaltypes.IsDefined(config.SubjectDns) && len(config.SubjectDns.Elements()) == 0) {
			resp.Diagnostics.AddError("if 'restrict_by_subject_dn' is true, subject DNs must be defined", "")
		}
		if !config.ClientCertAuthnEnabled.ValueBool() {
			resp.Diagnostics.AddError("'restrict_by_subject_dn' can only be true if 'client_cert_authn_enabled' is true", "")
		}
	}
}

func (state *serverSettingsWsTrustStsSettingsResourceModel) readClientResponseUsers(response *client.WsTrustStsSettings) diag.Diagnostics {
	var respDiags diag.Diagnostics
	usersAttrTypes := map[string]attr.Type{
		"password": types.StringType,
		"username": types.StringType,
	}
	usersElementType := types.ObjectType{AttrTypes: usersAttrTypes}
	var usersValues []attr.Value
	for _, usersResponseValue := range response.Users {
		var userPassword *string
		// Get password value from state, if it is set, since the PF API won't return the password
		if !state.Users.IsNull() && !state.Users.IsUnknown() {
			// Find the corresponding user in state, if it exists
			for _, user := range state.Users.Elements() {
				userAttrs := user.(types.Object).Attributes()
				password, ok := userAttrs["password"]
				if ok {
					userPassword = password.(types.String).ValueStringPointer()
					break
				}
			}
		}
		usersValue, diags := types.ObjectValue(usersAttrTypes, map[string]attr.Value{
			"password": types.StringPointerValue(userPassword),
			"username": types.StringPointerValue(usersResponseValue.Username),
		})
		respDiags.Append(diags...)
		usersValues = append(usersValues, usersValue)
	}
	usersValue, diags := types.SetValue(usersElementType, usersValues)
	respDiags.Append(diags...)

	state.Users = usersValue
	return respDiags
}

func (r *serverSettingsWsTrustStsSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource is singleton, so it can't be deleted from the service. Deleting this resource will remove it from Terraform state.
	// Instead this resource will be reset to the PingFederate default values.
	apiUpdateRequest := r.apiClient.ServerSettingsAPI.UpdateWsTrustStsSettings(config.AuthContext(ctx, r.providerConfig))
	apiUpdateRequest = apiUpdateRequest.Body(client.WsTrustStsSettings{})
	_, httpResp, err := r.apiClient.ServerSettingsAPI.UpdateWsTrustStsSettingsExecute(apiUpdateRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while resetting the serverSettingsWsTrustStsSettings", err, httpResp)
	}
}