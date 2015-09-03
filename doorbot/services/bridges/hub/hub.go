package hub

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"
)

// HubConfig represents hub configuration
type HubConfig struct {
	AccountID uint
	URL       string
	Token     string
}

// HubUser user data
type HubUser struct {
	ID          uint
	DisplayName *string
	Email       string
}

// Hub configuration data
type Hub struct {
	AccountID uint
	Token     string
	URL       string
}

// NewHub creates a new Hub instance from a given HubConfig
func New(c *HubConfig) *Hub {
	return &Hub{
		AccountID: c.AccountID,
		URL:       c.URL,
		Token:     c.Token,
	}
}

// GetUsers returns a list of users from the hub.
func (h *Hub) GetUsers() ([]*HubUser, error) {
	var users []*HubUser

	url := fmt.Sprintf("%s%s", h.URL, "/api/v1/people")

	log.WithFields(log.Fields{
		"account_id": h.AccountID,
		"url":        url,
	}).Info("Bridges::Hub::GetUsers requesting users", url)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("token %s", h.Token))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.WithFields(log.Fields{
			"account_id": h.AccountID,
			"url":        url,
			"error":      err,
		}).Error("Bridges::Hub::GetUsers http error")

		return users, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&users)

	if err != nil {
		log.WithFields(log.Fields{
			"account_id": h.AccountID,
			"url":        url,
			"error":      err,
		}).Error("Bridges::Hub::GetUsers JSON error")

		return users, err
	}

	log.WithFields(log.Fields{
		"account_id": h.AccountID,
		"url":        url,
		"users":      len(users),
	}).Info("Bridges::Hub::GetUsers Users received")

	return users, err
}
