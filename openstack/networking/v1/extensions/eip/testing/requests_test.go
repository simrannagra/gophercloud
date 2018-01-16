package testing

import (
"fmt"
"net/http"
"testing"

fake "github.com/gophercloud/gophercloud/openstack/networking/v1/common"
"github.com/gophercloud/gophercloud/openstack/networking/v1/extensions/eip"
th "github.com/gophercloud/gophercloud/testhelper"
)

func TestListEip(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/85636478b0bd8e67e89469c7749d4127/publicips", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "publicips": [
        {
            "id": "32e58fc8-1438-4a61-9a32-fdffe7284c88",
            "status": "ACTIVE",
            "profile": {
                "user_id": null,
                "product_id": null,
                "region_id": null
            },
            "type": "5_bgp",
            "port_id": "99163903-eb7b-4f1c-ae0f-610db955ab4e",
            "public_ip_address": "80.158.22.223",
            "private_ip_address": "10.0.0.214",
            "tenant_id": "85636478b0bd8e67e89469c7749d4127",
            "create_time": "2018-01-02 17:26:16",
            "bandwidth_id": "b179d110-8b4b-4184-b1b3-e24d32098b1c",
            "bandwidth_name": "ecs-nordea-test-bandwidth-4236",
            "bandwidth_share_type": "PER",
            "bandwidth_size": 5
        },
        {
            "id": "c9939fb4-48ac-4e6a-8ca5-54ba4b1f8105",
            "status": "DOWN",
            "profile": {
                "user_id": null,
                "product_id": null,
                "region_id": null
            },
            "type": "5_bgp",
            "public_ip_address": "80.158.18.202",
            "tenant_id": "85636478b0bd8e67e89469c7749d4127",
            "create_time": "2018-01-10 05:43:36",
            "bandwidth_id": "0cbb9fdb-96d9-4707-890b-43bd3b5cd640",
            "bandwidth_name": "bandwidth-2e9e",
            "bandwidth_share_type": "PER",
            "bandwidth_size": 5
        }
    ]
}
			`)
	})


	/*count := 0

	subnets.List(fake.ServiceClient(), subnets.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := subnets.ExtractSubnets(page)
		if err != nil {
			t.Errorf("Failed to extract subnets: %v", err)
			return false, err
		}

		expected := []subnets.Subnet{
			{
				Status:           "ACTIVE",
				CIDR:             "10.0.1.0/24",
				EnableDHCP:       true,
				Name:             "subnet-perf1",
				ID:               "249c7026-6fd3-4f3a-9613-1456f12f8e08",
				GatewayIP:        "10.0.1.1",
				PRIMARY_DNS:      "100.125.4.25",
				SECONDARY_DNS:    "8.8.8.8",
				VPC_ID:           "d4f2c817-d5df-4a66-994a-6571312b470e",
			},
			{
				Status:           "UNKNOWN",
				CIDR:             "192.168.199.0/24",
				EnableDHCP:       true,
				Name:             "tf_test_subnet",
				ID:               "404b11d4-6869-48c1-a359-da40b6c49dd7",
			},

		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}*/
	actual, err := eip.List(fake.ServiceClient(), eip.ListOpts{})
	if err != nil {
		t.Errorf("Failed to extract eips: %v", err)
	}

	expected := []eip.Eip{
		{
			ID: "32e58fc8-1438-4a61-9a32-fdffe7284c88",
			Status: "ACTIVE",
			Type: "5_bgp",
			PortId: "99163903-eb7b-4f1c-ae0f-610db955ab4e",
			PublicIpAddress: "80.158.22.223",
			PrivateIpAddress: "10.0.0.214",
			TenantID: "85636478b0bd8e67e89469c7749d4127",
			CreateTime: "2018-01-02 17:26:16",
			BandwidthId: "b179d110-8b4b-4184-b1b3-e24d32098b1c",
			Name: "ecs-nordea-test-bandwidth-4236",
			BandwidthShareType: "PER",
			BandwidthSize: 5,
		},
		{
			ID: "c9939fb4-48ac-4e6a-8ca5-54ba4b1f8105",
			Status: "DOWN",
			Type: "5_bgp",
			PublicIpAddress: "80.158.18.202",
			TenantID: "85636478b0bd8e67e89469c7749d4127",
			CreateTime: "2018-01-10 05:43:36",
			BandwidthId: "0cbb9fdb-96d9-4707-890b-43bd3b5cd640",
			Name: "bandwidth-2e9e",
			BandwidthShareType: "PER",
			BandwidthSize: 5,
		},

	}

	th.AssertDeepEquals(t, expected, actual)
}


func TestGetEip(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/85636478b0bd8e67e89469c7749d4127/publicips/32e58fc8-1438-4a61-9a32-fdffe7284c88", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "publicip": {
        "id": "32e58fc8-1438-4a61-9a32-fdffe7284c88",
        "status": "ACTIVE",
        "type": "5_bgp",
        "port_id": "99163903-eb7b-4f1c-ae0f-610db955ab4e",
        "public_ip_address": "80.158.22.223",
        "private_ip_address": "10.0.0.214",
        "tenant_id": "85636478b0bd8e67e89469c7749d4127",
        "create_time": "2018-01-02 17:26:16",
        "bandwidth_id": "b179d110-8b4b-4184-b1b3-e24d32098b1c",
        "bandwidth_share_type": "PER",
        "bandwidth_size": 5
    }
}
		`)
	})

	n, err := eip.Get(fake.ServiceClient(), "32e58fc8-1438-4a61-9a32-fdffe7284c88").Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "32e58fc8-1438-4a61-9a32-fdffe7284c88", n.ID)
	th.AssertEquals(t, "ACTIVE", n.Status)
	th.AssertEquals(t, "5_bgp", n.Type)
	th.AssertEquals(t, "99163903-eb7b-4f1c-ae0f-610db955ab4e", n.PortId)
	th.AssertEquals(t, "80.158.22.223", n.PublicIpAddress)
	th.AssertEquals(t, "10.0.0.214", n.PrivateIpAddress)
	th.AssertEquals(t, "85636478b0bd8e67e89469c7749d4127", n.TenantID)
	th.AssertEquals(t, "2018-01-02 17:26:16", n.CreateTime)
	th.AssertEquals(t, "b179d110-8b4b-4184-b1b3-e24d32098b1c", n.BandwidthId)
	th.AssertEquals(t, "PER", n.BandwidthShareType)
	th.AssertEquals(t, 5, n.BandwidthSize)

}

