package flavors

import "github.com/gophercloud/gophercloud"

func listURL(c *gophercloud.ServiceClient1, dataStoreID string, region string) string {

	return c.ResourceBaseURL() + "flavors?dbId=" + dataStoreID + "&region=" + region
}
