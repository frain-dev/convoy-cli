package net

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"testing"

	"github.com/frain-dev/convoy/config"

	"github.com/frain-dev/convoy/pkg/httpheader"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

var successBody = []byte("received webhook successfully")

func TestDispatcher_SendCliRequest(t *testing.T) {
	client := http.DefaultClient

	buf := make([]byte, config.MaxResponseSize*2)
	_, _ = rand.Read(buf)
	type args struct {
		headers  httpheader.HTTPHeader
		method   string
		url      string
		apiKey   string
		jsonData json.RawMessage
	}
	tests := []struct {
		name    string
		args    args
		want    *Response
		nFn     func() func()
		wantErr bool
	}{
		{
			name: "should_send_message",
			args: args{
				method:   http.MethodPost,
				jsonData: bytes.NewBufferString("testing").Bytes(),
			},
			want: &Response{
				Status:     "200",
				StatusCode: http.StatusOK,
				Method:     http.MethodPost,
				URL:        nil,
				RequestHeader: http.Header{
					"Content-Type":                         []string{"application/json"},
					"User-Agent":                           []string{defaultUserAgent()},
					config.DefaultSignatureHeader.String(): []string{"12345"}, // should equal hmac field above
				},
				ResponseHeader: nil,
				Body:           successBody,
				IP:             "",
				Error:          "",
			},
			nFn: func() func() {
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "https://google.com",
					httpmock.NewStringResponder(http.StatusOK, string(successBody)))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErr: false,
		},
		{
			name: "should_send_message_with_forwarded_headers",
			args: args{
				method:   http.MethodPost,
				jsonData: bytes.NewBufferString("testing").Bytes(),
				headers: map[string][]string{
					"X-Test-Sig": {"abcdef"},
				},
			},
			want: &Response{
				Status:     "200",
				StatusCode: http.StatusOK,
				Method:     http.MethodPost,
				URL:        nil,
				RequestHeader: http.Header{
					"Content-Type":                         []string{"application/json"},
					"User-Agent":                           []string{defaultUserAgent()},
					"X-Test-Sig":                           []string{"abcdef"},
					config.DefaultSignatureHeader.String(): []string{"12345"}, // should equal hmac field above
				},
				ResponseHeader: nil,
				Body:           successBody,
				IP:             "",
				Error:          "",
			},
			nFn: func() func() {
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "https://google.com",
					httpmock.NewStringResponder(http.StatusOK, string(successBody)))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErr: false,
		},
		{
			name: "should_cut_down_oversized_response_body",
			args: args{
				method:   http.MethodPost,
				jsonData: bytes.NewBufferString("testing").Bytes(),
			},
			want: &Response{
				Status:     "200",
				StatusCode: http.StatusOK,
				Method:     http.MethodPost,
				URL:        nil,
				RequestHeader: http.Header{
					"Content-Type":                         []string{"application/json"},
					"User-Agent":                           []string{defaultUserAgent()},
					config.DefaultSignatureHeader.String(): []string{"12345"}, // should equal hmac field above
				},
				ResponseHeader: nil,
				Body:           buf[:config.MaxResponseSize],
				IP:             "",
				Error:          "",
			},
			nFn: func() func() {
				httpmock.Activate()

				httpmock.RegisterResponder(http.MethodPost, "https://google.com",
					httpmock.NewBytesResponder(http.StatusOK, buf))

				return func() {
					httpmock.DeactivateAndReset()
				}
			},
			wantErr: false,
		},
		{
			name: "should_refuse_connection",
			args: args{
				method:   http.MethodPost,
				jsonData: bytes.NewBufferString("bossman").Bytes(),
			},
			want: &Response{
				Status:     "",
				StatusCode: 0,
				Method:     http.MethodPost,
				RequestHeader: http.Header{
					"Content-Type":                         []string{"application/json"},
					"User-Agent":                           []string{defaultUserAgent()},
					config.DefaultSignatureHeader.String(): []string{"12345"}, // should equal hmac field above
				},
				ResponseHeader: nil,
				Body:           nil,
				IP:             "",
				Error:          "connect: connection refused",
			},
			wantErr: true,
		},
		{
			name: "should_error_for_empty_signature_hmac",
			args: args{
				method:   http.MethodPost,
				jsonData: bytes.NewBufferString("bossman").Bytes(),
			},
			want: &Response{
				Status:         "",
				StatusCode:     0,
				Method:         "",
				RequestHeader:  nil,
				ResponseHeader: nil,
				Body:           nil,
				IP:             "",
				Error:          "signature header and hmac are required",
			},
			wantErr: true,
		},
		{
			name: "should_error_for_empty_signature_header",
			args: args{
				method:   http.MethodPost,
				jsonData: bytes.NewBufferString("bossman").Bytes(),
			},
			want: &Response{
				Status:         "",
				StatusCode:     0,
				Method:         "",
				RequestHeader:  nil,
				ResponseHeader: nil,
				Body:           nil,
				IP:             "",
				Error:          "signature header and hmac are required",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dispatcher{client: client}

			if tt.nFn != nil {
				deferFn := tt.nFn()
				defer deferFn()
			}

			got, err := d.SendCliRequest(tt.args.url, tt.args.method, tt.args.apiKey, tt.args.jsonData)
			if tt.wantErr {
				require.NotNil(t, err)
				require.Contains(t, err.Error(), tt.want.Error)
				require.Contains(t, got.Error, tt.want.Error)
			}

			require.Equal(t, tt.want.Status, got.Status)
			require.Equal(t, tt.want.StatusCode, got.StatusCode)
			require.Equal(t, tt.want.Method, got.Method)
			require.Equal(t, tt.want.IP, got.IP)
			require.Equal(t, tt.want.Body, got.Body)
			require.Equal(t, tt.want.RequestHeader, got.RequestHeader)
		})
	}
}