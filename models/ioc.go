package models

type MalwarebazaarApiBody struct {
	QueryStatus string `json:"query_status"`
	Data        []struct {
		Sha256Hash  string   `json:"sha256_hash"`
		Sha3384Hash string   `json:"sha3_384_hash"`
		Sha1Hash    string   `json:"sha1_hash"`
		Md5Hash     string   `json:"md5_hash"`
		FirstSeen   string   `json:"first_seen"`
		LastSeen    string   `json:"last_seen"`
		FileName    string   `json:"file_name"`
		FileType    string   `json:"file_type"`
		Signature   string   `json:"signature"`
		Tags        []string `json:"tags"`
	} `json:"data"`
}

type UrlHausApiBody struct {
	DateAdded string   `json:"dateadded"`
	URL       string   `json:"url"`
	Threat    string   `json:"threat"`
	Tags      []string `json:"tags"`
}

type Data struct {
	Threats []UrlHausApiBody `json:"threats"`
}

type ThreatMap map[string]Data
