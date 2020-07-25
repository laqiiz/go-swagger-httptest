package go_swagger_httptest

import (
	"github.com/go-openapi/loads"
	"github.com/laqiiz/go-swagger-httptest/gen/restapi"
	"github.com/laqiiz/go-swagger-httptest/gen/restapi/hello"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func HelloHandler() (http.Handler, error) {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return nil, err
	}

	api := hello.NewHelloAPI(swaggerSpec)
	h := restapi.ConfigureAPI(api)

	return h, nil
}

func TestHello(t *testing.T) {

	handler, err := HelloHandler()
	if err != nil {
		t.Fatal("api handler", err)
	}

	tests := []struct {
		name string
		handler   http.Handler
		path string
		wantCode int
		wantBody string
	}{
		{
			name: "get hello",
			handler:   handler,
			path: "hello",
			wantCode: http.StatusOK,
			wantBody: `{"message":"hello"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			// httptest request
			resp, err := http.Get(ts.URL +  "/v1/hello")
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			gotBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if tt.wantBody != strings.TrimSpace(string(gotBody)) {
				t.Errorf("want:%s but got:%v", tt.wantBody, string(gotBody))
			}
			if tt.wantCode != resp.StatusCode {
				t.Errorf("want:%d but got:%d", tt.wantCode, resp.StatusCode)
			}
		})
	}

}
