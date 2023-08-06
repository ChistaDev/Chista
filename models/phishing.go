package models

type PhishingDomain struct {
	ID     uint64 `json:"id"`
	Domain string `json:"domain" binding:"required"`
}
