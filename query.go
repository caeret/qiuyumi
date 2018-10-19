package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var client = &http.Client{Timeout: time.Second * 10}

type result struct {
	Status    bool   `json:"status"`
	Available bool   `json:"available"`
	From      string `json:"from"`
}

func query(name string) (n string, available bool, err error) {
	name = strings.Trim(name, ".")
	domain := "com"
	if idx := strings.LastIndex(name, "."); idx != -1 {
		domain = name[idx+1:]
		name = name[:idx]
	}

	n = name + "." + domain

	URL := fmt.Sprintf("http://www.qiuyumi.com/query/whois.%s.php", domain)
	query := make(url.Values)
	query.Set("name", name)
	req, _ := http.NewRequest(http.MethodPost, URL, strings.NewReader(query.Encode()))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.67 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://www.qiuyumi.com")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("fail to read response body with status %d: %s", resp.StatusCode, err.Error())
		return
	}

	var result result
	err = json.Unmarshal(b, &result)
	if err != nil {
		err = fmt.Errorf("fail to unmarshal result %s: %s", string(b), err.Error())
		return
	}

	if !result.Status {
		err = errors.New("invalid result status")
		return
	}

	available = result.Available

	return
}
