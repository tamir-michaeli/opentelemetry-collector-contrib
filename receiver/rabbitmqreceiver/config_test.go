// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package rabbitmqreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/rabbitmqreceiver"

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.uber.org/multierr"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/rabbitmqreceiver/internal/metadata"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		desc        string
		cfg         *Config
		expectedErr error
	}{
		{
			desc: "missing username, password, and invalid endpoint",
			cfg: &Config{
				HTTPClientConfig: confighttp.HTTPClientConfig{
					Endpoint: "invalid://endpoint:  12efg",
				},
			},
			expectedErr: multierr.Combine(
				errMissingUsername,
				errMissingPassword,
				fmt.Errorf("%w: %s", errInvalidEndpoint, `parse "invalid://endpoint:  12efg": invalid port ":  12efg" after host`),
			),
		},
		{
			desc: "missing password and invalid endpoint",
			cfg: &Config{
				Username: "otelu",
				HTTPClientConfig: confighttp.HTTPClientConfig{
					Endpoint: "invalid://endpoint:  12efg",
				},
			},
			expectedErr: multierr.Combine(
				errMissingPassword,
				fmt.Errorf("%w: %s", errInvalidEndpoint, `parse "invalid://endpoint:  12efg": invalid port ":  12efg" after host`),
			),
		},
		{
			desc: "missing username and invalid endpoint",
			cfg: &Config{
				Password: "otelp",
				HTTPClientConfig: confighttp.HTTPClientConfig{
					Endpoint: "invalid://endpoint:  12efg",
				},
			},
			expectedErr: multierr.Combine(
				errMissingUsername,
				fmt.Errorf("%w: %s", errInvalidEndpoint, `parse "invalid://endpoint:  12efg": invalid port ":  12efg" after host`),
			),
		},
		{
			desc: "invalid endpoint",
			cfg: &Config{
				Username: "otelu",
				Password: "otelp",
				HTTPClientConfig: confighttp.HTTPClientConfig{
					Endpoint: "invalid://endpoint:  12efg",
				},
			},
			expectedErr: multierr.Combine(
				fmt.Errorf("%w: %s", errInvalidEndpoint, `parse "invalid://endpoint:  12efg": invalid port ":  12efg" after host`),
			),
		},
		{
			desc: "valid config",
			cfg: &Config{
				Username: "otelu",
				Password: "otelp",
				HTTPClientConfig: confighttp.HTTPClientConfig{
					Endpoint: defaultEndpoint,
				},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actualErr := tc.cfg.Validate()
			if tc.expectedErr != nil {
				require.EqualError(t, actualErr, tc.expectedErr.Error())
			} else {
				require.NoError(t, actualErr)
			}

		})
	}
}

func TestLoadConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	sub, err := cm.Sub(component.NewIDWithName(metadata.Type, "").String())
	require.NoError(t, err)
	require.NoError(t, component.UnmarshalConfig(sub, cfg))

	expected := factory.CreateDefaultConfig().(*Config)
	expected.Endpoint = "http://localhost:15672"
	expected.Username = "otelu"
	expected.Password = "${env:RABBITMQ_PASSWORD}"
	expected.CollectionInterval = 10 * time.Second

	require.Equal(t, expected, cfg)
}
