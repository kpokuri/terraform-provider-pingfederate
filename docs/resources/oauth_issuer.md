---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingfederate_oauth_issuer Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Manages an OAuth Issuer.
---

# pingfederate_oauth_issuer (Resource)

Manages an OAuth Issuer.

## Example Usage

```terraform
resource "pingfederate_oauth_issuer" "example" {
  description = "example description"
  host        = "example"
  name        = "example"
  path        = "/example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `description` (String)
- `host` (String)
- `name` (String)
- `path` (String)

### Optional

- `id` (String) The persistent, unique ID for the virtual issuer. It can be any combination of [a-zA-Z0-9._-]. This property is system-assigned if not specified.

## Import

Import is supported using the following syntax:

```shell
# "oauthIssuerId" should be the id of the OAuth Issuer to be imported
terraform import pingfederate_oauth_issuer.myOauthIssuer oauthIssuerId
```