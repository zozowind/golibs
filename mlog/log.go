package mlog

import (
	"io"
	"runtime/debug"

	"github.com/rs/zerolog"
)

const (
	//On 开启
	On = "on"
	//Off 关闭
	Off = "off"
	//Auto 自动
	Auto = "auto"
)

//	按日志级别打印和控制日志
//	多渠道日志输出
//	文件行号定位 默认auto Info以下不输出
//  输出格式定义
//  支持自定义KV对输出
//  日志染色功能
// 堆栈输出 默认Error以下输出

//LogSetting 日志设定
type LogSetting struct {
	Write  io.Writer
	Level  zerolog.Level
	Caller string
	Stack  string
}

//Logger 日志生成器
type Logger struct {
	Agents map[*LogSetting]*zerolog.Logger
}

//New 新建Logger
func New(settings []*LogSetting) *Logger {
	logger := &Logger{
		Agents: map[*LogSetting]*zerolog.Logger{},
	}
	for _, s := range settings {
		logAgent := zerolog.New(s.Write).With().Timestamp().Logger().Level(s.Level)
		// 默认3
		zerolog.CallerSkipFrameCount = 3
		logger.Agents[s] = &logAgent
	}
	return logger
}

//Print Print
func (l *Logger) Print(i ...interface{}) {
	for _, a := range l.Agents {
		a.Print(i)
	}
}

//Printf Printf
func (l *Logger) Printf(msg string, v ...interface{}) {
	for _, a := range l.Agents {
		a.Printf(msg, v...)
	}
}

//Debug log Debug
func (l *Logger) Debug(msg string) {
	for s, a := range l.Agents {
		e := a.Debug()
		if s.Caller == On {
			e = e.Caller()
		}
		e.Msg(msg)
		if s.Stack == On {
			a.Debug().Msgf("%s", debug.Stack())
		}
	}
}

//Debugf log debugf
func (l *Logger) Debugf(msg string, v ...interface{}) {
	for s, a := range l.Agents {
		e := a.Debug()
		if s.Caller == On {
			e = e.Caller()
		}
		e.Msgf(msg, v...)
		if s.Stack == On {
			a.Debug().Msgf("%s", debug.Stack())
		}
	}
}

//DebugKv log with kv in debug
func (l *Logger) DebugKv(kv map[string]interface{}, msg string, v ...interface{}) {
	for s, a := range l.Agents {
		e := a.Debug()
		if s.Caller == On {
			e = e.Caller()
		}
		e.Fields(kv).Msgf(msg, v...)
		if s.Stack == On {
			a.Debug().Msgf("%s", debug.Stack())
		}
	}
}

//Info log Info
func (l *Logger) Info(msg string) {
	for s, a := range l.Agents {
		e := a.Info()
		if s.Caller == On {
			e = e.Caller()
		}
		e.Msg(msg)
		if s.Stack == On {
			a.Info().Msgf("%s", debug.Stack())
		}
	}
}

//Infof log Infof
func (l *Logger) Infof(msg string, v ...interface{}) {
	for s, a := range l.Agents {
		e := a.Info()
		if s.Caller == On {
			e = e.Caller()
		}
		e.Msgf(msg, v...)
		if s.Stack == On {
			a.Info().Msgf("%s", debug.Stack())
		}
	}
}

//InfoKv log with kv in Info
func (l *Logger) InfoKv(kv map[string]interface{}, msg string, v ...interface{}) {
	for s, a := range l.Agents {
		e := a.Info()
		if s.Caller == On {
			e = e.Caller()
		}
		e.Fields(kv).Msgf(msg, v...)
		if s.Stack == On {
			a.Info().Msgf("%s", debug.Stack())
		}
	}
}

//Warn log Warn
func (l *Logger) Warn(msg string) {
	for s, a := range l.Agents {
		e := a.Warn()
		if s.Caller != Off {
			e = e.Caller()
		}
		e.Msg(msg)
		if s.Stack == On {
			a.Warn().Msgf("%s", debug.Stack())
		}
	}
}

//Warnf log Warnf
func (l *Logger) Warnf(msg string, v ...interface{}) {
	for s, a := range l.Agents {
		e := a.Warn()
		if s.Caller != Off {
			e = e.Caller()
		}
		e.Msgf(msg, v...)
		if s.Stack == On {
			a.Warn().Msgf("%s", debug.Stack())
		}
	}
}

//WarnKv log with kv in Warn
func (l *Logger) WarnKv(kv map[string]interface{}, msg string, v ...interface{}) {
	for s, a := range l.Agents {
		e := a.Warn()
		if s.Caller == On {
			e = e.Caller()
		}
		e.Fields(kv).Msgf(msg, v...)
		if s.Stack == On {
			a.Warn().Msgf("%s", debug.Stack())
		}
	}
}

//Error log Error
func (l *Logger) Error(msg string) {
	for s, a := range l.Agents {
		e := a.Error()
		if s.Caller != Off {
			e = e.Caller()
		}
		e.Msg(msg)
		if s.Stack != Off {
			a.Error().Msgf("%s", debug.Stack())
		}
	}
}

//Errorf log Errorf
func (l *Logger) Errorf(msg string, v ...interface{}) {
	for s, a := range l.Agents {
		e := a.Error()
		if s.Caller != Off {
			e = e.Caller()
		}
		e.Msgf(msg, v...)
		if s.Stack != Off {
			a.Error().Msgf("%s", debug.Stack())
		}
	}
}

//ErrorKv log with kv in Error
func (l *Logger) ErrorKv(kv map[string]interface{}, msg string, v ...interface{}) {
	for s, a := range l.Agents {
		e := a.Error()
		if s.Caller == On {
			e = e.Caller()
		}
		e.Fields(kv).Msgf(msg, v...)
		if s.Stack != Off {
			a.Error().Msgf("%s", debug.Stack())
		}
	}
}

//Fatal log Fatal
func (l *Logger) Fatal(msg string) {
	for s, a := range l.Agents {
		e := a.Fatal()
		if s.Caller != Off {
			e = e.Caller()
		}
		e.Msg(msg)
		if s.Stack != Off {
			a.Error().Msgf("%s", debug.Stack())
		}
	}
}

//Fatalf log Fatalf
func (l *Logger) Fatalf(msg string, v ...interface{}) {
	for s, a := range l.Agents {
		e := a.Fatal()
		if s.Caller != Off {
			e = e.Caller()
		}
		e.Msgf(msg, v...)
		if s.Stack != Off {
			a.Fatal().Msgf("%s", debug.Stack())
		}
	}
}

//FatalKv log with kv in Fatal
func (l *Logger) FatalKv(kv map[string]interface{}, msg string, v ...interface{}) {
	for s, a := range l.Agents {
		e := a.Fatal()
		if s.Caller == On {
			e = e.Caller()
		}
		e.Fields(kv).Msgf(msg, v...)
		if s.Stack != Off {
			a.Fatal().Msgf("%s", debug.Stack())
		}
	}
}
