package openstack

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/gophercloud/gophercloud"
	tokens2 "github.com/gophercloud/gophercloud/openstack/identity/v2/tokens"
	tokens3 "github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/openstack/utils"
	"strings"
)

const (
	v20 = "v2.0"
	v30 = "v3.0"
)

// NewClient prepares an unauthenticated ProviderClient instance.
// Most users will probably prefer using the AuthenticatedClient function instead.
// This is useful if you wish to explicitly control the version of the identity service that's used for authentication explicitly,
// for example.
func NewClient(endpoint string) (*gophercloud.ProviderClient, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	hadPath := u.Path != ""
	u.Path, u.RawQuery, u.Fragment = "", "", ""
	base := u.String()

	endpoint = gophercloud.NormalizeURL(endpoint)
	base = gophercloud.NormalizeURL(base)

	if hadPath {
		return &gophercloud.ProviderClient{
			IdentityBase:     base,
			IdentityEndpoint: endpoint,
		}, nil
	}

	return &gophercloud.ProviderClient{
		IdentityBase:     base,
		IdentityEndpoint: "",
	}, nil
}

// AuthenticatedClient logs in to an OpenStack cloud found at the identity endpoint specified by options, acquires a token, and
// returns a Client instance that's ready to operate.
// It first queries the root identity endpoint to determine which versions of the identity service are supported, then chooses
// the most recent identity service available to proceed.
func AuthenticatedClient(options gophercloud.AuthOptions) (*gophercloud.ProviderClient, error) {
	client, err := NewClient(options.IdentityEndpoint)
	if err != nil {
		return nil, err
	}

	err = Authenticate(client, options)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Authenticate or re-authenticate against the most recent identity service supported at the provided endpoint.
func Authenticate(client *gophercloud.ProviderClient, options gophercloud.AuthOptions) error {
	versions := []*utils.Version{
		{ID: v20, Priority: 20, Suffix: "/v2.0/"},
		{ID: v30, Priority: 30, Suffix: "/v3/"},
	}

	chosen, endpoint, err := utils.ChooseVersion(client, versions)
	if err != nil {
		return err
	}

	switch chosen.ID {
	case v20:
		return v2auth(client, endpoint, options, gophercloud.EndpointOpts{})
	case v30:
		return v3auth(client, endpoint, &options, gophercloud.EndpointOpts{})
	default:
		// The switch statement must be out of date from the versions list.
		return fmt.Errorf("Unrecognized identity version: %s", chosen.ID)
	}
}

// AuthenticateV2 explicitly authenticates against the identity v2 endpoint.
func AuthenticateV2(client *gophercloud.ProviderClient, options gophercloud.AuthOptions, eo gophercloud.EndpointOpts) error {
	return v2auth(client, "", options, eo)
}

func v2auth(client *gophercloud.ProviderClient, endpoint string, options gophercloud.AuthOptions, eo gophercloud.EndpointOpts) error {
	v2Client, err := NewIdentityV2(client, eo)
	if err != nil {
		return err
	}

	if endpoint != "" {
		v2Client.Endpoint = endpoint
	}

	v2Opts := tokens2.AuthOptions{
		IdentityEndpoint: options.IdentityEndpoint,
		Username:         options.Username,
		Password:         options.Password,
		TenantID:         options.TenantID,
		TenantName:       options.TenantName,
		AllowReauth:      options.AllowReauth,
		TokenID:          options.TokenID,
	}

	result := tokens2.Create(v2Client, v2Opts)

	token, err := result.ExtractToken()
	if err != nil {
		return err
	}

	catalog, err := result.ExtractServiceCatalog()
	if err != nil {
		return err
	}

	if options.AllowReauth {
		client.ReauthFunc = func() error {
			client.TokenID = ""
			return v2auth(client, endpoint, options, eo)
		}
	}
	client.TokenID = token.ID
	client.EndpointLocator = func(opts gophercloud.EndpointOpts) (string, error) {
		return V2EndpointURL(catalog, opts)
	}

	return nil
}

// AuthenticateV3 explicitly authenticates against the identity v3 service.
func AuthenticateV3(client *gophercloud.ProviderClient, options tokens3.AuthOptionsBuilder, eo gophercloud.EndpointOpts) error {
	return v3auth(client, "", options, eo)
}

func v3auth(client *gophercloud.ProviderClient, endpoint string, opts tokens3.AuthOptionsBuilder, eo gophercloud.EndpointOpts) error {
	// Override the generated service endpoint with the one returned by the version endpoint.
	v3Client, err := NewIdentityV3(client, eo)
	if err != nil {
		return err
	}

	if endpoint != "" {
		v3Client.Endpoint = endpoint
	}

	result := tokens3.Create(v3Client, opts)

	token, err := result.ExtractToken()
	if err != nil {
		return err
	}

	catalog, err := result.ExtractServiceCatalog()
	if err != nil {
		return err
	}

	client.TokenID = token.ID

	if opts.CanReauth() {
		client.ReauthFunc = func() error {
			client.TokenID = ""
			return v3auth(client, endpoint, opts, eo)
		}
	}
	client.EndpointLocator = func(opts gophercloud.EndpointOpts) (string, error) {
		return V3EndpointURL(catalog, opts)
	}

	return nil
}

func GetProjectId(client *gophercloud.ProviderClient) (string, error) {
	versions := []*utils.Version{
		{ID: v20, Priority: 20, Suffix: "/v2.0/"},
		{ID: v30, Priority: 30, Suffix: "/v3/"},
	}

	chosen, endpoint, err := utils.ChooseVersion(client, versions)
	if err != nil {
		return "", err
	}

	switch chosen.ID {
	case v20:
		return getV2ProjectId(client, endpoint)
	case v30:
		return getV3ProjectId(client, endpoint)
	default:
		return "", fmt.Errorf("Unrecognized identity version: %s", chosen.ID)
	}
}

func getV2ProjectId(client *gophercloud.ProviderClient, endpoint string) (string, error) {
	v2Client, err := NewIdentityV2(client, gophercloud.EndpointOpts{})
	if err != nil {
		return "", err
	}

	if endpoint != "" {
		v2Client.Endpoint = endpoint
	}

	result := tokens2.Get(v2Client, client.TokenID)
	token, err := result.ExtractToken()
	if err != nil {
		return "", err
	}

	return token.Tenant.ID, nil
}

func getV3ProjectId(client *gophercloud.ProviderClient, endpoint string) (string, error) {
	v3Client, err := NewIdentityV3(client, gophercloud.EndpointOpts{})
	if err != nil {
		return "", err
	}

	if endpoint != "" {
		v3Client.Endpoint = endpoint
	}

	result := tokens3.Get(v3Client, client.TokenID)
	project, err := result.ExtractProject()
	if err != nil {
		return "", err
	}

	return project.ID, nil
}

// NewIdentityV2 creates a ServiceClient that may be used to interact with the v2 identity service.
func NewIdentityV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	endpoint := client.IdentityBase + "v2.0/"
	clientType := "identity"
	var err error
	if !reflect.DeepEqual(eo, gophercloud.EndpointOpts{}) {
		eo.ApplyDefaults(clientType)
		endpoint, err = client.EndpointLocator(eo)
		if err != nil {
			return nil, err
		}
	}

	return &gophercloud.ServiceClient{
		ProviderClient: client,
		Endpoint:       endpoint,
		Type:           clientType,
	}, nil
}

