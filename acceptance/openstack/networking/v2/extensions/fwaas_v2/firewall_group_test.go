// +build acceptance networking fwaas_v2

package fwaas_v2

import (
	"testing"
	//"fmt"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2/extensions/layer3"
	//"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	networking "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2"
	//compute "github.com/gophercloud/gophercloud/acceptance/openstack/compute/v2"
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

	rule, err := CreateRule(t, client)
	if err != nil {
		t.Fatalf("Unable to create rule: %v", err)
	}
	defer DeleteRule(t, client, rule.ID)
	//fmt.Printf("CreateRule finished, rule=%+v.\n", rule)

	tools.PrintResource(t, rule)

	policy, err := CreatePolicy(t, client, rule.ID)
	if err != nil {
		t.Fatalf("Unable to create policy: %v", err)
	}
	defer DeletePolicy(t, client, policy.ID)
	//fmt.Printf("CreatePolicy finished, policy=%+v.\n", policy)

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

// Problems on OTC
func TestFirewallGroupCRUDPort(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	/*choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	netid, err := networks.IDFromName(client, choices.NetworkName)
	if err != nil {
		t.Fatalf("Unable to find network id: %v", err)
	}*/

	// Create a network
	network, err := networking.CreateNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create network: %v", err)
	}
	defer networking.DeleteNetwork(t, client, network.ID)

	//tools.PrintResource(t, network)
	subnet, err := networking.CreateSubnet(t, client, network.ID)
	if err != nil {
		t.Fatalf("Unable to create subnet: %v", err)
	}
	defer networking.DeleteSubnet(t, client, subnet.ID)

	router, err := layer3.CreateExternalRouter(t, client)
	if err != nil {
		t.Fatalf("Unable to create router: %v", err)
	}
	defer layer3.DeleteRouter(t, client, router.ID)

	port, err := networking.CreatePort(t, client, network.ID, subnet.ID)
	if err != nil {
		t.Fatalf("Unable to create port: %v", err)
	}

	_, err = layer3.CreateRouterInterface(t, client, port.ID, router.ID)
	if err != nil {
		t.Fatalf("Unable to create router interface: %v", err)
	}
	defer layer3.DeleteRouterInterface(t, client, port.ID, router.ID)

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

	firewall_group, err := CreateFirewallGroupOnPort(t, client, policy.ID, port.ID)
	if err != nil {
		t.Fatalf("Unable to create firewall group: %v", err)
	}
	defer DeleteFirewallGroup(t, client, firewall_group.ID)

	tools.PrintResource(t, firewall_group)

	// Create second port
	/*port2, err := networking.CreatePort(t, client, network.ID, subnet.ID)
	if err != nil {
		t.Fatalf("Unable to create port 2: %v", err)
	}
	//deleting router interface will delete?
	//defer networking.DeletePort(t, client, port2.ID)

	// And interface it to router
	_, err = layer3.CreateRouterInterface(t, client, port2.ID, router.ID)
	if err != nil {
		t.Fatalf("Unable to create router interface 2: %v", err)
	}
	defer layer3.DeleteRouterInterface(t, client, port2.ID, router.ID)

	firewallGroupUpdateOpts := firewall_groups.UpdateOpts{
		IngressPolicyID:	policy.ID,
		EgressPolicyID:		policy.ID,
		Description: "Some firewall group description",
	}

	updateOpts := routerinsertion.UpdateOptsExt{
		firewallGroupUpdateOpts,
		[]string{port2.ID},
	}

	_, err = firewall_groups.Update(client, firewall_group.ID, updateOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update firewall group: %v", err)
	}

	newFirewallGroup, err := firewall_groups.Get(client, firewall_group.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get firewall group: %v", err)
	}

	tools.PrintResource(t, newFirewallGroup)*/
}

func TestFirewallGroupCRUDRemovePort(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	// Create Network
	network, err := networking.CreateNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create network: %v", err)
	}
	defer networking.DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := networking.CreateSubnet(t, client, network.ID)
	if err != nil {
		t.Fatalf("Unable to create subnet: %v", err)
	}
	defer networking.DeleteSubnet(t, client, subnet.ID)

        router, err := layer3.CreateExternalRouter(t, client)
        if err != nil {
                t.Fatalf("Unable to create router: %v", err)
        }
        defer layer3.DeleteRouter(t, client, router.ID)

	// Create port
	port, err := networking.CreatePort(t, client, network.ID, subnet.ID)
	if err != nil {
		t.Fatalf("Unable to create port: %v", err)
	}
	//defer networking.DeletePort(t, client, port.ID)

        _, err = layer3.CreateRouterInterface(t, client, port.ID, router.ID)
        if err != nil {
                t.Fatalf("Unable to create router interface: %v", err)
        }
        defer layer3.DeleteRouterInterface(t, client, port.ID, router.ID)

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

	firewall_group, err := CreateFirewallGroupOnPort(t, client, policy.ID, port.ID)
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
