// +build acceptance networking lbaas_v2 loadbalancers

package elbaas

import (
	"testing"
	"os"
    "fmt"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	//networking "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	//"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/listeners"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/loadbalancer_elbs"
	//"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/monitors"
)

func TestLoadbalancersList(t *testing.T) {
	client, err := clients.NewOtcV1Client("elb")
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}
    fmt.Printf("before  loadbalancer_elbs.List \n")
	allPages, err := loadbalancer_elbs.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list loadbalancers: %v", err)
	}
    fmt.Printf("after  loadbalancer_elbs.List \n")
	

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
	/*listener, err := CreateListener(t, clientlb, lb)
	if err != nil {
		t.Fatalf("Unable to create listener: %v", err)
	}
	defer DeleteListener(t, clientlb, lb.ID, listener.ID)

	updateListenerOpts := listeners.UpdateOpts{
		Description: "Some listener description",
	}
	_, err = listeners.Update(clientlb, listener.ID, updateListenerOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update listener")
	}

	if err := WaitForLoadBalancerState(clientlb, lb.ID, 1, loadbalancerActiveTimeoutSeconds); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newListener, err := listeners.Get(clientlb, listener.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get listener")
	}

	tools.PrintResource(t, newListener)
	*/

	/*
	// Pool

	pool, err := CreatePool(t, client, lb)
	if err != nil {
		t.Fatalf("Unable to create pool: %v", err)
	}
	defer DeletePool(t, client, lb.ID, pool.ID)

	updatePoolOpts := pools.UpdateOpts{
		Description: "Some pool description",
	}
	_, err = pools.Update(client, pool.ID, updatePoolOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update pool")
	}

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE", loadbalancerActiveTimeoutSeconds); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newPool, err := pools.Get(client, pool.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get pool")
	}

	tools.PrintResource(t, newPool)

	// Member
	member, err := CreateMember(t, client, lb, newPool, subnet.ID, subnet.CIDR)
	if err != nil {
		t.Fatalf("Unable to create member: %v", err)
	}
	defer DeleteMember(t, client, lb.ID, pool.ID, member.ID)

	newWeight := tools.RandomInt(11, 100)
	updateMemberOpts := pools.UpdateMemberOpts{
		Weight: newWeight,
	}
	_, err = pools.UpdateMember(client, pool.ID, member.ID, updateMemberOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update pool")
	}

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE", loadbalancerActiveTimeoutSeconds); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newMember, err := pools.GetMember(client, pool.ID, member.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get member")
	}

	tools.PrintResource(t, newMember)

	// Monitor
	monitor, err := CreateMonitor(t, client, lb, newPool)
	if err != nil {
		t.Fatalf("Unable to create monitor: %v", err)
	}
	defer DeleteMonitor(t, client, lb.ID, monitor.ID)

	newDelay := tools.RandomInt(20, 30)
	updateMonitorOpts := monitors.UpdateOpts{
		Delay: newDelay,
	}
	_, err = monitors.Update(client, monitor.ID, updateMonitorOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update monitor")
	}

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE", loadbalancerActiveTimeoutSeconds); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	newMonitor, err := monitors.Get(client, monitor.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get monitor")
	}

	tools.PrintResource(t, newMonitor)
	*/

}