func TestCreateEip(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/85636478b0bd8e67e89469c7749d4127/publicips", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "publicip": {
        "type": "5_bgp"
    },
    "bandwidth": {
        "name": "bandwidth123",
        "size": 10,
        "share_type": "PER"
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "publicip": {
        "id": "1d5fb72e-cc6f-4c86-b473-bf5f790bec56",
        "status": "PENDING_CREATE",
        "type": "5_bgp",
        "public_ip_address": "80.158.22.68",
        "tenant_id": "85636478b0bd8e67e89469c7749d4127",
        "create_time": "2018-01-15 10:36:38",
        "bandwidth_size": 0
    }
}	`)
	})

	options := eip.CreateOpts{
		PublicIp: eip.Publicip{Type:"5_bgp"},
		BandWidth: eip.Bandwidth{Name:"bandwidth123",Size:10,ShareType:"PER"},
	}
	n, err := eip.Create(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "1d5fb72e-cc6f-4c86-b473-bf5f790bec56", n.ID)
	th.AssertEquals(t, "PENDING_CREATE", n.Status)
	th.AssertEquals(t, "5_bgp", n.Type)
	th.AssertEquals(t, "80.158.22.68", n.PublicIpAddress)
	th.AssertEquals(t, "85636478b0bd8e67e89469c7749d4127", n.TenantID)
	th.AssertEquals(t, "2018-01-15 10:36:38", n.CreateTime)
	th.AssertEquals(t, 0, n.BandwidthSize)


}

func TestUpdateEip(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/85636478b0bd8e67e89469c7749d4127/publicips/f9e87316-9ca3-4fa8-a00d-7428dd619627", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "publicip": {
        "port_id": "99163903-eb7b-4f1c-ae0f-610db955ab4e"
    }
}
`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "publicip": {
        "id": "f9e87316-9ca3-4fa8-a00d-7428dd619627",
        "status": "DOWN",
        "type": "5_bgp",
        "port_id": "99163903-eb7b-4f1c-ae0f-610db955ab4e",
        "public_ip_address": "80.158.23.253",
        "tenant_id": "85636478b0bd8e67e89469c7749d4127",
        "create_time": "2018-01-15 11:35:11",
        "bandwidth_size": 5
    }
}
		`)
	})

	options := eip.UpdateOpts{PortId: "99163903-eb7b-4f1c-ae0f-610db955ab4e"}

	n, err := eip.Update(fake.ServiceClient(), "f9e87316-9ca3-4fa8-a00d-7428dd619627",options).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "f9e87316-9ca3-4fa8-a00d-7428dd619627", n.ID)
	th.AssertEquals(t, "DOWN", n.Status)
	th.AssertEquals(t, "5_bgp", n.Type)
	th.AssertEquals(t, "99163903-eb7b-4f1c-ae0f-610db955ab4e", n.PortId)
	th.AssertEquals(t, "80.158.23.253", n.PublicIpAddress)
	th.AssertEquals(t,"85636478b0bd8e67e89469c7749d4127", n.TenantID)
	th.AssertEquals(t, "2018-01-15 11:35:11", n.CreateTime)
	th.AssertEquals(t, 5, n.BandwidthSize)
}

func TestDeleteEip(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/85636478b0bd8e67e89469c7749d4127/publicips/d8a46410-3c33-48d6-884c-df692e445822", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := eip.Delete(fake.ServiceClient(), "d8a46410-3c33-48d6-884c-df692e445822")
	th.AssertNoErr(t, res.Err)
}
