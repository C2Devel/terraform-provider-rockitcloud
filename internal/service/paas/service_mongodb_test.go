package paas

import (
	"context"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	sdkpaas "github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/service/paas/services"
)

func TestResourceServiceMongoDBBackupUpdatesInPlace(t *testing.T) {
	t.Parallel()
	config := func(enabled bool) map[string]interface{} {
		return map[string]interface{}{
			"name": "tf-mongodb-backup", "instance_type": "m5.large", "ssh_key_name": "test-key",
			"root_volume": []interface{}{map[string]interface{}{"type": "gp2", "size": 32}},
			"data_volume": []interface{}{map[string]interface{}{"type": "gp2", "size": 32}},
			"subnet_ids":  []interface{}{"subnet-12345678"}, "security_group_ids": []interface{}{"sg-12345678"},
			"backup_settings": []interface{}{map[string]interface{}{
				"enabled": enabled, "bucket_name": "mongodb-backups", "user_login": "backup@example.com",
				"start_time": "03:00", "expiration_days": 7,
			}},
			"mongodb": []interface{}{map[string]interface{}{
				"version": "7.0", "database": []interface{}{map[string]interface{}{"name": "application", "backup_enabled": enabled}},
			}},
		}
	}
	for _, enabled := range []bool{true, false} {
		t.Run(strconv.FormatBool(enabled), func(t *testing.T) {
			resource := ResourceService()
			desired := terraform.NewResourceConfigRaw(config(enabled))
			if diagnostics := resource.Validate(desired); diagnostics.HasError() {
				t.Fatalf("MongoDB backup configuration rejected: %#v", diagnostics)
			}
			previous := schema.TestResourceDataRaw(t, resource.Schema, config(!enabled))
			previous.SetId("fm-cluster-12345678")
			diff, err := resource.Diff(context.Background(), previous.State(), desired, nil)
			if err != nil || diff == nil || diff.RequiresNew() {
				t.Fatalf("backup settings must update in place: diff=%#v err=%v", diff, err)
			}
			for _, key := range []string{"backup_settings.0.enabled", "mongodb.0.database.0.backup_enabled"} {
				change := diff.Attributes[key]
				if change == nil || change.New != strconv.FormatBool(enabled) {
					t.Errorf("expected %s to become %t, got %#v", key, enabled, change)
				}
			}
		})
	}
}

func TestMongoDBDatabaseBackupRoundTrip(t *testing.T) {
	t.Parallel()
	dataSource := DataSourceService()
	databaseSchema := dataSource.Schema["mongodb"].Elem.(*schema.Resource).Schema["database"].Elem.(*schema.Resource).Schema
	if field := databaseSchema["backup_enabled"]; field == nil || !field.Computed {
		t.Fatal("MongoDB data source must expose computed database.backup_enabled")
	}
	for _, enabled := range []bool{true, false} {
		t.Run(strconv.FormatBool(enabled), func(t *testing.T) {
			request := services.MongoDB.ExpandDatabase(map[string]interface{}{"name": "application", "backup_enabled": enabled})
			if request.BackupEnabled == nil || aws.BoolValue(request.BackupEnabled) != enabled {
				t.Fatalf("database backup flag lost in API request: %#v", request)
			}
			database := services.MongoDB.FlattenDatabase(&sdkpaas.DatabaseResponse{Name: request.Name, BackupEnabled: request.BackupEnabled})
			data := schema.TestResourceDataRaw(t, dataSource.Schema, nil)
			if err := data.Set("mongodb", []interface{}{map[string]interface{}{"database": []interface{}{database}}}); err != nil {
				t.Fatal(err)
			}
			databases := data.Get("mongodb.0.database").(*schema.Set).List()
			if got := databases[0].(map[string]interface{})["backup_enabled"]; got != enabled {
				t.Fatalf("database backup flag lost in Terraform state: got %#v, want %t", got, enabled)
			}
		})
	}
}
