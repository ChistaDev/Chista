package models

import "time"

type APTFeedConfigDB struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

// APTFeedData is a struct to hold the data from the APT Intrusion Set.
type APTFeedData struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	SpecVersion string `json:"spec_version"`
	Objects     []struct {
		Modified           time.Time `json:"modified"`
		Name               string    `json:"name"`
		Description        string    `json:"description"`
		Aliases            []string  `json:"aliases"`
		XMitreDeprecated   bool      `json:"x_mitre_deprecated"`
		XMitreVersion      string    `json:"x_mitre_version"`
		XMitreContributors []string  `json:"x_mitre_contributors"`
		Type               string    `json:"type"`
		ID                 string    `json:"id"`
		Created            time.Time `json:"created"`
		CreatedByRef       string    `json:"created_by_ref"`
		Revoked            bool      `json:"revoked"`
		ExternalReferences []struct {
			SourceName  string `json:"source_name"`
			URL         string `json:"url,omitempty"`
			ExternalID  string `json:"external_id,omitempty"`
			Description string `json:"description,omitempty"`
		} `json:"external_references"`
		ObjectMarkingRefs       []string `json:"object_marking_refs"`
		XMitreDomains           []string `json:"x_mitre_domains"`
		XMitreAttackSpecVersion string   `json:"x_mitre_attack_spec_version"`
		XMitreModifiedByRef     string   `json:"x_mitre_modified_by_ref"`
	} `json:"objects"`
}

// APTFeedAttackPattern is a struct to hold the data from the APT Technics.
type APTFeedTechnic struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	SpecVersion string `json:"spec_version"`
	Objects     []struct {
		XMitrePlatforms    []string  `json:"x_mitre_platforms"`
		XMitreDomains      []string  `json:"x_mitre_domains"`
		ObjectMarkingRefs  []string  `json:"object_marking_refs"`
		Type               string    `json:"type"`
		ID                 string    `json:"id"`
		Created            time.Time `json:"created"`
		XMitreVersion      string    `json:"x_mitre_version"`
		ExternalReferences []struct {
			SourceName  string `json:"source_name"`
			ExternalID  string `json:"external_id,omitempty"`
			URL         string `json:"url"`
			Description string `json:"description,omitempty"`
		} `json:"external_references"`
		XMitreDeprecated bool      `json:"x_mitre_deprecated"`
		Revoked          bool      `json:"revoked"`
		Description      string    `json:"description"`
		Modified         time.Time `json:"modified"`
		CreatedByRef     string    `json:"created_by_ref"`
		Name             string    `json:"name"`
		XMitreDetection  string    `json:"x_mitre_detection"`
		KillChainPhases  []struct {
			KillChainName string `json:"kill_chain_name"`
			PhaseName     string `json:"phase_name"`
		} `json:"kill_chain_phases"`
		XMitreIsSubtechnique    bool     `json:"x_mitre_is_subtechnique"`
		XMitreTacticType        []string `json:"x_mitre_tactic_type"`
		XMitreAttackSpecVersion string   `json:"x_mitre_attack_spec_version"`
		XMitreModifiedByRef     string   `json:"x_mitre_modified_by_ref"`
	} `json:"objects"`
}

// APTFeedTactic is a struct to hold the data from the APT Tactic.
type APTFeedTactic struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	SpecVersion string `json:"spec_version"`
	Objects     []struct {
		XMitreDomains      []string  `json:"x_mitre_domains"`
		ObjectMarkingRefs  []string  `json:"object_marking_refs"`
		ID                 string    `json:"id"`
		Type               string    `json:"type"`
		Created            time.Time `json:"created"`
		CreatedByRef       string    `json:"created_by_ref"`
		ExternalReferences []struct {
			ExternalID string `json:"external_id"`
			URL        string `json:"url"`
			SourceName string `json:"source_name"`
		} `json:"external_references"`
		Modified                time.Time `json:"modified"`
		Name                    string    `json:"name"`
		Description             string    `json:"description"`
		XMitreVersion           string    `json:"x_mitre_version"`
		XMitreAttackSpecVersion string    `json:"x_mitre_attack_spec_version"`
		XMitreModifiedByRef     string    `json:"x_mitre_modified_by_ref"`
		XMitreShortname         string    `json:"x_mitre_shortname"`
	} `json:"objects"`
}

// APTFeedRelationship is a struct to hold the data from the APT Relationship.
type APTFeedRelationship struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	SpecVersion string `json:"spec_version"`
	Objects     []struct {
		ObjectMarkingRefs  []string  `json:"object_marking_refs"`
		Type               string    `json:"type"`
		ID                 string    `json:"id"`
		Created            time.Time `json:"created"`
		XMitreVersion      string    `json:"x_mitre_version"`
		ExternalReferences []struct {
			SourceName  string `json:"source_name"`
			URL         string `json:"url"`
			Description string `json:"description"`
		} `json:"external_references"`
		XMitreDeprecated        bool      `json:"x_mitre_deprecated"`
		Revoked                 bool      `json:"revoked"`
		Description             string    `json:"description"`
		Modified                time.Time `json:"modified"`
		CreatedByRef            string    `json:"created_by_ref"`
		RelationshipType        string    `json:"relationship_type"`
		SourceRef               string    `json:"source_ref"`
		TargetRef               string    `json:"target_ref"`
		XMitreAttackSpecVersion string    `json:"x_mitre_attack_spec_version"`
		XMitreModifiedByRef     string    `json:"x_mitre_modified_by_ref"`
	} `json:"objects"`
}

// APTFeedMitigation is a struct to hold the data from the APT Mitigation.
type APTFeedMitigation struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	SpecVersion string `json:"spec_version"`
	Objects     []struct {
		XMitreDomains      []string  `json:"x_mitre_domains"`
		ObjectMarkingRefs  []string  `json:"object_marking_refs"`
		ID                 string    `json:"id"`
		Type               string    `json:"type"`
		Created            time.Time `json:"created"`
		CreatedByRef       string    `json:"created_by_ref"`
		ExternalReferences []struct {
			URL         string `json:"url"`
			SourceName  string `json:"source_name"`
			ExternalID  string `json:"external_id,omitempty"`
			Description string `json:"description,omitempty"`
		} `json:"external_references"`
		Modified            time.Time `json:"modified"`
		Name                string    `json:"name"`
		Description         string    `json:"description"`
		XMitreDeprecated    bool      `json:"x_mitre_deprecated"`
		XMitreVersion       string    `json:"x_mitre_version"`
		XMitreModifiedByRef string    `json:"x_mitre_modified_by_ref"`
	} `json:"objects"`
}
