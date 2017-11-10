// +build acceptance networking lbaas_v2 monitors

package elbaas

import (
	"testing"

	//"github.com/gophercloud/gophercloud/acceptance/clients"
	//"github.com/gophercloud/gophercloud/acceptance/tools"
	//"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/healthcheck"
)

func TestHealthList(t *testing.T) {
    /*
	client, err := clients.NewOtcV1Client("elb")
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	allPages, err := healthcheck.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list health: %v", err)
	}

	allHealth, err := healthcheck.ExtractHealth(allPages)
	if err != nil {
		t.Fatalf("Unable to extract health: %v", err)
	}

	for _, health := range allHealth {
		tools.PrintResource(t, health)
	}
    */
}
