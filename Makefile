BINARY = spider

PUBLISH = local
local:
	docker save brimstone/go-dht-spider > spider.tar

load:
	docker load < spider.tar

run:
	docker-1.10.0 run --rm -it brimstone/go-dht-spider

include ${PROJECTBUILDER}/Makefile
