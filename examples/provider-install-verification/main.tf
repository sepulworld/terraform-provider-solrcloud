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

/*
resource "solrcloud_collection" "test" {
  name = "test"
  num_shards = 1
  replication_factor = 1
  config_name = "myconfig"
  router_name = "compositeId"
  router_field = "id"
  shards = ["shard1"]
  max_shards_per_node = 1
  auto_add_replicas = false
  create_node_set = false
  rule = "shard:*,replica:1"
  snitch = "org.apache.solr.cloud.DefaultSnitch"
  collection_config_name = "myconfig"
  collection_config_set = "myconfig"
  collection_config_data = <<EOF
  <config>
    <luceneMatchVersion>8.8.2</luceneMatchVersion>
    <directoryFactory name="DirectoryFactory" class="${solr.directoryFactory:solr.NRTCachingDirectoryFactory}">
      <str name="solr.lock.type">${solr.lock.type:native}</str>
    </directoryFactory>
    <indexConfig>
      <lockType>${solr.lock.type:native}</lockType>
    </indexConfig>
    <query>
      <maxBooleanClauses>${solr.max.booleanClauses:1024}</maxBooleanClauses>
      <filterCache class="solr.FastLRUCache" size="512" initialSize="512" autowarmCount="0"/>
      <queryResultCache class="solr.LRUCache" size="512" initialSize="512" autowarmCount="0"/>
      <documentCache class="solr.LRUCache" size="512" initialSize="512" autowarmCount="0"/>
      <enableLazyFieldLoading>true</enableLazyFieldLoading>
      <queryResultWindowSize>20</queryResultWindowSize>
      <queryResultMaxDocsCached>200</queryResultMaxDocsCached>
      <useColdSearcher>false</useColdSearcher>
      <maxWarmingSearchers>2</maxWarmingSearchers>
    </query>
    <updateHandler class="solr.DirectUpdateHandler2">
      <updateLog>
        <str name="dir">${solr.ulog.dir:}</str>
        <int name="numVersionBuckets">${solr.ulog.numVersionBuckets:65536}</int>
      </updateLog>
    </updateHandler>
    <queryResponseWriter name="json" class="solr.JSONResponseWriter">
EOF
}
*/

resource "solrcloud_collection" "name" {
  name = "test" 
  num_shards = 1
  replication_factor = 1
  router = "compositeId"
}
