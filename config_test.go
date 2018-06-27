package main

import (
	"encoding/json"
	"testing"
)

func Test_ConfigUnmarshalJSON(t *testing.T) {
	c := Config{}

	j := `{"broker": "URI"}`
	err := json.Unmarshal([]byte(j), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Broker != "URI" || c.UserName != "" || c.Password != "" {
		t.Error("parse failed")
	}

	j = `{"broker": "URI", "username": "u"}`
	err = json.Unmarshal([]byte(j), &c)
	if err != nil {
		t.Error(err)
	}
	if c.Broker != "URI" || c.UserName != "u" || c.Password != "" {
		t.Error("parse failed with username set")
	}
}

func Test_ConfigCert(t *testing.T) {
	c := Config{}

	j := `{"caCert": "ca", "clientCert": "client", "privateKey": "key"}`
	err := json.Unmarshal([]byte(j), &c)
	if err != nil {
		t.Error(err)
	}
	if c.CaCert != "ca" || c.ClientCert != "client" || c.PrivateKey != "key" {
		t.Error("parse failed, %#v", c)
	}
}
