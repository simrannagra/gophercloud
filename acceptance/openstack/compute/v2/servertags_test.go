package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func TestServerTagsList(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	allPages, err := servers.List(client, servers.ListOpts{}).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve servers: %v", err)
	}

	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		t.Fatalf("Unable to extract servers: %v", err)
	}

	taglist := []string{"foo.bar", "name.value"}
	for _, server := range allServers {
		tools.PrintResource(t, server)
		CreateServerTags(t, client, server.ID, taglist)
		tags2, err := GetServerTags(t, client, server.ID)
		if err != nil {
			t.Fatalf("Unable to get tags: %v", err)
		}
		tools.PrintResource(t, tags2)
		err = DeleteServerTags(t, client, server.ID)
		if err != nil {
			t.Fatalf("Unable to delete tags: %v", err)
		}
	}
}


