package openstackobjectstorageservice

import (
	"context"
	"fmt"

	"github.com/openshift/managed-velero-operator/version"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"

	configv1 "github.com/openshift/api/config/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	openstackCredsSecretName = version.OperatorName + "-iam-credentials"
)

// OpenStackClient interact with Azure API
type OpenStackClient struct {
	region        string
	serviceClient gophercloud.ServiceClient
}

func getCredentialsSecret(kubeClient client.Client) (secret *corev1.Secret, err error) {
	namespace, err := k8sutil.GetOperatorNamespace()
	if err != nil {
		return nil, fmt.Errorf("failed to get operator namespace: %v", err)
	}

	secret = &corev1.Secret{}
	err = kubeClient.Get(context.TODO(),
		types.NamespacedName{
			Name:      openstackCredsSecretName,
			Namespace: namespace,
		},
		secret)

	return secret, err
}

func getStringFromSecret(secret *corev1.Secret, key string) (string, bool) {
	bytesVal, ok := secret.Data[key]
	if !ok {
		return "", false
	}
	return string(bytesVal), true
}

func getOpenStackCredentials(kubeClient client.Client) (authURL string, username string, password string, tenantID string, region string, err error) {

	secret, err := getCredentialsSecret(kubeClient)

	if err != nil {
		return "", "", "", "", "", err
	}

	authURL, ok := getStringFromSecret(secret, "os_auth_url")
	if !ok {
		return "", "", "", "", "", fmt.Errorf("os_auth_url is missing for secret: '%v', namespace: '%v'", secret.Name, secret.Namespace)
	}
	username, ok = getStringFromSecret(secret, "os_username")
	if !ok {
		return "", "", "", "", "", fmt.Errorf("os_username is missing for secret: '%v', namespace: '%v'", secret.Name, secret.Namespace)
	}
	password, ok = getStringFromSecret(secret, "os_password")
	if !ok {
		return "", "", "", "", "", fmt.Errorf("os_password is missing for secret: '%v', namespace: '%v'", secret.Name, secret.Namespace)
	}
	tenantID, ok = getStringFromSecret(secret, "os_tenant_id")
	if !ok {
		return "", "", "", "", "", fmt.Errorf("os_tenant_id is missing for secret: '%v', namespace: '%v'", secret.Name, secret.Namespace)
	}
	region, ok = getStringFromSecret(secret, "os_region_name")
	if !ok {
		return "", "", "", "", "", fmt.Errorf("os_region_name is missing for secret: '%v', namespace: '%v'", secret.Name, secret.Namespace)
	}

	return authURL, username, password, tenantID, region, nil
}

// NewAzureClient reads the credentials secret in the operator's namespace and uses
// them to create a new azure client.
func NewOpenStackClient(kubeClient client.Client, cfg *configv1.InfrastructureStatus) (*OpenStackClient, error) {
	var err error
	authURL, username, password, tenantID, region, err := getOpenStackCredentials(kubeClient)
	if err != nil {
		return nil, err
	}
	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: authURL,
		Username:         username,
		Password:         password,
		TenantID:         tenantID,
	}
	provider, err := openstack.AuthenticatedClient(authOpts)

	if err != nil {
		return nil, err
	}
	client, err := openstack.NewObjectStorageV1(provider, gophercloud.EndpointOpts{
		Region: region,
	})

	if err != nil {
		return nil, err
	}

	return &OpenStackClient{
		serviceClient: *client,
		region:        region,
	}, nil
}
