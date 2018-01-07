// Package v1 contains common functions for creating rds  instance
// for use in acceptance tests. See the `*_test.go` files for example usages.
package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/rds/v1/instances"
    "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
    "github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
    "github.com/gophercloud/gophercloud/openstack/rds/v1/datastores"
    "github.com/gophercloud/gophercloud/openstack/rds/v1/flavors"

	"fmt"
    "time"
    "errors"
)

// CreateRdsinstance creates a basic Rds instance with a randomly generated name.
func CreateRdsInstance(t *testing.T, client *gophercloud.ServiceClient,
    routerId string, networkId string, securityGroupId string) (*instances.Instance, error) {
	if testing.Short() {
		t.Skip("Skipping test that requires Rds instance creation in short mode.")
	}

    var instance *instances.Instance
    var datastoreId string
    datastoresList, err := datastores.List(client, "PostgreSQL").Extract()
    if err != nil {
        return instance, err
    }
    for _, datastore := range datastoresList {
        if datastore.Name == "9.5.5"{
            datastoreId = datastore.ID
        }
    }

    t.Logf("Attempting to get datastoreId: %s", datastoreId)

    flavorsList, err := flavors.List(client, datastoreId, "eu-de").Extract()
    if err != nil {
        return instance, err
    }

    if len(flavorsList) <1{
        t.Logf("Failed to get flavor for datastore : %s", datastoreId)
        return instance, err
    }
    t.Logf("Attempting to get flavors: %s", flavorsList[0].ID)

	name := tools.RandomString("ACPTTEST", 8)
	t.Logf("Attempting to create instance: %s", name)

	instance, err = instances.Create(client, instances.CreateOps{
		Name: name,
		DataStore: instances.DataStoreOps{Type:"PostgreSQL", Version:"9.5.5"},
        FlavorRef: flavorsList[0].ID,
		Volume: instances.VolumeOps{Type: "COMMON", Size: 100},
		Region: "eu-de",
		AvailabilityZone: "eu-de-01",
		Vpc: routerId,
		Nics: instances.NicsOps{SubnetId: networkId},
		SecurityGroup: instances.SecurityGroupOps{Id: securityGroupId},
		DbPort: "8635",
		BackupStrategy: instances.BackupStrategyOps{StartTime: "00:00:00", KeepDays: 0},
		DbRtPd: "Huangwei!120521",
        Ha: instances.HaOps{Enable: false},
	}).Extract()
	if err != nil {
		return instance, err
	}

	if err := WaitForInstanceStatus(t, client, instance.ID, "ACTIVE"); err != nil {
		return instance, err
	}

	return instance, nil
}

// WaitForInstanceStatus will poll an rds instance's status until it either matches
// the specified status or the status becomes ERROR.
func WaitForInstanceStatus(t *testing.T, client *gophercloud.ServiceClient, id string, status string) error {
    return WaitStatusFor(func() (bool, error) {
        latest, err := instances.Get(client, id).Extract()
        if err != nil {
            if _, ok := err.(gophercloud.ErrDefault404); ok {
                return true, nil
            }
            return false, err
        }
        t.Logf("Attempting to WaitForInstanceStatus: %s %s", latest.Status, status)
        if latest.Status == status {
            // Success!
            return true, nil
        }

        if latest.Status == "ERROR" {
            return false, fmt.Errorf("Rds Instance in ERROR state")
        }

        return false, nil
    })
}


func WaitForInstanceVolumeSize(t *testing.T, client *gophercloud.ServiceClient, id string, volumeSize int) error {
    return WaitStatusFor(func() (bool, error) {
        latest, err := instances.Get(client, id).Extract()
        if err != nil {
            return false, err
        }
        t.Logf("Attempting to WaitForInstanceVolumeSize: %s %s", latest.Volume.Size, volumeSize)
        if latest.Volume.Size == volumeSize {
            // Success!
            return true, nil
        }
        if latest.Status == "ERROR" {
            return false, fmt.Errorf("Rds Instance in ERROR state")
        }

        return false, nil
    })
}

func WaitStatusFor(predicate func() (bool, error)) error {
    for i := 0; i < 120; i++ {
        time.Sleep(10* time.Second)

        satisfied, err := predicate()
        if err != nil {
            return err
        }
        if satisfied {
            return nil
        }
    }
    return errors.New("Timed out")
}


