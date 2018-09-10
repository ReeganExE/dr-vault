package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

type MockClient struct {
	method string
	body []byte
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	m.method = req.Method
	m.body, _ = ioutil.ReadAll(req.Body)
	return &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString("{}")),
		StatusCode: 200,
		Status: "200 OK",
	}, nil
}

func TestWrite(t *testing.T) {
	muck :=  &MockClient{}
	var client = &VaultClient{
		http:   muck,
	}

	data := map[string]int{"one": 1,}

	client.Write("secret", data)

	if muck.method != "PUT" {
		t.Fatalf("Expected method to be PUT, but got %s", muck.method)
	}

	if string(muck.body) != `{"one":1}` {
		t.Fatalf("Body is not expected, got %s", string(muck.body))
	}
}

func TestWriteEmpty(t *testing.T) {
	muck :=  &MockClient{}
	var client = &VaultClient{
		http:   muck,
	}

	data := map[string]int{}

	client.Write("secret", data)

	if muck.method != "PUT" {
		t.Fatalf("Expected method to be PUT, but got %s", muck.method)
	}

	if string(muck.body) != `{"":""}` {
		t.Fatalf("Expected an empty body, but got %s", string(muck.body))
	}
}

func TestNewVaultClient(t *testing.T) {
	client := NewVaultClient("1.11.2.3:88", "")
	if client.address != "http://1.11.2.3:88" {
		t.Fatalf("Expected a vault address, but got %s", client.address)
	}
}
