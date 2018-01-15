package eip

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"reflect"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the floating IP attributes you want to see returned. SortKey allows you to
// sort by a particular network attribute. SortDir sets the direction, and is
// either `asc' or `desc'. Marker and Limit are used for pagination.

type ListOpts struct {
	// ID is the unique identifier for the eip.
	ID string `json:"id,omitempty"`

	// Name is the human readable name for the eip. It does not have to be
	// unique.
	Name string `json:"bandwidth_name,omitempty"`

	//Specifies the elastic IP address type in the eip.
	Type string `json:"type,omitempty"`

	// Status indicates whether or not a eip is currently operational.
	Status string `json:"status,omitempty"`

	//Specifies the elastic IP address obtained in the eip.
	PublicIpAddress string `json:"public_ip_address,omitempty"`

	//Specifies the private IP address bound to the elastic IP address in the eip.
	PrivateIpAddress string `json:"private_ip_address,omitempty"`

	//Specifies the ID of the VM NIC in the eip.
	PortId string `json:"port_id,omitempty"`

	//Specifies the tenant ID of the operator.
	TenantID string `json:"tenant_id,omitempty"`

	//Specifies the time for applying for the elastic IP address.
	CreateTime string `json:"create_time,omitempty"`

	//Specifies the bandwidth ID of the elastic IP address.
	BandwidthId string `json:"bandwidth_id,omitempty"`

	//Specifies the bandwidth capacity in the eip.
	BandwidthSize int `json:"bandwidth_size,omitempty"`

	//Specifies the bandwidth capacity in the eip.
	BandwidthShareType string `json:"bandwidth_share_type,omitempty"`
	}

// List returns collection of
// eips. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those eips that are owned by the
// tenant who submits the request, unless an admin user submits the request.

func List(c *gophercloud.ServiceClient, opts ListOpts) ([]Eip, error) {

	u := rootURL(c)

	pages, err := pagination.NewPager(c, u, func(r pagination.PageResult) pagination.Page {
		return EipPage{pagination.LinkedPageBase{PageResult: r}}
	}).AllPages()

	allEip, err := ExtractEips(pages)
	if err != nil {
		return nil, err
	}

	return FilterEIPs(allEip, opts)
}

func FilterEIPs(eips []Eip, opts ListOpts) ([]Eip, error) {

	var refinedEIPs []Eip
	var matched bool
	m := map[string]interface{}{}

	if opts.ID != "" {
		m["ID"] = opts.ID
	}
	if opts.Name != "" {
		m["Name"] = opts.Name
	}
	if opts.Status != "" {
		m["Status"] = opts.Status
	}
	if opts.Type != "" {
		m["Type"] = opts.Type
	}
	if opts.PublicIpAddress != "" {
		m["PublicIpAddress"] = opts.PublicIpAddress
	}
	if opts.TenantID != "" {
		m["TenantID"] = opts.TenantID
	}
	if opts.CreateTime != "" {
		m["CreateTime"] = opts.CreateTime
	}
	if opts.BandwidthId != "" {
		m["BandwidthId"] = opts.BandwidthId
	}
	if opts.TenantID != "" {
		m["BandwidthSize"] = opts.BandwidthSize
	}
	if opts.BandwidthId != "" {
		m["BandwidthShareType"] = opts.BandwidthShareType
	}

	if len(m) > 0 && len(eips) > 0 {
		for _, eip := range eips {
			matched = true

			for key, value := range m {
				if sVal := getStructField(&eip, key); !(sVal == value) {
					matched = false
				}
			}

			if matched {
				refinedEIPs = append(refinedEIPs, eip)
			}
		}

	} else {
		refinedEIPs = eips
	}

	return refinedEIPs, nil
}

func getStructField(v *Eip, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return string(f.String())
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToEipCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains all the values needed to create a new eip.
type CreateOpts struct {
	PublicIp  Publicip  `json:"publicip,omitempty" required:"true"`
	BandWidth Bandwidth `json:"bandwidth,omitempty" required:"true"`
}

// ToEipCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToEipCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// logical eip.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToEipCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &gophercloud.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular vpc based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToEipUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts contains the values used when updating a eip.
type UpdateOpts struct {
	PortId string `json:"port_id,omitempty"`
}

// ToEipUpdateMap builds an update body based on UpdateOpts.
func (opts UpdateOpts) ToEipUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "publicip")
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
func Update(c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToEipUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular eip based on its unique ID.
func Delete(c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
