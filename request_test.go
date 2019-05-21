package httpbox

import (
	"net/http"
	"reflect"
	"testing"
	"io"
	"strings"
	"io/ioutil"
)


type getHandler struct {}

func (gh *getHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.Copy(w,strings.NewReader(r.URL.RawQuery))
}

func testGetHandler() http.Handler {
	return &getHandler{}
}

type postHandler struct {}

func (gh *postHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqBody,_ := ioutil.ReadAll(r.Body)
	io.Copy(w,strings.NewReader(string(reqBody)))
}

func testPostHandler() http.Handler {
	return &postHandler{}
}


func TestGet(t *testing.T) {
	type args struct {
		url        string
		params     map[string]string
		newHandler func() http.Handler
	}
	tests := []struct {
		name string
		args  args
		want  []byte
		want1 int
	}{
		{
			args:  args{
				url:        "/test/get",
				params:     map[string]string{"k":"v"},
				newHandler: testGetHandler,
			},
			want:  []byte("k=v"),
			want1: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Get(tt.args.url, tt.args.params, tt.args.newHandler)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPost(t *testing.T) {
	type args struct {
		url        string
		params     map[string]string
		newHandler func() http.Handler
	}
	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 int
	}{
		{
			args:  args{
				url:        "/test/post",
				params:     map[string]string{"k":"v"},
				newHandler: testPostHandler,
			},
			want:  []byte(`{"k":"v"}`),
			want1: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Post(tt.args.url, tt.args.params, tt.args.newHandler)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Post() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Post() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
