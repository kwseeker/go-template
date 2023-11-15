package net

import (
	"bytes"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common/protocol"
	"log"
	"net"
	"strings"
)

type AddressFamily byte

const (
	// AddressFamilyIPv4 represents address as IPv4
	AddressFamilyIPv4 = AddressFamily(0)

	// AddressFamilyIPv6 represents address as IPv6
	AddressFamilyIPv6 = AddressFamily(1)

	// AddressFamilyDomain represents address as Domain
	AddressFamilyDomain = AddressFamily(2)
)

type Address interface {
	IP() net.IP     // IP of this Address
	Domain() string // Domain of this Address
	Family() AddressFamily
	String() string // String representation of this Address
}

type ipv4Address [4]byte

func (a ipv4Address) IP() net.IP {
	return net.IP(a[:])
}

func (ipv4Address) Domain() string {
	panic("Calling Domain() on an IPv4Address.")
}

func (ipv4Address) Family() AddressFamily {
	return AddressFamilyIPv4
}

func (a ipv4Address) String() string {
	return a.IP().String()
}

type ipv6Address [16]byte

func (a ipv6Address) IP() net.IP {
	return net.IP(a[:])
}

func (ipv6Address) Domain() string {
	panic("Calling Domain() on an IPv6Address.")
}

func (ipv6Address) Family() AddressFamily {
	return AddressFamilyIPv6
}

func (a ipv6Address) String() string {
	return "[" + a.IP().String() + "]"
}

func isAlphaNum(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

var bytes0 = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// IPAddress creates an Address with given IP.
func IPAddress(ip []byte) Address {
	switch len(ip) {
	case net.IPv4len:
		var addr ipv4Address = [4]byte{ip[0], ip[1], ip[2], ip[3]}
		return addr
	case net.IPv6len:
		if bytes.Equal(ip[:10], bytes0) && ip[10] == 0xff && ip[11] == 0xff {
			return IPAddress(ip[12:16])
		}
		var addr ipv6Address = [16]byte{
			ip[0], ip[1], ip[2], ip[3],
			ip[4], ip[5], ip[6], ip[7],
			ip[8], ip[9], ip[10], ip[11],
			ip[12], ip[13], ip[14], ip[15],
		}
		return addr
	default:
		log.Println("invalid IP format: ", ip)
		return nil
	}
}

type domainAddress string

func (domainAddress) IP() net.IP {
	panic("Calling IP() on a DomainAddress.")
}

func (a domainAddress) Domain() string {
	return string(a)
}

func (domainAddress) Family() AddressFamily {
	return AddressFamilyDomain
}

func (a domainAddress) String() string {
	return a.Domain()
}

// DomainAddress creates an Address with given domain.
func DomainAddress(domain string) Address {
	return domainAddress(domain)
}

func ParseAddress(addr string) Address {
	// Handle IPv6 address in form as "[2001:4860:0:2001::68]"
	lenAddr := len(addr)
	if lenAddr > 0 && addr[0] == '[' && addr[lenAddr-1] == ']' {
		addr = addr[1 : lenAddr-1]
		lenAddr -= 2
	}

	if lenAddr > 0 && (!isAlphaNum(addr[0]) || !isAlphaNum(addr[len(addr)-1])) {
		addr = strings.TrimSpace(addr)
	}

	ip := net.ParseIP(addr)
	if ip != nil {
		return IPAddress(ip)
	}
	return DomainAddress(addr)
}

// NewIPOrDomain translates Address to IPOrDomain
func NewIPOrDomain(addr Address) *protocol.IPOrDomain {
	switch addr.Family() {
	case AddressFamilyDomain:
		return &protocol.IPOrDomain{
			Address: &protocol.IPOrDomain_Domain{
				Domain: addr.Domain(),
			},
		}
	case AddressFamilyIPv4, AddressFamilyIPv6:
		return &protocol.IPOrDomain{
			Address: &protocol.IPOrDomain_Ip{
				Ip: addr.IP(),
			},
		}
	default:
		panic("Unknown Address type.")
	}
}
