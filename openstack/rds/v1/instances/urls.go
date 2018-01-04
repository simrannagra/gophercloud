package instances

import "github.com/gophercloud/gophercloud"

func createURL(c *gophercloud.ServiceClient1) string {
	return c.ServiceURL("instances")
}

func deleteURL(c *gophercloud.ServiceClient1, id string) string {
	return c.ServiceURL("instances", id)
}

func getURL(c *gophercloud.ServiceClient1, id string) string {
	return c.ServiceURL("instances", id)
}

func listURL(c *gophercloud.ServiceClient1) string {
	return c.ServiceURL("instances")
}

func updateURL(c *gophercloud.ServiceClient1, id string) string {
	return c.ServiceURL("instances", id, "action")
}
