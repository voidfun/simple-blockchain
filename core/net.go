package main

/*
import (
	"bytes"
	"io/ioutil"
	"net/http"
)


import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func (bc *Blockchain) syncBlocks() {
	longestChain := bc.Blocks
	for url, _ := range bc.Nodes {
		listBlocksUrl := url + "/blocks"
		resp, err := httpGet(listBlocksUrl)
		if err != nil {
			log.Fatal(err)
			continue
		}
		var blocks []Block
		json.Unmarshal(resp, &blocks)
		if len(blocks) > len(longestChain) {
			longestChain = blocks
		}
	}
	blockchain.replaceChain(longestChain)
}

func (bc *Blockchain) syncAll() {
	bc.syncNodes()

	for url, _ := range bc.Nodes {
		syncBlocksUrl := url + "/syncBlocks"
		resp, err := httpGet(syncBlocksUrl)
		if err != nil {
			log.Fatal(err)
			continue
		}
		var blocks []Block
		json.Unmarshal(resp, &blocks)
	}
}

func (bc *Blockchain) syncNodes() {
	urls := make([]string, 0, len(bc.Nodes))
	for k := range bc.Nodes {
		urls = append(urls, k)
	}
	m := Message{}
	m.URLs = urls
	jsonBytes, _ := json.Marshal(m)

	for url, _ := range bc.Nodes {
		nodesUrl := url + "/nodes"
		resp, err := httpPost(nodesUrl, string(jsonBytes))
		if err != nil {
			return
		}
		log.Println(resp)
	}
}

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return []byte(""), err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatal(err)
		return []byte(""), err
	}
	return body, nil
}

func httpPost(url string, data string) (bool, error) {
	log.Println("[HTTP POST] url:" + url)
	log.Println("[HTTP POST] data:" + data)
	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
	return true, nil
}
*/
