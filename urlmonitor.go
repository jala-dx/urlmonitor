/*
 * Implements a simple endpoint monitor.
 * Refer README.md for the deployment on K8s
 */

package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"time"
	"os"
	// "sync"
	"encoding/json"
	"net/http"

)

/*
 * Config to bootstrap the ep monitor
 */

type Config struct {
    Address           string        `json:"address"`
    LogFile           string        `json:"logfile"`
    ExternalUrls      []ExternalUrl `json:"externalurls"`
}

type ExternalUrl struct {
   Host string     `json:"host"`
   Type string     `json:"type"`

}

type Monitor struct {
    Cfg    *Config
    Client *http.Client
}

type MonitorIfc interface {
    StartMonitor()
    GetCurrentMetrics()
    StopMonitor()
}


func (m *Monitor) GetCurrentMetrics() string{

	respStr := ""
	for _, v := range m.Cfg.ExternalUrls {
		fmt.Println("==============", v.Host)
		start := time.Now()
		resp, err := m.Client.Get(v.Host)
		end := time.Now()
		elapsed := end.Sub(start)
		if err != nil {
			log.Println(err)
			return respStr
		}

		htmlData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return respStr
		}
		defer resp.Body.Close()
		fmt.Printf("%v\n", resp.Status)
		fmt.Printf(string(htmlData))
		if resp.StatusCode == http.StatusOK {
			fmt.Println("Received 200")
			respStr = respStr + BuildResponse(v.Host, 1, elapsed)
		}else {
			fmt.Println("Received something else ")
			respStr = respStr + BuildResponse(v.Host, 0, elapsed)

		}
	}
	return respStr

}

func (m *Monitor) StartMonitor() {

	go func() {
		for {
			for _, v := range m.Cfg.ExternalUrls {
				fmt.Println("==============", v.Host)
				resp, err := m.Client.Get(v.Host)
				if err != nil {
					log.Println(err)
					return
				}

				htmlData, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer resp.Body.Close()
				fmt.Printf("%v\n", resp.Status)
				fmt.Printf(string(htmlData))
			}
		}
	}()


}

func (m *Monitor) StopMonitor() {

}

func MonitorEP(m MonitorIfc) {

	// Not needed yet, for future
	//m.StartMonitor()

	// Not needed yet, for future
	m.StopMonitor()

}

func ParseConfig() (error, *Config) {

	jsonFile, err := os.Open("/tmp/config.json")
	if err != nil {
		//TODO log
		return err, nil
	}
	defer jsonFile.Close()
	
	raw, err := ioutil.ReadAll(jsonFile)
	cfg := &Config{}
	err = json.Unmarshal(raw, cfg)
	if err != nil {
		fmt.Println("Unmarshal failed", err)
		return err, nil
	}
	return  nil, cfg
	

}


/* Sample Scrape

# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 12720
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0


# HELP myapp_processed_ops_total The total number of processed events
# TYPE myapp_processed_ops_total counter
myapp_processed_ops_total 2240
*/

func BuildResponse(url string, status int, d time.Duration) string {

	connStr := fmt.Sprintf("# HELP external_url_up Connectivity status of an endpoint url." + "\n" + 
                   "# TYPE external_url_up counter" + "\n" +
                   "external_url_up{url=\"%s\"} %d\n", url, status)
	respStr := fmt.Sprintf("# HELP external_url_response_ms Latency to reach endpoint url." + "\n" + 
                   "# TYPE external_url_response_ms counter" + "\n" +
                   "external_url_response_ms{url=\"%s\"} %v\n", url, d.Milliseconds())
	result := connStr + respStr
	fmt.Println(result)
	return result

}

func main() {
	// TODO fix the log file
	log.SetFlags(log.Lshortfile)

        err, cfg := ParseConfig()
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
		// TODO make it configurable
		Timeout: 5  * time.Second,
	}

	m := &Monitor{cfg, client}
	//MonitorEP(m)

        http.HandleFunc("/metrics", MetricsHandler(m))
        http.ListenAndServe(":2112", nil)

}

/*
 * Metrics Handler that handles the incoming requests
 */
func MetricsHandler(m *Monitor) http.HandlerFunc {

	return func (w http.ResponseWriter, r *http.Request) {
		fmt.Println("Inside the Handler")
		if r.URL.Path != "/metrics" {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
			case "GET": 
				//w.Write([]byte("Received a GET request\n"))
				//resp := BuildResponse("http://jalaja.com", 1)
				//w.Write([]byte(resp))
				mResp := m.GetCurrentMetrics()
				w.Write([]byte(mResp))
			default:
				w.WriteHeader(http.StatusNotImplemented)
				v := http.StatusText(http.StatusNotImplemented) + "\n"
				w.Write([]byte(v))

		}
	}

}
