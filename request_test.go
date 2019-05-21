package httpbox

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

type getHandler struct{}

type handlerResp struct {
	Msg string `json:"msg"`
}

func (gh *getHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := &handlerResp{}
	if rawQuery := r.URL.RawQuery; rawQuery != "" {
		resp.Msg = rawQuery
	} else {
		resp.Msg = "no params"
	}

	b, _ := json.Marshal(resp)
	io.Copy(w, bytes.NewReader(b))
}

type postHandler struct{}

func (gh *postHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	resp := &handlerResp{Msg: string(reqBody)}
	b, _ := json.Marshal(resp)

	io.Copy(w, bytes.NewReader(b))
}

func TestGet(t *testing.T) {
	type args struct {
		testReq    Requester
		newHandler http.Handler
		resp       *handlerResp
	}
	tests := []struct {
		name     string
		args     args
		want     int
		wantResp *handlerResp
		wantErr  bool
	}{
		{
			args: args{
				testReq:    DefaultRequester("/test/get", map[string]string{"k": "v"}),
				newHandler: &getHandler{},
				resp:       &handlerResp{},
			},
			want:     http.StatusOK,
			wantResp: &handlerResp{Msg: "k=v"},
			wantErr:  false,
		},

		{
			args: args{
				testReq:    DefaultRequester("/test/get"),
				newHandler: &getHandler{},
				resp:       &handlerResp{},
			},
			want:     http.StatusOK,
			wantResp: &handlerResp{Msg: "no params"},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.testReq, tt.args.newHandler, tt.args.resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
			if tt.args.resp.Msg != tt.wantResp.Msg {
				t.Errorf("Get() = %v, want %v", tt.args.resp.Msg, tt.wantResp.Msg)
			}
		})
	}
}

func TestPost(t *testing.T) {
	type args struct {
		testReq    Requester
		newHandler http.Handler
		resp       *handlerResp
	}
	tests := []struct {
		name     string
		args     args
		want     int
		wantResp handlerResp
		wantErr  bool
	}{
		{
			args: args{
				testReq:    DefaultRequester("/test/get", map[string]string{"k": "v"}),
				newHandler: &postHandler{},
				resp:       &handlerResp{},
			},
			want:     http.StatusOK,
			wantResp: handlerResp{Msg: `{"k":"v"}`},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Post(tt.args.testReq, tt.args.newHandler, tt.args.resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Post() = %v, want %v", got, tt.want)
			}
			if tt.args.resp.Msg != tt.wantResp.Msg {
				t.Errorf("Post() = %v, want %v", tt.args.resp.Msg, tt.wantResp.Msg)
			}
		})
	}
}
