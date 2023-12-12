package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Collection struct {
	Name              string   `json:"name"`
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
	ReplicationFactor int      `json:"replicationFactor,omitempty"`
	Shards            []string `json:"shards,omitempty"`
}

// CreateCollection sends a request to SolrCloud to create a new collection.
// return resp and error
func (c *Client) CreateCollection(ctx context.Context, name string, numShards int, replicationFactor int, shards []string) (http.Response, error) {
	// Construct the request payload
	// remove nil values from shards
	requestData := CollectionCreationRequest{
		Name:              name,
		NumShards:         numShards,
		ReplicationFactor: replicationFactor,
		Shards:            shards,
	}

	// tflog
	tflog.Info(ctx, fmt.Sprintf("Creating collection: %s", name))

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return http.Response{}, fmt.Errorf("error marshalling request data: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Request data: %s", jsonData))

	// Create the request
	url := fmt.Sprintf("%s/api/collections", c.HostURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return http.Response{}, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return http.Response{}, fmt.Errorf("error executing request: %w", err)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return http.Response{}, fmt.Errorf("error reading response body: %w", err)
	}

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return http.Response{}, fmt.Errorf("error creating collection: %s", body)
	}

	return *resp, nil
}

type CollectionStatusResponse struct {
	ResponseHeader struct {
		Status int `json:"status"`
		QTime  int `json:"QTime"`
	} `json:"responseHeader"`
	CollectionProperties map[string]struct {
		ZnodeVersion       int               `json:"znodeVersion"`
		Properties         map[string]string `json:"properties"`
		ActiveShards       int               `json:"activeShards"`
		InactiveShards     int               `json:"inactiveShards"`
		SchemaNonCompliant []string          `json:"schemaNonCompliant"`
		Shards             map[string]struct {
			State    string `json:"state"`
			Range    string `json:"range"`
			Replicas map[string]struct {
				Total          int `json:"total"`
				Active         int `json:"active"`
				Down           int `json:"down"`
				Recovering     int `json:"recovering"`
				RecoveryFailed int `json:"recovery_failed"`
			} `json:"replicas"`
			Leader struct {
				CoreNode string `json:"coreNode"`
			} `json:"leader"`
		} `json:"shards"`
	} `json:"collectionName"`
}

type CollectionStatusResponse2 struct {
	ResponseHeader ResponseHeader `json:"responseHeader"`
	Cluster        ClusterInfo    `json:"cluster"`
}

type ResponseHeader struct {
	Status int `json:"status"`
	QTime  int `json:"QTime"`
}

type ClusterInfo struct {
	Collections map[string]CollectionInfo `json:"collections"`
	LiveNodes   []string                  `json:"live_nodes"`
}

type CollectionInfo struct {
	PullReplicas      string               `json:"pullReplicas"`
	ConfigName        string               `json:"configName"`
	ReplicationFactor int                  `json:"replicationFactor"`
	Router            RouterInfo           `json:"router"`
	NrtReplicas       int                  `json:"nrtReplicas"`
	TlogReplicas      string               `json:"tlogReplicas"`
	Shards            map[string]ShardInfo `json:"shards"`
	Health            string               `json:"health"`
	ZnodeVersion      int                  `json:"znodeVersion"`
}

type RouterInfo struct {
	Name string `json:"name"`
}

type ShardInfo struct {
	Range    string                 `json:"range"`
	State    string                 `json:"state"`
	Replicas map[string]ReplicaInfo `json:"replicas"`
	Health   string                 `json:"health"`
}

type ReplicaInfo struct {
	Core          string `json:"core"`
	NodeName      string `json:"node_name"`
	Type          string `json:"type"`
	State         string `json:"state"`
	Leader        string `json:"leader"`
	ForceSetState string `json:"force_set_state"`
	BaseURL       string `json:"base_url"`
}

func (c *Client) GetCollectionStatus(collectionName string) (CollectionInfo, error) {
	var collectionStatus CollectionStatusResponse2

	// Construct the URL for the collection status API
	url := fmt.Sprintf("%s/api/collections/%s", c.HostURL, collectionName)

	// Create an HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return CollectionInfo{}, fmt.Errorf("error creating request: %w", err)
	}

	// Perform the HTTP request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return CollectionInfo{}, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	// Read and unmarshal the response
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &collectionStatus)
	if err != nil {
		return CollectionInfo{}, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return collectionStatus.Cluster.Collections[collectionName], nil
}
