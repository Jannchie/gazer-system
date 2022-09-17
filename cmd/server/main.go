package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/jannchie/gazer-system/internal/gs"
)

type TemporaryError struct {
	error
}

func (t *TemporaryError) Temporary() bool {
	return true
}

func main() {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config := gs.NewDefaultConfig()
	config.CollectHandle = func(c *gs.Collector, targetURL string) ([]byte, error) {
		// set user agent

		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			log.Fatalln(err)
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36 Edg/105.0.1343.42")
		resp, err := c.Client.Do(req)
		if err != nil {
			return nil, err
		}
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			// if is 4XX error, should wait for proxy refresh.
			err := c.RefreshClient()
			if err != nil {
				return nil, TemporaryError{err}
			}
			return nil, TemporaryError{fmt.Errorf("status code error: %d", resp.StatusCode)}
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	gss := gs.NewServer(config)
	gss.Run()
}
