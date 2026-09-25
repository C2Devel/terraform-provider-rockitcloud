package eks

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceClusterK2Schema(t *testing.T) {
	resource := ResourceCluster()
	if err := resource.InternalValidate(nil, true); err != nil {
		t.Fatalf("unexpected schema validation error: %s", err)
	}

	for _, name := range []string{"enabled_cluster_log_types", "encryption_config"} {
		if field := resource.Schema[name]; !field.Computed || field.Optional {
			t.Fatalf("%s must be computed-only", name)
		}
	}

	vpc := resource.Schema["vpc_config"].Elem.(*schema.Resource).Schema
	if vpc["security_group_ids"].ForceNew {
		t.Fatal("security_group_ids must update in place")
	}
	for _, name := range []string{"endpoint_private_access", "endpoint_public_access", "public_access_cidrs"} {
		if field := vpc[name]; !field.Computed || field.Optional {
			t.Fatalf("vpc_config.%s must be computed-only", name)
		}
	}

	legacy := resource.Schema["legacy_cluster_params"].Elem.(*schema.Resource).Schema
	if !legacy["master_config"].Required {
		t.Fatal("master_config must be required when legacy_cluster_params is configured")
	}
	for blockName, flagName := range map[string]string{
		"docker_registry_config": "docker_registry_required",
		"ebs_provider_config":    "ebs_provider_required",
		"ingress_config":         "ingress_required",
		"nlb_provider_config":    "nlb_provider_required",
	} {
		block := legacy[blockName].Elem.(*schema.Resource).Schema
		if _, ok := block[flagName]; ok {
			t.Fatalf("%s must remain implicitly enabled for backward compatibility", flagName)
		}
	}
	userData := legacy["user_data_config"].Elem.(*schema.Resource).Schema
	if userData["user_data"].ForceNew || userData["user_data_content_type"].ForceNew {
		t.Fatal("user_data_config fields must update in place")
	}

	if got := resource.Timeouts.Delete; got == nil || *got != 60*time.Minute {
		t.Fatalf("cluster delete timeout must match the 60 minute delete retry window, got %v", got)
	}
}

