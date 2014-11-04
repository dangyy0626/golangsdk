package lbs

import (
	"github.com/mitchellh/mapstructure"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/lb/v1/nodes"
	"github.com/rackspace/gophercloud/rackspace/lb/v1/vips"
)

// Protocol represents the network protocol which the load balancer accepts.
type Protocol string

// The constants below represent all the compatible load balancer protocols.
const (
	// DNSTCP is a protocol that works with IPv6 and allows your DNS server to
	// receive traffic using TCP port 53.
	DNSTCP = "DNS_TCP"

	// DNSUDP is a protocol that works with IPv6 and allows your DNS server to
	// receive traffic using UDP port 53.
	DNSUDP = "DNS_UDP"

	// TCP is one of the core protocols of the Internet Protocol Suite. It
	// provides a reliable, ordered delivery of a stream of bytes from one
	// program on a computer to another program on another computer. Applications
	// that require an ordered and reliable delivery of packets use this protocol.
	TCP = "TCP"

	// TCPCLIENTFIRST is a protocol similar to TCP, but is more efficient when a
	// client is expected to write the data first.
	TCPCLIENTFIRST = "TCP_CLIENT_FIRST"

	// UDP provides a datagram service that emphasizes speed over reliability. It
	// works well with applications that provide security through other measures.
	UDP = "UDP"

	// UDPSTREAM is a protocol designed to stream media over networks and is
	// built on top of UDP.
	UDPSTREAM = "UDP_STREAM"
)

// Algorithm defines how traffic should be directed between back-end nodes.
type Algorithm string

const (
	// LC directs traffic to the node with the lowest number of connections.
	LC = "LEAST_CONNECTIONS"

	// RAND directs traffic to nodes at random.
	RAND = "RANDOM"

	// RR directs traffic to each of the nodes in turn.
	RR = "ROUND_ROBIN"

	// WLC directs traffic to a node based on the number of concurrent
	// connections and its weight.
	WLC = "WEIGHTED_LEAST_CONNECTIONS"

	// WRR directs traffic to a node according to the RR algorithm, but with
	// different proportions of traffic being directed to the back-end nodes.
	// Weights must be defined as part of the node configuration.
	WRR = "WEIGHTED_ROUND_ROBIN"
)

// Status represents the potential state of a load balancer resource.
type Status string

const (
	// ACTIVE indicates that the LB is configured properly and ready to serve
	// traffic to incoming requests via the configured virtual IPs.
	ACTIVE = "ACTIVE"

	// BUILD indicates that the LB is being provisioned for the first time and
	// configuration is being applied to bring the service online. The service
	// cannot yet serve incoming requests.
	BUILD = "BUILD"

	// PENDINGUPDATE indicates that the LB is online but configuration changes
	// are being applied to update the service based on a previous request.
	PENDINGUPDATE = "PENDING_UPDATE"

	// PENDINGDELETE indicates that the LB is online but configuration changes
	// are being applied to begin deletion of the service based on a previous
	// request.
	PENDINGDELETE = "PENDING_DELETE"

	// SUSPENDED indicates that the LB has been taken offline and disabled.
	SUSPENDED = "SUSPENDED"

	// ERROR indicates that the system encountered an error when attempting to
	// configure the load balancer.
	ERROR = "ERROR"

	// DELETED indicates that the LB has been deleted.
	DELETED = "DELETED"
)

// Datetime represents the structure of a Created or Updated field.
type Datetime struct {
	Time string
}

