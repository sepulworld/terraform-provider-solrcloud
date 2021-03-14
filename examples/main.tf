terraform {
  required_providers {
    solrcluster = {
      version = "0.1"
      source  = "hashicorp.com/edu/solrcluster"
    }
  }
}

provider "hashicups" {}

module "psl" {
  source = "./solrcluster"

  collection_name = "collection1"
}

output "psl" {
  value = module.psl.collection
}
