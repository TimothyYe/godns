package utils

const (
	// PanicMax is the max allowed panic times.
	PanicMax = 5
	// DNSPOD for dnspod.cn.
	DNSPOD = "DNSPod"
	// HE for he.net.
	HE = "HE"
	// CLOUDFLARE for cloudflare.com.
	CLOUDFLARE = "Cloudflare"
	// ALIDNS for AliDNS.
	ALIDNS = "AliDNS"
	// GOOGLE for Google Domains.
	GOOGLE = "Google"
	// DUCK for Duck DNS.
	DUCK = "DuckDNS"
	// DREAMHOST for Dreamhost.
	DREAMHOST = "Dreamhost"
	// DYNV6 for Dynv6.
	DYNV6 = "Dynv6"
	// DYNU for Dynu.
	DYNU = "Dynu"
	// NOIP for NoIP.
	NOIP = "NoIP"
	// SCALEWAY for Scaleway.
	SCALEWAY = "Scaleway"
	// LINODE for Linode.
	LINODE = "Linode"
	// STRATO for Strato.
	STRATO = "Strato"
	// LOOPIASE for LoopiaSE.
	LOOPIASE = "LoopiaSE"
	// INFOMANIAK for Infomaniak.
	INFOMANIAK = "Infomaniak"
	// HETZNER for Hetzner.
	HETZNER = "Hetzner"
	// OVH for OVH.
	OVH = "OVH"
	// IPV4 for IPV4 mode.
	IPV4 = "IPV4"
	// IPV6 for IPV6 mode.
	IPV6 = "IPV6"
	// IPTypeA.
	IPTypeA = "A"
	// IPTypeAAAA.
	IPTypeAAAA = "AAAA"
	// RootDomain.
	RootDomain = "@"
	// Regex pattern to match IPV4 address.
	IPv4Pattern = `((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)`
	// Regex pattern to match IPV6 address.
	IPv6Pattern = `(([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|` +
		`(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|` +
		`(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|` +
		`(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|` +
		`(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|` +
		`(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|` +
		`(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|` +
		`(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))`
	// Regex pattern to match IPV4 and IPV6 address.
	IPPattern = "(" + IPv4Pattern + ")|(" + IPv6Pattern + ")"

	// DefaultTimeout is the default timeout value, in seconds.
	DefaultTimeout = 10
)
