package main

var Validators = map[string][]string{
	"meta": {
		"-", "omitempty", "omitnil", "structonly", "nostructlevel",
		"dive", "keys", "endkeys", "|",
	},
	"required": {
		"required", "required_if", "required_unless",
		"required_with", "required_with_all",
		"required_without", "required_without_all",
		"excluded_if", "excluded_unless",
		"isdefault",
	},
	"length_numeric": {
		"len", "max", "min",
		"eq", "ne", "oneof", "oneofci",
		"gt", "gte", "lt", "lte",
	},
	"cross_field": {
		"eqfield", "nefield", "gtfield", "gtefield", "ltfield", "ltefield",
		"eqcsfield", "necsfield", "gtcsfield", "gtecsfield", "ltcsfield", "ltecsfield",
	},
	"strings": {
		"alpha", "alphanum", "alphaunicode", "alphanumunicode",
		"ascii", "boolean", "contains", "containsany", "containsrune",
		"endswith", "endsnotwith", "startswith", "startsnotwith",
		"excludes", "excludesall", "excludesrune",
		"lowercase", "uppercase", "multibyte",
		"number", "numeric", "printascii",
	},
	"formats": {
		"base64", "base64url", "base64rawurl",
		"bic", "bcp47_language_tag", "btc_addr", "btc_addr_bech32",
		"credit_card", "mongodb", "mongodb_connection_string", "cron",
		"datetime", "e164", "email", "eth_addr",
		"hexadecimal", "hexcolor", "hsl", "hsla",
		"html", "html_encoded", "isbn", "isbn10", "isbn13",
		"issn", "iso3166_1_alpha2", "iso3166_1_alpha3",
		"iso3166_1_alpha_numeric", "iso3166_2", "iso4217",
		"json", "jwt", "latitude", "longitude", "luhn_checksum",
		"postcode_iso3166_alpha2", "postcode_iso3166_alpha2_field",
		"rgb", "rgba", "ssn", "timezone",
		"uri", "url", "http_url", "url_encoded", "urn_rfc2141",
	},
	"network": {
		"cidr", "cidrv4", "cidrv6", "fqdn",
		"hostname", "hostname_port", "hostname_rfc1123",
		"ip", "ip4_addr", "ip6_addr", "ip_addr",
		"ipv4", "ipv6", "mac",
		"tcp_addr", "tcp4_addr", "tcp6_addr",
		"udp_addr", "udp4_addr", "udp6_addr",
		"unix_addr",
	},
	"other": {
		"unique",
	},
	"aliases": {
		"iscolor", "country_code",
	},
}