// NewIdentityV3 creates a ServiceClient that may be used to access the v3 identity service.
func NewIdentityV3(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	endpoint := client.IdentityBase + "v3/"
	clientType := "identity"
	var err error
	if !reflect.DeepEqual(eo, gophercloud.EndpointOpts{}) {
		eo.ApplyDefaults(clientType)
		endpoint, err = client.EndpointLocator(eo)
		if err != nil {
			return nil, err
		}
	}

	return &gophercloud.ServiceClient{
		ProviderClient: client,
		Endpoint:       endpoint,
		Type:           clientType,
	}, nil
}

func initClientOpts(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts, clientType string) (*gophercloud.ServiceClient, error) {
	sc := new(gophercloud.ServiceClient)
	eo.ApplyDefaults(clientType)
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return sc, err
	}
	sc.ProviderClient = client
	sc.Endpoint = url
	sc.Type = clientType
	return sc, nil
}

func initClientOpts1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts, clientType string) (*gophercloud.ServiceClient1, error) {
	pid, e := GetProjectId(client)
	if e != nil {
		return nil, e
	}

	c, e := initClientOpts(client, eo, clientType)
	if e != nil {
		return nil, e
	}

	sc := new(gophercloud.ServiceClient1)
	sc.ServiceClient = c
	sc.ProjectID = pid
	return sc, nil
}

