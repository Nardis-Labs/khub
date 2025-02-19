package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
}

func (suite *ConfigSuite) SetupTest() {
	viper.Reset()
}

func createConfigFile(t *testing.T, cfgFile string) {
	content := []byte(`listen_port: "8080"
timeout: "5000"
`)
	assert.NoError(t, os.WriteFile(cfgFile, content, 0644))
}

func createErrorConfigFile(t *testing.T, cfgFile string) {
	content := []byte(`
	listen_port: "8080"
timeout: "3000"
`)

	assert.NoError(t, os.WriteFile(cfgFile, content, 0644))
}

func (suite *ConfigSuite) TestEnvConfig() {
	cfgFile := "./not-found.yaml"
	os.Setenv("khub_LISTEN_PORT", "8080")
	os.Setenv("khub_TIMEOUT", "2000")
	os.Setenv("khub_ENVIRONMENT", "Production")
	os.Setenv("khub_OIDC_ISSUER", "test")
	c := Load("1.2.3", cfgFile)
	suite.Equal(8080, c.ListenPort)
	suite.Equal(2000, c.Timeout)
	suite.Equal("test", c.OIDCIssuer)
	suite.Equal("Production", c.Environment)
	os.Unsetenv("khub_LISTEN_PORT")
	os.Unsetenv("khub_TIMEOUT")
	os.Unsetenv("khub_ENVIRONMENT")
}

func (suite *ConfigSuite) TestIsProduction() {
	cases := []struct {
		Environment string
		Expected    bool
	}{
		{"Production", true},
		{"", false},
		{"Dev", false},
		{"PRODUCTION", true},
		{"PrOdUcTiOn", true},
	}

	for _, c := range cases {
		os.Unsetenv("khub_ENVIRONMENT")
		if c.Environment != "" {
			os.Setenv("khub_ENVIRONMENT", c.Environment)
		}
		cfg := Load("test", "not-found")
		suite.Equal(c.Expected, cfg.IsProduction(), "IsProduction should be %t", c.Expected)
	}
	os.Unsetenv("khub_ENVIRONMENT")
}

func (suite *ConfigSuite) TestConfigFile() {
	viper.Reset()
	cfgFile := "./test-config-file.yml"
	createConfigFile(suite.T(), cfgFile)
	defer os.RemoveAll(cfgFile)

	c := Load("1.2.3", cfgFile)
	suite.Equal(8080, c.ListenPort, "Listen Port should be '8080'")
	suite.Equal(5000, c.Timeout, "Timeout should be '5000'")
	suite.Equal("1.2.3", c.Version, "Version should be '1.2.3'")
}

func (suite *ConfigSuite) TestDefaultConfigFile() {
	viper.Reset()
	cfgFile := "./config.yml"
	createConfigFile(suite.T(), cfgFile)
	defer os.RemoveAll(cfgFile)

	c := Load("1.2.3", "")
	suite.Equal(8080, c.ListenPort, "Listen Port should be '8080'")
	suite.Equal(5000, c.Timeout, "Timeout should be '5000'")
	suite.Equal("1.2.3", c.Version, "Version should be '1.2.3'")
}

func (suite *ConfigSuite) TestConfigFileErr() {
	cfgFile := "./test-config-file-error.yml"
	createErrorConfigFile(suite.T(), cfgFile)
	defer os.RemoveAll(cfgFile)
	c := Load("1.2.3", cfgFile)
	suite.Equal(8080, c.ListenPort, "Listen Port should be '8080'")
	suite.Equal(2000, c.Timeout, "Timeout should be '2000'")
	suite.Equal("1.2.3", c.Version, "Version should be '1.2.3'")
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
