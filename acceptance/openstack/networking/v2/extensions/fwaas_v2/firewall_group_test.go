// +build acceptance networking fwaas_v2

package fwaas_v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	layer3 "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2/extensions/layer3"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/firewall_groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/routerinsertion"
)

func TestFirewallGroupList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	allPages, err := firewall_groups.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list firewall groups: %v", err)
	}

	allFirewallGroups, err := firewall_groups.ExtractFirewallGroups(allPages)
	if err != nil {
		t.Fatalf("Unable to extract firewall groups: %v", err)
	}

	for _, firewall_group := range allFirewallGroups {
		tools.PrintResource(t, firewall_group)
	}
}

func TestFirewallGroupCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	router, err := layer3.CreateExternalRouter(t, client)
	if err != nil {
		t.Fatalf("Unable to create router: %v", err)
	}
	defer layer3.DeleteRouter(t, client, router.ID)

	rule, err := CreateRule(t, client)
	if err != nil {
		t.Fatalf("Unable to create rule: %v", err)
	}
	defer DeleteRule(t, client, rule.ID)

	tools.PrintResource(t, rule)

	policy, err := CreatePolicy(t, client, rule.ID)
	if err != nil {
		t.Fatalf("Unable to create policy: %v", err)
	}
	defer DeletePolicy(t, client, policy.ID)

	tools.PrintResource(t, policy)

	firewall_group, err := CreateFirewallGroup(t, client, policy.ID)
	if err != nil {
		t.Fatalf("Unable to create firewall group: %v", err)
	}
	defer DeleteFirewallGroup(t, client, firewall_group.ID)

	tools.PrintResource(t, firewall_group)

	updateOpts := firewall_groups.UpdateOpts{
		IngressPolicyID:	policy.ID,
		EgressPolicyID:		policy.ID,
		Description: "Some firewall group description",
	}

	_, err = firewall_groups.Update(client, firewall_group.ID, updateOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update firewall group: %v", err)
	}

	newFirewallGroup, err := firewall_groups.Get(client, firewall_group.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get firewall group: %v", err)
	}

	tools.PrintResource(t, newFirewallGroup)
}

func TestFirewallGroupCRUDRouter(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	router, err := layer3.CreateExternalRouter(t, client)
	if err != nil {
		t.Fatalf("Unable to create router: %v", err)
	}
	defer layer3.DeleteRouter(t, client, router.ID)

	rule, err := CreateRule(t, client)
	if err != nil {
		t.Fatalf("Unable to create rule: %v", err)
	}
	defer DeleteRule(t, client, rule.ID)

	tools.PrintResource(t, rule)

	policy, err := CreatePolicy(t, client, rule.ID)
	if err != nil {
		t.Fatalf("Unable to create policy: %v", err)
	}
	defer DeletePolicy(t, client, policy.ID)

	tools.PrintResource(t, policy)

	firewall_group, err := CreateFirewallGroupOnRouter(t, client, policy.ID, router.ID)
	if err != nil {
		t.Fatalf("Unable to create firewall group: %v", err)
	}
	defer DeleteFirewallGroup(t, client, firewall_group.ID)

	tools.PrintResource(t, firewall_group)

	router2, err := layer3.CreateExternalRouter(t, client)
	if err != nil {
		t.Fatalf("Unable to create router: %v", err)
	}
	defer layer3.DeleteRouter(t, client, router2.ID)

	firewallGroupUpdateOpts := firewall_groups.UpdateOpts{
		IngressPolicyID:	policy.ID,
		EgressPolicyID:		policy.ID,
		Description: "Some firewall group description",
	}

	updateOpts := routerinsertion.UpdateOptsExt{
		firewallGroupUpdateOpts,
		[]string{router2.ID},
	}

	_, err = firewall_groups.Update(client, firewall_group.ID, updateOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update firewall group: %v", err)
	}

	newFirewallGroup, err := firewall_groups.Get(client, firewall_group.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get firewall group: %v", err)
	}

	tools.PrintResource(t, newFirewallGroup)
}

func TestFirewallGroupCRUDRemoveRouter(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	router, err := layer3.CreateExternalRouter(t, client)
	if err != nil {
		t.Fatalf("Unable to create router: %v", err)
	}
	defer layer3.DeleteRouter(t, client, router.ID)

	rule, err := CreateRule(t, client)
	if err != nil {
		t.Fatalf("Unable to create rule: %v", err)
	}
	defer DeleteRule(t, client, rule.ID)

	tools.PrintResource(t, rule)

	policy, err := CreatePolicy(t, client, rule.ID)
	if err != nil {
		t.Fatalf("Unable to create policy: %v", err)
	}
	defer DeletePolicy(t, client, policy.ID)

	tools.PrintResource(t, policy)

	firewall_group, err := CreateFirewallGroupOnRouter(t, client, policy.ID, router.ID)
	if err != nil {
		t.Fatalf("Unable to create firewall group: %v", err)
	}
	defer DeleteFirewallGroup(t, client, firewall_group.ID)

	tools.PrintResource(t, firewall_group)

	firewallGroupUpdateOpts := firewall_groups.UpdateOpts{
		IngressPolicyID:    policy.ID,
		EgressPolicyID:		policy.ID,
		Description: "Some firewall group description",
	}

	updateOpts := routerinsertion.UpdateOptsExt{
		firewallGroupUpdateOpts,
		[]string{},
	}

	_, err = firewall_groups.Update(client, firewall_group.ID, updateOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update firewall group: %v", err)
	}

	newFirewallGroup, err := firewall_groups.Get(client, firewall_group.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get firewall group: %v", err)
	}

	tools.PrintResource(t, newFirewallGroup)
}
