---
subcategory: "Route 53"
layout: "aws"
page_title: "aws_route53_record"
description: |-
  Manages a Route53 record.
---

[fqdn]: https://en.wikipedia.org/wiki/Fully_qualified_domain_name

# Resource: aws_route53_record

Manages a Route53 record.

## Example usage

### Simple routing policy

```terraform
resource "aws_route53_record" "www" {
  zone_id = aws_route53_zone.primary.zone_id
  name    = "www.example.com"
  type    = "A"
  ttl     = "300"
  records = [aws_eip.lb.public_ip]
}
```

### NS record management

When creating Route 53 zones, the `NS` records for the zone are automatically created.
Enabling the `allow_overwrite` argument will allow managing these records in a single Terraform run without the requirement for `terraform import`.

```terraform
resource "aws_route53_zone" "example" {
  name = "test.example.com"
}

resource "aws_route53_record" "example" {
  allow_overwrite = true
  name            = "test.example.com"
  ttl             = 172800
  type            = "NS"
  zone_id         = aws_route53_zone.example.zone_id

  records = [
    aws_route53_zone.example.name_servers[0],
    aws_route53_zone.example.name_servers[1],
  ]
}
```

## Argument reference

The following arguments are required:

* `name` - (Required, Forces new resource, String) The name of the record.
* `type` - (Required, Editable, String) The type of the record.
    * _Valid values:_ `A`, `AAAA`, `CNAME`, `MX`, `NS`, `PTR`, `SRV` or `TXT`
* `zone_id` - (Required, Forces new resource, String) The ID of the hosted zone to contain this record.

The following arguments are optional:

* `allow_overwrite` - (Optional, Editable, Boolean) Indicates whether to allow creation of this record in Terraform to overwrite an existing record, if any.
    * _Default value:_ `false`
* `records` - (Optional, Editable, Set of strings) A list of DNS records.
  To specify a single record value longer than 255 characters such as a TXT record for DKIM, add `\" \"` inside the Terraform configuration string to split characters into multiple text strings (for example, `"first255characters\" \"next255characters"`).
* `ttl` - (Optional, Editable, Integer) The TTL of the record.

    ~> **Note** This does not affect the ability to update the record in Terraform and does not prevent other resources within Terraform or manual Route 53 changes outside Terraform from overwriting this record.

    !> **Warning** This configuration is not recommended for most environments.

* `failover_routing_policy` - (Optional, Editable, [Block](#failover_routing_policy)) A failover routing policy block.
    * _Constraints:_ Conflicts with `geolocation_routing_policy`, `latency_routing_policy`, `weighted_routing_policy` and `multivalue_answer_routing_policy` arguments.
* `geolocation_routing_policy` - (Optional, Editable, [Block](#geolocation_routing_policy)) A geolocation routing policy block.
    * _Constraints:_ Conflicts with `failover_routing_policy`, `latency_routing_policy`, `weighted_routing_policy` and `multivalue_answer_routing_policy` arguments.
* `latency_routing_policy` - (Optional, Editable, [Block](#latency_routing_policy)) A latency routing policy block.
    * _Constraints:_ Conflicts with `failover_routing_policy`, `geolocation_routing_policy`, `weighted_routing_policy` and `multivalue_answer_routing_policy` arguments.
* `multivalue_answer_routing_policy` - (Optional, Editable, Boolean) Indicates whether to route traffic to all the records in the group, allowing multiple responses to a DNS query.
    * _Default value:_ `false`
    * _Constraints:_ Conflicts with `failover_routing_policy`, `geolocation_routing_policy`, `latency_routing_policy` and `weighted_routing_policy` arguments.
* `set_identifier` - (Optional, Editable, String) Unique identifier to differentiate records with the same `name` and `type`.
* `weighted_routing_policy` - (Optional, Editable, [Block](#weighted_routing_policy)) A weighted routing policy block.
    * _Constraints:_ Conflicts with `failover_routing_policy`, `geolocation_routing_policy`, `latency_routing_policy` and `multivalue_answer_routing_policy` arguments.

### failover_routing_policy

The following arguments are required:

* `type` - (Required, Editable, String) The failover type.
    * _Valid values:_ `PRIMARY` or `SECONDARY`.

### geolocation_routing_policy

The following arguments are optional:

* `continent` - (Optional, Editable, String) A two-letter continent code.
    * _Example:_ `EU`, `NA`, `AS`
* `country` - (Optional, Editable, String) A two-letter country code.
    * _Example:_ `US`, `DE`, `RU`
* `subdivision` - (Optional, Editable, String) A subdivision code.
    * _Example:_ `US-CA`, `US-NY`

### latency_routing_policy

The following arguments are required:

* `region` - (Required, Editable, String) The AWS region for the latency-based routing.

### weighted_routing_policy

The following arguments are required:

* `weight` - (Required, Editable, Integer) The weight for the weighted routing policy.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `fqdn` - (String) [FQDN][fqdn] built using the zone domain and `name`.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`alias`, `health_check_id`.

## Timeouts

Timeouts usage for records is not currently supported.

## Import

Route53 records can be imported using the record ID, which consists of the zone identifier, record name, and record type separated by underscores (`_`), for example:

```
$ terraform import aws_route53_record.myrecord z-xxxxxxxx_dev.example.com_NS
```