func TestExpandAndFlattenK2ClusterConfiguration(t *testing.T) {
	network := expandNetworkConfigRequest([]interface{}{map[string]interface{}{
		"ip_family":         "ipv4",
		"pod_ipv4_cidr":     "10.10.0.0/16",
		"service_ipv4_cidr": "10.20.0.0/16",
	}})
	if got := aws.StringValue(network.PodIpv4Cidr); got != "10.10.0.0/16" {
		t.Fatalf("unexpected pod CIDR: %q", got)
	}
	flattenedNetwork := flattenNetworkConfig(&eks.KubernetesNetworkConfigResponse{
		PodIpv4Cidr: aws.String("10.10.0.0/16"),
	})
	if got := flattenedNetwork[0].(map[string]interface{})["pod_ipv4_cidr"]; got != "10.10.0.0/16" {
		t.Fatalf("unexpected flattened pod CIDR: %q", got)
	}

	securityGroups := schema.NewSet(schema.HashString, []interface{}{"sg-1", "sg-2"})
	vpcUpdate := expandVPCSecurityGroupUpdateRequest(securityGroups)
	if len(vpcUpdate.SecurityGroupIds) != 2 || vpcUpdate.SubnetIds != nil ||
		vpcUpdate.EndpointPrivateAccess != nil || vpcUpdate.EndpointPublicAccess != nil ||
		vpcUpdate.PublicAccessCidrs != nil {
		t.Fatalf("security group update contains unsupported fields: %#v", vpcUpdate)
	}

	remoteAccess := expandClusterRemoteAccessConfig([]interface{}{map[string]interface{}{
		"ec2_ssh_key": "key-name",
	}})
	if got := aws.StringValue(remoteAccess.Ec2SshKey); got != "key-name" {
		t.Fatalf("unexpected SSH key: %q", got)
	}

	userData := expandUserDataConfig([]interface{}{map[string]interface{}{
		"user_data":              "#cloud-config",
		"user_data_content_type": "cloud-config",
	}})
	if got := aws.StringValue(userData.UserData); got != "#cloud-config" {
		t.Fatalf("unexpected user data: %q", got)
	}
	if got := aws.StringValue(userData.UserDataContentType); got != "cloud-config" {
		t.Fatalf("unexpected user data content type: %q", got)
	}
	flattenedUserData := flattenUserDataConfig(userData)
	if got := flattenedUserData[0].(map[string]interface{})["user_data"]; got != "#cloud-config" {
		t.Fatalf("unexpected flattened user data: %q", got)
	}

	legacy := expandLegacyClusterParams([]interface{}{map[string]interface{}{
		"cluster_autoscaler_config": []interface{}{map[string]interface{}{
			"cluster_autoscaler_required": false,
			"cluster_autoscaler_user":     "autoscaler",
		}},
		"docker_registry_config": []interface{}{map[string]interface{}{
			"volume_size": 10,
			"volume_type": "ssd",
		}},
		"ebs_provider_config": []interface{}{map[string]interface{}{
			"ebs_user": "ebs",
		}},
		"ingress_config": []interface{}{map[string]interface{}{
			"instance_type": "small",
			"volume_size":   10,
			"volume_type":   "ssd",
		}},
		"master_config": []interface{}{map[string]interface{}{
			"high_availability": false,
			"instance_type":     "small",
			"volume_size":       10,
			"volume_type":       "ssd",
		}},
		"nlb_provider_config": []interface{}{map[string]interface{}{
			"nlb_user": "nlb",
		}},
	}})

	if aws.BoolValue(legacy.ClusterAutoscalerConfig.ClusterAutoscalerRequired) {
		t.Fatal("explicit false Cluster Autoscaler flag was not preserved")
	}
	if !aws.BoolValue(legacy.DockerRegistryConfig.DockerRegistryRequired) ||
		!aws.BoolValue(legacy.EbsProviderConfig.EbsProviderRequired) ||
		!aws.BoolValue(legacy.IngressConfig.IngressRequired) ||
		!aws.BoolValue(legacy.NlbProviderConfig.NlbProviderRequired) {
		t.Fatal("existing integrated services must remain implicitly enabled")
	}

	flattened := flattenLegacyClusterParams(&eks.LegacyClusterParamsResponse{
		ClusterAutoscalerConfig: &eks.ClusterAutoscalerConfig{
			ClusterAutoscalerRequired: aws.Bool(true),
			ClusterAutoscalerUserName: aws.String("managed-autoscaler"),
		},
		EbsProviderConfig: &eks.EbsProviderConfigResponse{
			EbsProviderRequired: aws.Bool(true),
			EbsUser:             aws.String("configured-ebs"),
			EbsUserName:         aws.String("managed-ebs"),
		},
	})
	values := flattened[0].(map[string]interface{})
	autoscaler := values["cluster_autoscaler_config"].([]interface{})[0].(map[string]interface{})
	if got := autoscaler["cluster_autoscaler_user"]; got != "managed-autoscaler" {
		t.Fatalf("unexpected autoscaler user: %q", got)
	}
	ebs := values["ebs_provider_config"].([]interface{})[0].(map[string]interface{})
	if got := ebs["ebs_user"]; got != "configured-ebs" {
		t.Fatalf("unexpected EBS user: %q", got)
	}

	emptyLegacy := flattenLegacyClusterParams(&eks.LegacyClusterParamsResponse{
		DockerRegistryConfig: &eks.DockerRegistryConfig{
			DockerRegistryRequired: aws.Bool(false),
		},
		EbsProviderConfig: &eks.EbsProviderConfigResponse{
			EbsProviderRequired: aws.Bool(false),
		},
		IngressConfig: &eks.IngressConfig{
			IngressRequired: aws.Bool(false),
		},
		NlbProviderConfig: &eks.NlbProviderConfigResponse{
			NlbProviderRequired: aws.Bool(false),
		},
		UserDataConfig: &eks.UserDataConfig{},
		PlacementConfig: &eks.PlacementConfig{
			Tenancy: aws.String("default"),
		},
	})
	if len(emptyLegacy) != 0 {
		t.Fatalf("empty API legacy blocks must not create Terraform drift: %#v", emptyLegacy)
	}
	if got := flattenClusterRemoteAccessConfig(&eks.RemoteAccessConfig{}); len(got) != 0 {
		t.Fatalf("empty remote access block must not create Terraform drift: %#v", got)
	}
}

// A disabled Cluster Autoscaler must survive the read: dropping it made an explicit
// cluster_autoscaler_required = false replace the cluster on every apply.
func TestDisabledClusterAutoscalerIsReadBack(t *testing.T) {
	flattened := flattenClusterAutoscalerConfig(&eks.ClusterAutoscalerConfig{
		ClusterAutoscalerRequired: aws.Bool(false),
	})
	if len(flattened) != 1 {
		t.Fatalf("disabled Cluster Autoscaler was dropped: %#v", flattened)
	}
	if got := flattened[0].(map[string]interface{})["cluster_autoscaler_required"]; got != false {
		t.Fatalf("unexpected Cluster Autoscaler flag: %#v", got)
	}

	// Configurations that do not declare the block must keep reading it from the API.
	legacy := ResourceCluster().Schema["legacy_cluster_params"]
	autoscaler := legacy.Elem.(*schema.Resource).Schema["cluster_autoscaler_config"]
	if !legacy.Computed || !autoscaler.Computed {
		t.Fatal("cluster_autoscaler_config must be computed to stay out of the diff when it is not configured")
	}
}

