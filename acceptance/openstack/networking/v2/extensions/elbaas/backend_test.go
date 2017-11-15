// +build acceptance networking lbaas_v2 listeners

package elbaas

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/backendmember"
)

func TestBackendList(t *testing.T) {
    client, err := clients.NewOtcV1Client("elb")
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}
	
	allPages, err := backendmember.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list backend: %v", err)
	}

	allBackend, err := backendmember.ExtractBackend(allPages)
	if err != nil {
		t.Fatalf("Unable to extract Backend: %v", err)
	}

	for _, backend := range allBackend {
		tools.PrintResource(t, backend)
	}
}
