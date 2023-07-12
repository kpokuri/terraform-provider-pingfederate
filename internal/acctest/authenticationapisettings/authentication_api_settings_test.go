package acctest_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/acctest"
	"github.com/pingidentity/terraform-provider-pingfederate/internal/provider"
)

// Attributes to test with. Add optional properties to test here if desired.
type authenticationApiSettingsResourceModel struct {
	apiEnabled                       bool
	enableApiDescriptions            bool
	restrictAccessToRedirectlessMode bool
	includeRequestContext            bool
}

func TestAccAuthenticationApiSettings(t *testing.T) {
	resourceName := "myAuthenticationApiSettings"
	initialResourceModel := authenticationApiSettingsResourceModel{
		apiEnabled:                       false,
		enableApiDescriptions:            false,
		restrictAccessToRedirectlessMode: false,
		includeRequestContext:            false,
	}

	updatedResourceModel := authenticationApiSettingsResourceModel{
		apiEnabled:                       true,
		enableApiDescriptions:            true,
		restrictAccessToRedirectlessMode: false,
		includeRequestContext:            false,
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAuthenticationApiSettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedAuthenticationApiSettingsAttributes(initialResourceModel),
			},
			{
				Config: testAccAuthenticationApiSettings(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedAuthenticationApiSettingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccAuthenticationApiSettings(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_authentication_api_settings." + resourceName,
				ImportStateId:     resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAuthenticationApiSettings(resourceName string, resourceModel authenticationApiSettingsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_authentication_api_settings" "%[1]s" {
  api_enabled                          = %[2]t
  enable_api_descriptions              = %[3]t
  restrict_access_to_redirectless_mode = %[4]t
  include_request_context              = %[5]t
}`, resourceName,
		resourceModel.apiEnabled,
		resourceModel.enableApiDescriptions,
		resourceModel.restrictAccessToRedirectlessMode,
		resourceModel.includeRequestContext,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedAuthenticationApiSettingsAttributes(config authenticationApiSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "AuthenticationApiSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.AuthenticationApiApi.GetAuthenticationApiSettings(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchBool(resourceType, nil, "api_enabled",
			config.apiEnabled, *response.ApiEnabled)
		if err != nil {
			return err
		}

		err = acctest.TestAttributesMatchBool(resourceType, nil, "enable_api_descriptions",
			config.enableApiDescriptions, *response.EnableApiDescriptions)
		if err != nil {
			return err
		}
		return nil
	}
}