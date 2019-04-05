// Package jsonconfig provides ...
package jsonconfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//Configuration  Store the Port for the application
type Configuration struct {
	ServerPort string
}

var config Configuration
var err error

//ReadConfig will read the config json file
//if read json file success ,it will return config and nil
func ReadConfig(filename string) (Configuration, error) {
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Print("Unable to Read Config file")
		return Configuration{}, err
	}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Print(err)
		return Configuration{}, err
	}
	return config, nil

}
