package eip

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// Bandwidth contains the Name,the Size,the ShareType and the ChargeMode.
type Bandwidth struct {
	Name string `json:"name" required:"true"`
	Size int `json:"size" required:"true"`
	ShareType string `json:"share_type" required:"true"`
	ChargeMode string `json:"charge_mode,omitempty"`
}

//Publicip contains the Type and the IpAddress.
type Publicip struct {
	Type string `json:"type" required:"true"`
	IpAddress string `json:"ip_address,omitempty"`
}


// Eip represents a Neutron eip.
type Eip struct {
	// ID is the unique identifier for the eip.
	ID string `json:"id"`

	// Name is the human readable name for the eip. It does not have to be
	// unique.
	Name string `json:"bandwidth_name"`

	//Specifies the elastic IP address type in the eip.
	Type string `json:"type"`

	// Status indicates whether or not a eip is currently operational.
	Status string `json:"status"`

	//Specifies the elastic IP address obtained in the eip.
	PublicIpAddress string `json:"public_ip_address"`

	//Specifies the private IP address bound to the elastic IP address in the eip.
	PrivateIpAddress string `json:"private_ip_address"`

	//Specifies the ID of the VM NIC in the eip.
	PortId string `json:"port_id"`

	//Specifies the tenant ID of the operator.
	TenantID string `json:"tenant_id"`

	//Specifies the time for applying for the elastic IP address.
	CreateTime string `json:"create_time"`

	//Specifies the bandwidth ID of the elastic IP address.
	BandwidthId string `json:"bandwidth_id"`

	//Specifies the bandwidth capacity in the eip.
	BandwidthSize int `json:"bandwidth_size"`

	//Specifies the bandwidth capacity in the eip.
	BandwidthShareType string `json:"bandwidth_share_type"`

}

// EipPage is the page returned by a pager when traversing over a
// collection of eips.
type EipPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of eips has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r EipPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"publicips_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a EipPage struct is empty.
func (r EipPage) IsEmpty() (bool, error) {
	is, err := ExtractEips(r)
	return len(is) == 0, err
}

// ExtractEips accepts a Page struct, specifically a EipPage struct,
// and extracts the elements into a slice of Eip structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractEips(r pagination.Page) ([]Eip, error) {
	var s struct {
		Eips []Eip `json:"publicips"`
	}
	err := (r.(EipPage)).ExtractInto(&s)
	return s.Eips, err
}

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a eip.
func (r commonResult) Extract() (*Eip, error) {
	var s struct {
		Eip *Eip `json:"publicip"`
	}
	err := r.ExtractInto(&s)
	return s.Eip, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Eip.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Eip.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Eip.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}
