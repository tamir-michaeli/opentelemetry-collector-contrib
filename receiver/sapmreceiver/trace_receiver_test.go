// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package sapmreceiver

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/binary"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jaegertracing/jaeger/model"
	"github.com/klauspost/compress/zstd"
	splunksapm "github.com/signalfx/sapm-proto/gen"
	"github.com/signalfx/sapm-proto/sapmprotocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/receivertest"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/common/testutil"
	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/splunk"
)

func expectedTraceData(t1, t2, t3 time.Time) ptrace.Traces {
	traceID := pcommon.TraceID(
		[16]byte{0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF, 0x80})
	parentSpanID := pcommon.SpanID([8]byte{0x1F, 0x1E, 0x1D, 0x1C, 0x1B, 0x1A, 0x19, 0x18})
	childSpanID := pcommon.SpanID([8]byte{0xAF, 0xAE, 0xAD, 0xAC, 0xAB, 0xAA, 0xA9, 0xA8})

	traces := ptrace.NewTraces()
	rs := traces.ResourceSpans().AppendEmpty()
	rs.Resource().Attributes().PutStr(conventions.AttributeServiceName, "issaTest")
	rs.Resource().Attributes().PutBool("bool", true)
	rs.Resource().Attributes().PutStr("string", "yes")
	rs.Resource().Attributes().PutInt("int64", 10000000)
	spans := rs.ScopeSpans().AppendEmpty().Spans()

	span0 := spans.AppendEmpty()
	span0.SetSpanID(childSpanID)
	span0.SetParentSpanID(parentSpanID)
	span0.SetTraceID(traceID)
	span0.SetName("DBSearch")
	span0.SetStartTimestamp(pcommon.NewTimestampFromTime(t1))
	span0.SetEndTimestamp(pcommon.NewTimestampFromTime(t2))
	span0.Status().SetCode(ptrace.StatusCodeError)
	span0.Status().SetMessage("Stale indices")

	span1 := spans.AppendEmpty()
	span1.SetSpanID(parentSpanID)
	span1.SetTraceID(traceID)
	span1.SetName("ProxyFetch")
	span1.SetStartTimestamp(pcommon.NewTimestampFromTime(t2))
	span1.SetEndTimestamp(pcommon.NewTimestampFromTime(t3))
	span1.Status().SetCode(ptrace.StatusCodeError)
	span1.Status().SetMessage("Frontend crash")

	return traces
}

func grpcFixture(t1 time.Time) *model.Batch {
	traceID := model.TraceID{}
	_ = traceID.Unmarshal([]byte{0xF1, 0xF2, 0xF3, 0xF4, 0xF5, 0xF6, 0xF7, 0xF8, 0xF9, 0xFA, 0xFB, 0xFC, 0xFD, 0xFE, 0xFF, 0x80})
	parentSpanID := model.NewSpanID(binary.BigEndian.Uint64([]byte{0x1F, 0x1E, 0x1D, 0x1C, 0x1B, 0x1A, 0x19, 0x18}))
	childSpanID := model.NewSpanID(binary.BigEndian.Uint64([]byte{0xAF, 0xAE, 0xAD, 0xAC, 0xAB, 0xAA, 0xA9, 0xA8}))

	return &model.Batch{
		Process: &model.Process{
			ServiceName: "issaTest",
			Tags: []model.KeyValue{
				model.Bool("bool", true),
				model.String("string", "yes"),
				model.Int64("int64", 1e7),
			},
		},
		Spans: []*model.Span{
			{
				TraceID:       traceID,
				SpanID:        childSpanID,
				OperationName: "DBSearch",
				StartTime:     t1,
				Duration:      10 * time.Minute,
				Tags: []model.KeyValue{
					model.String(conventions.OtelStatusDescription, "Stale indices"),
					model.String(conventions.OtelStatusCode, "ERROR"),
					model.Bool("error", true),
				},
				References: []model.SpanRef{
					{
						TraceID: traceID,
						SpanID:  parentSpanID,
						RefType: model.SpanRefType_CHILD_OF,
					},
				},
			},
			{
				TraceID:       traceID,
				SpanID:        parentSpanID,
				OperationName: "ProxyFetch",
				StartTime:     t1.Add(10 * time.Minute),
				Duration:      2 * time.Second,
				Tags: []model.KeyValue{
					model.String(conventions.OtelStatusDescription, "Frontend crash"),
					model.String(conventions.OtelStatusCode, "ERROR"),
					model.Bool("error", true),
				},
			},
		},
	}
}

