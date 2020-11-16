/*
 * Implements a simple url monitor.
 * Refer README.md for the deployment on K8s
 */

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

/*
 * Config to bootstrap the url monitor
 */
type Config struct {
	// Address to listen on
	Address      string        `json:"address"`
	// Logfile for stdout/debug
	LogFile      string        `json:"logfile"`
        //List of URLs to be monitored
	ExternalUrls []ExternalUrl `json:"externalurls"`
}

/* 
 * External URL to be monitored
 */
type ExternalUrl struct {
	// Host name or IP to be monitored
	Host string `json:"host"`
	// http or https
	Type string `json:"type"`
}

/*
 * Monitor struct implements MonitorIfc interface
 */
type Monitor struct {
	Cfg    *Config
	Client *http.Client
}

/*
 * Interface for future enhancements to stop/start monitor
 */
type MonitorIfc interface {
	GetCurrentMetrics()
}

/*
 * GetCurrentMetrics queries the external URLs and responds with the Prometheus output format
 */
func (m *Monitor) GetCurrentMetrics() string {

	respStr := ""
	for _, v := range m.Cfg.ExternalUrls {
		fmt.Println("==============", v.Host)
		start := time.Now()
		resp, err := m.Client.Get(v.Host)
		end := time.Now()
		elapsed := end.Sub(start)
		if err != nil {
			log.Println(err)
			// Return connectivity unreached on failure cases
			respStr = respStr + BuildResponse(v.Host, 0, elapsed)
			return respStr
		}

		htmlData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			respStr = respStr + BuildResponse(v.Host, 0, elapsed)
			return respStr
		}
		defer resp.Body.Close()
		fmt.Printf("%v\n", resp.Status)
		fmt.Printf(string(htmlData))
		if resp.StatusCode == http.StatusOK {
			fmt.Println("Received 200")
			respStr = respStr + BuildResponse(v.Host, 1, elapsed)
		} else {
			fmt.Println("Received something else ")
			respStr = respStr + BuildResponse(v.Host, 0, elapsed)

		}
	}
	return respStr

}

/*
 * ParseConfig parses the json config file and picks the URLs to be monitored
 */
func ParseConfig(fileName string) (*Config, error) {

	jsonFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	raw, err := ioutil.ReadAll(jsonFile)
	cfg := &Config{}
	err = json.Unmarshal(raw, cfg)
	if err != nil {
		fmt.Println("Unmarshal failed", err)
		return nil, err
	}
	return cfg, nil

}

/*
 * Builds the Prometheus format output if the query hits /metrics
 */
func BuildResponse(url string, status int, d time.Duration) string {

	connStr := fmt.Sprintf("# HELP external_url_up Connectivity status of an endpoint url."+"\n"+
		"# TYPE external_url_up counter"+"\n"+
		"external_url_up{url=\"%s\"} %d\n", url, status)
	respStr := fmt.Sprintf("# HELP external_url_response_ms Latency to reach endpoint url."+"\n"+
		"# TYPE external_url_response_ms counter"+"\n"+
		"external_url_response_ms{url=\"%s\"} %v\n", url, d.Milliseconds())
	result := connStr + respStr
	fmt.Println(result)
	return result

}

func main() {

	cfg, err := ParseConfig("/tmp/config.json")
	if err != nil {
		fmt.Println("ParseConfig returned error", err)
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: http.ProxyFromEnvironment,
		},
		Timeout: 5 * time.Second,
	}

	m := &Monitor{cfg, client}

	http.HandleFunc("/metrics", MetricsHandler(m))
	http.ListenAndServe(cfg.Address, nil)

}

/*
 * Metrics Handler that handles the incoming requests
 */
func MetricsHandler(m *Monitor) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Inside the Handler")
		if r.URL.Path != "/metrics" {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case "GET":
			mResp := m.GetCurrentMetrics()
			w.Write([]byte(mResp))
		default:
			w.WriteHeader(http.StatusNotImplemented)
			v := http.StatusText(http.StatusNotImplemented) + "\n"
			w.Write([]byte(v))

		}
	}

}
