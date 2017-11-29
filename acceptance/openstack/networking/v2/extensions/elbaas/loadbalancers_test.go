// +build acceptance networking lbaas_v2 loadbalancers

package elbaas

import (
	"testing"
	"os"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	//networking "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/listeners"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/loadbalancer_elbs"
	//"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/backendmember"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/healthcheck"
	//"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	//compute "github.com/gophercloud/gophercloud/acceptance/openstack/compute/v2"
)

func TestLoadbalancersList(t *testing.T) {
	client, err := clients.NewOtcV1Client("elb")
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}
	//fmt.Printf("before  loadbalancer_elbs.List \n")
	allPages, err := loadbalancer_elbs.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list loadbalancers: %v", err)
	}
	// fmt.Printf("after  loadbalancer_elbs.List \n")


	allLoadbalancers, err := loadbalancer_elbs.ExtractLoadBalancers(allPages)
	if err != nil {
		t.Fatalf("Unable to extract loadbalancers: %v", err)
	}

	for _, lb := range allLoadbalancers {
		tools.PrintResource(t, lb)
	}
}

func TestLoadbalancersCRUD(t *testing.T) {
	clientlb, err := clients.NewOtcV1Client("elb")
	if err != nil {
		t.Fatalf("Unable to create an elb client: %v", err)
	}
	tenantID := os.Getenv("OS_TENANT_ID")
	lb, err := CreateLoadBalancer(t, clientlb, "", tenantID, os.Getenv("OS_VPC_ID"), "External")
	if err != nil {
		t.Fatalf("Unable to create loadbalancer: %v", err)
	}
	defer DeleteLoadBalancer(t, clientlb, lb.ID)

	newLB, err := loadbalancer_elbs.Get(clientlb, lb.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get loadbalancer: %v", err)
	}

	tools.PrintResource(t, newLB)

	// Because of the time it takes to create a loadbalancer,
	// this test will include some other resources.

	// Listener
	listener, err := CreateListener(t, clientlb, lb)
	if err != nil {
		t.Fatalf("Unable to create listener: %v", err)
	}
	//fmt.Printf("Listener created: %+v.\n", listener)
	defer DeleteListener(t, clientlb, listener.ID)

	updateListenerOpts := listeners.UpdateOpts{
		Description: "Some listener description",
	}
	_, err = listeners.Update(clientlb, listener.ID, updateListenerOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update listener")
	}

	newListener, err := listeners.Get(clientlb, listener.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get listener")
	}

	tools.PrintResource(t, newListener)

	// Health check
	health, err := CreateHealth(t, clientlb, lb, listener)
	if err != nil {
		t.Fatalf("Unable to create health: %v", err)
	}
	//fmt.Printf("######   HEALTH before DeleteHealth !!!! lb=%v+ health=%v+ \n", lb, health)
	defer DeleteHealth(t, clientlb, lb.ID, health.ID)

	newInterval:= tools.RandomInt(1, 5)
	updateHealthOpts := healthcheck.UpdateOpts{
		HealthcheckInterval: newInterval,
	}
	_, err = healthcheck.Update(clientlb, health.ID, updateHealthOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update health")
	}

	newHealth, err := healthcheck.Get(clientlb, health.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get health")
	}

	tools.PrintResource(t, newHealth)

	backend, err := AddBackend(t, clientlb, lb, listener, os.Getenv("OS_SERVER_ID"), os.Getenv("OS_SERVER_ADDRESS"))
	if err != nil {
		t.Fatalf("Unable to create backend: %v", err)
	}
	//fmt.Printf("backend: %+v.\n", backend)
	tools.PrintResource(t, backend)

	//fmt.Printf("######   BackEnd before DeleteBackend !!!! lb=%v+ backend=%v+ \n", lb, backend)
	RemoveBackend(t, clientlb, lb.ID, backend.ID)

}