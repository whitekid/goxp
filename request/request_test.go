package request

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/goex/log"
)

func TestSimple(t *testing.T) {
	resp, err := Get("https://www.google.co.kr").
		WithClient(http.DefaultClient).
		Do(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFormContentType(t *testing.T) {
	req, err := Post("http:....").
		Form("Key", "Vaue").
		Form("Key1", "Value").
		makeRequest()
	require.NoError(t, err)

	require.Equal(t, ContentTypeForm, req.Header.Get(headerContentType))
}

func TestPapagoSMT(t *testing.T) {
	if _, ok := os.LookupEnv("NAVER_CLIENT_ID"); !ok {
		t.Skip()
	}

	resp, err := Post("https://openapi.naver.com/v1/papago/n2mt").
		Header("X-Naver-Client-Id", os.Getenv("NAVER_CLIENT_ID")).
		Header("X-Naver-Client-Secret", os.Getenv("NAVER_CLIENT_SECRET")).
		Forms(map[string]string{
			"source": "ko",
			"target": "en",
			"text":   "만나서 반갑습니다.",
		}).
		Do(context.Background())
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var r struct {
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
	require.NoError(t, resp.JSON(&r))
	defer resp.Body.Close()
	require.Equal(t, "Nice to meet you.", r.Message.Result.TranslatedText)
}

func TestGithubGet(t *testing.T) {
	resp, err := Get("https://api.github.com").Do(context.Background())
	require.NoError(t, err)

	r := make(map[string]string)
	require.NoError(t, resp.JSON(&r))
	defer resp.Body.Close()
	require.Equal(t, "https://api.github.com/hub", r["hub_url"])
}

func TestGoogleCustomSearch(t *testing.T) {
	key, ok := os.LookupEnv("GOOGLE_API_KEY")
	if !ok {
		t.Skip("GOOGLE_API_KEY required")
	}

	resp, err := Get("https://www.googleapis.com/customsearch/v1").
		Query("key", key).
		Query("cx", os.Getenv("GOOGLE_cx")).
		Query("q", "request").
		Do(context.Background())

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// https://developers.google.com/custom-search/json-api/v1/reference/cse/list#response
	var response struct {
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
	require.NoError(t, resp.JSON(&response))
	defer resp.Body.Close()
	require.True(t, len(response.Items) > 0)
	log.Infof("link: %s", response.Items[0].Link)
}

func TestRedirect(t *testing.T) {
	type args struct {
		url            string
		followRedirect bool
	}

	tests := [...]struct {
		name           string
		args           args
		wantErr        bool
		wantStatusCode int
	}{
		{"", args{"http://google.com", false}, false, http.StatusMovedPermanently},
		{"", args{"http://google.com", true}, false, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := Get(tt.args.url).
				FollowRedirect(tt.args.followRedirect).
				Do(context.Background())
			require.NoError(t, err)
			require.Equal(t, tt.wantStatusCode, resp.StatusCode)
		})
	}
}

func TestBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		io.Copy(w, req.Body)
	}))
	defer ts.Close()

	message := "hello world"
	want := message
	resp, err := Post(ts.URL).
		Body(strings.NewReader(message)).
		Do(context.Background())
	require.NoError(t, err)
	require.True(t, resp.Success())
	require.Equal(t, want, resp.String())
}
