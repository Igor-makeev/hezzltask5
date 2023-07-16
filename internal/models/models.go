package models

import "time"

type CreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

type DeleteResponse struct {
	Id         int  `json:"id" `
	CampaignId int  `json:"campaignId" `
	Removed    bool `json:"removed,omitempty" `
}

type Item struct {
	Id          int       `json:"id" `
	CampaignId  int       `json:"campaignId" `
	Name        string    `json:"name" `
	Description string    `json:"description" `
	Priority    int       `json:"priority" `
	Removed     bool      `json:"removed" `
	CreatedAt   time.Time `json:"createdAt" `
}

type Event struct {
	Id          int       `json:"id" `
	CampaignId  int       `json:"campaignId" `
	Name        string    `json:"name" `
	Description string    `json:"description" `
	Priority    int       `json:"priority" `
	Removed     bool      `json:"removed" `
	EventTime   time.Time `json:"eventTime" `
}
