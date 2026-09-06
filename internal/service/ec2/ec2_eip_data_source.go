package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func DataSourceEIP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceEIPRead,

		Schema: map[string]*schema.Schema{
			"association_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter": CustomFiltersSchema(),
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_interface_owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"public_dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ipv4_pool": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"carrier_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_owned_ipv4_pool": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_owned_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourceEIPRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	req := &ec2.DescribeAddressesInput{}

	if v, ok := d.GetOk("id"); ok {
		req.AllocationIds = []*string{aws.String(v.(string))}
	}

	if v, ok := d.GetOk("public_ip"); ok {
		req.PublicIps = []*string{aws.String(v.(string))}
	}

	req.Filters = []*ec2.Filter{}

	req.Filters = append(req.Filters, BuildCustomFilterList(
		d.Get("filter").(*schema.Set),
	)...)

	if tags, tagsOk := d.GetOk("tags"); tagsOk {
		req.Filters = append(req.Filters, BuildTagFilterList(
			Tags(tftags.New(tags.(map[string]interface{}))),
		)...)
	}

	if len(req.Filters) == 0 {
		// Don't send an empty filters list; the EC2 API won't accept it.
		req.Filters = nil
	}

	resp, err := conn.DescribeAddresses(req)
	if err != nil {
		return fmt.Errorf("error describing EC2 Address: %w", err)
	}
	if resp == nil || len(resp.Addresses) == 0 {
		return fmt.Errorf("no matching Elastic IP found")
	}
	if len(resp.Addresses) > 1 {
		return fmt.Errorf("multiple Elastic IPs matched; use additional constraints to reduce matches to a single Elastic IP")
	}

	eip := resp.Addresses[0]

	d.SetId(aws.StringValue(eip.AllocationId))

	d.Set("association_id", eip.AssociationId)
	d.Set("domain", eip.Domain)
	d.Set("instance_id", eip.InstanceId)
	d.Set("network_interface_id", eip.NetworkInterfaceId)
	d.Set("network_interface_owner_id", eip.NetworkInterfaceOwnerId)

	d.Set("private_ip", eip.PrivateIpAddress)
	d.Set("public_ip", eip.PublicIp)

	privateDNS, publicDNS, err := eipDNSNames(conn, eip.NetworkInterfaceId)
	if err != nil {
		return fmt.Errorf("error reading EC2 Network Interface (%s) for EIP (%s): %w", aws.StringValue(eip.NetworkInterfaceId), d.Id(), err)
	}
	d.Set("private_dns", privateDNS)
	d.Set("public_dns", publicDNS)
	d.Set("public_ipv4_pool", eip.PublicIpv4Pool)
	d.Set("carrier_ip", eip.CarrierIp)
	d.Set("customer_owned_ipv4_pool", eip.CustomerOwnedIpv4Pool)
	d.Set("customer_owned_ip", eip.CustomerOwnedIp)

	if err := d.Set("tags", KeyValueTags(eip.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	return nil
}
