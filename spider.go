package main

import (
	"encoding/hex"
	"log"

	"github.com/shiyanhui/dht"
)

func main() {
	downloader := dht.NewWire()
	go func() {
		// once we got the request result
		for resp := range downloader.Response() {
			log.Println(string(resp.InfoHash), string(resp.MetadataInfo))
		}
	}()
	go downloader.Run()

	config := dht.NewCrawlConfig()
	log.Println("Spider configured")
	/*
		config.OnGetPeers = func(infohash, ip string, port int) {
			log.Println("OnGetPeers:", hex.EncodeToString([]byte(infohash)), ip, port)
		}
	*/
	config.OnAnnouncePeer = func(infoHash, ip string, port int) {
		log.Println("Attempting to download", hex.EncodeToString([]byte(infoHash)), "from", ip, port)
		// request to download the metadata info
		downloader.Request([]byte(infoHash), ip, port)
	}
	d := dht.New(config)

	log.Println("Spider started")
	d.Run()
}
