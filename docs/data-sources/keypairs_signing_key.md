---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingfederate_keypairs_signing_key Data Source - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Data source to retrieve a signing key pair.
---

# pingfederate_keypairs_signing_key (Data Source)

Data source to retrieve a signing key pair.

## Example Usage

```terraform
data "pingfederate_keypairs_signing_key" "signingKey" {
  key_id = "signingkey"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `key_id` (String) The persistent, unique ID for the certificate.

### Read-Only

- `crypto_provider` (String) Cryptographic Provider. This is only applicable if Hybrid HSM mode is true. Supported values are `LOCAL` and `HSM`.
- `expires` (String) The end date up until which the item is valid, in ISO 8601 format (UTC)
- `id` (String) ID of this resource.
- `issuer_dn` (String) The issuer's distinguished name
- `key_algorithm` (String) The public key algorithm.
- `key_size` (Number) The public key size, in bits.
- `rotation_settings` (Attributes) The local identity profile data store configuration. (see [below for nested schema](#nestedatt--rotation_settings))
- `serial_number` (String) The serial number assigned by the CA
- `sha1_fingerprint` (String) SHA-1 fingerprint in Hex encoding
- `sha256_fingerprint` (String) SHA-256 fingerprint in Hex encoding
- `signature_algorithm` (String) The signature algorithm.
- `status` (String) Status of the item.
- `subject_alternative_names` (Set of String) The subject alternative names (SAN).
- `subject_dn` (String) The subject's distinguished name
- `valid_from` (String) The start date from which the item is valid, in ISO 8601 format (UTC).
- `version` (Number) The X.509 version to which the item conforms

<a id="nestedatt--rotation_settings"></a>
### Nested Schema for `rotation_settings`

Read-Only:

- `activation_buffer_days` (Number) Buffer days before key pair expiration for activation of the new key pair.
- `creation_buffer_days` (Number) Buffer days before key pair expiration for creation of a new key pair.
- `id` (String) The base DN to search from. If not specified, the search will start at the LDAP's root.
- `key_algorithm` (String) Key algorithm to be used while creating a new key pair. If this property is unset, the key algorithm of the original key pair will be used. Supported algorithms are available through the /keyPairs/keyAlgorithms endpoint.
- `key_size` (Number) Key size, in bits. If this property is unset, the key size of the original key pair will be used. Supported key sizes are available through the /keyPairs/keyAlgorithms endpoint.
- `signature_algorithm` (String) Required if the original key pair used SHA1 algorithm. If this property is unset, the default signature algorithm of the original key pair will be used. Supported signature algorithms are available through the /keyPairs/keyAlgorithms endpoint.
- `valid_days` (Number) Valid days for the new key pair to be created. If this property is unset, the validity days of the original key pair will be used.