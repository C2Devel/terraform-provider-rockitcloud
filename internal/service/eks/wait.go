package eks

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const (
	addonCreatedTimeout = 20 * time.Minute
	addonUpdatedTimeout = 20 * time.Minute
	addonDeletedTimeout = 40 * time.Minute
)

func waitAddonCreated(ctx context.Context, conn *eks.EKS, clusterName, addonName string) (*eks.Addon, error) {
	stateConf := resource.StateChangeConf{
		Pending: []string{eks.AddonStatusCreating, eks.AddonStatusDegraded},
		Target:  []string{eks.AddonStatusActive},
		Refresh: statusAddon(ctx, conn, clusterName, addonName),
		Timeout: addonCreatedTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*eks.Addon); ok {
		if status, health := aws.StringValue(output.Status), output.Health; status == eks.AddonStatusCreateFailed && health != nil {
			tfresource.SetLastError(err, AddonIssuesError(health.Issues))
		}

		return output, err
	}

	return nil, err
}

func waitAddonDeleted(ctx context.Context, conn *eks.EKS, clusterName, addonName string) (*eks.Addon, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{eks.AddonStatusActive, eks.AddonStatusDeleting},
		Target:  []string{},
		Refresh: statusAddon(ctx, conn, clusterName, addonName),
		Timeout: addonDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*eks.Addon); ok {
		if status, health := aws.StringValue(output.Status), output.Health; status == eks.AddonStatusDeleteFailed && health != nil {
			tfresource.SetLastError(err, AddonIssuesError(health.Issues))
		}

		return output, err
	}

	return nil, err
}

func waitAddonUpdateSuccessful(ctx context.Context, conn *eks.EKS, clusterName, addonName, id string) (*eks.Update, error) {
	stateConf := resource.StateChangeConf{
		Pending: []string{eks.UpdateStatusInProgress},
		Target:  []string{eks.UpdateStatusSuccessful},
		Refresh: statusAddonUpdate(ctx, conn, clusterName, addonName, id),
		Timeout: addonUpdatedTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*eks.Update); ok {
		if status := aws.StringValue(output.Status); status == eks.UpdateStatusCancelled || status == eks.UpdateStatusFailed {
			tfresource.SetLastError(err, ErrorDetailsError(output.Errors))
		}

		return output, err
	}

	return nil, err
}

func waitClusterCreated(conn *eks.EKS, name string, timeout time.Duration) (*eks.Cluster, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{eks.ClusterStatusCreating, eks.ClusterStatusPending, eks.ClusterStatusClaimed, eks.ClusterStatusProvisioning},
		Target:  []string{eks.ClusterStatusReady},
		Refresh: statusCluster(conn, name),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*eks.Cluster); ok {
		if status, health := aws.StringValue(output.Status), output.Health; status == eks.ClusterStatusFailed && health != nil {
			tfresource.SetLastError(err, ClusterIssuesError(health.Issues))
		}

		return output, err
	}

	return nil, err
}

func waitClusterDeleted(conn *eks.EKS, name string, timeout time.Duration) (*eks.Cluster, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			eks.ClusterStatusActive,
			eks.ClusterStatusClaimed,
			eks.ClusterStatusCreating,
			eks.ClusterStatusDeleting,
			eks.ClusterStatusPending,
			eks.ClusterStatusProvisioning,
			eks.ClusterStatusReady,
			eks.ClusterStatusUpdating,
		},
		Target:         []string{eks.ClusterStatusDeleted},
		Refresh:        statusCluster(conn, name),
		Timeout:        timeout,
		NotFoundChecks: 1,
	}

	outputRaw, err := stateConf.WaitForState()

	if tfresource.NotFound(err) {
		return nil, nil
	}

	if output, ok := outputRaw.(*eks.Cluster); ok {
		if status, health := aws.StringValue(output.Status), output.Health; status == eks.ClusterStatusFailed && health != nil {
			tfresource.SetLastError(err, ClusterIssuesError(health.Issues))
		}

		return output, err
	}

	return nil, err
}

