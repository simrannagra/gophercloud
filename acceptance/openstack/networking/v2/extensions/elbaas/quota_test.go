// +build acceptance networking lbaas_v2 listeners

package elbaas

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/listeners"
)

func TestQuotaList(t *testing.T) {
    client, err := clients.NewOtcV1Client("elb")
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}
	
	allPages, err := listeners.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list listeners: %v", err)
	}

	allListeners, err := listeners.ExtractListeners(allPages)
	if err != nil {
		t.Fatalf("Unable to extract listeners: %v", err)
	}

	for _, listener := range allListeners {
		tools.PrintResource(t, listener)
	}
}
