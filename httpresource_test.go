package webjacker_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/kitd/webjacker"
)

type testType struct {
	*webjacker.HttpResource
}

func TestPlainGet(t *testing.T) {

	test := testType{
		webjacker.NewHttpResource("test"),
	}

	method := http.MethodGet

	test.On(method, func(w http.ResponseWriter, r *http.Request, params url.Values) {
		checkMethod(t, r, method)
		w.Write([]byte("OK"))
	})

	mux := http.NewServeMux()
	webjacker.RegisterHttpResourceOnPath(test.HttpResource, mux, "")
	ts := httptest.NewServer(mux)
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL + test.Path())
	checkResponse(t, resp, err, method, "OK")
}

func TestGetWithParams(t *testing.T) {

	test := testType{
		webjacker.NewHttpResource("test"),
	}

	method := http.MethodGet

	param := "food"
	value := "cheese"

	test.On(method, func(w http.ResponseWriter, r *http.Request, params url.Values) {
		checkMethod(t, r, method)
		checkParam(t, params, param, value)
		w.Write([]byte("OK"))
	})

	mux := http.NewServeMux()
	webjacker.RegisterHttpResourceOnPath(test.HttpResource, mux, "")
	ts := httptest.NewServer(mux)
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL + fmt.Sprintf("%s?%s=%s", test.Path(), param, value))
	checkResponse(t, resp, err, method, "OK")
}

func TestPostWithParams(t *testing.T) {

	test := testType{
		webjacker.NewHttpResource("test"),
	}

	method := http.MethodPost

	param := "food"
	value := "cheese"

	test.On(method, func(w http.ResponseWriter, r *http.Request, params url.Values) {
		checkMethod(t, r, method)
		checkParam(t, params, param, value)
		w.Write([]byte("OK"))
	})

	mux := http.NewServeMux()
	webjacker.RegisterHttpResourceOnPath(test.HttpResource, mux, "")
	ts := httptest.NewServer(mux)
	defer ts.Close()

	resp, err := ts.Client().PostForm(ts.URL+test.Path(),
		url.Values{param: {value}})
	checkResponse(t, resp, err, method, "OK")
}

func TestCustomEvent(t *testing.T) {

	test := testType{
		webjacker.NewHttpResource("test"),
	}

	event := "my_event"

	test.On(event, func(w http.ResponseWriter, r *http.Request, params url.Values) {
		checkMethod(t, r, http.MethodGet)
		checkParam(t, params, webjacker.EventParam, event)
		w.Write([]byte("OK"))
	})

	mux := http.NewServeMux()
	webjacker.RegisterHttpResourceOnPath(test.HttpResource, mux, "")
	ts := httptest.NewServer(mux)
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL + test.EventPath(event))
	checkResponse(t, resp, err, http.MethodGet, "OK")
}

func TestNotImplemented(t *testing.T) {

	test := testType{
		webjacker.NewHttpResource("test"),
	}

	event := "my_event"

	test.On(event, func(w http.ResponseWriter, r *http.Request, params url.Values) {
		checkMethod(t, r, http.MethodGet)
		checkParam(t, params, webjacker.EventParam, event)
		w.Write([]byte("OK"))
	})

	mux := http.NewServeMux()
	webjacker.RegisterHttpResourceOnPath(test.HttpResource, mux, "")
	ts := httptest.NewServer(mux)
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL + test.EventPath("dummy_event"))
	if err != nil {
		t.Fatalf("Error running %s: %v", http.MethodGet, err)
	}
	if resp.StatusCode != http.StatusNotImplemented {
		t.Fatalf("Error running GET: expected status %d, received %d", http.StatusNotImplemented, resp.StatusCode)
	}
}

func checkParam(t *testing.T, params url.Values, name, expectedValue string) {
	if !params.Has(name) {
		t.Errorf("Expecting param %s. Was not present", name)
	}

	value := params.Get(name)
	if value != expectedValue {
		t.Errorf("Expecting param %s to be %s. Was %s", name, expectedValue, value)
	}
}

func checkMethod(t *testing.T, r *http.Request, expectedMethod string) {
	if r.Method != expectedMethod {
		t.Errorf("Expecting method %s. Was %s", expectedMethod, r.Method)
	}
}

func checkResponse(t *testing.T, resp *http.Response, err error, expectedMethod, expectedBody string) {
	if err != nil {
		t.Fatalf("Error running %s: %v", expectedMethod, err)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response: %v", err)
	}
	if string(bytes) != "OK" {
		t.Errorf("Expecting response body %s. Was %s", "OK", resp.Body)
	}
}
