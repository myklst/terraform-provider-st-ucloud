terraform-provider-st-ucloud
============================

This Terraform custom provider is designed for own use case scenario.123

Supported Versions
------------------

| Terraform version | minimum provider version |maxmimum provider version
| ---- | ---- | ----|
| >= 1.3.x	| 0.1.1	| latest |

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 1.3.x
-	[Go](https://golang.org/doc/install) 1.19 (to build the provider plugin)

Local Installation
------------------

1. Run make file `make install-local-custom-provider` to install the provider under ~/.terraform.d/plugins.

2. The provider source should be change to the path that configured in the *Makefile*:

    ```
    terraform {
      required_providers {
        st-ucloud = {
          source = "example.local/myklst/st-ucloud"
        }
      }
    }
    ```

Why Custom Provider
-------------------

This custom provider exists due to UCloud doesn't support Terraform officially.

### Resources

- **st-ucloud_cdn_domain**

  Configure acl, origin, cache control of a domain.

- **st-ucloud_ssl_certificate**

  Manage ssl certificates.

### Data Sources

- **st-ucloud_ssl_certificate**

  Query all ssl certificates in UCloud.

References
----------

- Website: https://www.terraform.io
- Terraform Plugin Framework: https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework
- UCloud official Terraform provider: https://github.com/ucloud/terraform-provider-ucloud
