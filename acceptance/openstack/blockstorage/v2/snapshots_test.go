// +build acceptance blockstorage

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/snapshots"
	"reflect"
)

func TestSnapshotsList(t *testing.T) {
	client, err := clients.NewBlockStorageV2Client()
	if err != nil {
		t.Fatalf("Unable to create a blockstorage client: %v", err)
	}

	allPages, err := snapshots.List(client, snapshots.ListOpts{}).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve snapshots: %v", err)
	}

	allSnapshots, err := snapshots.ExtractSnapshots(allPages)
	if err != nil {
		t.Fatalf("Unable to extract snapshots: %v", err)
	}

	for _, snapshot := range allSnapshots {
		tools.PrintResource(t, snapshot)
	}
}

func TestSnapshotsCreateDelete(t *testing.T) {
	client, err := clients.NewBlockStorageV2Client()
	if err != nil {
		t.Fatalf("Unable to create a blockstorage client: %v", err)
	}

	volume, err := CreateVolume(t, client)
	if err != nil {
		t.Fatalf("Unable to create volume: %v", err)
	}
	defer DeleteVolume(t, client, volume)

	snapshot, err := CreateSnapshot(t, client, volume)
	if err != nil {
		t.Fatalf("Unable to create snapshot: %v", err)
	}
	defer DeleteSnapshot(t, client, snapshot)

	newSnapshot, err := snapshots.Get(client, snapshot.ID).Extract()
	if err != nil {
		t.Errorf("Unable to retrieve snapshot: %v", err)
	}

	tools.PrintResource(t, newSnapshot)
}

func TestSnapshotsTags(t *testing.T) {
	client, err := clients.NewBlockStorageV2Client()
	if err != nil {
		t.Fatalf("Unable to create blockstorage client: %v", err)
	}

	volume, err := CreateVolume(t, client)
	if err != nil {
		t.Fatalf("Unable to create volume: %v", err)
	}
	defer DeleteVolume(t, client, volume)

	snapshot, err := CreateSnapshot(t, client, volume)
	if err != nil {
		t.Fatalf("Unable to create snapshot: %v", err)
	}
	defer DeleteSnapshot(t, client, snapshot)

	tagmap := map[string]string{"foo" : "bar", "name" : "value"}
	tagmap2, err := CreateVolumeTags(t, client, "snapshots", snapshot.ID, tagmap)
	if err != nil {
		t.Errorf("Unable to create tags for snapshot: %v", err)
	}
	tagmap3, err := GetVolumeTags(t, client, "snapshots", snapshot.ID)
	if err != nil {
		t.Errorf("Unable to get tags from snapshot: %v", err)
	}
	if !reflect.DeepEqual(tagmap3.Tags, tagmap) {
		t.Errorf("Tags aren't equal after set/get: %v != %v", tagmap3.Tags, tagmap)
	}
	tools.PrintResource(t, tagmap2)

	tagmap4 := map[string]string{"foo2" : "bar2", "name2" : "value2"}
	tagmap5, err := CreateVolumeTags(t, client, "snapshots", snapshot.ID, tagmap4)
	if err != nil {
		t.Errorf("Unable to create tags for snapshot: %v", err)
	}
	tagmap6, err := GetVolumeTags(t, client, "snapshots", snapshot.ID)
	if err != nil {
		t.Errorf("Unable to get tags from snapshot: %v", err)
	}
	if !reflect.DeepEqual(tagmap6.Tags, tagmap4) {
		t.Errorf("Tags aren't equal after set/get: %v != %v", tagmap6.Tags, tagmap4)
	}
	tools.PrintResource(t, tagmap5)

	tagmap0a := map[string]string{}
	tagmap0b, err := CreateVolumeTags(t, client, "snapshots", snapshot.ID, tagmap0a)
	if err != nil {
		t.Errorf("Unable to delete tags for snapshot: %v", err)
	}
	tagmap0c, err := GetVolumeTags(t, client, "snapshots", snapshot.ID)
	if err != nil {
		t.Errorf("Unable to get empty tags from snapshot: %v", err)
	}
	if !reflect.DeepEqual(tagmap0c.Tags, tagmap0a) {
		t.Errorf("Tags aren't equal after set/get: %v != %v", tagmap0c.Tags, tagmap0a)
	}
	tools.PrintResource(t, tagmap0b)

	newSnapshot, err := snapshots.Get(client, snapshot.ID).Extract()
	if err != nil {
		t.Errorf("Unable to retrieve snapshot: %v", err)
	}

	tools.PrintResource(t, newSnapshot)
}
