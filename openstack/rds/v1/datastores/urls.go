package datastores

import "github.com/gophercloud/gophercloud"

func listURL(c *gophercloud.ServiceClient1, dataStoreName string) string {
	return c.ServiceURL("datastores", dataStoreName, "versions")
}
