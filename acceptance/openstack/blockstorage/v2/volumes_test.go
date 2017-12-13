// +build acceptance blockstorage

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/volumes"
	"reflect"
)

func TestVolumesList(t *testing.T) {
	client, err := clients.NewBlockStorageV2Client()
	if err != nil {
		t.Fatalf("Unable to create a blockstorage client: %v", err)
	}

	allPages, err := volumes.List(client, volumes.ListOpts{}).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve volumes: %v", err)
	}

	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		t.Fatalf("Unable to extract volumes: %v", err)
	}

	for _, volume := range allVolumes {
		tools.PrintResource(t, volume)
	}
}

func TestVolumesCreateDestroy(t *testing.T) {
	client, err := clients.NewBlockStorageV2Client()
	if err != nil {
		t.Fatalf("Unable to create blockstorage client: %v", err)
	}

	volume, err := CreateVolume(t, client)
	if err != nil {
		t.Fatalf("Unable to create volume: %v", err)
	}
	defer DeleteVolume(t, client, volume)

	newVolume, err := volumes.Get(client, volume.ID).Extract()
	if err != nil {
		t.Errorf("Unable to retrieve volume: %v", err)
	}

	tools.PrintResource(t, newVolume)
}


func TestVolumesTags(t *testing.T) {
	client, err := clients.NewBlockStorageV2Client()
	if err != nil {
		t.Fatalf("Unable to create blockstorage client: %v", err)
	}

	volume, err := CreateVolume(t, client)
	if err != nil {
		t.Fatalf("Unable to create volume: %v", err)
	}
	defer DeleteVolume(t, client, volume)

	tagmap0, err := GetVolumeTags(t, client, "volumes", volume.ID)
	if err != nil {
		t.Errorf("Unable to get initial tags from volume: %v", err)
	}
	tools.PrintResource(t, tagmap0)

	tagmap := map[string]string{"foo" : "bar", "name" : "value"}
	tagmap2, err := CreateVolumeTags(t, client, "volumes", volume.ID, tagmap)
	if err != nil {
		t.Errorf("Unable to create tags for volume: %v", err)
	}
	tagmap3, err := GetVolumeTags(t, client, "volumes", volume.ID)
	if err != nil {
		t.Errorf("Unable to get tags from volume: %v", err)
	}
	if !reflect.DeepEqual(tagmap3.Tags, tagmap) {
		t.Errorf("Tags aren't equal after set/get: %v != %v", tagmap3.Tags, tagmap)
	}
	tools.PrintResource(t, tagmap2)

	tagmap4 := map[string]string{"foo2" : "bar2", "name2" : "value2"}
	tagmap5, err := CreateVolumeTags(t, client, "volumes", volume.ID, tagmap4)
	if err != nil {
		t.Errorf("Unable to create tags for volume: %v", err)
	}
	tagmap6, err := GetVolumeTags(t, client, "volumes", volume.ID)
	if err != nil {
		t.Errorf("Unable to get tags from volume: %v", err)
	}
	if !reflect.DeepEqual(tagmap6.Tags, tagmap4) {
		t.Errorf("Tags aren't equal after set/get: %v != %v", tagmap6.Tags, tagmap4)
	}
	tools.PrintResource(t, tagmap5)

	tagmap0a := map[string]string{}
	err = DeleteVolumeTags(t, client, "volumes", volume.ID)
	if err != nil {
		t.Errorf("Unable to delete tags for volume: %v", err)
	}
	tagmap0c, err := GetVolumeTags(t, client, "volumes", volume.ID)
	if err != nil {
		t.Errorf("Unable to get empty tags from volume: %v", err)
	}
	if !reflect.DeepEqual(tagmap0c.Tags, tagmap0a) {
		t.Errorf("Tags aren't equal after set/get: %v != %v", tagmap0c.Tags, tagmap0a)
	}
	tools.PrintResource(t, tagmap0a)

	//err = DeleteVolumeTags(t, client, "volumes", volume.ID)
	newVolume, err := volumes.Get(client, volume.ID).Extract()
	if err != nil {
		t.Errorf("Unable to retrieve volume: %v", err)
	}

	tools.PrintResource(t, newVolume)
}