func TestResourceNodeGroupK2SchemaAndLabelValidation(t *testing.T) {
	resource := ResourceNodeGroup()
	if err := resource.InternalValidate(nil, true); err != nil {
		t.Fatalf("unexpected schema validation error: %s", err)
	}

	for _, name := range []string{"force_update_version", "launch_template", "release_version", "version"} {
		if field := resource.Schema[name]; !field.Computed || field.Optional {
			t.Fatalf("%s must be computed-only", name)
		}
	}

	minSize := resource.Schema["scaling_config"].Elem.(*schema.Resource).Schema["min_size"]
	if _, errors := minSize.ValidateFunc(0, "min_size"); len(errors) == 0 {
		t.Fatal("min_size=0 must be rejected")
	}

	valid := map[string]interface{}{
		"example.com/name": "value_1",
		"empty":            "",
	}
	if _, errors := validateKubernetesLabels(valid, "labels"); len(errors) != 0 {
		t.Fatalf("valid labels rejected: %v", errors)
	}

	invalid := map[string]interface{}{
		"Bad Prefix/name":  strings.Repeat("x", 64),
		"bad..prefix/name": "value",
	}
	if _, errors := validateKubernetesLabels(invalid, "labels"); len(errors) < 3 {
		t.Fatalf("invalid labels were not fully rejected: %v", errors)
	}
}

func TestClusterAuthIsExplicitlyUnsupported(t *testing.T) {
	err := dataSourceClusterAuthRead(nil, nil)
	if err == nil || !strings.Contains(err.Error(), "aws_eks_cluster_kubeconfig") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestK2FailureDetailsAreActionable(t *testing.T) {
	clusterErr := ClusterIssuesError([]*eks.ClusterIssue{{
		Code:        aws.String("InsufficientCapacity"),
		Message:     aws.String("control plane could not be provisioned"),
		ResourceIds: aws.StringSlice([]string{"cluster/test"}),
	}})
	if clusterErr == nil ||
		!strings.Contains(clusterErr.Error(), "InsufficientCapacity") ||
		!strings.Contains(clusterErr.Error(), "control plane could not be provisioned") {
		t.Fatalf("cluster health details were lost: %v", clusterErr)
	}

	updateErr := ErrorDetailsError([]*eks.ErrorDetail{{
		ErrorCode:    aws.String("InvalidParameter"),
		ErrorMessage: aws.String("security group update was rejected"),
		ResourceIds:  aws.StringSlice([]string{"sg-test"}),
	}})
	if updateErr == nil ||
		!strings.Contains(updateErr.Error(), "InvalidParameter") ||
		!strings.Contains(updateErr.Error(), "security group update was rejected") {
		t.Fatalf("update failure details were lost: %v", updateErr)
	}
}

func TestClusterUpdateFallsBackWhenDescribeUpdateIsUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/clusters/test/updates/update-1":
			w.Header().Set("X-Amzn-Errortype", "PathNotFoundError")
			http.Error(w, `{"message":"Specified path does not exist."}`, http.StatusBadRequest)
		case "/clusters/test":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"cluster":{"name":"test","status":"READY"}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	conn := eks.New(session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("test", "test", ""),
		Endpoint:    aws.String(server.URL),
		Region:      aws.String("ru-msk"),
	})))

	update, err := waitClusterUpdateSuccessful(conn, "test", "update-1", 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected fallback waiter error: %s", err)
	}
	if got := aws.StringValue(update.Status); got != eks.UpdateStatusSuccessful {
		t.Fatalf("unexpected fallback update status: %q", got)
	}

	update, err = waitClusterUpdateSuccessful(conn, "test", "", 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected empty update ID fallback error: %s", err)
	}
	if got := aws.StringValue(update.Status); got != eks.UpdateStatusSuccessful {
		t.Fatalf("unexpected empty ID fallback update status: %q", got)
	}
}

