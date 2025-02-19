package logging

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type LoggingSuite struct {
	suite.Suite
}

func TestLoggingSuite(t *testing.T) {
	suite.Run(t, new(LoggingSuite))
}

func (suite *LoggingSuite) TestInitLogger() {
	InitLogger("debug")
	suite.Equal(logLevels["debug"], zerolog.GlobalLevel())
}
