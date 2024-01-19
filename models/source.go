package models

type Source struct {
	Market []struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		Category string `json:"category,omitempty"`
	} `json:"market,omitempty"`
	Telegram []struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		Category string `json:"category,omitempty"`
	} `json:"telegram,omitempty"`
	Ransom []struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		Category string `json:"category,omitempty"`
	} `json:"ransom,omitempty"`
	Exploit []struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		Category string `json:"category,omitempty"`
	} `json:"exploit,omitempty"`
	Forum []struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		Category string `json:"category,omitempty"`
	} `json:"forum,omitempty"`
	Discord []struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		Category string `json:"category,omitempty"`
	} `json:"discord,omitempty"`
}
