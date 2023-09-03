package http

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	type input struct {
		mockHTTPServer mockHTTPServer
	}

	type want struct {
		err      bool
		contains string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles unexpected status code",
			input{
				mockHTTPServer: newMockHTTPServer(t, true),
			},
			want{
				err:      true,
				contains: "unexpected status code",
			},
		},
		{
			"is successful",
			input{
				mockHTTPServer: newMockHTTPServer(t, false),
			},
			want{
				err: false,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.input.mockHTTPServer.stop()
			tt.input.mockHTTPServer.start()
			time.Sleep(time.Second)

			u, err := url.Parse("http://127.0.0.1:33333/macd")
			if !assert.Nil(t, err) {
				t.FailNow()
			}

			_, err = Get(u, http.DefaultClient)
			if tt.want.err {
				if assert.NotNil(t, err) {
					assert.Contains(t, err.Error(), tt.want.contains)
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

type mockHTTPServer struct {
	t      *testing.T
	server *http.Server
}

func newMockHTTPServer(t *testing.T, wantErr bool) mockHTTPServer {
	return mockHTTPServer{
		server: &http.Server{Addr: "127.0.0.1:33333", Handler: &handler{
			wantErr: wantErr,
		}},
		t: t,
	}
}

func (m mockHTTPServer) start() {
	go func() {
		if err := m.server.ListenAndServe(); err != nil {
			m.t.Log("Failed to setup server: ", err.Error())
		}
	}()
}

func (m mockHTTPServer) stop() {
	if err := m.server.Shutdown(context.Background()); err != nil {
		m.t.Fatal("Failed to shutdown mock http server: ", err.Error())
	}
}

type handler struct {
	wantErr bool
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.wantErr {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Failed!")); err != nil {
			fmt.Println("Failed to respond in server test.")
		}
		return
	}

	if _, err := w.Write([]byte("Success")); err != nil {
		fmt.Println("Failed to respond in server test.")
	}
}
