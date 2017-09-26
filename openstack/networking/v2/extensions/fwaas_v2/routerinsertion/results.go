package routerinsertion

// FirewallExt is an extension to the base Firewall group object
type FirewallGroupExt struct {
	// RouterIDs are the routers that the firewall is attached to.
	PortIDs []string `json:"ports"`
}
