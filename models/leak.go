package models

import "time"

type Email struct {
	Email string `json:"email"`
}

type LeakResponse struct {
	Success  bool `json:"success"`
	Breaches []struct {
		Name        string    `json:"Name"`
		Title       string    `json:"Title"`
		Domain      string    `json:"Domain"`
		BreachDate  time.Time `json:"BreachDate"`
		DataClasses []string  `json:"DataClasses"`
	} `json:"breaches"`
	Total int `json:"total"`
}
