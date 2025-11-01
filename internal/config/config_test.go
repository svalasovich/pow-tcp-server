package config

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	type testFixtures struct {
		action  func(*cobra.Command, *testing.T)
		asserts func(t *testing.T, cfg Config, err error)
	}

	testCases := []struct {
		name     string
		fixtures testFixtures
	}{
		{
			name: "default values",
			fixtures: testFixtures{
				action: func(_ *cobra.Command, t *testing.T) {
					t.Setenv("PRIVATE_KEY", "44a533d657e872b9dd8e435c8cc75c6f6da58518258c819819aba92a2aa4243f")
				},
				asserts: func(t *testing.T, cfg Config, err error) {
					expected := Config{
						Log: Log{
							Level:      "info",
							ShowSource: false,
							JSONFormat: false,
						},
						Monitoring: Monitoring{
							Server: HTTPServer{
								Port: 9000,
							},
						},
					}

					require.NoError(t, err)
					assert.Equal(t, expected, cfg)
				},
			},
		},
		{
			name: "environment values",
			fixtures: testFixtures{
				action: func(_ *cobra.Command, t *testing.T) {
					t.Setenv("LOG__LEVEL", "error")
				},
				asserts: func(t *testing.T, cfg Config, err error) {
					expected := Config{
						Log: Log{
							Level:      "error",
							ShowSource: false,
							JSONFormat: false,
						},
						Monitoring: Monitoring{
							Server: HTTPServer{
								Port: 9000,
							},
						},
					}

					require.NoError(t, err)
					assert.Equal(t, expected, cfg)
				},
			},
		},
		{
			name: "cli values",
			fixtures: testFixtures{
				action: func(cmd *cobra.Command, t *testing.T) {
					require.NoError(t, cmd.Flags().Set("log.level", "error"))
					t.Setenv("LOG__LEVEL", "envLevel")
				},
				asserts: func(t *testing.T, cfg Config, err error) {
					expected := Config{
						Log: Log{
							Level:      "error",
							ShowSource: false,
							JSONFormat: false,
						},
						Monitoring: Monitoring{
							Server: HTTPServer{
								Port: 9000,
							},
						},
					}

					require.NoError(t, err)
					assert.Equal(t, expected, cfg)
				},
			},
		},
		{
			name: "invalid values",
			fixtures: testFixtures{
				action: func(cmd *cobra.Command, t *testing.T) {
					require.NoError(t, cmd.Flags().Set("log.level", ""))
					t.Setenv("MONITORING__SERVER__PORT", "99999")
				},
				asserts: func(t *testing.T, _ Config, err error) {
					assert.ErrorContains(t, err, "log.level")
					assert.ErrorContains(t, err, "monitoring.server.port")
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// given
			cmd := &cobra.Command{}
			InitConfigFileFlags(cmd)
			InitLogFlags(cmd)
			InitMonitoringFlags(cmd)

			testCase.fixtures.action(cmd, t)

			// when
			cfg := &Config{}
			err := Init(cfg, cmd)

			// then
			testCase.fixtures.asserts(t, *cfg, err)
		})
	}
}
