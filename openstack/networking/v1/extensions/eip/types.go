package eip

//ShareType indicating that the bandwidth is exclusive.
type ShareType string

const (
	//The value is PER, indicating that the bandwidth is exclusive.
	PER = ShareType("PER")
)

//ElasticIpType specifies the elastic IP address type.
type ElasticIpType string

const (
	//The value is 5_bgp.The value must be a type supported by the system.
	Type= ElasticIpType("5_bgp")
)

