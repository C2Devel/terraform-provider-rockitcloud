package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func DataSourceVPNGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPNGatewayRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attached_vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": CustomFiltersSchema(),
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceVPNGatewayRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	input := &ec2.DescribeVpnGatewaysInput{}

	if id, ok := d.GetOk("id"); ok {
		input.VpnGatewayIds = aws.StringSlice([]string{id.(string)})
	}

	input.Filters = BuildAttributeFilterList(
		map[string]string{
			"state":             d.Get("state").(string),
			"availability-zone": d.Get("availability_zone").(string),
		},
	)
	if id, ok := d.GetOk("attached_vpc_id"); ok {
		input.Filters = append(input.Filters, BuildAttributeFilterList(
			map[string]string{
				"attachment.state":  "attached",
				"attachment.vpc-id": id.(string),
			},
		)...)
	}
	input.Filters = append(input.Filters, BuildCustomFilterList(
		d.Get("filter").(*schema.Set),
	)...)
	if len(input.Filters) == 0 {
		// Don't send an empty filters list; the EC2 API won't accept it.
		input.Filters = nil
	}

	vgw, err := FindVPNGateway(conn, input)

	if err != nil {
		return tfresource.SingularDataSourceFindError("EC2 VPN Gateway", err)
	}

	d.SetId(aws.StringValue(vgw.VpnGatewayId))

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   ec2.ServiceName,
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("vpn-gateway/%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	for _, attachment := range vgw.VpcAttachments {
		if aws.StringValue(attachment.State) == ec2.AttachmentStatusAttached {
			d.Set("attached_vpc_id", attachment.VpcId)
			break
		}
	}
	d.Set("state", vgw.State)

	return nil
}
