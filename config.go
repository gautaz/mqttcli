package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	simpleJson "github.com/bitly/go-simplejson"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const DefaultConfigFile = ".mqttcli.cfg" // Under HOME

type Config struct {
  Broker   string `json:"broker"`
	UserName string `json:"username"`
	Password string `json:"password"`

	CaCert     string `json:"caCert"`
	ClientCert string `json:"clientCert"`
	PrivateKey string `json:"privateKey"`
}

func (c *Config) UnmarshalJSON(data []byte) error {
	js, err := simpleJson.NewJson(data)
	if err != nil {
		return err
	}
	if c.Broker, err = js.Get("broker").String(); err != nil {
		c.Broker = ""
	}
	if c.UserName, err = js.Get("username").String(); err != nil {
		c.UserName = ""
	}
	if c.Password, err = js.Get("password").String(); err != nil {
		c.Password = ""
	}
	if c.CaCert, err = js.Get("caCert").String(); err != nil {
		c.CaCert = ""
	}
	if c.ClientCert, err = js.Get("clientCert").String(); err != nil {
		c.ClientCert = ""
	}
	if c.PrivateKey, err = js.Get("privateKey").String(); err != nil {
		c.PrivateKey = ""
	}
	return nil
}

func readFromConfigFile(path string) (Config, error) {
	ret := Config{}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ret, err
	}

	err = json.Unmarshal(b, &ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
func getSettingsFromFile(p string, opts *MQTT.ClientOptions) error {
	confPath := ""
	home := UserHomeDir()
	// replace home to ~ in order to match
	p = strings.Replace(p, home, "~", 1)
	if p == "~/.mqttcli.cfg" || p == "" {
		confPath = path.Join(home, DefaultConfigFile)
		_, err := os.Stat(confPath)
		if os.IsNotExist(err) {
			return err
		}
	} else {
		confPath = p
	}

	ret, err := readFromConfigFile(confPath)
	if err != nil {
		log.Error(err)
		return err
	}

	tlsConfig, ok, err := makeTlsConfig(ret.CaCert, ret.ClientCert, ret.PrivateKey, false)
	if err != nil {
		return err
	}
	if ok {
		opts.SetTLSConfig(tlsConfig)
	}

	if ret.Broker != "" {
		log.Infof("Broker URI: %s", ret.Broker)
		opts.AddBroker(ret.Broker)
	}

	if ret.UserName != "" {
		opts.SetUsername(ret.UserName)
	}
	if ret.Password != "" {
		opts.SetPassword(ret.Password)
	}
	return nil
}
