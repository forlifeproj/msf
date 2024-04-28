package config

/*
@version v1.0
@author declan
*/

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	LogLevelNull    = 0
	LogLevelTrace   = 1
	LogLevelDebug   = 2
	LogLevelInfo    = 3
	LogLevelWarning = 4
	LogLevelError   = 5
	LogLevelFatal   = 6
)

// ConfPath 默认配置文件地址，可修改
var ConfPath = "../conf/config.toml"

var logLevelMap = map[string]uint8{
	"trace": LogLevelTrace,
	"debug": LogLevelDebug,
	"info":  LogLevelInfo,
	"warn":  LogLevelWarning,
	"error": LogLevelError,
	"fatal": LogLevelFatal,
}
var logLevelStrMap = map[uint8]string{
	LogLevelTrace:   "trace",
	LogLevelDebug:   "debug",
	LogLevelInfo:    "info",
	LogLevelWarning: "warn",
	LogLevelError:   "error",
	LogLevelFatal:   "fatal",
}

// LogLevel 日志级别
type LogLevel uint8

// UnmarshalText 通过字符串解析日志级别
func (l *LogLevel) UnmarshalText(text []byte) error {
	level, ok := logLevelMap[strings.ToLower(string(text))]
	if !ok {
		return fmt.Errorf("not support log level %v", string(text))
	}
	*l = LogLevel(level)
	return nil
}

// String 日志级别字符串展示
func (l LogLevel) String() string {
	name, ok := logLevelStrMap[uint8(l)]
	if ok {
		return name
	}

	return "unknown"
}

// Level return uint8 level
func (l LogLevel) Level() uint8 {
	return uint8(l)
}

// Value return uint8 level
func (l LogLevel) Value() uint8 {
	return uint8(l)
}

// LogSize 日志文件大小 B K M G
type LogSize int64

// UnmarshalText 通过字符串解析日志大小
func (l *LogSize) UnmarshalText(text []byte) error {
	if len(text) < 2 { //至少两个字节
		return fmt.Errorf("not support log size %v", string(text))
	}
	c := strings.ToLower(string(text[len(text)-1:])) //最后一个字符
	n, e := strconv.ParseInt(string(text[:len(text)-1]), 10, 64)
	if e != nil {
		return e
	}
	if c == "k" {
		n *= 1024
	} else if c == "m" {
		n *= 1024 * 1024
	} else if c == "g" {
		n *= 1024 * 1024 * 1024
	}
	*l = LogSize(n)
	return nil
}

// String 日志大小字符串展示
func (l LogSize) String() string {
	if l < 1024 {
		return fmt.Sprintf("%dB", int64(l))
	} else if l < 1024*1024 {
		return fmt.Sprintf("%dK", int64(l)/1024)
	} else if l < 1024*1024*1024 {
		return fmt.Sprintf("%dM", int64(l)/1024/1024)
	} else if l < 1024*1024*1024*1024 {
		return fmt.Sprintf("%dG", int64(l)/1024/1024/1024)
	}

	return "unknown"
}

// Size return int64 size
func (l LogSize) Size() int64 {
	return int64(l)
}

// Value return int64 size
func (l LogSize) Value() int64 {
	return int64(l)
}

// Duration duration for config parse
type Duration time.Duration

func (d Duration) String() string {
	dd := time.Duration(d)
	return dd.String()
}

// GoString  duration go string
func (d Duration) GoString() string {
	dd := time.Duration(d)
	return dd.String()
}

// Duration duration
func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

// Value duration
func (d Duration) Value() time.Duration {
	return time.Duration(d)
}

// UnmarshalText 字符串解析时间
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	dd, err := time.ParseDuration(string(text))
	if err == nil {
		*d = Duration(dd)
	}
	return err
}

// Parse parse config with default and config file ../conf/config.toml
func Parse(c interface{}) error {
	return ParseConfigWithoutDefaults(c)
}

// ParseConfig same as Parse
func ParseConfig(c interface{}) error {
	return Parse(c)
}

// ParseConfigWithPath 自己定义配置文件路径
func ParseConfigWithPath(c interface{}, path string) error {
	if _, err := toml.DecodeFile(path, c); err != nil {
		return err
	}
	return nil
}

// ParseConfigWithoutDefaults no default value
func ParseConfigWithoutDefaults(c interface{}) error {
	if _, err := toml.DecodeFile(ConfPath, c); err != nil {
		return err
	}
	return nil
}
