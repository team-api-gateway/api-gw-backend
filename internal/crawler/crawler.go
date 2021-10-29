package crawler

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
)

type Crawler struct {
	openAPIUrls []string
	httpClient  http.Client
}

func New(urls []string) *Crawler {
	return &Crawler{
		openAPIUrls: urls,
		httpClient:  http.Client{Transport: &http.Transport{Proxy: nil, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}},
	}
}

func (c *Crawler) LoadSpecs() ([]*openapi3.T, error) {
	var specs []*openapi3.T
	for _, url := range c.openAPIUrls {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Printf("skipping '%s': %s\n", url, err.Error())
			continue
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			fmt.Printf("skipping '%s': %s\n", url, err.Error())
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("skipping '%s': %s\n", url, err.Error())
			continue
		}
		doc, err := loadOpenApi3(body)
		if err != nil {
			fmt.Printf("skipping '%s': %s\n", url, err.Error())
			continue
		}
		if doc.OpenAPI == "" {
			doc2, err := loadOpenApi2(body)
			if err != nil {
				fmt.Printf("skipping '%s': %s\n", url, err.Error())
				continue
			}
			doc, err = openapi2conv.ToV3(doc2)
			if err != nil {
				fmt.Printf("skipping '%s': %s\n", url, err.Error())
				continue
			}
		}
		specs = append(specs, doc)
	}
	return specs, nil
}
func loadOpenApi3(data []byte) (*openapi3.T, error) {
	return openapi3.NewLoader().LoadFromData(data)
}
func loadOpenApi2(data []byte) (*openapi2.T, error) {
	var doc openapi2.T
	if err := json.Unmarshal(data, &doc); err != nil {
		if err = yaml.Unmarshal(data, &doc); err != nil {
			return nil, err
		}
	}
	return &doc, nil
}