// NewObjectStorageV1 creates a ServiceClient that may be used with the v1 object storage package.
func NewObjectStorageV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(client, eo, "object-store")
}

// NewComputeV2 creates a ServiceClient that may be used with the v2 compute package.
func NewComputeV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(client, eo, "compute")
}

// NewNetworkV2 creates a ServiceClient that may be used with the v2 network package.
func NewNetworkV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "network")
	sc.ResourceBase = sc.Endpoint + "v2.0/"
	return sc, err
}

// NewVpcV1 creates a ServiceClient that may be used with the v1 VPC for OTC.
func NewVpcV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	pid, e := GetProjectId(client)
	if e != nil {
		return nil, e
	}

	sc, err := initClientOpts(client, eo, "network")
	sc.ResourceBase = sc.Endpoint + "v1/" + pid + "/"
	return sc, err
}

// NewOtcV1 creates a ServiceClient that may be used with the v1 network package.
func NewOtcV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts, otctype string) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "compute")
	//fmt.Printf("client=%+v.\n", sc)
	sc.Endpoint = strings.Replace(strings.Replace(sc.Endpoint, "ecs", otctype, 1), "/v2/", "/v1.0/", 1)
	//fmt.Printf("url=%s.\n", sc.Endpoint)
	sc.ResourceBase = sc.Endpoint
	sc.Type = otctype
	return sc, err
}

// NewBlockStorageV1 creates a ServiceClient that may be used to access the v1 block storage service.
func NewBlockStorageV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(client, eo, "volume")
}

// NewBlockStorageV2 creates a ServiceClient that may be used to access the v2 block storage service.
func NewBlockStorageV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(client, eo, "volumev2")
}

// NewSharedFileSystemV2 creates a ServiceClient that may be used to access the v2 shared file system service.
func NewSharedFileSystemV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(client, eo, "sharev2")
}

// NewCDNV1 creates a ServiceClient that may be used to access the OpenStack v1
// CDN service.
func NewCDNV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(client, eo, "cdn")
}

// NewOrchestrationV1 creates a ServiceClient that may be used to access the v1 orchestration service.
func NewOrchestrationV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(client, eo, "orchestration")
}

// NewDBV1 creates a ServiceClient that may be used to access the v1 DB service.
func NewDBV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(client, eo, "database")
}

// NewDNSV2 creates a ServiceClient that may be used to access the v2 DNS service.
func NewDNSV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "dns")
	sc.ResourceBase = sc.Endpoint + "v2/"
	return sc, err
}

// NewImageServiceV2 creates a ServiceClient that may be used to access the v2 image service.
func NewImageServiceV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(client, eo, "image")
	sc.ResourceBase = sc.Endpoint + "v2/"
	return sc, err
}

func NewCESClient(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient1, error) {
	sc, err := initClientOpts1(client, eo, "ces")
	if err != nil {
		return nil, err
	}
	sc.ResourceBase = sc.Endpoint
	return sc, err
}

// NewSmnServiceV2 creates a ServiceClient that may be used to access the v2 Simple Message Notification service.
func NewSmnServiceV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {

	sc, err := initClientOpts(client, eo, "compute")
	sc.Endpoint = strings.Replace(sc.Endpoint, "ecs", "smn", 1)
	sc.ResourceBase = sc.Endpoint + "notifications/"
	sc.Type = "smn"
	return sc, err
}

//NewRdsServiceV1 creates the a ServiceClient that may be used to access the v1
//rds service which is a service of db instances management.
func NewRdsServiceV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	newsc, err := initClientOpts(client, eo, "compute")
	rdsendpoint := strings.Replace(strings.Replace(newsc.Endpoint, "ecs", "rds", 1), "/v2/", "/rds/v1/", 1)
	newsc.Endpoint = rdsendpoint
	newsc.ResourceBase = rdsendpoint
	newsc.Type = "rds"
	return newsc, err
}