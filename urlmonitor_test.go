/*
 * Test cases to verify the functionality of Endpoint Monitor
 */

package main

import "fmt"
import "testing"


// TODO can create a MOCK http server

func TestBuildResponse(t *testing.T) {

    resp := BuildResponse("http://testurl.com", 200, 15) 

    if resp == "" {
        t.Errorf("Received incorrect response")
    }

}