// LoadBalancer represents a load balancer API resource.
type LoadBalancer struct {
	// Human-readable name for the load balancer.
	Name string

	// The unique ID for the load balancer.
	ID int

	// Represents the service protocol being load balanced. See Protocol type for
	// a list of accepted values.
	Protocol Protocol

	// Defines how traffic should be directed between back-end nodes. The default
	// algorithm is RANDOM. See Algorithm type for a list of accepted values.
	Algorithm Algorithm

	// The current status of the load balancer.
	Status Status

	// The number of load balancer nodes.
	NodeCount int `mapstructure:"nodeCount"`

	// Slice of virtual IPs associated with this load balancer.
	VIPs []vips.VIP `mapstructure:"virtualIps"`

	// Datetime when the LB was created.
	Created Datetime

	// Datetime when the LB was created.
	Updated Datetime

	// Port number for the service you are load balancing.
	Port int

	// HalfClosed provides the ability for one end of the connection to
	// terminate its output while still receiving data from the other end. This
	// is only available on TCP/TCP_CLIENT_FIRST protocols.
	HalfClosed bool

	// Timeout represents the timeout value between a load balancer and its
	// nodes. Defaults to 30 seconds with a maximum of 120 seconds.
	Timeout int

	// TODO
	Cluster Cluster

	// Nodes shows all the back-end nodes which are associated with the load
	// balancer. These are the devices which are delivered traffic.
	Nodes []nodes.Node

	// TODO
	ConnectionLogging ConnectionLogging

	// SessionPersistence specifies whether multiple requests from clients are
	// directed to the same node.
	SessionPersistence SessionPersistence

	// ConnectionThrottle specifies a limit on the number of connections per IP
	// address to help mitigate malicious or abusive traffic to your applications.
	ConnectionThrottle ConnectionThrottle

	// TODO
	SourceAddrs SourceAddrs `mapstructure:"sourceAddresses"`
}

// SourceAddrs - temp
type SourceAddrs struct {
	IPv4Public  string `json:"ipv4Public" mapstructure:"ipv4Public"`
	IPv4Private string `json:"ipv4Servicenet" mapstructure:"ipv4Servicenet"`
	IPv6Public  string `json:"ipv6Public" mapstructure:"ipv6Public"`
	IPv6Private string `json:"ipv6Servicenet" mapstructure:"ipv6Servicenet"`
}

// SessionPersistence - temp
type SessionPersistence struct {
	Type string `json:"persistenceType" mapstructure:"persistenceType"`
}

// ConnectionThrottle - temp
type ConnectionThrottle struct {
	MinConns     int `json:"minConnections" mapstructure:"minConnections"`
	MaxConns     int `json:"maxConnections" mapstructure:"maxConnections"`
	MaxConnRate  int `json:"maxConnectionRate" mapstructure:"maxConnectionRate"`
	RateInterval int `json:"rateInterval" mapstructure:"rateInterval"`
}

// ConnectionLogging - temp
type ConnectionLogging struct {
	Enabled bool
}

// Cluster - temp
type Cluster struct {
	Name string
}

// LBPage is the page returned by a pager when traversing over a collection of
// LBs.
type LBPage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks whether a NetworkPage struct is empty.
func (p LBPage) IsEmpty() (bool, error) {
	is, err := ExtractLBs(p)
	if err != nil {
		return true, nil
	}
	return len(is) == 0, nil
}

// ExtractLBs accepts a Page struct, specifically a LBPage struct, and extracts
// the elements into a slice of LoadBalancer structs. In other words, a generic
// collection is mapped into a relevant slice.
func ExtractLBs(page pagination.Page) ([]LoadBalancer, error) {
	var resp struct {
		LBs []LoadBalancer `mapstructure:"loadBalancers" json:"loadBalancers"`
	}

	err := mapstructure.Decode(page.(LBPage).Body, &resp)

	return resp.LBs, err
}

type commonResult struct {
	gophercloud.Result
}

// Extract interprets any commonResult as a LB, if possible.
func (r commonResult) Extract() (*LoadBalancer, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		LB LoadBalancer `mapstructure:"loadBalancer"`
	}

	err := mapstructure.Decode(r.Body, &response)

	return &response.LB, err
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// UpdateResult represents the result of an update operation.
type UpdateResult struct {
	gophercloud.ErrResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}
