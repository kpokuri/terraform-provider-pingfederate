---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingfederate_authentication_api_settings Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Manages a AuthenticationApiSettings.
---

# pingfederate_authentication_api_settings (Resource)

Manages a AuthenticationApiSettings.

## Example Usage

```terraform
# this resource does not support import as the PF API only supports PUT Method
resource "pingfederate_authentication_api_settings" "authenticationApiSettingsExample" {
  api_enabled                          = true
  enable_api_descriptions              = false
  restrict_access_to_redirectless_mode = false
  include_request_context              = true
  # To remove a previously added default application ref, change id and location values to empty strings
  default_application_ref = {
    id       = ""
    location = ""
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `api_enabled` (Boolean) Enable Authentication API
- `default_application_ref` (Attributes) Enable API descriptions (see [below for nested schema](#nestedatt--default_application_ref))
- `enable_api_descriptions` (Boolean) Enable API descriptions
- `include_request_context` (Boolean) Includes request context in API responses
- `restrict_access_to_redirectless_mode` (Boolean) Enable restrict access to redirectless mode

### Read-Only

- `id` (String) Placeholder name of this object required by Terraform.

<a id="nestedatt--default_application_ref"></a>
### Nested Schema for `default_application_ref`

Optional:

- `id` (String)
- `location` (String)