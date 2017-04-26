//
//  randomorg.go
//  Random input - random.org web request
//
//  Created by Adam Hosier on 2017-0309.
//  Copyright © 2017 Adam Hosier. All rights reserved.
//

package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"fmt"
)

type Request struct {
	JsonRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  *RequestParams  `json:"params"`
	Id int                  `json:"id"`
}

type RequestParams struct {
	ApiKey string `json:"apiKey"`
	N      int    `json:"n"`
	Size   int    `json:"size"`
}

type Response struct {
	Error *struct {
		Code    int
		Message string
	}
	Result *struct {
		Random *struct {
			Data []string
		}
	}
}

func main() {
	url := "https://api.random.org/json-rpc/1/invoke"
	method := "generateBlobs"
	params := RequestParams{"c463485b-1a4d-4f26-af04-792affc469bb", 1, 256}
	request := &Request{JsonRPC: "2.0", Method: method, Params: &params, Id: 1}

	// Build request
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(request)

	// Send request
	resp, err := http.Post(url, "application/json-rpc", b)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("Invalid return code: %d", resp.StatusCode)
	}

	// Extract body
	var data Response
	json.NewDecoder(resp.Body).Decode(&data)

	// Handle json-rpc errors
	if data.Error != nil {
		log.Fatalf("JSON response error: %s", data.Error.Message)
	}

	// Collect result
	result := data.Result.Random.Data[0]
	fmt.Println(result[:len(result) - 1])
}
