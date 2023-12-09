package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Collection struct {
	Name string `json:"name"`
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
