package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Collection struct {
	Name              string   `json:"name"`
	RouterName        string   `json:"router.name"`
	NumShards         int      `json:"numShards"`
	ReplicationFactor int      `json:"replicationFactor"`
	Shards            []string `json:"shards"`
}

type SolrResponseCollectionList struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
	} `json:"responseHeader"`
	Collections []string `json:"collections"`
}

func (c *Client) GetCollections() (SolrResponseCollectionList, error) {
	var response SolrResponseCollectionList
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/solr/admin/collections?action=LIST", c.HostURL), nil)
	if err != nil {
		return response, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

// CollectionCreationRequest represents the JSON payload for creating a collection.
type CollectionCreationRequest struct {
	Name              string   `json:"name"`
	NumShards         int      `json:"numShards,omitempty"`
	RouterName        string   `json:"router.name,omitempty"`
	ReplicationFactor int      `json:"replicationFactor,omitempty"`
	Shards            []string `json:"shards,omitempty"`
}

// CreateCollection sends a request to SolrCloud to create a new collection.
func (c *Client) CreateCollection(name string, numShards int, routerName string, replicationFactor int, shards []string) error {
	// Construct the request payload
	requestData := CollectionCreationRequest{
		Name:              name,
		NumShards:         numShards,
		RouterName:        routerName,
		ReplicationFactor: replicationFactor,
		Shards:            shards,
	}

	// Marshal the request data into JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	// Create the request
	url := fmt.Sprintf("%s/api/collections", c.HostURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed: status code %d", resp.StatusCode)
	}

	return nil
}