func waitClusterUpdateSuccessful(conn *eks.EKS, name, id string, timeout time.Duration) (*eks.Update, error) {
	if id == "" {
		if err := waitClusterReadyAfterUpdate(conn, name, timeout); err != nil {
			return nil, err
		}

		return &eks.Update{
			Status: aws.String(eks.UpdateStatusSuccessful),
		}, nil
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{eks.UpdateStatusInProgress},
		Target:  []string{eks.UpdateStatusSuccessful},
		Refresh: statusClusterUpdate(conn, name, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForState()
	if tfawserr.ErrCodeEquals(err, "PathNotFoundError") {
		if fallbackErr := waitClusterReadyAfterUpdate(conn, name, timeout); fallbackErr != nil {
			return nil, fallbackErr
		}

		return &eks.Update{
			Id:     aws.String(id),
			Status: aws.String(eks.UpdateStatusSuccessful),
		}, nil
	}

	if output, ok := outputRaw.(*eks.Update); ok {
		if status := aws.StringValue(output.Status); status == eks.UpdateStatusCancelled || status == eks.UpdateStatusFailed {
			tfresource.SetLastError(err, ErrorDetailsError(output.Errors))
		}

		return output, err
	}

	return nil, err
}

func waitClusterReadyAfterUpdate(conn *eks.EKS, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			eks.ClusterStatusActive,
			eks.ClusterStatusClaimed,
			eks.ClusterStatusCreating,
			eks.ClusterStatusPending,
			eks.ClusterStatusProvisioning,
			eks.ClusterStatusUpdating,
		},
		Target:  []string{eks.ClusterStatusReady},
		Refresh: statusCluster(conn, name),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*eks.Cluster); ok {
		if status, health := aws.StringValue(output.Status), output.Health; status == eks.ClusterStatusFailed && health != nil {
			tfresource.SetLastError(err, ClusterIssuesError(health.Issues))
		}

		return err
	}

	return err
}

func waitFargateProfileCreated(conn *eks.EKS, clusterName, fargateProfileName string, timeout time.Duration) (*eks.FargateProfile, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{eks.FargateProfileStatusCreating},
		Target:  []string{eks.FargateProfileStatusActive},
		Refresh: statusFargateProfile(conn, clusterName, fargateProfileName),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*eks.FargateProfile); ok {
		return output, err
	}

	return nil, err
}

func waitFargateProfileDeleted(conn *eks.EKS, clusterName, fargateProfileName string, timeout time.Duration) (*eks.FargateProfile, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{eks.FargateProfileStatusActive, eks.FargateProfileStatusDeleting},
		Target:  []string{},
		Refresh: statusFargateProfile(conn, clusterName, fargateProfileName),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*eks.FargateProfile); ok {
		return output, err
	}

	return nil, err
}

func waitNodegroupCreated(ctx context.Context, conn *eks.EKS, clusterName, nodeGroupName string, timeout time.Duration) (*eks.Nodegroup, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{eks.NodegroupStatusClaimed, eks.NodegroupStatusPending, eks.NodegroupStatusCreating, eks.NodegroupStatusProvisioning},
		Target:  []string{eks.NodegroupStatusActive},
		Refresh: statusNodegroup(conn, clusterName, nodeGroupName),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*eks.Nodegroup); ok {
		if status, health := aws.StringValue(output.Status), output.Health; (status == eks.NodegroupStatusCreateFailed || status == eks.NodegroupStatusDegraded) && health != nil {
			tfresource.SetLastError(err, IssuesError(health.Issues))
		}

		return output, err
	}

	return nil, err
}

