package restserver

import (
	"bytes"
	"encoding/json"
	"golangbufferedsocket/logging"
	"io/ioutil"
	"net"
	"net/http"
)

// Config configuration for the server(s)
type Config struct {
	ServerType       string `yaml:"ServerType"`
	ServerAddress    string `yaml:"ServerAddress"`
	URL              string `yaml:"URL"`
	ExtraHeader      string `yaml:"ExtraHeader"`
	ExtraHeaderValue string `yaml:"ExtraHeaderValue"`
	MaxRetries       int    `yaml:"MaxRetries"`
}

// Server to "buffer" socket data an send it as http request
func Server(c net.Conn, serverConfig Config) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		logging.GetInfoLogger().Printf("Received data %s", string(data))

		// http api handling
		reqBody, err := json.Marshal(string(data))
		if err != nil {
			logging.GetErrorLogger().Printf("Could nod create body for %s", string(data))
			continue
		}
		if serverConfig.MaxRetries == 0 {
			serverConfig.MaxRetries = 5
		}

		var resp *http.Response

		for i := 0; i < serverConfig.MaxRetries; i++ {
			logging.GetInfoLogger().Println("Sending post request")

			if serverConfig.URL == "" {
				logging.GetInfoLogger().Println("No URL not trying to send data")
				break
			}

			resp, err = http.Post(serverConfig.URL,
				"application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				print(err)
			}

			if resp.StatusCode != 201 {
				continue
			}

			defer resp.Body.Close()
			_, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				print(err)
			}

			break
		}

		if resp != nil && resp.StatusCode != 201 {
			logging.GetErrorLogger().Printf("Could not send the following data: %s", string(data))
		}
		logging.GetInfoLogger().Println("Sending done")
	}
}
