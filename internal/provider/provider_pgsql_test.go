package provider

import "testing"

func TestProviderRegistersPostgreSQLUserAndDatabase(t *testing.T) {
	t.Parallel()

	provider := Provider()
	for _, name := range []string{"aws_paas_pgsql_database", "aws_paas_pgsql_user"} {
		if resource, ok := provider.ResourcesMap[name]; !ok || resource == nil {
			t.Fatalf("expected %s to be registered", name)
		}
	}

	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("provider schema validation failed: %s", err)
	}
}
