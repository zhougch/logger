package logger

import (
	"go.uber.org/zap"
	"testing"
)

func TestLogger(t *testing.T) {
	cfg := new(LogConfig)
	cfg.Develop = true
	cfg.Level = "debug"
	cfg.Path = "./test.log"
	cfg.ErrorPath = "./test-error.log"

	log := NewRotateLogger(*cfg)
	defer log.Sync()

	log.Debug("debug", zap.String("aaa", "bbb"), zap.String("ccc", "ddd"))
	log.Info("info")
	log.Warn("warn")
	log.Error("error")

	log1 := log.Sugar()
	log1.Debugf("debug: %s", "aaa")
	log1.Infof("info: %s", "bbb")
	log1.Warnf("warn: %s", "ccc")
	log1.Errorf("error: %s", "ddd")
}
