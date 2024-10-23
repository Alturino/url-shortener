package log

import (
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

const (
	KeyEndTime        = "endTime"
	KeyHashcode       = "hashcode"
	KeyID             = "id"
	KeyProcess        = "process"
	KeyProcessingTime = "processingTime"
	KeyRequestBody    = "requestBody"
	KeyRequestHeader  = "requestHeader"
	KeyRequestIp      = "requesterIP"
	KeyRequestMethod  = "requestMethod"
	KeyRequestURI     = "requestURI"
	KeyRequestURL     = "requestURL"
	KeyShortUrl       = "shortUrl"
	KeyStartTime      = "startTime"
)

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
		Out:          os.Stdout,
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

	logger.Info().
		Str(KeyProcess, "InitLogger").
		Time(KeyStartTime, startTime).
		Time(KeyEndTime, time.Now()).
		Dur(KeyProcessingTime, time.Since(startTime)).
		Msg("finish initiating logging")

	return &logger
}
