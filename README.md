terraform-provider-st-ucloud
============================

This Terraform custom provider is designed for own use case scenario.

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

- **st-ucloud_alidns_gtm_instance**

  The original reason to write this resource is official UCloud Terraform
  provider does not support creating GTM (Global Traffic Manager) instance on
  UCloud international account using Terraform. As we developing on the
  resource, we added few more features which are useful in our use case, which
  includes:

    - setting the renewal status to *NotRenewal* when destroying the resource.

    - allowing changing of renewal period and status without recreating the GTM instsance.

- **st-ucloud_alidns_record_weight**

  Official UCloud Terraform provider does not have the resource to modify DNS
  records weight.

- **st-ucloud_ram_user_group_attachment**

  The official UCloud Terraform provider's resource
  [*ucloud_ram_group_membership*](https://registry.terraform.io/providers/aliyun/ucloud/latest/docs/resources/ram_group_membership)
  will remove all other attached users for the target group, which may cause a
  problem where Terraform may delete those users attached outside from Terraform.

- **st-ucloud_ram_policy**

  This resource is designed to handle policy content that exceeds the limit of 6144 characters.
  It provides functionality to create policies by splitting the content into smaller segments that fit within the limit,
  enabling the management and combination of these segments to form the complete policy. Finally, the policy will be attached to the relevant user.

- **st-ucloud_cms_alarm_rule**

  The official UCloud Terraform provider's resource
  [*ucloud_cms_group_metric_rule*](https://registry.terraform.io/providers/aliyun/ucloud/latest/docs/resources/cms_group_metric_rule)
  does not support adding alarm rules into application groups based on expression-based creation.

  For namespaces and metric inputs, please refer to: [*Alicloud Alarm Metric List*](https://cms.console.aliyun.com/metric-meta)

**st-ucloud_alidns_instance**

   The official UCloud Terraform provider's resource
   [*ucloud_alidns_instance*](https://registry.terraform.io/providers/aliyun/ucloud/latest/docs/resources/alidns_instance)
   will destroy and create a new instance everytime when upgrading or downgrading.

**st-ucloud_alidns_domain_attachment**

   The official UCloud Terraform provider's resource
   [*ucloud_dns_domain_attachment*](https://registry.terraform.io/providers/aliyun/ucloud/latest/docs/resources/dns_domain_attachment)
   accept input of a list of domains. There will be an issue when upgrading a batch of domains when the existing attachment
   is more than 100 domains. The official resources will first destroy all the domains and re-add the new one together with
   the existing one. The resources will hit timeout during adding of new domains and make some of the domains not re-add back.

- **st-ucloud_cms_system_event_contact_group_attachment**

  The official UCloud Terraform provider's resource [*ucloud_cms_event_rule*](https://registry.terraform.io/providers/aliyun/ucloud/latest/docs/resources/cms_event_rule) does not bind the created system event rule to the contact group itself.
  This may cause system event rule could create as usual but with an empty target contact group.


- **st-ucloud_ddoscoo_webconfig_ssl_attachment**

  This resource is designed to associate a SSL certificate to a website/domain before being added
  into Anti-DDoS as UCloud Terraform Provider does not support the SSL binding operation.

### Data Sources

- **st-ucloud_ddoscoo_domain_resources**

  Official UCloud Terraform provider does not support querying the CNAME of
  AntiDDoS domain resources through
  [*ucloud_ddoscoo_domain_resources*](https://registry.terraform.io/providers/aliyun/ucloud/latest/docs/data-sources/ddoscoo_domain_resources).

- **st-ucloud_ddoscoo_instances**

  Official UCloud Terraform provider does not support querying the EIP of
  AntiDDoS instance resources through
  [*ucloud_ddoscoo_instances*](https://registry.terraform.io/providers/aliyun/ucloud/latest/docs/data-sources/ddoscoo_instances).

- **st-ucloud_cdn_domain**

  Official UCloud Terraform provider does not have the data source to query
  the CNAME of CDN domain.

- **st-ucloud_slb_load_balancers**

  The tags parameter of UCloud API
  [*DescribeLoadBalancers*](https://www.alibabacloud.com/help/en/server-load-balancer/latest/describeloadbalancers)
  will return all load balancers when any one of the tags are matched. This may
  be a problem when the user wants to match exactly all given tags, therefore
  this data source will filter once more after listing the load balancers
  from UCloud API to match all the given tags.

  The example bahaviors of UCloud API *DescribeLoadBalancers*:

  | Load Balancer   | Tags                                            | Given tags: { "location": "office" "env": "test" }          |
  |-----------------|-------------------------------------------------|-------------------------------------------------------------|
  | load-balancer-A | { "location": "office" "env" : "test" }         | Matched (work as expected)                                  |
  | load-balancer-B | { "location": "office" "env" : "prod" }         | Matched (should not be matched as the `env` is prod)          |

References
----------

- Website: https://www.terraform.io
- Terraform Plugin Framework: https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework
- UCloud official Terraform provider: https://github.com/aliyun/terraform-provider-ucloud
