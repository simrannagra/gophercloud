package elbaas

import (
	"fmt"
	// "strings"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/listeners"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/loadbalancer_elbs"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/monitors"
    "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/healthcheck"
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
	fmt.Printf("Extracted listener: %+v.\n", listener)
	if err != nil {

        t.Logf("Attempting to create listener %s on port %d failed err=%v", listenerName, listenerPort, err)
        return listener, err
	}

	t.Logf("Successfully created listener %s", listenerName)

	if err := WaitForLoadBalancerState(client, lb.ID, 1, loadbalancerActiveTimeoutSeconds); err != nil {
		return listener, fmt.Errorf("Timed out waiting for loadbalancer to become active")
	}

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

	fmt.Printf("job=%+v.\n", job)

	t.Logf("Successfully created loadbalancer %s on subnet %s", lbName, subnetID)
	t.Logf("Waiting for loadbalancer %s to become active", lbName)

	if err := WaitForJobSuccess(client, job.URI, loadbalancerActiveTimeoutSeconds); err != nil {
		return nil, err
	}

	mlb, err := GetJobEntity(client, job.URI,"elb")
	fmt.Printf("mlb=%+v.\n", mlb)
	t.Logf("LoadBalancer %s is active", lbName)

	if vid, ok := mlb["id"]; ok {
		fmt.Printf("vid=%s.\n", vid)
		if id, ok := vid.(string); ok {
			fmt.Printf("id=%s.\n", id)
			lb, err := loadbalancer_elbs.Get(client, id).Extract()
			if err != nil {
				fmt.Printf("Error: %s.\n", err.Error())
				return nil, err
			}
			fmt.Printf("lb=%+v.\n", lb)
			return lb, err
		}
	}

	return nil, err
}

// CreateMember will create a member with a random name, port, address, and
// weight. An error will be returned if the member could not be created.
func CreateMember(t *testing.T, client *gophercloud.ServiceClient, lb *loadbalancer_elbs.LoadBalancer, subnetID, subnetCIDR string) ( error) {
	return nil
	/*
	memberName := tools.RandomString("TESTACCT-", 8)
	memberPort := tools.RandomInt(100, 1000)
	memberWeight := tools.RandomInt(1, 10)

	cidrParts := strings.Split(subnetCIDR, "/")
	subnetParts := strings.Split(cidrParts[0], ".")
	memberAddress := fmt.Sprintf("%s.%s.%s.%d", subnetParts[0], subnetParts[1], subnetParts[2], tools.RandomInt(10, 100))


	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE", loadbalancerActiveTimeoutSeconds); err != nil {
		return member, fmt.Errorf("Timed out waiting for loadbalancer to become active")
	}

	return member, nil
	*/
}

// CreateMonitor will create a monitor with a random name for a specific pool.
// An error will be returned if the monitor could not be created.
func CreateMonitor(t *testing.T, client *gophercloud.ServiceClient, lb *loadbalancer_elbs.LoadBalancer/* pool *pools.Pool*/) (*monitors.Monitor, error) {
	monitorName := tools.RandomString("TESTACCT-", 8)

	t.Logf("Attempting to create monitor %s", monitorName)

	createOpts := monitors.CreateOpts{
		//PoolID:     pool.ID,
		Name:       monitorName,
		Delay:      10,
		Timeout:    5,
		MaxRetries: 5,
		Type:       "PING",
	}

	monitor, err := monitors.Create(client, createOpts).Extract()
	if err != nil {
		return monitor, err
	}

	t.Logf("Successfully created monitor: %s", monitorName)

	if err := WaitForLoadBalancerState(client, lb.ID, 1, loadbalancerActiveTimeoutSeconds); err != nil {
		return monitor, fmt.Errorf("Timed out waiting for loadbalancer to become active")
	}

	return monitor, nil
}

// CreatePool will create a pool with a random name with a specified listener
// and loadbalancer. An error will be returned if the pool could not be
// created.
func CreatePool(t *testing.T, client *gophercloud.ServiceClient, lb *loadbalancer_elbs.LoadBalancer) (/**pools.Pool, */error) {
	return nil
/*
	poolName := tools.RandomString("TESTACCT-", 8)
    t.Logf("Attempting to create pool %s", poolName)

	createOpts := pools.CreateOpts{
		Name:           poolName,
		Protocol:       pools.ProtocolTCP,
		LoadbalancerID: lb.ID,
		LBMethod:       pools.LBMethodLeastConnections,
	}

	pool, err := pools.Create(client, createOpts).Extract()
	if err != nil {
		return pool, err
	}

	t.Logf("Successfully created pool %s", poolName)

	if err := WaitForLoadBalancerState(client, lb.ID, "ACTIVE", loadbalancerActiveTimeoutSeconds); err != nil {
		return pool, fmt.Errorf("Timed out waiting for loadbalancer to become active")
	}

	return pool, nil
*/
}

// CreateMonitor will create a monitor with a random name for a specific pool.
// An error will be returned if the monitor could not be created.
func CreateHealth(t *testing.T, client *gophercloud.ServiceClient, lb *loadbalancer_elbs.LoadBalancer, listener *listeners.Listener) (*healthcheck.Health, error) {
	healthName := tools.RandomString("TESTACCT-", 8)

	t.Logf("Attempting to create health %s", healthName)

    fmt.Printf("######    before  health.CreateOpts listener.ID=%v+  \n", listener.ID)

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
    
    fmt.Printf("#######    after  health.CreateOpts %v+ \n", createOpts)

	health, err := healthcheck.Create(client, createOpts).Extract()
	if err != nil {
		return health, err
	}

	t.Logf("Successfully created health: %s", healthName)

	if err := WaitForLoadBalancerState(client, lb.ID, 1, loadbalancerActiveTimeoutSeconds); err != nil {
		return health, fmt.Errorf("Timed out waiting for loadbalancer to become active")
	}

	return health, nil
}

