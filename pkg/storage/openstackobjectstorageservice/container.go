package openstackobjectstorageservice

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
)

var (
	veleroContainerName = "managed-velero-backup-container"
)

func checkExistingContainer(ctx context.Context, reqLogger logr.Logger, client *OpenStackClient, containerName string) bool {

	Container := containers.Get(&client.serviceClient, veleroContainerName, nil)

	if Container.Result.Err != nil {
		return false
	}

	return true
}

func createContainer(ctx context.Context, client *OpenStackClient) (*string, error) {
	res := containers.Create(&client.serviceClient, veleroContainerName, nil)

	if res.Result.Err != nil {
		return nil, fmt.Errorf(
			"Error creating  container: %v . Error: %w",
			veleroContainerName,
			res.Result.Err,
		)
	}

	return &veleroContainerName, nil
}