func waitNodegroupDeleted(ctx context.Context, conn *eks.EKS, clusterName, nodeGroupName string, timeout time.Duration) (*eks.Nodegroup, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			eks.NodegroupStatusActive,
			eks.NodegroupStatusClaimed,
			eks.NodegroupStatusCreateFailed,
			eks.NodegroupStatusDegraded,
			eks.NodegroupStatusDeleting,
			eks.NodegroupStatusPending,
			eks.NodegroupStatusProvisioning,
			eks.NodegroupStatusUpdating,
		},
		Target:         []string{eks.NodegroupStatusDeleted},
		Refresh:        statusNodegroup(conn, clusterName, nodeGroupName),
		Timeout:        timeout,
		NotFoundChecks: 1,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if tfresource.NotFound(err) {
		return nil, nil
	}

	if output, ok := outputRaw.(*eks.Nodegroup); ok {
		if status, health := aws.StringValue(output.Status), output.Health; status == eks.NodegroupStatusDeleteFailed && health != nil {
			tfresource.SetLastError(err, IssuesError(health.Issues))
		}

		return output, err
	}

	return nil, err
}

func waitNodegroupUpdateSuccessful(ctx context.Context, conn *eks.EKS, clusterName, nodeGroupName, id string, timeout time.Duration) (*eks.Update, error) {
	if id == "" {
		if _, err := waitNodegroupActiveAfterUpdate(ctx, conn, clusterName, nodeGroupName, timeout); err != nil {
			return nil, err
		}

		return &eks.Update{
			Status: aws.String(eks.UpdateStatusSuccessful),
		}, nil
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{eks.UpdateStatusInProgress},
		Target:  []string{eks.UpdateStatusSuccessful},
		Refresh: statusNodegroupUpdate(conn, clusterName, nodeGroupName, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if tfawserr.ErrCodeEquals(err, "PathNotFoundError") {
		if _, fallbackErr := waitNodegroupActiveAfterUpdate(ctx, conn, clusterName, nodeGroupName, timeout); fallbackErr != nil {
			return nil, fallbackErr
		}

		return &eks.Update{
			Id:     aws.String(id),
			Status: aws.String(eks.UpdateStatusSuccessful),
		}, nil
	}

	if output, ok := outputRaw.(*eks.Update); ok {
		if status := aws.StringValue(output.Status); status == eks.UpdateStatusCancelled || status == eks.UpdateStatusFailed {
			tfresource.SetLastError(err, ErrorDetailsError(output.Errors))
		}

		return output, err
	}

	return nil, err
}

func waitNodegroupActiveAfterUpdate(ctx context.Context, conn *eks.EKS, clusterName, nodeGroupName string, timeout time.Duration) (*eks.Nodegroup, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			eks.NodegroupStatusClaimed,
			eks.NodegroupStatusPending,
			eks.NodegroupStatusProvisioning,
			eks.NodegroupStatusUpdating,
		},
		Target:  []string{eks.NodegroupStatusActive},
		Refresh: statusNodegroup(conn, clusterName, nodeGroupName),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*eks.Nodegroup); ok {
		if status, health := aws.StringValue(output.Status), output.Health; (status == eks.NodegroupStatusCreateFailed || status == eks.NodegroupStatusDegraded || status == eks.NodegroupStatusDeleteFailed) && health != nil {
			tfresource.SetLastError(err, IssuesError(health.Issues))
		}

		return output, err
	}

	return nil, err
}

func waitOIDCIdentityProviderConfigCreated(ctx context.Context, conn *eks.EKS, clusterName, configName string, timeout time.Duration) (*eks.OidcIdentityProviderConfig, error) {
	stateConf := resource.StateChangeConf{
		Pending: []string{eks.ConfigStatusCreating},
		Target:  []string{eks.ConfigStatusActive},
		Refresh: statusOIDCIdentityProviderConfig(ctx, conn, clusterName, configName),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*eks.OidcIdentityProviderConfig); ok {
		return output, err
	}

	return nil, err
}

func waitOIDCIdentityProviderConfigDeleted(ctx context.Context, conn *eks.EKS, clusterName, configName string, timeout time.Duration) (*eks.OidcIdentityProviderConfig, error) {
	stateConf := resource.StateChangeConf{
		Pending: []string{eks.ConfigStatusActive, eks.ConfigStatusDeleting},
		Target:  []string{},
		Refresh: statusOIDCIdentityProviderConfig(ctx, conn, clusterName, configName),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*eks.OidcIdentityProviderConfig); ok {
		return output, err
	}

	return nil, err
}
