package secretsmanager

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func DataSourceSecrets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecretsRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"description",
								"name",
								"all",
							}, false),
						},
						"values": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"include_planned_deletion": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"names": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sort_by": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"created-date",
					"last-accessed-date",
					"last-changed-date",
					"name",
				}, false),
			},
		},
	}
}

func dataSourceSecretsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).SecretsManagerConn

	input := &secretsmanager.ListSecretsInput{}

	if v, ok := d.GetOk("filter"); ok {
		var filters []*secretsmanager.Filter
		for _, f := range v.(*schema.Set).List() {
			m := f.(map[string]interface{})
			filter := &secretsmanager.Filter{
				Key: aws.String(m["key"].(string)),
			}
			for _, val := range m["values"].(*schema.Set).List() {
				filter.Values = append(filter.Values, aws.String(val.(string)))
			}
			filters = append(filters, filter)
		}
		input.Filters = filters
	}

	if v, ok := d.GetOk("include_planned_deletion"); ok {
		input.IncludePlannedDeletion = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("sort_by"); ok {
		input.SortBy = aws.String(v.(string))
	}

	var results []*secretsmanager.SecretListEntry

	err := conn.ListSecretsPages(input, func(page *secretsmanager.ListSecretsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, secretListEntry := range page.SecretList {
			if secretListEntry == nil {
				continue
			}

			results = append(results, secretListEntry)
		}

		return !lastPage
	})

	if err != nil {
		return fmt.Errorf("listing Secrets Manager Secrets: %w", err)
	}

	var names []string

	for _, r := range results {
		names = append(names, aws.StringValue(r.Name))
	}

	d.SetId(meta.(*conns.AWSClient).Region)
	d.Set("names", names)

	return nil
}
