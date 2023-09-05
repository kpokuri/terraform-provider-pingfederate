---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingfederate_authentication_policy_contracts Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Manages a AuthenticationPolicyContracts.
---

# pingfederate_authentication_policy_contracts (Resource)

Manages a AuthenticationPolicyContracts.

## Example Usage

```terraform
resource "pingfederate_authentication_policy_contracts" "authenticationPolicyContractsExample" {
  core_attributes     = [{ name = "subject" }]
  extended_attributes = [{ name = "extended_attribute" }, { name = "extended_attribute2" }]
  name                = "example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `core_attributes` (Attributes Set) A list of read-only assertion attributes (for example, subject) that are automatically populated by PingFederate. (see [below for nested schema](#nestedatt--core_attributes))
- `extended_attributes` (Attributes Set) A list of additional attributes as needed. (see [below for nested schema](#nestedatt--extended_attributes))
- `name` (String) The Authentication Policy Contract Name. Name is unique.

### Optional

- `id` (String) The persistent, unique ID for the authentication policy contract. It can be any combination of [a-zA-Z0-9._-]. This property is system-assigned if not specified.

<a id="nestedatt--core_attributes"></a>
### Nested Schema for `core_attributes`

Required:

- `name` (String)


<a id="nestedatt--extended_attributes"></a>
### Nested Schema for `extended_attributes`

Required:

- `name` (String)