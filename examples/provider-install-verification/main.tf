terraform {
  required_providers {
    solrcloud = {
      source = "hashicorp.com/edu/solrcloud"
    }
  }
}

provider "solrcloud" {
  host     = "http://localhost:8983"
  username = "solr"
  password = "SolrRocks"
}

data "solrcloud_collections" "example" {}
