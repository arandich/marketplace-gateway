package handlers

import (
	"bytes"
	imagecdn "github.com/arandich/marketplace-sdk/image-cdn"
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"strings"
)

type CdnClient struct {
	client imagecdn.IClient
}

func NewCdnClientHandler(client imagecdn.IClient) CdnClient {
	return CdnClient{client: client}
}

type Response struct {
	Url   string `json:"url,omitempty"`
	Error string `json:"error,omitempty"`
}

func (c *CdnClient) UploadImg(w http.ResponseWriter, r *http.Request) {
	var resp Response
	var err error
	var buf bytes.Buffer
	defer func() {
		if err != nil {
			resp.Error = err.Error()
			respB, _ := json.Marshal(resp)
			w.Write(respB)
		}
	}()
	w.Header().Set("Content-Type", "application/json")

	err = r.ParseMultipartForm(32 << 18)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	name := strings.Split(header.Filename, ".")
	_, err = io.Copy(&buf, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	url, err := c.client.UploadImage(r.Context(), buf.Bytes(), name[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.Url = url

	respB, _ := json.Marshal(resp)
	w.Write(respB)

}
