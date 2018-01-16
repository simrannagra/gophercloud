package common

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

const TokenID = client.TokenID

// Fake project id to use.
const ProjectID = "85636478b0bd8e67e89469c7749d4127"

func ServiceClient() *gophercloud.ServiceClient {
	sc := client.ServiceClient()
	sc.ResourceBase = sc.Endpoint + "v1/" + ProjectID + "/"
	return sc
}
