package keypairs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingfederate-go-client"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingfederate/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &keyPairsSigningImportDataSource{}
	_ datasource.DataSourceWithConfigure = &keyPairsSigningImportDataSource{}
)

// Create a Administrative Account data source
func NewKeyPairsSigningImportDataSource() datasource.DataSource {
	return &keyPairsSigningImportDataSource{}
}

// keyPairsSigningImportDataSource is the datasource implementation.
type keyPairsSigningImportDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

type keyPairsSigningImportDataSourceModel struct {
	Id             types.String `tfsdk:"id"`
	FileData       types.String `tfsdk:"file_data"`
	Format         types.String `tfsdk:"format"`
	Password       types.String `tfsdk:"password"`
	CryptoProvider types.String `tfsdk:"crypto_provider"`
}

// GetSchema defines the schema for the datasource.
func (r *keyPairsSigningImportDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a KeyPairsSigningImport.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The persistent, unique ID for the certificate. It can be any combination of [a-z0-9._-]. This property is system-assigned if not specified.",
				Computed:    true,
				Optional:    true,
			},
			"file_data": schema.StringAttribute{
				Description: "Base-64 encoded PKCS12 or PEM file data. In the case of PEM, the raw (non-base-64) data is also accepted. In BCFIPS mode, only PEM with PBES2 and AES or Triple DES encryption is accepted and 128-bit salt is required.",
				Required:    true,
			},
			"format": schema.StringAttribute{
				Description: "Key pair file format. If specified, this field will control what file format is expected, otherwise the format will be auto-detected. In BCFIPS mode, only PEM is supported. (PKCS12, PEM)",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for the file. In BCFIPS mode, the password must be at least 14 characters.",
				Required:    true,
				Sensitive:   true,
			},
			"crypto_provider": schema.StringAttribute{
				Description: "Cryptographic Provider. This is only applicable if Hybrid HSM mode is true. (LOCAL, HSM)",
				Computed:    true,
				Optional:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Metadata returns the data source type name.
func (r *keyPairsSigningImportDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key_pair_signing_import"
}

// Configure adds the provider configured client to the data source.
func (r *keyPairsSigningImportDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Read a DseeCompatAdministrativeAccountResponse object into the model struct
func readKeyPairsSigningImportResponseDataSource(ctx context.Context, r *client.KeyPairView, state *keyPairsSigningImportDataSourceModel, expectedValues *keyPairsSigningImportDataSourceModel, planFileData string, planFormat string, planPassword string) {
	state.Id = internaltypes.StringTypeOrNil(r.Id, false)
	state.FileData = internaltypes.StringTypeOrNil(&planFileData, false)
	state.Format = internaltypes.StringTypeOrNil(&planFormat, false)
	state.Password = types.StringValue(planPassword)
	state.CryptoProvider = internaltypes.StringTypeOrNil(r.CryptoProvider, false)
}

// Read resource information
func (r *keyPairsSigningImportDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state keyPairsSigningImportDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReadKeyPairsSigningImport, httpResp, err := r.apiClient.KeyPairsSigningApi.GetSigningKeyPair(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the KeyPair Signing Import", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, responseErr := apiReadKeyPairsSigningImport.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	} else {
		diags.AddError("There was an issue retrieving the response of the KeyPair Signing Import: %s", responseErr.Error())
	}

	// Read the response into the state
	stateFileData := state.FileData.ValueString()
	stateFormat := state.Format.ValueString()
	statePassword := state.Password.ValueString()
	readKeyPairsSigningImportResponseDataSource(ctx, apiReadKeyPairsSigningImport, &state, &state, stateFileData, stateFormat, statePassword)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
