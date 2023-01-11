package request

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp/log"
)

func TestSimple(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := Get("https://www.google.co.kr").
		WithClient(http.DefaultClient).
		ContentType(MIMEApplicationJSON).
		Do(ctx)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestMakeRequest(t *testing.T) {
	{
		req, err := Get("http:....").
			Header("Key", "Value").
			Headers(map[string]string{
				"Key1": "Value1",
				"Key2": "Value2",
			}).
			Query("queryKey", "queryValue").
			Queries(map[string]string{
				"queryKey1": "queryValue1",
				"queryKey2": "queryValue2",
			}).
			makeRequest()
		require.NoError(t, err)

		query, _ := url.ParseQuery(req.URL.RawQuery)
		require.Equal(t, "queryValue", query.Get("queryKey"))
		require.Equal(t, "queryValue2", query.Get("queryKey2"))
		require.Equal(t, "Value", req.Header.Get("Key"))
		require.Equal(t, "Value2", req.Header.Get("Key2"))
	}

	{
		req, err := Post("http:....").
			Form("Key", "Vaue").
			Form("Key1", "Value").
			makeRequest()
		require.NoError(t, err)

		require.Equal(t, MIMEApplicationForm, req.Header.Get(HeaderContentType))
	}
}

func TestRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

	}))
	defer ts.Close()

	ctx := context.Background()
	{
		resp, err := Get(ts.URL).Do(ctx)
		require.NoError(t, err)
		require.Truef(t, resp.Success(), "status=%d", resp.StatusCode)
	}

	{
		var param = map[string]string{
			"key": "value",
		}
		resp, err := Post(ts.URL).JSON(&param).JSON(&param).Do(ctx)
		require.NoError(t, err)
		require.True(t, resp.Success())
	}
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
	defer resp.Body.Close()
	require.NoError(t, resp.JSON(&r))
	require.Equal(t, "Good to meet you.", r.Message.Result.TranslatedText)
}

func TestGithubGet(t *testing.T) {
	resp, err := Get("https://api.github.com").Do(context.Background())
	require.NoError(t, err)
	require.Truef(t, resp.Success(), "failed with status %d", resp.StatusCode)

	r := make(map[string]string)
	defer resp.Body.Close()
	require.NoError(t, resp.JSON(&r))
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
	defer resp.Body.Close()
	require.NoError(t, resp.JSON(&response))
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

func TestResponseBody(t *testing.T) {
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
