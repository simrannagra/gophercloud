package elbaas

import (
	//"fmt"
	// "strings"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/listeners"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/loadbalancer_elbs"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/backendmember"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/healthcheck"
	//"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/quota"
	"fmt"
)

const loadbalancerActiveTimeoutSeconds = 300
const loadbalancerDeleteTimeoutSeconds = 300

// CreateListener will create a listener for a given load balancer on a random
// port with a random name. An error will be returned if the listener could not
// be created.
func CreateListener(t *testing.T, client *gophercloud.ServiceClient, lb *loadbalancer_elbs.LoadBalancer) (*listeners.Listener, error) {
	listenerName := tools.RandomString("TESTACCT-", 8)
	listenerPort := tools.RandomInt(1, 100)

	t.Logf("Attempting to create listener %s on port %d", listenerName, listenerPort)

	// fmt.Printf("*******    before  listeners.CreateOpts  \n")

	createOpts := listeners.CreateOpts{
		Name:           listenerName,
		LoadbalancerID: lb.ID,
		Protocol:       "TCP",
		ProtocolPort:   listenerPort,
		BackendProtocol: "TCP",
		BackendProtocolPort: listenerPort,
		Algorithm:		 "roundrobin",
	}
	// fmt.Printf("*******    after  listeners.CreateOpts %v+ \n", createOpts)

	listener, err := listeners.Create(client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Successfully created listener %s", listener.ID)

	return listener, nil
}

// CreateLoadBalancer will create a load balancer with a random name on a given
// subnet. An error will be returned if the loadbalancer could not be created.
func CreateLoadBalancer(t *testing.T, client *gophercloud.ServiceClient, subnetID string, tenantID string, vpcID string, lb_type string) (*loadbalancer_elbs.LoadBalancer, error) {
	lbName := tools.RandomString("TESTACCT-", 8)

	t.Logf("Attempting to create loadbalancer %s on subnet %s", lbName, subnetID)

	createOpts := loadbalancer_elbs.CreateOpts{
		//Tenant_ID: 	  tenantID,
		VpcID: 		  vpcID,
		Name:         lbName,
		Bandwidth:	  5,	// Must not be passed in for Internal
		Type:		  lb_type,
		VipSubnetID:  subnetID,	// Must be blank for External, required for Internal
		AdminStateUp: gophercloud.Enabled,
	}

	job, err := loadbalancer_elbs.Create(client, createOpts).ExtractJobResponse()
	if err != nil {
		return nil, err
	}

	//fmt.Printf("job=%+v.\n", job)
	t.Logf("Waiting for loadbalancer %s to become active", lbName)

	if err := gophercloud.WaitForJobSuccess(client, job.URI, loadbalancerActiveTimeoutSeconds); err != nil {
		return nil, err
	}

	entity, err := gophercloud.GetJobEntity(client, job.URI,"elb")
	//fmt.Printf("mlb=%+v.\n", mlb)
	t.Logf("LoadBalancer %s is active", lbName)

	if mlb, ok := entity.(map[string]interface{}); ok {
		if vid, ok := mlb["id"]; ok {
			//fmt.Printf("vid=%s.\n", vid)
			if id, ok := vid.(string); ok {
				//fmt.Printf("id=%s.\n", id)
				lb, err := loadbalancer_elbs.Get(client, id).Extract()
				if err != nil {
					//fmt.Printf("Error: %s.\n", err.Error())
					return nil, err
				}
				//fmt.Printf("lb=%+v.\n", lb)
				return lb, err
			}
		}
	}

	return nil, fmt.Errorf("Unexpected conversion error in CreateLoadBalancer.")
}

// CreateHealth will create a monitor with a random name for a specific pool.
// An error will be returned if the monitor could not be created.
func CreateHealth(t *testing.T, client *gophercloud.ServiceClient, lb *loadbalancer_elbs.LoadBalancer, listener *listeners.Listener) (*healthcheck.Health, error) {
	//fmt.Printf("######    before  health.CreateOpts listener.ID=%v+  \n", listener.ID)

	createOpts := healthcheck.CreateOpts{
		HealthcheckConnectPort:  80,
		HealthcheckInterval: 5,
		HealthcheckProtocol: "HTTP",
		HealthcheckTimeout: 10,
		HealthcheckUri: "/",
		HealthyThreshold: 3,
		ListenerID: listener.ID,
		UnhealthyThreshold: 3,
	}

	//fmt.Printf("#######    after  health.CreateOpts %v+ \n", createOpts)

	health, err := healthcheck.Create(client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Successfully created healthcheck %s.", health.ID)

	//return health, nil
	return health, nil
}

// CreateBackend will create a listener backend for a given load balancer on a random
// port with a random name. An error will be returned if the listener could not
// be created.
func AddBackend(t *testing.T, client *gophercloud.ServiceClient, lb *loadbalancer_elbs.LoadBalancer, listener *listeners.Listener, server_id string, address string) (*backendmember.Backend, error) {
	addOpts := backendmember.AddOpts{
		ServerId: server_id,
		Address:   address,
	}
	//fmt.Printf("*******    after  backendmember.AddOpts %v+ \n", addOpts)

	job, err := backendmember.Add(client, listener.ID, addOpts).ExtractJobResponse()
	if err != nil {
		return nil, err
	}

	//fmt.Printf("job=%+v.\n", job)
	t.Logf("Waiting for backend to become active")

	if err := gophercloud.WaitForJobSuccess(client, job.URI, loadbalancerActiveTimeoutSeconds); err != nil {
		return nil, err
	}

	entity, err := gophercloud.GetJobEntity(client, job.URI,"members")
	//fmt.Printf("mlb=%+v.\n", mlb)
	t.Logf("Backend for listener %s, lb %s is active", listener.ID, lb.ID)

	if members, ok := entity.([]interface{}); ok {
		if len(members) > 0 {
			vmember := members[0]
			if member, ok := vmember.(map[string]interface{}); ok {
				//return member, nil
				if vid, ok := member["id"]; ok {
					//fmt.Printf("vid=%s.\n", vid)
					if id, ok := vid.(string); ok {
						//fmt.Printf("id=%s.\n", id)
						backend, err := backendmember.Get(client, listener.ID, id).Extract()
						if err != nil {
							//fmt.Printf("Error: %s.\n", err.Error())
							return nil, err
						}
						//fmt.Printf("lb=%+v.\n", lb)
						return backend, err
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("Unexpected conversion error in AddBackend.")
}

// DeleteListener will delete a specified listener. A fatal error will occur if
// the listener could not be deleted. This works best when used as a deferred
// function.
func DeleteListener(t *testing.T, client *gophercloud.ServiceClient, id string) {
	t.Logf("Attempting to delete listener %s", id)

	err := listeners.Delete(client, id).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete listner: %v", err)
	}

	t.Logf("Successfully deleted listener %s", id)
}


// DeleteLoadBalancer will delete a specified loadbalancer. A fatal error will
// occur if the loadbalancer could not be deleted. This works best when used
// as a deferred function.
func DeleteLoadBalancer(t *testing.T, client *gophercloud.ServiceClient, lbID string) {
	t.Logf("Attempting to delete loadbalancer %s", lbID)

	job, err := loadbalancer_elbs.Delete(client, lbID).ExtractJobResponse()
	//fmt.Printf("delete job: %+v.\n", job)
	if err != nil {
		t.Fatalf("Unable to delete loadbalancer: %v", err)
	}

	t.Logf("Waiting for loadbalancer %s to delete", lbID)

	if err := gophercloud.WaitForJobSuccess(client, job.URI, loadbalancerActiveTimeoutSeconds); err != nil {
		t.Fatalf("Loadbalancer did not delete in time.")
	}

	t.Logf("Successfully deleted loadbalancer %s", lbID)
}

// DeleteHealth will delete a specified monitor. A fatal error will occur if
// the monitor could not be deleted. This works best when used as a deferred
// function.
func DeleteHealth(t *testing.T, client *gophercloud.ServiceClient, lbID, healthID string) {
	t.Logf("Attempting to delete health %s", healthID)

	if err := healthcheck.Delete(client, healthID).ExtractErr(); err != nil {
		t.Fatalf("Unable to delete health: %v", err)
	}

	t.Logf("Successfully deleted health %s", healthID)
}

// DeleteBackend will delete a specified listener. A fatal error will occur if
// the listener could not be deleted. This works best when used as a deferred
// function.
func RemoveBackend(t *testing.T, client *gophercloud.ServiceClient, listener_id, id string) {
	t.Logf("Attempting to delete backend member %s", id)

	job, err := backendmember.Remove(client, listener_id, id).ExtractJobResponse()
	if err != nil {
		t.Fatalf("Unable to delete backend: %v", err)
	}

	t.Logf("Waiting for backend member %s to delete", id)

	if err := gophercloud.WaitForJobSuccess(client, job.URI, loadbalancerActiveTimeoutSeconds); err != nil {
		t.Fatalf("backend member did not delete in time.")
	}

	t.Logf("Successfully deleted backend member %s", id)
}