// sendSapm acts as a client for sending sapm to the receiver.  This could be replaced with a sapm exporter in the future.
func sendSapm(
	endpoint string,
	sapm *splunksapm.PostSpansRequest,
	compression string,
	tlsEnabled bool,
	token string,
) (*http.Response, error) {
	// marshal the sapm
	reqBytes, err := sapm.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sapm %w", err)
	}

	switch compression {
	case "gzip":
		reqBytes, err = compressGzip(reqBytes)
		if err != nil {
			return nil, err
		}
	case "zstd":
		reqBytes, err = compressZstd(reqBytes)
		if err != nil {
			return nil, err
		}
	case "":
		// no compression
	default:
		return nil, fmt.Errorf("unknown compression %q", compression)
	}

	// build the request
	url := fmt.Sprintf("http://%s%s", endpoint, sapmprotocol.TraceEndpointV2)
	if tlsEnabled {
		url = fmt.Sprintf("https://%s%s", endpoint, sapmprotocol.TraceEndpointV2)
	}
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBytes))
	req.Header.Set(sapmprotocol.ContentTypeHeaderName, sapmprotocol.ContentTypeHeaderValue)

	// set headers for gzip
	if compression != "" {
		req.Header.Set(sapmprotocol.ContentEncodingHeaderName, compression)
		req.Header.Set(sapmprotocol.AcceptEncodingHeaderName, sapmprotocol.GZipEncodingHeaderValue)
	}

	if token != "" {
		req.Header.Set("x-sf-token", token)
	}

	// send the request
	client := &http.Client{}

	if tlsEnabled {
		tlscs := configtls.TLSClientSetting{
			TLSSetting: configtls.TLSSetting{
				CAFile:   "./testdata/ca.crt",
				CertFile: "./testdata/client.crt",
				KeyFile:  "./testdata/client.key",
			},
			ServerName: "localhost",
		}
		tls, errTLS := tlscs.LoadTLSConfig()
		if errTLS != nil {
			return nil, fmt.Errorf("failed to send request to receiver %w", err)
		}
		client.Transport = &http.Transport{
			TLSClientConfig: tls,
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return resp, fmt.Errorf("failed to send request to receiver %w", err)
	}

	return resp, nil
}

func compressGzip(reqBytes []byte) ([]byte, error) {
	// create a gzip writer
	var buff bytes.Buffer
	writer := gzip.NewWriter(&buff)

	// run the request bytes through the gzip writer
	_, err := writer.Write(reqBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to write gzip sapm %w", err)
	}

	// close the writer
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close the gzip writer %w", err)
	}

	return buff.Bytes(), nil
}

func compressZstd(reqBytes []byte) ([]byte, error) {
	// create a gzip writer
	var buff bytes.Buffer
	writer, err := zstd.NewWriter(&buff)
	if err != nil {
		return nil, fmt.Errorf("failed to write zstd sapm %w", err)
	}

	// run the request bytes through the gzip writer
	_, err = writer.Write(reqBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to write zstd sapm %w", err)
	}

	// close the writer
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close the zstd writer %w", err)
	}

	return buff.Bytes(), nil
}

func setupReceiver(t *testing.T, config *Config, sink *consumertest.TracesSink) receiver.Traces {
	params := receivertest.NewNopCreateSettings()
	sr, err := newReceiver(params, config, sink)
	assert.NoError(t, err, "should not have failed to create the SAPM receiver")
	t.Log("Starting")

	require.NoError(t, sr.Start(context.Background(), componenttest.NewNopHost()), "should not have failed to start trace reception")
	require.NoError(t, sr.Start(context.Background(), componenttest.NewNopHost()), "should not fail to start log on second Start call")

	// If there are errors reported through host.ReportFatalError() this will retrieve it.
	<-time.After(500 * time.Millisecond)
	t.Log("Trace Reception Started")
	return sr
}

