package v1

import (
	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v1/vpcs"
	"testing"
)

func TestVpcList(t *testing.T) {
	client, err := clients.NewVpcV1Client()
	if err != nil {
		t.Fatalf("Unable to create a vpc client: %v", err)
	}

	listOpts := vpcs.ListOpts{}
	allVpcs, err := vpcs.List(client, listOpts)
	if err != nil {
		t.Fatalf("Unable to list routers: %v", err)
	}
	for _, router := range allVpcs {
		tools.PrintResource(t, router)
	}
}

func TestVpcsCRUD(t *testing.T) {
	client, err := clients.NewVpcV1Client()
	if err != nil {
		t.Fatalf("Unable to create a vpc client: %v", err)
	}

	// Create a vpc
	vpc, err := CreateVpc(t, client)
	if err != nil {
		t.Fatalf("Unable to create create: %v", err)
	}
	defer DeleteVpc(t, client, vpc.ID)

	tools.PrintResource(t, vpc)

	newName := tools.RandomString("TESTACC-", 8)
	updateOpts := &vpcs.UpdateOpts{
		Name: newName,
	}

	_, err = vpcs.Update(client, vpc.ID, updateOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update vpc: %v", err)
	}

	newVpc, err := vpcs.Get(client, vpc.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to retrieve vpc: %v", err)
	}

	tools.PrintResource(t, newVpc)
}
