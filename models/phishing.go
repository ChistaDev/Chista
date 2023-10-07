package models

type PhishingDomain struct {
	ID       uint64 `json:"id"`
	Domain   string `json:"domain" binding:"required"`
	Hostname string
	TLD      string
	NS       []string
}

type PhishingDomainAsArray struct {
	PhishingDomains []PhishingDomain
}

type ResponseDomain struct {
	Domain string `json:"domain"`
}

type DnsTwisterToHexResponse struct {
	Domain               string `json:"domain"`
	DomainAsHexadecimal  string `json:"domain_as_hexadecimal"`
	FuzzURL              string `json:"fuzz_url"`
	HasMXURL             string `json:"has_mx_url"`
	ParkedScoreURL       string `json:"parked_score_url"`
	ResolveIPURL         string `json:"resolve_ip_url"`
	SafeBrowsingCheckURL string `json:"safebrowsing_check_url"`
	URL                  string `json:"url"`
}

type DnsTwisterFuzzResponse struct {
	Domain              string                  `json:"domain"`
	DomainAsHexadecimal string                  `json:"domain_as_hexadecimal"`
	FuzzyDomains        []DnsTwisterFuzzyDomain `json:"fuzzy_domains"`
	HasMXURL            string                  `json:"has_mx_url"`
	ParkedScoreURL      string                  `json:"parked_score_url"`
	ResolveIPURL        string                  `json:"resolve_ip_url"`
	URL                 string                  `json:"url"`
}

type DnsTwisterFuzzyDomain struct {
	Domain              string `json:"domain"`
	DomainAsHexadecimal string `json:"domain_as_hexadecimal"`
	FuzzURL             string `json:"fuzz_url"`
	Fuzzer              string `json:"fuzzer"`
	HasMXURL            string `json:"has_mx_url"`
	ParkedScoreURL      string `json:"parked_score_url"`
	ResolveIPURL        string `json:"resolve_ip_url"`
}
