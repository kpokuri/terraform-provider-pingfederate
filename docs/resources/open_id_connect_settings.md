---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingfederate_open_id_connect_settings Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Manages OpenID Connect configuration settings
---

# pingfederate_open_id_connect_settings (Resource)

Manages OpenID Connect configuration settings

## Example Usage

```terraform
resource "pingfederate_open_id_connect_settings" "openIdConnectSettingsExample" {
  default_policy_ref = {
    id = "oauth_example_policy"
  }
  session_settings = {
    track_user_sessions_for_logout = true
    revoke_user_session_on_logout  = false
    session_revocation_lifetime    = 180
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `default_policy_ref` (Attributes) Reference to the default policy. (see [below for nested schema](#nestedatt--default_policy_ref))
- `session_settings` (Attributes) The session settings (see [below for nested schema](#nestedatt--session_settings))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--default_policy_ref"></a>
### Nested Schema for `default_policy_ref`

Required:

- `id` (String) The ID of the resource.

Read-Only:

- `location` (String, Deprecated) A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.


<a id="nestedatt--session_settings"></a>
### Nested Schema for `session_settings`

Optional:

- `revoke_user_session_on_logout` (Boolean) Determines whether the user's session is revoked on logout. This property is now available under /session/settings and should be accessed through that resource.
- `session_revocation_lifetime` (Number) How long a session revocation is tracked and stored, in minutes. This property is now available under /session/settings and should be accessed through that resource.
- `track_user_sessions_for_logout` (Boolean) Determines whether user sessions are tracked for logout. This property is now available under /oauth/authServerSettings and should be accessed through that resource.

## Import

Import is supported using the following syntax:

```shell
# This resource is singleton, so the value of "openIdConnectSettingsId" doesn't matter - it is just a placeholder, and required by Terraform
terraform import pingfederate_open_id_connect_settings.openIdConnectSettings openIdConnectSettingsId
```