func UpdateRdsInstance(t *testing.T, client *gophercloud.ServiceClient, id string,
    volumeSize int) (*instances.Instance, error) {
    if testing.Short() {
        t.Skip("Skipping test that requires Rds instance creation in short mode.")
    }

    var instance *instances.Instance

    t.Logf("Attempting to update instance: %s", id)

    var updateOpts instances.UpdateOps
    volume := make(map[string]interface{})
    volume["size"] = volumeSize
    updateOpts.Volume = volume

    instance, err := instances.UpdateVolumeSize(client, updateOpts, id).Extract()
    if err != nil {
        return instance, err
    }

    if err := WaitForInstanceVolumeSize(t, client, id, volumeSize); err != nil {
        return instance, err
    }

    return instance, nil
}


// DeleteRdsInstance deletes an instance via its UUID.
func DeleteRdsInstance(t *testing.T, client *gophercloud.ServiceClient, instance *instances.Instance) {
    result := instances.Delete(client, instance.ID)
    if result.Err != nil {
        t.Fatalf("Unable to delete Rds instance %s: %s", instance.ID, result.Err)
    }

    if err := WaitForInstanceStatus(t, client, instance.ID, "DELETED"); err != nil {
        t.Fatalf("Unable to delete Rds instance err status %s: %s", instance.ID, err)
    }
    time.Sleep(80* time.Second)
    t.Logf("Deleted Rds instance: %s", instance.ID)
}


func CreateRdsRouter(t *testing.T, client *gophercloud.ServiceClient) (*routers.Router, error) {
    var router *routers.Router

    routerName := tools.RandomString("TESTACC-", 8)

    t.Logf("Attempting to create router: %s", routerName)

    adminStateUp := true

    createOpts := routers.CreateOpts{
        Name:         routerName,
        AdminStateUp: &adminStateUp,
    }

    router, err := routers.Create(client, createOpts).Extract()
    if err != nil {
        return router, err
    }

    if err := WaitForRdsRouterToCreate(client, router.ID, 60); err != nil {
        return router, err
    }

    t.Logf("Created router: %s", routerName)

    return router, nil
}

func WaitForRdsRouterToCreate(client *gophercloud.ServiceClient, routerID string, secs int) error {
    return gophercloud.WaitFor(secs, func() (bool, error) {
        r, err := routers.Get(client, routerID).Extract()
        if err != nil {
            return false, err
        }

        if r.Status == "ACTIVE" {
            return true, nil
        }

        return false, nil
    })
}

// CreateRdsRouterInterface will attach a subnet to a router. An error will be
// returned if the operation fails.
func CreateRdsRouterInterface(t *testing.T, client *gophercloud.ServiceClient, subnetId, routerID string) (*routers.InterfaceInfo, error) {
    t.Logf("Attempting to add subnetId %s to router %s", subnetId, routerID)

    aiOpts := routers.AddInterfaceOpts{
        SubnetID: subnetId,
    }

    iface, err := routers.AddInterface(client, routerID, aiOpts).Extract()
    if err != nil {
        return iface, err
    }

    if err := WaitForRdsRouterInterfaceToAttach(client, iface.PortID, 60); err != nil {
        return iface, err
    }

    t.Logf("Successfully added port %s to router %s", iface.PortID, routerID)
    return iface, nil
}

func WaitForRdsRouterInterfaceToAttach(client *gophercloud.ServiceClient, routerInterfaceID string, secs int) error {
    return gophercloud.WaitFor(secs, func() (bool, error) {
        r, err := ports.Get(client, routerInterfaceID).Extract()
        fmt.Printf("WaitForRouterInterfaceToAttach got port=%+v\n.", r)
        if err != nil {
            fmt.Printf("WaitForRouterInterfaceToAttach returning error=%s.\n", err.Error())
            return false, err
        }

        return true, nil
    })
}

// DeleteRdsRouterInterface will detach a subnet to a router. A fatal error will
// occur if the deletion failed. This works best when used as a deferred
// function.
func DeleteRdsRouterInterface(t *testing.T, client *gophercloud.ServiceClient, subnetId, routerID string) {
    t.Logf("Attempting to detach subnet %s from router %s", subnetId, routerID)

    riOpts := routers.RemoveInterfaceOpts{
        SubnetID: subnetId,
    }

    _, err := routers.RemoveInterface(client, routerID, riOpts).Extract()
    if err != nil {
        if _, ok := err.(gophercloud.ErrDefault404); ok {
            t.Logf("Successfully detached subnetId %s from router %s", subnetId, routerID)
            return
        }
        if errCode, ok := err.(gophercloud.ErrUnexpectedResponseCode); ok {
            if errCode.Actual == 409 {
                t.Fatalf("Router Interface is still in use.subnetId:%s  routerID:%s", subnetId, routerID)
            }
        }
        t.Fatalf("Failed to detach subnetId %s from router %s with error %s", subnetId, routerID, err.Error())
    }

    t.Logf("Successfully detached subnetId %s from router %s", subnetId, routerID)
}
