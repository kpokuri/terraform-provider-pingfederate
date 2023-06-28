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

const virtualHostNamesId = "id"

// Attributes to test with. Add optional properties to test here if desired.
type virtualHostNamesResourceModel struct {
	id               string
	virtualHostNames []string
}

func TestAccVirtualHostNames(t *testing.T) {
	resourceName := "myVirtualHostNames"
	initialResourceModel := virtualHostNamesResourceModel{
		virtualHostNames: []string{"example1", "example2"},
	}
	updatedResourceModel := virtualHostNamesResourceModel{
		virtualHostNames: []string{"example1", "example2", "example3"},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.ConfigurationPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"pingfederate": providerserver.NewProtocol6WithError(provider.New()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualHostNames(resourceName, initialResourceModel),
				Check:  testAccCheckExpectedVirtualHostNamesAttributes(initialResourceModel),
			},
			{
				// Test updating some fields
				Config: testAccVirtualHostNames(resourceName, updatedResourceModel),
				Check:  testAccCheckExpectedVirtualHostNamesAttributes(updatedResourceModel),
			},
			{
				// Test importing the resource
				Config:            testAccVirtualHostNames(resourceName, updatedResourceModel),
				ResourceName:      "pingfederate_virtual_host_names." + resourceName,
				ImportStateId:     virtualHostNamesId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccVirtualHostNames(resourceName string, resourceModel virtualHostNamesResourceModel) string {
	return fmt.Sprintf(`
resource "pingfederate_virtual_host_names" "%[1]s" {
  virtual_host_names = %[2]s
}`, resourceName,
		acctest.StringSliceToTerraformString(resourceModel.virtualHostNames),
	)
}

// Test that the expected attributes are set on the PingFederate server
func testAccCheckExpectedVirtualHostNamesAttributes(config virtualHostNamesResourceModel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceType := "VirtualHostNames"
		testClient := acctest.TestClient()
		ctx := acctest.TestBasicAuthContext()
		response, _, err := testClient.VirtualHostNamesApi.GetVirtualHostNamesSettings(ctx).Execute()
		if err != nil {
			return err
		}

		// Verify that attributes have expected values
		err = acctest.TestAttributesMatchStringSlice(resourceType, &config.id, "virtual_host_names",
			config.virtualHostNames, response.GetVirtualHostNames())
		if err != nil {
			return err
		}

		return nil
	}
}