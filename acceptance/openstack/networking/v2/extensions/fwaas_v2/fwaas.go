package fwaas_v2

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/firewall_groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/policies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/routerinsertion"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/rules"
)

// CreateFirewallGroup will create a Firewall group with a random name and a specified
// policy ID. An error will be returned if the firewall could not be created.
func CreateFirewallGroup(t *testing.T, client *gophercloud.ServiceClient, policyID string) (*firewall_groups.FirewallGroup, error) {
	firewallName := tools.RandomString("TESTACC-", 8)

	t.Logf("Attempting to create firewall group %s", firewallName)

	iTrue := true
	createOpts := firewall_groups.CreateOpts{
		Name:         		firewallName,
		IngressPolicyID:    policyID,
		EgressPolicyID:     policyID,
		AdminStateUp: 		&iTrue,
	}

	firewall_group, err := firewall_groups.Create(client, createOpts).Extract()
	if err != nil {
		return firewall_group, err
	}
	//fmt.Printf("Created firewall_group=%+v.\n", firewall_group)

	/*
	// NOT ACTIVE if not connected to subnet, so don't wait?
	t.Logf("Waiting for firewall group to become active.")
	if err := WaitForFirewallGroupState(client, firewall_group.ID, "ACTIVE", 60); err != nil {
		return firewall_group, err
	} */

	t.Logf("Successfully created firewall group %s", firewallName)

	return firewall_group, nil
}

// CreateFirewallGroupOnPort will create a Firewall group with a random name and a
// specified policy ID attached to a specified Port. An error will be
// returned if the firewall group could not be created.
func CreateFirewallGroupOnPort(t *testing.T, client *gophercloud.ServiceClient, policyID string, portID string) (*firewall_groups.FirewallGroup, error) {
	firewallName := tools.RandomString("TESTACC-", 8)

	t.Logf("Attempting to create firewall group %s", firewallName)

	firewallGroupCreateOpts := firewall_groups.CreateOpts{
		Name:     			firewallName,
		IngressPolicyID:	policyID,
		EgressPolicyID:		policyID,
	}

	createOpts := routerinsertion.CreateOptsExt{
		CreateOptsBuilder: firewallGroupCreateOpts,
		PortIDs:         []string{portID},
	}

	firewall_group, err := firewall_groups.Create(client, createOpts).Extract()
	if err != nil {
		return firewall_group, err
	}

	t.Logf("Waiting for firewall group to become active.")
	if err := WaitForFirewallGroupState(client, firewall_group.ID, "ACTIVE", 60); err != nil {
		return firewall_group, err
	}

	t.Logf("Successfully created firewall group %s", firewallName)

	return firewall_group, nil
}

// CreatePolicy will create a Firewall Policy with a random name and given
// rule. An error will be returned if the rule could not be created.
func CreatePolicy(t *testing.T, client *gophercloud.ServiceClient, ruleID string) (*policies.Policy, error) {
	policyName := tools.RandomString("TESTACC-", 8)

	t.Logf("Attempting to create policy %s", policyName)

	createOpts := policies.CreateOpts{
		Name: policyName,
		Rules: []string{
			ruleID,
		},
	}

	policy, err := policies.Create(client, createOpts).Extract()
	if err != nil {
		return policy, err
	}

	t.Logf("Successfully created policy %s", policyName)

	return policy, nil
}

// CreateRule will create a Firewall Rule with a random source address and
//source port, destination address and port. An error will be returned if
// the rule could not be created.
func CreateRule(t *testing.T, client *gophercloud.ServiceClient) (*rules.Rule, error) {
	ruleName := tools.RandomString("TESTACC-", 8)
	sourceAddress := fmt.Sprintf("192.168.1.%d", tools.RandomInt(1, 100))
	sourcePort := strconv.Itoa(tools.RandomInt(1, 100))
	destinationAddress := fmt.Sprintf("192.168.2.%d", tools.RandomInt(1, 100))
	destinationPort := strconv.Itoa(tools.RandomInt(1, 100))

	t.Logf("Attempting to create rule %s with source %s:%s and destination %s:%s",
		ruleName, sourceAddress, sourcePort, destinationAddress, destinationPort)

	createOpts := rules.CreateOpts{
		Name:                 ruleName,
		Protocol:             rules.ProtocolTCP,
		Action:               "allow",
		SourceIPAddress:      sourceAddress,
		SourcePort:           sourcePort,
		DestinationIPAddress: destinationAddress,
		DestinationPort:      destinationPort,
	}

	rule, err := rules.Create(client, createOpts).Extract()
	if err != nil {
		return rule, err
	}

	t.Logf("Rule %s successfully created", ruleName)

	return rule, nil
}

// DeleteFirewallGroup will delete a firewall group with a specified ID. A fatal error
// will occur if the delete was not successful. This works best when used as
// a deferred function.
func DeleteFirewallGroup(t *testing.T, client *gophercloud.ServiceClient, firewallID string) {
	t.Logf("Attempting to delete firewall group: %s", firewallID)

	err := firewall_groups.Delete(client, firewallID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete firewall group %s: %v", firewallID, err)
	}

	t.Logf("Waiting for firewall group to delete.")
	if err := WaitForFirewallGroupState(client, firewallID, "DELETED", 60); err != nil {
		t.Logf("Unable to delete firewall group: %s", firewallID)
	}

	t.Logf("Firewall group deleted: %s", firewallID)
}

// DeletePolicy will delete a policy with a specified ID. A fatal error will
// occur if the delete was not successful. This works best when used as a
// deferred function.
func DeletePolicy(t *testing.T, client *gophercloud.ServiceClient, policyID string) {
	t.Logf("Attempting to delete policy: %s", policyID)

	err := policies.Delete(client, policyID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete policy %s: %v", policyID, err)
	}

	t.Logf("Deleted policy: %s", policyID)
}

// DeleteRule will delete a rule with a specified ID. A fatal error will occur
// if the delete was not successful. This works best when used as a deferred
// function.
func DeleteRule(t *testing.T, client *gophercloud.ServiceClient, ruleID string) {
	t.Logf("Attempting to delete rule: %s", ruleID)

	err := rules.Delete(client, ruleID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete rule %s: %v", ruleID, err)
	}

	t.Logf("Deleted rule: %s", ruleID)
}

// WaitForFirewallGroupState will wait until a firewall reaches a given state.
func WaitForFirewallGroupState(client *gophercloud.ServiceClient, firewallID, status string, secs int) error {
	return gophercloud.WaitFor(secs, func() (bool, error) {
		current, err := firewall_groups.Get(client, firewallID).Extract()
		if err != nil {
			if httpStatus, ok := err.(gophercloud.ErrDefault404); ok {
				if httpStatus.Actual == 404 {
					if status == "DELETED" {
						return true, nil
					}
				}
			}
			return false, err
		}

		if current.Status == status {
			return true, nil
		}

		return false, nil
	})
}
