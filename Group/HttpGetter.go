package group

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type HttpGetter struct {
	base_url string
}

func NewHttpGetter(base_url string) *HttpGetter {
	return &HttpGetter{
		base_url: base_url,
	}
}

func (h *HttpGetter) Get(group, key string) ([]byte, error) {
	client := http.Client{}
	log.Println("ask", h.base_url+group+"/"+key)
	req, err := http.NewRequest(http.MethodGet, h.base_url+group+"/"+key, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", resp.Status)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	return bytes, err
}
