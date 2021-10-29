package global

import (
	"go.uber.org/zap"

	"v6tool/config"

	"github.com/spf13/viper"
)

var (
	V6TOOL_CONFIG config.Server
	V6TOOL_VP     *viper.Viper
	V6TOOL_LOG    *zap.Logger
)
