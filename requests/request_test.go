package requests

import (
	"compress/flate"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/require"

	"github.com/whitekid/goxp"
	"github.com/whitekid/goxp/log"
)

const uaFirefox = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:132.0) Gecko/20100101 Firefox/132.0"

func TestSimple(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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
		require.NoError(t, resp.Success(), "status=%d", resp.StatusCode)
	}

	{
		var param = map[string]string{
			"key": "value",
		}
		resp, err := Post(ts.URL).JSON(&param).JSON(&param).Do(ctx)
		require.NoError(t, err)
		require.NoError(t, resp.Success())
	}
}

func TestAuth(t *testing.T) {
	type args struct {
		auth func(r *Request) *Request
	}
	tests := [...]struct {
		name       string
		args       args
		wantErr    bool
		wantHeader map[string]string
	}{
		{"bearer", args{func(r *Request) *Request { return r.AuthBearer("token") }},
			false, map[string]string{"Authorization": "Bearer token"}},
		{"token", args{func(r *Request) *Request { return r.AuthToken("token") }},
			false, map[string]string{"Authorization": "Token token"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := tt.args.auth(Get("https://google.com")).makeRequest()
			require.NoError(t, err)
			for k, v := range tt.wantHeader {
				require.Equal(t, v, req.Header.Get(k))
			}
		})
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

	type Response struct {
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
	r, err := goxp.ReadJSON[Response](resp.Body)
	require.NoError(t, err)
	require.Equal(t, "Good to meet you.", r.Message.Result.TranslatedText)
}

func TestGithubGet(t *testing.T) {
	resp, err := Get("https://api.github.com").Do(context.Background())
	require.NoError(t, err)
	require.NoError(t, resp.Success(), "failed with status %d", resp.StatusCode)

	defer resp.Body.Close()
	r, err := goxp.ReadJSON[map[string]string](resp.Body)
	require.NoError(t, err)
	require.Equal(t, "https://api.github.com/hub", (*r)["hub_url"])
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
	type Response struct {
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
	r, err := goxp.ReadJSON[Response](resp.Body)
	require.NoError(t, err)
	require.True(t, len(r.Items) > 0)
	log.Infof("link: %s", r.Items[0].Link)
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
		resolvedURL    string
	}{
		{"", args{"http://google.com", false}, false, http.StatusMovedPermanently, ""},
		{"", args{"http://google.com", true}, false, http.StatusOK, "http://www.google.com/"},
		{"", args{"https://code.facebook.com/posts/rss", true}, false, http.StatusOK, "https://engineering.fb.com/feed/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := Get(tt.args.url).
				FollowRedirect(tt.args.followRedirect).
				Do(context.Background())
			require.NoError(t, err)
			require.Equal(t, tt.wantStatusCode, resp.StatusCode)

			if !tt.args.followRedirect {
				return
			}
			require.Equal(t, tt.resolvedURL, resp.Request.URL.String())
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
	require.NoError(t, resp.Success())
	require.Equal(t, want, resp.String())
}

func TestRequestEncodingPure(t *testing.T) {
	type args struct {
		url     string
		headers map[string]string
	}
	tests := [...]struct {
		name       string
		args       args
		wantHeader map[string]string
		wantText   string
	}{
		{`gzip`, args{
			url: "https://www.google.com",
			headers: map[string]string{
				HeaderUserAgent:      uaFirefox,
				headerAcceptEncoding: "gzip",
			},
		}, map[string]string{
			HeaderContentEncoding: "gzip",
		}, "googleg"},
		{`br`, args{
			url: "https://www.google.com",
			headers: map[string]string{
				HeaderUserAgent:      uaFirefox,
				headerAcceptEncoding: "br",
			},
		}, map[string]string{
			HeaderContentEncoding: "br",
		}, "googleg"},
		{`zstd`, args{
			url: "https://www.facebook.com",
			headers: map[string]string{
				HeaderUserAgent:      uaFirefox,
				headerAcceptEncoding: "zstd",
			},
		}, map[string]string{
			HeaderContentEncoding: "zstd",
		}, "facebook"},
		{`deflate`, args{
			url: "https://www.facebook.com",
			headers: map[string]string{
				HeaderUserAgent:      uaFirefox,
				headerAcceptEncoding: "deflate",
			},
		}, map[string]string{
			HeaderContentEncoding: "", // deflate를 보내면 인코딩을 하지 않고 응답한다.
		}, "facebook"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			req, err := http.NewRequest(http.MethodGet, tt.args.url, nil)
			require.NoError(t, err)

			for k, v := range tt.args.headers {
				req.Header.Set(k, v)
			}

			client := &http.Client{}

			req = req.WithContext(ctx)
			resp, err := client.Do(req)
			require.NoError(t, err)

			for k, v := range tt.wantHeader {
				require.Equal(t, v, resp.Header.Get(k))
			}

			var reader io.ReadCloser
			switch enc := resp.Header.Get(HeaderContentEncoding); enc {
			case "gzip":
				reader, err = gzip.NewReader(resp.Body)
				if err != nil {
					panic(err)
				}
			case "br":
				reader = io.NopCloser(brotli.NewReader(resp.Body))
			case "zstd":
				decoder, err := zstd.NewReader(resp.Body)
				if err != nil {
					panic(err)
				}
				defer decoder.Close()
				reader = io.NopCloser(decoder)
			case "deflate":
				reader = io.NopCloser(flate.NewReader(resp.Body))
			case "":
				reader = resp.Body
			default:
				require.Failf(t, "unsupported encoding: %s", enc)
			}

			body, err := io.ReadAll(reader)
			require.NoError(t, err)
			require.NotEmpty(t, body)

			if tt.wantText != "" {
				require.Contains(t, string(body), tt.wantText)
			}
		})
	}
}

func TestRequestEncoding(t *testing.T) {
	type args struct {
		url            string
		acceptEncoding string
	}
	tests := [...]struct {
		name                string
		args                args
		wantContentEncoding string
		wantText            string
	}{
		{`default`, args{
			url:            "https://www.google.com",
			acceptEncoding: "",
		}, "br", "googleg"},
		{`gzip`, args{
			url:            "https://www.google.com",
			acceptEncoding: "gzip",
		}, "gzip", "googleg"},
		{`br`, args{
			url:            "https://www.google.com",
			acceptEncoding: "br",
		}, "br", "googleg"},
		{`zstd`, args{
			url:            "https://www.facebook.com",
			acceptEncoding: "zstd",
		}, "zstd", "facebook"},
		{`deflate`, args{
			url:            "https://www.facebook.com",
			acceptEncoding: "deflate",
		}, "", "facebook"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			req := Get(tt.args.url).Header(HeaderUserAgent, uaFirefox).Header(HeaderContentType, "text/html")
			if tt.args.acceptEncoding != "" {
				req.Header(headerAcceptEncoding, tt.args.acceptEncoding)
			}

			resp, err := req.Do(ctx)
			require.NoError(t, err)

			require.Equal(t, tt.wantContentEncoding, resp.Header.Get(HeaderContentEncoding))

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.NotEmpty(t, body)

			if tt.wantText != "" {
				require.Contains(t, string(body), tt.wantText)
			}
		})
	}
}
