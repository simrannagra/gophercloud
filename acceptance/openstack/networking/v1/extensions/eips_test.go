package extensions

import (
	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v1/extensions/eip"
	"testing"
)

func TestEipList(t *testing.T) {
	client, err := clients.NewVpcV1Client()
	if err != nil {
		t.Fatalf("Unable to create a eip client: %v", err)
	}

	listOpts := eip.ListOpts{}
	allEips, err := eip.List(client, listOpts)
	if err != nil {
		t.Fatalf("Unable to list eips: %v", err)
	}
	for _, router := range allEips {
		tools.PrintResource(t, router)
	}
}

func TestEipsCRUD(t *testing.T) {
	client, err := clients.NewVpcV1Client()
	if err != nil {
		t.Fatalf("Unable to create a eip client: %v", err)
	}

	// Create a eip
	Eip, err := CreateEip(t, client)
	if err != nil {
		t.Fatalf("Unable to create eip: %v", err)
	}

	// Delete a eip
	defer DeleteEip(t, client, Eip.ID)

	tools.PrintResource(t, Eip)
	updateOpts := eip.UpdateOpts{
		PortId: "3f2e210a-d2f0-4275-a0d5-79c69a571df8",
	}

	_, err = eip.Update(client, Eip.ID, updateOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update eip: %v", err)
	}

	newEip, err := eip.Get(client, Eip.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to retrieve eip: %v", err)
	}

	tools.PrintResource(t, newEip)
}

