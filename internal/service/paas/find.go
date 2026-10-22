package paas

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/paas"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func FindServiceByID(conn *paas.PaaS, id string) (*paas.Service, error) {
	return findServiceByID(context.Background(), conn, id)
}

func findServiceByID(ctx context.Context, conn *paas.PaaS, id string) (*paas.Service, error) {
	input := &paas.DescribeServiceInput{
		ServiceId: aws.String(id),
	}

	output, err := conn.DescribeServiceWithContext(
		ctx,
		input,
		request.WithLogLevel(aws.LogOff),
	)

	if tfawserr.ErrCodeEquals(err, ServiceNotFoundCode) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.Service == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Service, nil
}

func FindPostgreSQLUserByID(conn *paas.PaaS, serviceID, id string) (*paas.UserResponse, error) {
	return findPostgreSQLUserByID(context.Background(), conn, serviceID, id)
}

// findPostgreSQLUserByID checks that the user exists with ListUsers, because the
// error code returned by DescribeUser for a missing user is not documented,
// and then reads the user parameters with DescribeUser.
func findPostgreSQLUserByID(ctx context.Context, conn *paas.PaaS, serviceID, id string) (*paas.UserResponse, error) {
	listInput := &paas.ListUsersInput{
		ServiceId: aws.String(serviceID),
	}

	listOutput, err := conn.ListUsersWithContext(ctx, listInput)

	if tfawserr.ErrCodeEquals(err, ServiceNotFoundCode) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: listInput,
		}
	}

	if err != nil {
		return nil, err
	}

	if listOutput == nil {
		return nil, tfresource.NewEmptyResultError(listInput)
	}

	found := false
	for _, user := range listOutput.Users {
		if user != nil && aws.StringValue(user.Id) == id {
			found = true
			break
		}
	}

	if !found {
		return nil, &resource.NotFoundError{
			LastRequest: listInput,
		}
	}

	input := &paas.DescribeUserInput{
		ServiceId: aws.String(serviceID),
		UserId:    aws.String(id),
	}

	output, err := conn.DescribeUserWithContext(
		ctx,
		input,
		request.WithLogLevel(aws.LogOff),
	)

	if tfawserr.ErrCodeEquals(err, ServiceNotFoundCode, ResourceDoesNotExistCode) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.User == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.User, nil
}

func FindPostgreSQLDatabaseByID(conn *paas.PaaS, serviceID, id string) (*paas.DatabaseResponse, error) {
	return findPostgreSQLDatabaseByID(context.Background(), conn, serviceID, id)
}

// findPostgreSQLDatabaseByID checks that the database exists with ListDatabases, because the
// error code returned by DescribeDatabase for a missing database is not documented,
// and then reads the database parameters with DescribeDatabase.
func findPostgreSQLDatabaseByID(ctx context.Context, conn *paas.PaaS, serviceID, id string) (*paas.DatabaseResponse, error) {
	listInput := &paas.ListDatabasesInput{
		ServiceId: aws.String(serviceID),
	}

	listOutput, err := conn.ListDatabasesWithContext(ctx, listInput)

	if tfawserr.ErrCodeEquals(err, ServiceNotFoundCode) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: listInput,
		}
	}

	if err != nil {
		return nil, err
	}

	if listOutput == nil {
		return nil, tfresource.NewEmptyResultError(listInput)
	}

	found := false
	for _, database := range listOutput.Databases {
		if database != nil && aws.StringValue(database.Id) == id {
			found = true
			break
		}
	}

	if !found {
		return nil, &resource.NotFoundError{
			LastRequest: listInput,
		}
	}

	input := &paas.DescribeDatabaseInput{
		DatabaseId: aws.String(id),
		ServiceId:  aws.String(serviceID),
	}

	output, err := conn.DescribeDatabaseWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, ServiceNotFoundCode, ResourceDoesNotExistCode) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.Database == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Database, nil
}

func FindClickHouseProfileByID(conn *paas.PaaS, serviceID, id string) (*paas.ClickHouseProfile, error) {
	return findClickHouseProfileByID(context.Background(), conn, serviceID, id)
}

func findClickHouseProfileByID(ctx context.Context, conn *paas.PaaS, serviceID, id string) (*paas.ClickHouseProfile, error) {
	input := &paas.ListClickhouseProfilesInput{
		ServiceId: aws.String(serviceID),
	}

	output, err := conn.ListClickhouseProfilesWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, ServiceNotFoundCode) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	for _, profile := range output.Profiles {
		if profile != nil && aws.StringValue(profile.Id) == id {
			return profile, nil
		}
	}

	return nil, &resource.NotFoundError{
		LastRequest: input,
	}
}

func FindClickHouseQuotaByID(conn *paas.PaaS, serviceID, id string) (*paas.ClickHouseQuota, error) {
	return findClickHouseQuotaByID(context.Background(), conn, serviceID, id)
}

func findClickHouseQuotaByID(ctx context.Context, conn *paas.PaaS, serviceID, id string) (*paas.ClickHouseQuota, error) {
	input := &paas.ListClickhouseQuotasInput{
		ServiceId: aws.String(serviceID),
	}

	output, err := conn.ListClickhouseQuotasWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, ServiceNotFoundCode) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	for _, quota := range output.Quotas {
		if quota != nil && aws.StringValue(quota.Id) == id {
			return quota, nil
		}
	}

	return nil, &resource.NotFoundError{
		LastRequest: input,
	}
}

