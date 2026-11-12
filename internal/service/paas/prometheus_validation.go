package paas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/service/paas/services"
)

func validateServiceConfiguration(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if err := validateMySQLServiceConfiguration(d); err != nil {
		return err
	}
	if err := validatePostgreSQLServiceConfiguration(d); err != nil {
		return err
	}
	if err := validateManagedExternally(d); err != nil {
		return err
	}
	if err := validateKafkaServiceConfiguration(ctx, d, meta); err != nil {
		return err
	}
	if err := validateGitLabServiceConfiguration(d); err != nil {
		return err
	}

	if err := validateAutoDiskExpansionConfiguration(ctx, d, meta); err != nil {
		return err
	}

	return validatePrometheusServiceConfiguration(ctx, d, meta)
}

// validatePrometheusServiceConfiguration rejects an unsupported cluster shape
// before the provider sends a billable CreateService request. K2 Cloud exposes
// Prometheus only as a single-node service.
func validatePrometheusServiceConfiguration(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	_, isPrometheus := d.GetOk(services.Prometheus.ServiceType())
	highAvailability := false
	if v, ok := d.GetOk("high_availability"); ok {
		highAvailability = v.(bool)
	}

	return validatePrometheusServiceConfig(isPrometheus, highAvailability)
}

func validatePrometheusServiceConfig(isPrometheus, highAvailability bool) error {
	if isPrometheus && highAvailability {
		return fmt.Errorf("prometheus supports only high_availability = false in K2 Cloud")
	}

	return nil
}
