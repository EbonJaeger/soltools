package cmd

import (
	log2 "log"
	"os"

	"github.com/DataDrake/waterlog"
	"github.com/DataDrake/waterlog/format"
	"github.com/DataDrake/waterlog/level"
)

// NewLogger initializes a Waterlog logging instance and returns it.
func NewLogger() *waterlog.WaterLog {
	logger := waterlog.New(os.Stdout, "", log2.Ltime)
	logger.SetLevel(level.Info)
	logger.SetFormat(format.Min)

	return logger
}