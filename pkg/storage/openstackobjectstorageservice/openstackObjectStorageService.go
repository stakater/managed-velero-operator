package openstackobjectstorageservice

import (
	"context"

	"github.com/go-logr/logr"
	configv1 "github.com/openshift/api/config/v1"
	veleroInstallCR "github.com/openshift/managed-velero-operator/pkg/apis/managed/v1alpha2"
	storageBase "github.com/openshift/managed-velero-operator/pkg/storage/base"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type driver struct {
	storageBase.Driver
	client *OpenStackClient
}

// NewDriver creates a new OpenStackObjectStorageService driver
// Used during bootstrapping
func NewDriver(ctx context.Context, cfg *configv1.InfrastructureStatus, kubeClient client.Client) (*driver, error) {
	client, err := NewOpenStackClient(kubeClient, cfg)

	if err != nil {
		return nil, err
	}
	drv := driver{
		client: client,
	}
	drv.Context = ctx
	drv.KubeClient = kubeClient

	return &drv, nil
}

// GetPlatformType returns the platform type of this driver
func (d *driver) GetPlatformType() configv1.PlatformType {
	return configv1.AzurePlatformType
}

// CreateStorage attempts to create a Azure Blob Service Container with relevant tags
func (d *driver) CreateStorage(reqLogger logr.Logger, instance *veleroInstallCR.VeleroInstall) error {

	return nil // yet to be implemented
}

// StorageExists checks that the blob exists, and that we have access to it.
func (d *driver) StorageExists(status *veleroInstallCR.VeleroInstallStatus) (bool, error) {

	return false, nil // yet to be implemented
}
