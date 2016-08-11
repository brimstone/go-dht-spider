DHT Spider
==========

Listens to DHT announcements and requests, saving what it can find to elasticsearch.

Usage
-----

First start Elasticsearch:
```
docker run -d --name elasticsearch elasticsearch:2
```

Then start the spider:
```
docker run -d --name spider --link elasticsearch brimstone/go-dht-spider
```

If you'd like, also attach kibana:
```
docker run -d --name kibana --link elasticsearch kibana:4
```


The spider populates the `dht-*` indexes with `infohash`, `request`, and `announce` types.
