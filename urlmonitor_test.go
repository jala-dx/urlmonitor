/*
 * Test cases to verify the functionality of Endpoint Monitor
 */

package main

import "testing"

// TODO can create a MOCK http server and run HTTP get

/*
 * Test the REsponse builder
 */
func TestBuildResponse(t *testing.T) {

	resp := BuildResponse("http://testurl.com", 200, 15)

	if resp == "" {
		t.Errorf("Received incorrect response")
	}

}

/*
 * Test if the config file is parsed correctly
 */
func TestParseConfig(t *testing.T) {

	config, err := ParseConfig("./config_test.json")
	if err != nil {
		t.Errorf("Failed to parse the config file")
	}

	if config.Address !=  "0.0.0.0:2112" {
		t.Errorf("incorrect address in the config file")
	}

	if len(config.ExternalUrls) != 3 {
		t.Errorf("Expected 3 URLs, but found %d", len(config.ExternalUrls) )
	}

}
