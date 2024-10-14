package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	KeyProcess        = "process"
	KeyProcessingTime = "processingTime"
	KeyStartTime      = "startTime"
	KeyEndTime        = "endTime"
)

var Logger *zerolog.Logger

func InitLogger() *zerolog.Logger {
	startTime := time.Now()

	zerolog.ErrorFieldName = "error"
	zerolog.ErrorStackFieldName = "stack-trace"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.DurationFieldUnit = time.Microsecond

	fileWriter := &lumberjack.Logger{
		Filename:   "url_shortener.jsonl",
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     10,
		Compress:   true,
	}
	consoleWriter := zerolog.ConsoleWriter{
		Out:          os.Stderr,
		TimeFormat:   time.RFC3339Nano,
		NoColor:      false,
		TimeLocation: time.UTC,
	}
	output := zerolog.MultiLevelWriter(consoleWriter, fileWriter)
	logger := zerolog.New(output).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Int("gid", os.Getegid()).
		Int("uid", os.Getuid()).
		Logger()

	Logger = &logger

	Logger.Info().
		Str(KeyProcess, "InitLogger").
		Time(KeyStartTime, startTime).
		Time(KeyEndTime, time.Now()).
		Dur("duration", time.Since(startTime)).
		Msg("finish initiating logging")

	return &logger
}
