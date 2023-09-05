---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pingfederate_key_pairs_signing_import Resource - terraform-provider-pingfederate"
subcategory: ""
description: |-
  Manages a KeyPairsSigningImport.
---

# pingfederate_key_pairs_signing_import (Resource)

Manages a KeyPairsSigningImport.

## Example Usage

```terraform
# WARNING! You will need to secure your state file properly when using this resource! #
# Please refer to the link below on how to best store state files and data within. #
# https://developer.hashicorp.com/terraform/plugin/best-practices/sensitive-state #
resource "pingfederate_key_pairs_signing_import" "keyPairsSigningImportExample" {
  file_data = "example"
  format    = "PKCS12"
  # This value will be stored into your state file 
  password = "example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `file_data` (String) Base-64 encoded PKCS12 or PEM file data. In the case of PEM, the raw (non-base-64) data is also accepted. In BCFIPS mode, only PEM with PBES2 and AES or Triple DES encryption is accepted and 128-bit salt is required.
- `format` (String) Key pair file format. If specified, this field will control what file format is expected, otherwise the format will be auto-detected. In BCFIPS mode, only PEM is supported. (PKCS12, PEM)
- `password` (String, Sensitive) Password for the file. In BCFIPS mode, the password must be at least 14 characters.

### Optional

- `crypto_provider` (String) Cryptographic Provider. This is only applicable if Hybrid HSM mode is true. (LOCAL, HSM)
- `id` (String) The persistent, unique ID for the certificate. It can be any combination of [a-z0-9._-]. This property is system-assigned if not specified.