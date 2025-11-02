package config

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/svalasovich/pow-tcp-server/internal/log"
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
				action: func(_ *cobra.Command, _ *testing.T) {
				},
				asserts: func(t *testing.T, cfg Config, err error) {
					expected := Config{
						Log: log.Config{
							Level:      "info",
							ShowSource: false,
							JSONFormat: false,
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
						Log: log.Config{
							Level:      "error",
							ShowSource: false,
							JSONFormat: false,
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
						Log: log.Config{
							Level:      "error",
							ShowSource: false,
							JSONFormat: false,
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
					require.NoError(t, cmd.Flags().Set("log.level", "tt"))
				},
				asserts: func(t *testing.T, _ Config, err error) {
					assert.ErrorContains(t, err, "log.level")
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

			testCase.fixtures.action(cmd, t)

			// when
			cfg := &Config{}
			err := Init(cfg, cmd)

			// then
			testCase.fixtures.asserts(t, *cfg, err)
		})
	}
}