func TestClusterUpdateFallbackAcceptsModifyingStatus(t *testing.T) {
	var describes int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/clusters/test/updates/update-1":
			w.Header().Set("X-Amzn-Errortype", "PathNotFoundError")
			http.Error(w, `{"message":"Specified path does not exist."}`, http.StatusBadRequest)
		case "/clusters/test":
			status := eks.ClusterStatusReady
			if atomic.AddInt32(&describes, 1) == 1 {
				status = clusterStatusModifying
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"cluster":{"name":"test","status":%q}}`, status)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	conn := eks.New(session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("test", "test", ""),
		Endpoint:    aws.String(server.URL),
		Region:      aws.String("ru-msk"),
	})))

	update, err := waitClusterUpdateSuccessful(conn, "test", "update-1", 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected fallback waiter error: %s", err)
	}
	if got := aws.StringValue(update.Status); got != eks.UpdateStatusSuccessful {
		t.Fatalf("unexpected fallback update status: %q", got)
	}
}

func TestClusterDeleteAcceptsModifyingStatus(t *testing.T) {
	var describes int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&describes, 1) == 1 {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"cluster":{"name":"test","status":%q}}`, clusterStatusModifying)

			return
		}

		w.Header().Set("X-Amzn-Errortype", ErrCodeClusterNotFound)
		http.Error(w, `{"message":"Cluster not found."}`, http.StatusNotFound)
	}))
	defer server.Close()

	conn := eks.New(session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("test", "test", ""),
		Endpoint:    aws.String(server.URL),
		Region:      aws.String("ru-msk"),
	})))

	cluster, err := waitClusterDeleted(conn, "test", 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected delete waiter error: %s", err)
	}
	if cluster != nil {
		t.Fatalf("deleted cluster was reported as existing: %#v", cluster)
	}
}

// The cluster resource and data source are filled by the same flatten functions,
// so every block the data source exposes must accept the whole resource block.
func TestDataSourceClusterBlocksMatchResource(t *testing.T) {
	block := func(fields map[string]*schema.Schema, name string) map[string]*schema.Schema {
		if field, ok := fields[name]; ok {
			if nested, ok := field.Elem.(*schema.Resource); ok {
				return nested.Schema
			}
		}

		return nil
	}

	var compare func(path string, resource, dataSource map[string]*schema.Schema)
	compare = func(path string, resource, dataSource map[string]*schema.Schema) {
		for name := range resource {
			if _, ok := dataSource[name]; !ok {
				t.Errorf("data source cannot store %s%s", path, name)
				continue
			}

			if nested := block(resource, name); nested != nil {
				compare(path+name+".", nested, block(dataSource, name))
			}
		}
	}

	resource := ResourceCluster().Schema
	dataSource := DataSourceCluster().Schema

	// Only the blocks the data source declares are compared: the resource also has
	// top-level attributes the data source deliberately does not read.
	for name := range dataSource {
		if nested := block(dataSource, name); nested != nil {
			compare(name+".", block(resource, name), nested)
		}
	}
}

func TestNodegroupUpdateFallsBackWhenDescribeUpdateIsUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/clusters/test/updates/update-1":
			w.Header().Set("X-Amzn-Errortype", "PathNotFoundError")
			http.Error(w, `{"message":"Specified path does not exist."}`, http.StatusBadRequest)
		case "/clusters/test/node-groups/general":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"nodegroup":{"clusterName":"test","nodegroupName":"general","status":"ACTIVE"}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	conn := eks.New(session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("test", "test", ""),
		Endpoint:    aws.String(server.URL),
		Region:      aws.String("ru-msk"),
	})))

	ctx := context.Background()
	update, err := waitNodegroupUpdateSuccessful(ctx, conn, "test", "general", "update-1", nil, 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected fallback waiter error: %s", err)
	}
	if got := aws.StringValue(update.Status); got != eks.UpdateStatusSuccessful {
		t.Fatalf("unexpected fallback update status: %q", got)
	}

	update, err = waitNodegroupUpdateSuccessful(ctx, conn, "test", "general", "", nil, 5*time.Second)
	if err != nil {
		t.Fatalf("unexpected empty update ID fallback error: %s", err)
	}
	if got := aws.StringValue(update.Status); got != eks.UpdateStatusSuccessful {
		t.Fatalf("unexpected empty ID fallback update status: %q", got)
	}
}

func TestNodegroupUpdateFallbackWaitsForRequestedDesiredSize(t *testing.T) {
	var describes int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/clusters/test/updates/update-1":
			w.Header().Set("X-Amzn-Errortype", "PathNotFoundError")
			http.Error(w, `{"message":"Specified path does not exist."}`, http.StatusBadRequest)
		case "/clusters/test/node-groups/general":
			desiredSize := 3
			if atomic.AddInt32(&describes, 1) == 1 {
				desiredSize = 2
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"nodegroup":{"clusterName":"test","nodegroupName":"general","status":"ACTIVE","scalingConfig":{"desiredSize":%d,"maxSize":3,"minSize":1}}}`, desiredSize)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	conn := eks.New(session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("test", "test", ""),
		Endpoint:    aws.String(server.URL),
		Region:      aws.String("ru-msk"),
	})))

	if _, err := waitNodegroupUpdateSuccessful(context.Background(), conn, "test", "general", "update-1", aws.Int64(3), 5*time.Second); err != nil {
		t.Fatalf("unexpected fallback waiter error: %s", err)
	}
	if got := atomic.LoadInt32(&describes); got < 2 {
		t.Fatalf("fallback waiter returned on the stale desired size after %d describes", got)
	}
}
