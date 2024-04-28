package fllog

import (
	"fmt"

	"github.com/forlifeproj/msf/config"
)

type LogCfg struct {
	LogConf struct {
		Name       string `default:"../log/fllog.log"`
		Level      string `default:"INFO"`
		MaxSize    int    `default:"1073741824"`
		MaxAge     int    `default:"30"` //默认最多保存10个日志文件
		MaxBackups int    `default:"10"` // 最大保存日志数量
	}
}

func Init(cfg string) error {
	logCfg := LogCfg{}

	if err := config.ParseConfigWithPath(&logCfg, cfg); err != nil {
		fmt.Printf("load logcfg failed. err:%+v cfg:%s", err, cfg)
		return err
	}
	fmt.Printf("cfg:%s logCfg:%+v", cfg, logCfg)
	builder := NewLogUtilsBuilder(
		logCfg.LogConf.Level,
		logCfg.LogConf.Name,
		logCfg.LogConf.MaxSize,
		logCfg.LogConf.MaxAge,
		logCfg.LogConf.MaxBackups,
		false,
		true,
	)
	logUtils := NewLogUtils().SetBuilder(builder)
	err := logUtils.Init()
	if err != nil {
		fmt.Printf("logUtils init failed. err:%+v", err)
		return err
	}
	fmt.Printf("log init succ. cfg:%s logCfg:%+v", cfg, logCfg)
	Debug("log init succ. cfg:%s logCfg:%+v", cfg, logCfg)
	return nil
}

func Debug(f string, p ...interface{}) {
	if Allow("DEBUG") {
		msg := fmt.Sprintf(f, p...)
		utils := NewLogUtils()
		utils.getLogsUtils().Debug(msg)
	}
}

func Info(f string, p ...interface{}) {
	if Allow("INFO") {
		msg := fmt.Sprintf(f, p...)
		utils := NewLogUtils()
		utils.getLogsUtils().Info(msg)
	}
}

func Warn(f string, p ...interface{}) {
	if Allow("WARN") {
		msg := fmt.Sprintf(f, p...)
		utils := NewLogUtils()
		utils.getLogsUtils().Warn(msg)
	}
}

func Error(f string, p ...interface{}) {
	if Allow("ERROR") {
		msg := fmt.Sprintf(f, p...)
		utils := NewLogUtils()
		utils.getLogsUtils().Error(msg)
	}
}
