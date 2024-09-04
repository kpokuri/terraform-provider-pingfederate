---
page_title: "pingfederate_openid_connect_settings Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Manages OpenID Connect configuration settings
---

# pingfederate_openid_connect_settings (Resource)

Manages OpenID Connect configuration settings

## Example Usage

```terraform
resource "pingfederate_oauth_open_id_connect_policy" "oauthOIDCPolicyExample" {
  policy_id = "oidcPolicy"
  name      = "oidcPolicy"
  access_token_manager_ref = {
    id = pingfederate_oauth_access_token_manager.example.manager_id
  }
  attribute_contract = {
    extended_attributes = []
  }
  attribute_mapping = {
    attribute_contract_fulfillment = {
      "sub" = {
        source = {
          type = "TOKEN"
        }
        value = "Username"
      }
    }
  }
  return_id_token_on_refresh_grant = false
  include_sri_in_id_token          = false
  include_s_hash_in_id_token       = false
  include_user_info_in_id_token    = false
  reissue_id_token_in_hybrid_flow  = false
  id_token_lifetime                = 5
}

resource "pingfederate_openid_connect_settings" "openIdConnectSettingsExample" {
  default_policy_ref = {
    id = pingfederate_oauth_open_id_connect_policy.oauthOIDCPolicyExample.policy_id
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `default_policy_ref` (Attributes) Reference to the default policy. (see [below for nested schema](#nestedatt--default_policy_ref))

### Optional

- `session_settings` (Attributes) The session settings (see [below for nested schema](#nestedatt--session_settings))

### Read-Only

- `id` (String, Deprecated) The ID of this resource.

<a id="nestedatt--default_policy_ref"></a>
### Nested Schema for `default_policy_ref`

Required:

- `id` (String) The ID of the resource.


<a id="nestedatt--session_settings"></a>
### Nested Schema for `session_settings`

Optional:

- `revoke_user_session_on_logout` (Boolean, Deprecated) Determines whether the user's session is revoked on logout. The default is `true`.
- `session_revocation_lifetime` (Number, Deprecated) How long a session revocation is tracked and stored, in minutes. The default is `490`. Value must be between `1` and `432001`, inclusive.
- `track_user_sessions_for_logout` (Boolean, Deprecated) Determines whether user sessions are tracked for logout. The default is `false`.

## Import

Import is supported using the following syntax:

~> This resource is singleton, so the value of "id" doesn't matter - it is just a placeholder, and required by Terraform

```shell
terraform import pingfederate_openid_connect_settings.openIdConnectSettings id
```