// DeleteListener will delete a specified listener. A fatal error will occur if
// the listener could not be deleted. This works best when used as a deferred
// function.
func DeleteListener(t *testing.T, client *gophercloud.ServiceClient, lbID, listenerID string) {
	t.Logf("Attempting to delete listener %s", listenerID)

	if err := listeners.Delete(client, listenerID).ExtractErr(); err != nil {
		t.Fatalf("Unable to delete listener: %v", err)
	}

	if err := WaitForLoadBalancerState(client, lbID, 1, loadbalancerActiveTimeoutSeconds); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	t.Logf("Successfully deleted listener %s", listenerID)
}

// DeleteMember will delete a specified member. A fatal error will occur if the
// member could not be deleted. This works best when used as a deferred
// function.
func DeleteMember(t *testing.T, client *gophercloud.ServiceClient, lbID, poolID, memberID string) {
	/*
	t.Logf("Attempting to delete member %s", memberID)


	if err := pools.DeleteMember(client, poolID, memberID).ExtractErr(); err != nil {
		t.Fatalf("Unable to delete member: %s", memberID)
	}

	if err := WaitForLoadBalancerState(client, lbID, "ACTIVE", loadbalancerActiveTimeoutSeconds); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	t.Logf("Successfully deleted member %s", memberID)
	*/
}

// DeleteLoadBalancer will delete a specified loadbalancer. A fatal error will
// occur if the loadbalancer could not be deleted. This works best when used
// as a deferred function.
func DeleteLoadBalancer(t *testing.T, client *gophercloud.ServiceClient, lbID string) {
	t.Logf("Attempting to delete loadbalancer %s", lbID)

	if err := loadbalancer_elbs.Delete(client, lbID).ExtractErr(); err != nil {
		t.Fatalf("Unable to delete loadbalancer: %v", err)
	}

	t.Logf("Waiting for loadbalancer %s to delete", lbID)

	//if err := WaitForLoadBalancerState(client, lbID, "DELETED", loadbalancerActiveTimeoutSeconds); err != nil {
	//	t.Fatalf("Loadbalancer did not delete in time.")
	//}

	t.Logf("Successfully deleted loadbalancer %s", lbID)
}

// DeleteMonitor will delete a specified monitor. A fatal error will occur if
// the monitor could not be deleted. This works best when used as a deferred
// function.
func DeleteMonitor(t *testing.T, client *gophercloud.ServiceClient, lbID, monitorID string) {
	t.Logf("Attempting to delete monitor %s", monitorID)

	if err := monitors.Delete(client, monitorID).ExtractErr(); err != nil {
		t.Fatalf("Unable to delete monitor: %v", err)
	}

	if err := WaitForLoadBalancerState(client, lbID, 1, loadbalancerActiveTimeoutSeconds); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	t.Logf("Successfully deleted monitor %s", monitorID)
}

// DeleteHealth will delete a specified monitor. A fatal error will occur if
// the monitor could not be deleted. This works best when used as a deferred
// function.
func DeleteHealth(t *testing.T, client *gophercloud.ServiceClient, lbID, healthID string) {
	t.Logf("Attempting to delete health %s", healthID)

	if err := healthcheck.Delete(client, healthID).ExtractErr(); err != nil {
		t.Fatalf("Unable to delete health: %v", err)
	}

	if err := WaitForLoadBalancerState(client, lbID, 1, loadbalancerActiveTimeoutSeconds); err != nil {
		t.Fatalf("Timed out waiting for loadbalancer to become active")
	}

	t.Logf("Successfully deleted health %s", healthID)
}

func WaitForJobSuccess(client *gophercloud.ServiceClient, uri string, secs int) error {
	return gophercloud.WaitFor(secs, func() (bool, error) {
		job := new(loadbalancer_elbs.JobStatus)
		_, err := client.Get("https://elb.eu-de.otc.t-systems.com" + uri, &job, nil)
		if err != nil {
			return false, err
		}
		fmt.Printf("JobStatus: %+v.\n", job)

		if job.Status == "SUCCESS" {
			return true, nil
		}
		if job.Status == "FAIL" {
			err = fmt.Errorf("Job failed with code %s: %s.\n", job.ErrorCode, job.FailReason)
			return false, err
		}

		return false, nil
	})
}

func GetJobEntity(client *gophercloud.ServiceClient, uri string, label string) (map[string]interface{}, error) {
	job := new(loadbalancer_elbs.JobStatus)
	_, err := client.Get("https://elb.eu-de.otc.t-systems.com" + uri, &job, nil)
	if err != nil {
		return nil, err
	}
	fmt.Printf("JobStatus: %+v.\n", job)

	if job.Status == "SUCCESS" {
		if e := job.Entities[label]; e != nil {
			if m, ok := e.(map[string]interface{}); ok {
				return m, nil
			}
		}
	}

	return nil, nil
}

// WaitForLoadBalancerState will wait until a loadbalancer reaches a given state.
func WaitForLoadBalancerState(client *gophercloud.ServiceClient, lbID string, status int, secs int) error {
	return gophercloud.WaitFor(secs, func() (bool, error) {
		current, err := loadbalancer_elbs.Get(client, lbID).Extract()
		if err != nil {
			if httpStatus, ok := err.(gophercloud.ErrDefault404); ok {
				if httpStatus.Actual == 404 {
					//if status == "DELETED" {
					//	return true, nil
					//}
				}
			}
			return false, err
		}

		if current.AdminStateUp == status {
			return true, nil
		}

		return false, nil
	})
}
