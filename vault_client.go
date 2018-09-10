package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Client interface {
	Write(path string, data interface{}) error
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type VaultClientInterface interface {
	Write(path string, data interface{}) error
}

type VaultClient struct {
	address string
	token   string
	http    HttpClient
}

// NewVaultClient creates a new Vault fate with an associated token
func NewVaultClient(address, token string) *VaultClient {
	host, port := parseAddr(address)

	// url.Parse didn't work with IP:PORT
	// parse, e := url.Parse(address)

	return &VaultClient{
		address: fmt.Sprintf("http://%s:%s", host, port), // parse.String(),
		token:   token,
		http:    &http.Client{},
	}
}

func parseAddr(addr string) (string, string) {
	arr := strings.Split(addr, ":")
	return arr[0], arr[1]
}

func (c *VaultClient) Write(path string, data interface{}) error {
	bits, e := json.Marshal(data)
	if e != nil {
		panic(e)
	}

	if len(bits) == 2 {
		// Pass an empty object to avoid complaint from Vault
		bits = []byte(`{"":""}`)
	}

	request, _ := http.NewRequest("PUT", fmt.Sprintf("%s/v1/%s", c.address, path), bytes.NewBuffer(bits))
	request.Header.Add("x-vault-token", c.token)
	request.Header.Add("content-type", "application/json")

	res, e := c.http.Do(request)
	if e != nil {
		return e
	}

	if res.StatusCode >= 300 {
		all, _ := ioutil.ReadAll(res.Body)
		log.Fatalf("Vault error code: %s. Body: %s", res.Status, string(all))
		return errors.New("failed to write")
	}

	return nil
}