func TestReception(t *testing.T) {
	now := time.Unix(1542158650, 536343000).UTC()
	nowPlus10min := now.Add(10 * time.Minute)
	nowPlus10min2sec := now.Add(10 * time.Minute).Add(2 * time.Second)
	tlsAddress := testutil.GetAvailableLocalAddress(t)

	type args struct {
		config      *Config
		sapm        *splunksapm.PostSpansRequest
		compression string
		useTLS      bool
	}
	tests := []struct {
		name string
		args args
		want ptrace.Traces
	}{
		{
			name: "receive uncompressed sapm",
			args: args{
				// 1. Create the SAPM receiver aka "server"
				config: &Config{
					HTTPServerConfig: confighttp.HTTPServerConfig{
						Endpoint: "0.0.0.0:7226",
					},
				},
				sapm:        &splunksapm.PostSpansRequest{Batches: []*model.Batch{grpcFixture(now)}},
				compression: "",
				useTLS:      false,
			},
			want: expectedTraceData(now, nowPlus10min, nowPlus10min2sec),
		},
		{
			name: "receive compressed sapm",
			args: args{
				config: &Config{
					HTTPServerConfig: confighttp.HTTPServerConfig{
						Endpoint: "0.0.0.0:7226",
					},
				},
				sapm:        &splunksapm.PostSpansRequest{Batches: []*model.Batch{grpcFixture(now)}},
				compression: "gzip",
				useTLS:      false,
			},
			want: expectedTraceData(now, nowPlus10min, nowPlus10min2sec),
		},
		{
			name: "connect via TLS zstd compressed sapm",
			args: args{
				config: &Config{
					HTTPServerConfig: confighttp.HTTPServerConfig{
						Endpoint: tlsAddress,
						TLSSetting: &configtls.TLSServerSetting{
							TLSSetting: configtls.TLSSetting{
								CAFile:   "./testdata/ca.crt",
								CertFile: "./testdata/server.crt",
								KeyFile:  "./testdata/server.key",
							},
						},
					},
				},
				sapm:        &splunksapm.PostSpansRequest{Batches: []*model.Batch{grpcFixture(now)}},
				compression: "zstd",
				useTLS:      true,
			},
			want: expectedTraceData(now, nowPlus10min, nowPlus10min2sec),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			sink := new(consumertest.TracesSink)
			sr := setupReceiver(t, tt.args.config, sink)
			defer func() {
				require.NoError(t, sr.Shutdown(context.Background()))
			}()

			t.Log("Sending Sapm Request")
			var resp *http.Response
			resp, err := sendSapm(tt.args.config.Endpoint, tt.args.sapm, tt.args.compression, tt.args.useTLS, "")
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			t.Log("SAPM Request Received")

			// retrieve received traces
			got := sink.AllTraces()
			assert.Equal(t, 1, len(got))

			// compare what we got to what we wanted
			t.Log("Comparing expected data to trace data")
			assert.EqualValues(t, tt.want, got[0])
		})
	}
}

func TestAccessTokenPassthrough(t *testing.T) {
	tests := []struct {
		name                   string
		accessTokenPassthrough bool
		token                  string
	}{
		{
			name:                   "no passthrough and no token",
			accessTokenPassthrough: false,
			token:                  "",
		},
		{
			name:                   "no passthrough and token",
			accessTokenPassthrough: false,
			token:                  "MyAccessToken",
		},
		{
			name:                   "passthrough and no token",
			accessTokenPassthrough: true,
			token:                  "",
		},
		{
			name:                   "passthrough and token",
			accessTokenPassthrough: true,
			token:                  "MyAccessToken",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				HTTPServerConfig: confighttp.HTTPServerConfig{
					Endpoint: "0.0.0.0:7226",
				},
				AccessTokenPassthroughConfig: splunk.AccessTokenPassthroughConfig{
					AccessTokenPassthrough: tt.accessTokenPassthrough,
				},
			}

			sapm := &splunksapm.PostSpansRequest{
				Batches: []*model.Batch{grpcFixture(time.Now().UTC())},
			}

			sink := new(consumertest.TracesSink)
			sr := setupReceiver(t, config, sink)
			defer func() {
				require.NoError(t, sr.Shutdown(context.Background()))
			}()

			var resp *http.Response
			resp, err := sendSapm(config.Endpoint, sapm, "gzip", false, tt.token)
			require.NoErrorf(t, err, "should not have failed when sending sapm %v", err)
			assert.Equal(t, 200, resp.StatusCode)

			got := sink.AllTraces()
			assert.Equal(t, 1, len(got))

			received := got[0].ResourceSpans()
			for i := 0; i < received.Len(); i++ {
				rspan := received.At(i)
				attrs := rspan.Resource().Attributes()
				amap, contains := attrs.Get("com.splunk.signalfx.access_token")
				if tt.accessTokenPassthrough && tt.token != "" {
					assert.Equal(t, tt.token, amap.Str())
				} else {
					assert.False(t, contains)
				}
			}
		})
	}
}
