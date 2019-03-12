package request

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormContentType(t *testing.T) {
	req, err := Post("http:....").Form("Key", "Vaue").Form("Key1", "Value").makeRequest()
	assert.NoError(t, err)

	assert.Equal(t, ContentTypeForm, req.Header.Get(headerContentType))
}

func TestPapagoSMT(t *testing.T) {
	type papagoSMTResp struct {
		Message struct {
			Type    string `json:"@type"`
			Service string `json:"@service"`
			Version string `json:"@version"`
			Result  struct {
				TranslatedText string `json:"translatedText"`
				SrcLangType    string `json:"srcLangType"`
			} `json:"result"`
		} `json:"message"`
	}

	resp, err := Post("https://openapi.naver.com/v1/language/translate").
		Headers(map[string]string{
			"X-Naver-Client-Id":     os.Getenv("NAVER_CLIENT_ID"),
			"X-Naver-Client-Secret": os.Getenv("NAVER_CLIENT_SECRET"),
		}).
		Forms(map[string]string{
			"source": "ko",
			"target": "en",
			"text":   "만나서 반갑습니다.",
		}).
		Do()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	r := &papagoSMTResp{}
	assert.NoError(t, resp.JSON(r))
	defer resp.Body.Close()
	assert.Equal(t, "Nice to meet you.", r.Message.Result.TranslatedText)
}

func TestGithubGet(t *testing.T) {
	resp, err := Get("https://api.github.com").Do()
	assert.NoError(t, err)

	r := map[string]string{}
	assert.NoError(t, resp.JSON(&r))
	defer resp.Body.Close()
	assert.Equal(t, "https://api.github.com/hub", r["hub_url"])
}

func TestGoogleCustomSearch(t *testing.T) {
	resp, err := Get("https://www.googleapis.com/customsearch/v1").
		Param("key", os.Getenv("GOOGLE_API_KEY")).
		Param("cx", os.Getenv("GOOGLE_cx")).
		Param("q", "request").
		Do()

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// https://developers.google.com/custom-search/json-api/v1/reference/cse/list#response
	type GoogleCustomSearchResponse struct {
		Kind string `json:"kind"`
		URL  struct {
			Type     string `json:"type"`
			Template string `json:"template"`
		} `json:"url"`
		Items []struct {
			Kind  string `json:"kind"`
			Title string `json:"title"`
			Link  string `json:"link"`
		} `json:"items"`
	}
	result := &GoogleCustomSearchResponse{}
	assert.NoError(t, resp.JSON(result))
	defer resp.Body.Close()
	assert.True(t, len(result.Items) > 0)
	log.Printf("link: %s", result.Items[0].Link)
}
