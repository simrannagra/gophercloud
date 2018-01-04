package v1

import (
	"testing"
    "github.com/gophercloud/gophercloud/acceptance/clients"
    "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2/extensions/layer3"
    networking "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2"
    "github.com/gophercloud/gophercloud/acceptance/tools"
    extensions "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2/extensions"
    "github.com/gophercloud/gophercloud/openstack/rds/v1/instances"
    "github.com/gophercloud/gophercloud/openstack/rds/v1/datastores"
    "github.com/gophercloud/gophercloud/openstack/rds/v1/flavors"
)

func TestRdsInstanceCreateDestroy(t *testing.T) {
    client, err := clients.NewNetworkV2Client()
    if err != nil {
        t.Fatalf("Unable to create a network client: %v", err)
    }

    client1, err := clients.NewRdsV1Client()
    if err != nil {
        t.Fatalf("Unable to create a rds client: %v", err)
    }
    // Create a Security Group
    group, err := extensions.CreateSecurityGroup(t, client)
    if err != nil {
        t.Fatalf("Unable to create security group: %v", err)
    }
    defer extensions.DeleteSecurityGroup(t, client, group.ID)

    // Create a network
    network, err := networking.CreateNetwork(t, client)
    if err != nil {
        t.Fatalf("Unable to create network: %v", err)
    }
    defer networking.DeleteNetwork(t, client, network.ID)

    subnet, err := networking.CreateSubnet(t, client, network.ID)
    if err != nil {
        t.Fatalf("Unable to create subnet: %v", err)
    }
    defer networking.DeleteSubnet(t, client, subnet.ID)

    router, err := CreateRdsRouter(t, client)
    if err != nil {
        t.Fatalf("Unable to create router: %v", err)
    }
    defer layer3.DeleteRouter(t, client, router.ID)

    _, err = CreateRdsRouterInterface(t, client, subnet.ID, router.ID)
    if err != nil {
        t.Fatalf("Unable to create router interface: %v", err)
    }
    defer DeleteRdsRouterInterface(t, client, subnet.ID, router.ID)

    instance, err := CreateRdsInstance(t, client1, router.ID, network.ID, group.ID)
    if err != nil {
        t.Fatalf("Unable to create rds instance: %v", err)
    }
    defer DeleteRdsInstance(t, client1, instance)

    newinstance, err := instances.Get(client1, instance.ID).Extract()
    if err != nil {
        t.Errorf("Unable to retrieve rds instance: %v", err)
    }
    tools.PrintResource(t, newinstance)
}

func TestRdsInstanceList(t *testing.T) {
    client, err := clients.NewRdsV1Client()
    if err != nil {
        t.Fatalf("Unable to create a Rds client: %v", err)
    }

    instancesList, err := instances.List(client).Extract()
    if err != nil {
        t.Fatalf("Unable to retrieve Rds instances: %v", err)
    }

    for _, instance := range instancesList {
        tools.PrintResource(t, instance)
    }
}

func TestRdsDatastoreList(t *testing.T) {
    client, err := clients.NewRdsV1Client()
    if err != nil {
        t.Fatalf("Unable to create a Rds client: %v", err)
    }

    datastoresList, err := datastores.List(client, "PostgreSQL").Extract()
    if err != nil {
        t.Fatalf("Unable to retrieve Rds datastores: %v", err)
    }

    for _, datastore := range datastoresList {
        tools.PrintResource(t, datastore)
    }
}

func TestRdsFlavorList(t *testing.T) {
    client, err := clients.NewRdsV1Client()
    if err != nil {
        t.Fatalf("Unable to create a Rds client: %v", err)
    }

    datastoresList, err := datastores.List(client, "PostgreSQL").Extract()
    if err != nil {
        t.Fatalf("Unable to retrieve Rds datastores: %v", err)
    }
    var datastoreId string
    for _, datastore := range datastoresList {
        if datastore.Name == "9.5.5"{
            datastoreId = datastore.ID
        }
    }
    flavorsList, err := flavors.List(client, datastoreId, "eu-de").Extract()
    if err != nil {
        t.Fatalf("Unable to retrieve Rds flavors: %v", err)
    }
    for _, flavor := range flavorsList {
        tools.PrintResource(t, flavor)
    }
}