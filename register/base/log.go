package base

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
	"waho/register"
)

type LogRegister struct {
	register.BaseRegister
}

func init() {

	logFileDir, _ := filepath.Abs(register.GetConf.Section("logs").GetDefault("path", "./logs"))
	os.MkdirAll(logFileDir, os.ModePerm)

	prefix := register.GetConf.Section("logs").GetDefault("name", "waho")
	logFilePath := path.Join(logFileDir, prefix + ".Info.%Y%m%d%H.log")
	writer, _ := rotatelogs.New(
		logFilePath,
		//rotatelogs.WithMaxAge(time.Hour*24),
		//rotatelogs.WithLinkName(logFileDir), // 软链
		rotatelogs.WithRotationTime(time.Hour*1), // 1小时切割一次
	)


	formatter := &log.TextFormatter{}
	formatter.CallerPrettyfier = func(frame *runtime.Frame) (function string, file string) {
		function = frame.Function
		dir, filename := path.Split(frame.File)
		f := path.Base(dir)
		return function, fmt.Sprintf("%s/%s:%d", f, filename, frame.Line)
	}
	formatter.DisableTimestamp = false
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05.000000"
	log.SetFormatter(formatter)

	logFilePathWarn := path.Join(logFileDir, prefix + ".Warn.%Y%m%d%H.log")
	writerWarn, _ := rotatelogs.New(logFilePathWarn,rotatelogs.WithRotationTime(time.Hour*1),)
	logFilePathError := path.Join(logFileDir, prefix + ".Error.%Y%m%d%H.log")
	writerError, _ := rotatelogs.New(logFilePathError,rotatelogs.WithRotationTime(time.Hour*1),)
	logFilePathFatal := path.Join(logFileDir, prefix + ".Fatal.%Y%m%d%H.log")
	writerFatal, _ := rotatelogs.New(logFilePathFatal,rotatelogs.WithRotationTime(time.Hour*1),)
	logFilePathPanic := path.Join(logFileDir, prefix + ".Panic.%Y%m%d%H.log")
	writerPanic, _ := rotatelogs.New(logFilePathPanic,rotatelogs.WithRotationTime(time.Hour*1),)
	logFilePathDebug := path.Join(logFileDir, prefix + ".Debug.%Y%m%d%H.log")
	writerDebug, _ := rotatelogs.New(logFilePathDebug,rotatelogs.WithRotationTime(time.Hour*1),)
	logFilePathTrace := path.Join(logFileDir, prefix + ".Trace.%Y%m%d%H.log")
	writerTrace, _ := rotatelogs.New(logFilePathTrace,rotatelogs.WithRotationTime(time.Hour*1),)
	log.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			log.InfoLevel:  writer,
			log.WarnLevel:  writerWarn,
			log.ErrorLevel: writerError,
			log.FatalLevel: writerFatal,
			log.PanicLevel: writerPanic,
			log.DebugLevel: writerDebug,
			log.TraceLevel: writerTrace,
		},
		formatter,
	))
	//log.SetReportCaller(true)
}

func (logRegister *LogRegister) Init() {
	log.Info("log init")
}