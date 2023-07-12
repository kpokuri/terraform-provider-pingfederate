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

const protocolMetadataLifetimeSettingsId = "id"

// Attributes to test with. Add optional properties to test here if desired.
type protocolMetadataLifetimeSettingsResourceModel struct {
	id            string
	cacheDuration int64
	reloadDelay   int64
}

func TestAccProtocolMetadataLifetimeSettings(t *testing.T) {
	resourceName := "myProtocolMetadataLifetimeSettings"
	initialResourceModel := protocolMetadataLifetimeSettingsResourceModel{
		cacheDuration: 1,
		reloadDelay:   1,
	}
	updatedResourceModel := protocolMetadataLifetimeSettingsResourceModel{
		cacheDuration: 1440,
		reloadDelay:   1440,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccProtocolMetadataLifetimeSettings(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedProtocolMetadataLifetimeSettingsAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccProtocolMetadataLifetimeSettings(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedProtocolMetadataLifetimeSettingsAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccProtocolMetadataLifetimeSettings(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_protocol_metadata_lifetime_settings." + resourceName,
				ImportStateId:     protocolMetadataLifetimeSettingsId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProtocolMetadataLifetimeSettings(resourceName string, resourceModel protocolMetadataLifetimeSettingsResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_protocol_metadata_lifetime_settings" "%[1]s" {
  cache_duration = %[2]d
  reload_delay   = %[3]d
}`, resourceName,
		resourceModel.cacheDuration,
		resourceModel.reloadDelay,
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedProtocolMetadataLifetimeSettingsAttributes(config protocolMetadataLifetimeSettingsResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "ProtocolMetadataLifetimeSettings"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.ProtocolMetadataApi.GetLifetimeSettings(ctx).Execute()

		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "cache_duration",
			config.cacheDuration, *response.CacheDuration)
		if err != nil {
			return err
		}
		err = acctest.TestAttributesMatchInt(resourceType, &config.id, "reload_delay",
			config.reloadDelay, *response.ReloadDelay)
		if err != nil {
			return err
		}

		return nil
	}
}