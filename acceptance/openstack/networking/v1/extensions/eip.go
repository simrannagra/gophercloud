package extensions

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v1/extensions/eip"
)

func CreateEip(t *testing.T, client *gophercloud.ServiceClient) (*eip.Eip, error) {

	eipName := tools.RandomString("TESTACC-", 8)

	bandwidth:=eip.Bandwidth{Name:"bandwidth123",Size:10,ShareType:"PER"}

	publicip:=eip.Publicip{Type:"5_bgp"}

	createOpts := eip.CreateOpts{
		PublicIp: publicip,
		BandWidth: bandwidth,
	}

	t.Logf("Attempting to create eip: %s", eipName)

	Eip, err := eip.Create(client, createOpts).Extract()
	if err != nil {
		return Eip, err
	}
	t.Logf("Created eip: %s", eipName)

	return Eip, nil
}

func DeleteEip(t *testing.T, client *gophercloud.ServiceClient, EipID string) {
	t.Logf("Attempting to delete eip: %s", EipID)

	err := eip.Delete(client, EipID).ExtractErr()
	if err != nil {
		t.Fatalf("Error deleting eip: %v", err)
	}

	t.Logf("Deleted eip: %s", EipID)
}
