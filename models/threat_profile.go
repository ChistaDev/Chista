package models

type CombinedProfilesData struct {
	RansomProfileData RansomProfileData
	AptData           AptData
}

// Main Ransomware Profile data container
type RansomProfileData struct {
	Name      string     `json:"name"`
	Meta      string     `json:"meta,omitempty"`
	Locations []Location `json:"locations"`
	Profile   []string   `json:"profile"`
}

type Location struct {
	FQDN       string `json:"fqdn"`
	Title      string `json:"title,omitempty"`
	Version    int    `json:"version,omitempty"`
	Slug       string `json:"slug,omitempty"`
	Available  bool   `json:"available,omitempty"`
	Updated    string `json:"updated,omitempty"`
	LastScrape string `json:"lastscrape,omitempty"`
	Enabled    bool   `json:"enabled,omitempty"`
}

// Main APT Profile data container
type AptDataContainer struct {
	Authors      []string  `json:"authors"`
	Category     string    `json:"category"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Source       string    `json:"source"`
	Description  string    `json:"description"`
	Tlp          string    `json:"tlp"`
	License      string    `json:"license"`
	UUID         string    `json:"uuid"`
	LastDbChange string    `json:"last-db-change"`
	AptData      []AptData `json:"values"`
}

type AptData struct {
	Actor string `json:"actor"`
	Names []struct {
		Name      string `json:"name"`
		NameGiver string `json:"name-giver"`
	} `json:"names"`
	Country           []string `json:"country"`
	Description       string   `json:"description"`
	Information       []string `json:"information,omitempty"`
	UUID              string   `json:"uuid"`
	LastCardChange    string   `json:"last-card-change"`
	Motivation        []string `json:"motivation,omitempty"`
	FirstSeen         string   `json:"first-seen,omitempty"`
	ObservedSectors   []string `json:"observed-sectors,omitempty"`
	ObservedCountries []string `json:"observed-countries,omitempty"`
	Tools             []string `json:"tools,omitempty"`
	Operations        []struct {
		Date     string `json:"date"`
		Activity string `json:"activity"`
	} `json:"operations,omitempty"`
	MitreAttack       []string `json:"mitre-attack,omitempty"`
	Sponsor           string   `json:"sponsor,omitempty"`
	CounterOperations []struct {
		Date     string `json:"date"`
		Activity string `json:"activity"`
	} `json:"counter-operations,omitempty"`
	Playbook      []string `json:"playbook,omitempty"`
	AlienvaultOtx []string `json:"alienvault-otx,omitempty"`
}
