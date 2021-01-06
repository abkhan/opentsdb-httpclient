package tsdbh

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// HttpClient struct represents open-tsdb http client
type HttpClient struct {

	//Number of datapoints to buffer before flushing
	BufferSize int

	//Number of seconds before writitng what we've buffered
	Interval int

	base_url string
	client   *http.Client
	putChan  chan DataPoint
}

// NewHttpClient takes conf and returns an HttpClient
func NewHttpClient(conf Conf) *HttpClient {
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: conf.HTTPConf.DialTimeout,
		}).Dial,
		TLSHandshakeTimeout: conf.HTTPConf.TLSHandshakeTimeout,
		MaxIdleConnsPerHost: conf.HTTPConf.MaxIdleConnsPerHost,
		MaxIdleConns:        conf.HTTPConf.MaxIdleConns,
		IdleConnTimeout:     conf.HTTPConf.IdleConnTimeout,
	}

	t := &HttpClient{
		BufferSize: conf.DefaultBuffer,
		Interval:   conf.DefaultInterval,
		client: &http.Client{
			Transport: netTransport,
			Timeout:   conf.HTTPConf.ClientTimeout,
		},
		base_url: fmt.Sprintf("http://%s:%d/api", conf.Host, conf.Port),
		putChan:  make(chan DataPoint, 1000),
	}

	if t.Interval != 0 && t.BufferSize != 0 {
		go t.writer()
	} else {
		log.Warningf("Not starting writer for tsdb client, if you want one please set interval and buffer")
	}

	return t
}

func (t *HttpClient) PutMany(dps []DataPoint) error {

	var (
		buf  bytes.Buffer
		err  error
		resp *http.Response
		req  *http.Request
	)

	url := fmt.Sprintf("%s/%s", t.base_url, "put")

	gw := gzip.NewWriter(&buf)

	if err = json.NewEncoder(gw).Encode(dps); err != nil {
		return err
	}

	if err = gw.Close(); err != nil {
		return err
	}

	if req, err = http.NewRequest("POST", url, &buf); err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")

	len, start := buf.Len(), time.Now()
	if resp, err = t.client.Do(req); err != nil {
		return err
	}
	log.Debugf("%d bytes in %+v", len, time.Since(start))

	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}

	rv := map[string]interface{}{}
	if err := json.NewDecoder(resp.Body).Decode(&rv); err != nil {
		return err
	}

	if _, ok := rv["error"]; ok {
		errmap := rv["error"].(map[string]interface{})
		return errors.New(errmap["message"].(string))
	}

	return nil
}

// PutOne sends directly
func (t *HttpClient) PutOne(dp DataPoint) {
	buffer := []DataPoint{dp}
	if err := t.PutMany(buffer); err != nil {
		log.Errorf("Failed to putOne data point: %+v", err)
	}
}

// Put sends to a put chan that is put with many when buffer is filled or timeout
func (t *HttpClient) Put(dp DataPoint) {
	t.putChan <- dp
}

func (t *HttpClient) writer() {

	buffer := []DataPoint{}

	ticker := time.NewTicker(time.Duration(t.Interval) * time.Second)
	for {
		select {
		case m := <-t.putChan:
			buffer = append(buffer, m)
			if len(buffer) < t.BufferSize {
				continue
			}
		case <-ticker.C:
			if len(buffer) == 0 {
				continue
			}
		}

		if err := t.PutMany(buffer); err != nil {
			log.Errorf("Failed to put many data points: %+v", err)
		}

		buffer = []DataPoint{}
	}
}

// CloseIdleConnections to close idle transport connections
func (t *HttpClient) CloseIdleConnections() {
	t.client.CloseIdleConnections()
}

//SetClientTimeout to update timeout
func (t *HttpClient) SetClientTimeout(val time.Duration) {
	t.client.Timeout = val
}
