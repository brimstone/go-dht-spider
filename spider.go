package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shiyanhui/dht"
)

type File struct {
	Path   []interface{} `json:"path"`
	Length int           `json:"length"`
}
type BitTorrent struct {
	Timestamp int64  `json:"@timestamp"`
	InfoHash  string `json:"infohash"`
	Name      string `json:"name"`
	Files     []File `json:"files,omitempty"`
	Length    int    `json:"length,omitempty"`
}

type Peer struct {
	Timestamp int64  `json:"@timestamp"`
	InfoHash  string `json:"infohash"`
	Address   string `json:"address"`
	Port      int    `json:"port"`
}

func Nowish() string {
	return time.Now().UTC().Format("20060102")
}

func main() {
	w := dht.NewWire()

	es := "http://" + os.Getenv("ELASTICSEARCH_PORT_9200_TCP_ADDR") + ":" + os.Getenv("ELASTICSEARCH_PORT_9200_TCP_PORT")
	DeleteTemplate(es)
	SetTemplate(es)

	go func() {
		for resp := range w.Response() {
			metadata, err := dht.Decode(resp.MetadataInfo)
			if err != nil {
				continue
			}
			info := metadata.(map[string]interface{})
			if _, ok := info["name"]; !ok {
				continue
			}
			bt := BitTorrent{
				Timestamp: time.Now().Unix(),
				InfoHash:  hex.EncodeToString(resp.InfoHash),
				Name:      info["name"].(string),
			}
			if v, ok := info["files"]; ok {
				files := v.([]interface{})
				bt.Files = make([]File, len(files))
				for i, item := range files {
					file := item.(map[string]interface{})
					bt.Files[i] = File{
						Path:   file["path"].([]interface{}),
						Length: file["length"].(int),
					}
				}
			} else if _, ok := info["length"]; ok {
				bt.Length = info["length"].(int)
			}
			data, err := json.Marshal(bt)
			if err == nil {
				fmt.Printf("%s\n\n", data)
			}

			err = Index(es+"/dht-"+Nowish()+"/infohash/"+string(bt.InfoHash), data)
			if err != nil {
				log.Println("Error saving to ES:", err)
			}

		}
	}()
	go w.Run()

	config := dht.NewCrawlConfig()
	log.Println("Spider configured")

	config.OnAnnouncePeer = func(infoHash, ip string, port int) {
		log.Println("Attempting to download", hex.EncodeToString([]byte(infoHash)), "from", ip, port)

		peer, _ := json.Marshal(Peer{
			Timestamp: time.Now().Unix(),
			InfoHash:  hex.EncodeToString([]byte(infoHash)),
			Address:   ip,
			Port:      port,
		})

		err := Index(es+"/dht-"+Nowish()+"/announce/", peer)
		if err != nil {
			log.Println("Error saving to ES:", err)
		}
		// request to download the metadata info
		w.Request([]byte(infoHash), ip, port)
	}

	config.OnGetPeers = func(infoHash, ip string, port int) {
		peer, _ := json.Marshal(Peer{
			Timestamp: time.Now().Unix(),
			InfoHash:  hex.EncodeToString([]byte(infoHash)),
			Address:   ip,
			Port:      port,
		})

		err := Index(es+"/dht-"+Nowish()+"/request/", peer)
		if err != nil {
			log.Println("Error saving to ES:", err)
		}
	}

	d := dht.New(config)

	log.Println("Spider started")
	d.Run()
}