func FindClickHouseShardByIndex(conn *paas.PaaS, serviceID string, index int64) (*paas.ClickHouseShard, error) {
	return findClickHouseShardByIndex(context.Background(), conn, serviceID, index)
}

func findClickHouseShardByIndex(ctx context.Context, conn *paas.PaaS, serviceID string, index int64) (*paas.ClickHouseShard, error) {
	input := &paas.ListClickhouseShardsInput{
		ServiceId: aws.String(serviceID),
	}

	output, err := conn.ListClickhouseShardsWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, ServiceNotFoundCode) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	for _, shard := range output.Shards {
		if shard != nil && shard.Index != nil && aws.Int64Value(shard.Index) == index {
			return shard, nil
		}
	}

	return nil, &resource.NotFoundError{
		LastRequest: input,
	}
}

func FindBackupUsers(conn *paas.PaaS) ([]*paas.BackupUser, error) {
	input := &paas.ListBackupUsersInput{}

	output, err := conn.ListBackupUsers(input)

	if err != nil {
		return nil, err
	}

	if output == nil || output.Users == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Users, nil
}

func FindBackupById(conn *paas.PaaS, id string) (*paas.Backup, error) {
	input := &paas.DescribeBackupInput{
		BackupId: aws.String(id),
	}

	output, err := conn.DescribeBackup(input)

	if err != nil {
		return nil, err
	}

	if output == nil || output.Backup == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Backup, nil
}

func FindBackups(conn *paas.PaaS, serviceClass, serviceId, serviceType string) ([]*paas.Backup, error) {
	input := &paas.ListBackupsInput{}

	if serviceClass != "" {
		input.ServiceClass = aws.String(serviceClass)
	}

	if serviceId != "" {
		input.ServiceId = aws.String(serviceId)
	}

	if serviceType != "" {
		input.ServiceType = aws.String(serviceType)
	}

	output, err := conn.ListBackups(input)

	if err != nil {
		return nil, err
	}

	if output == nil || output.Backups == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Backups, nil
}

func FindPrometheusNotificationChannelByID(ctx context.Context, conn *paas.PaaS, serviceID, channelID string) (*paas.NotificationChannel, error) {
	input := &paas.ListNotificationChannelsInput{
		ServiceId: aws.String(serviceID),
	}
	output, err := conn.ListNotificationChannelsWithContext(ctx, input)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, ServiceNotFoundCode) || isPrometheusNotFoundError(err) {
			return nil, newPrometheusNotFoundError(input, err)
		}
		return nil, err
	}
	if output == nil || output.NotificationChannels == nil {
		return nil, newPrometheusNotFoundError(input, fmt.Errorf("notification channel %q not found", channelID))
	}

	for _, item := range output.NotificationChannels {
		if item != nil && aws.StringValue(item.Id) == channelID {
			return item, nil
		}
	}

	return nil, newPrometheusNotFoundError(input, fmt.Errorf("notification channel %q not found", channelID))
}

func FindPrometheusRouteByID(ctx context.Context, conn *paas.PaaS, serviceID, routeID string) (*paas.PrometheusRoute, error) {
	input := &paas.ListPrometheusRoutesInput{
		ServiceId: aws.String(serviceID),
	}
	output, err := conn.ListPrometheusRoutesWithContext(ctx, input)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, ServiceNotFoundCode) || isPrometheusNotFoundError(err) {
			return nil, newPrometheusNotFoundError(input, err)
		}
		return nil, err
	}
	if output == nil || output.PrometheusRoutes == nil {
		return nil, newPrometheusNotFoundError(input, fmt.Errorf("Prometheus route %q not found", routeID))
	}

	for _, item := range output.PrometheusRoutes {
		if item != nil && aws.StringValue(item.Id) == routeID {
			return item, nil
		}
	}

	return nil, newPrometheusNotFoundError(input, fmt.Errorf("Prometheus route %q not found", routeID))
}

func FindPrometheusScrapeJobByID(ctx context.Context, conn *paas.PaaS, serviceID, jobID string) (*paas.PrometheusScrapeJob, error) {
	input := &paas.ListPrometheusScrapeJobsInput{
		ServiceId: aws.String(serviceID),
	}
	output, err := conn.ListPrometheusScrapeJobsWithContext(ctx, input)
	if err != nil {
		if tfawserr.ErrCodeEquals(err, ServiceNotFoundCode) || isPrometheusNotFoundError(err) {
			return nil, newPrometheusNotFoundError(input, err)
		}
		return nil, err
	}
	if output == nil || output.PrometheusScrapeJobs == nil {
		return nil, newPrometheusNotFoundError(input, fmt.Errorf("Prometheus scrape job %q not found", jobID))
	}

	for _, item := range output.PrometheusScrapeJobs {
		if item != nil && aws.StringValue(item.Id) == jobID {
			return item, nil
		}
	}

	return nil, newPrometheusNotFoundError(input, fmt.Errorf("Prometheus scrape job %q not found", jobID))
}

func newPrometheusNotFoundError(input interface{}, err error) error {
	return &resource.NotFoundError{
		LastError:   err,
		LastRequest: input,
	}
}

func isPrometheusNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	var requestFailure awserr.RequestFailure
	if errors.As(err, &requestFailure) && requestFailure.StatusCode() == http.StatusNotFound {
		return true
	}

	var awsErr awserr.Error
	if errors.As(err, &awsErr) {
		return strings.Contains(strings.ToLower(awsErr.Code()), "notfound")
	}

	return false
}
