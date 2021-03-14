terraform {
  required_providers {
    solrcluster = {
      version = "0.1"
      source  = "hashicorp.com/edu/solrcluster"
    }
  }
}

variable "collection_name" {
  type    = string
  default = "collection1"
}

data "solrcluster_collections" "all" {}

# Returns all collections
output "all_collections" {
  value = data.solrcluster_collections.all.collections
}

# Only returns collection1 
output "collection" {
  value = {
    for collection in data.solrcluster_collections.all.collections :
    collection.id => collection
    if collection.name == var.collection_name
  }
}
