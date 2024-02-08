package models

import (
	"time"
)

type Response interface {
	ToJSON() ([]byte, error)
	FromJSON([]byte) error
}

type CrtShResponseModel []struct {
	IssuerCaID     int    `json:"issuer_ca_id"`
	IssuerName     string `json:"issuer_name"`
	CommonName     string `json:"common_name"`
	NameValue      string `json:"name_value"`
	ID             int64  `json:"id"`
	EntryTimestamp string `json:"entry_timestamp"`
	NotBefore      string `json:"not_before"`
	NotAfter       string `json:"not_after"`
	SerialNumber   string `json:"serial_number"`
}

type CensysAccountEndpointResponseModel struct {
	Email      string `json:"email"`
	Login      string `json:"login"`
	FirstLogin string `json:"first_login"`
	LastLogin  string `json:"last_login"`
	Quota      struct {
		Used      int    `json:"used"`
		Allowance int    `json:"allowance"`
		ResetsAt  string `json:"resets_at"`
	} `json:"quota"`
}

type CensysCTSearchEndpointResponseModel struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Result struct {
		Query      string  `json:"query"`
		Total      float64 `json:"total"`
		DurationMs int     `json:"duration_ms"`
		Hits       []struct {
			Parsed struct {
				IssuerDn       string `json:"issuer_dn"`
				SubjectDn      string `json:"subject_dn"`
				ValidityPeriod struct {
					NotAfter  time.Time `json:"not_after"`
					NotBefore time.Time `json:"not_before"`
				} `json:"validity_period"`
			} `json:"parsed"`
			FingerprintSha256 string   `json:"fingerprint_sha256"`
			Names             []string `json:"names"`
		} `json:"hits"`
		Links struct {
			Next string `json:"next"`
			Prev string `json:"prev"`
		} `json:"links"`
	} `json:"result"`
}

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
