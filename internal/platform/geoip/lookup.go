package geoip

import (
	"net"
	"strings"
)

type Result struct {
	IP          string
	Country     string
	CountryCode string
	City        string
}

func Lookup(ip string) Result {
	ip = strings.TrimSpace(ip)
	out := Result{IP: ip, Country: "Unknown", CountryCode: "XX", City: "Unknown"}
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return out
	}
	if parsed.IsLoopback() {
		return Result{IP: ip, Country: "Local", CountryCode: "LO", City: "localhost"}
	}
	if parsed.IsPrivate() {
		return Result{IP: ip, Country: "Private", CountryCode: "PR", City: "lan"}
	}
	v4 := parsed.To4()
	if v4 == nil {
		return out
	}
	first := int(v4[0])
	switch {
	case first >= 1 && first < 50:
		out.Country, out.CountryCode, out.City = "United States", "US", "Unknown"
	case first >= 50 && first < 100:
		out.Country, out.CountryCode, out.City = "Germany", "DE", "Unknown"
	case first >= 100 && first < 150:
		out.Country, out.CountryCode, out.City = "Iran", "IR", "Unknown"
	case first >= 150 && first < 200:
		out.Country, out.CountryCode, out.City = "United Kingdom", "GB", "Unknown"
	default:
		out.Country, out.CountryCode, out.City = "Other", "OT", "Unknown"
	}
	return out
}
