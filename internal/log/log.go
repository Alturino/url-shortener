package log

import (
	"context"
	"io"
	"os"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

const (
	KeyHashcode           = "hashcode"
	KeyUrlID              = "urlId"
	KeyUrl                = "url"
	KeyProcess            = "process"
	KeyRequestBody        = "requestBody"
	KeyRequestHeader      = "requestHeader"
	KeyRequestHost        = "host"
	KeyRequestIp          = "requesterIP"
	KeyRequestMethod      = "requestMethod"
	KeyRequestProcessedAt = "requestProcessedAt"
	KeyNewUrl             = "newUrl"
	KeyOldUrl             = "oldUrl"
	KeyRequestURI         = "requestURI"
	KeyRequestURL         = "requestURL"
	KeyShortUrl           = "shortUrl"
	KeyConfig             = "config"
)

type hashcode struct{}

func HashcodeFromContext(c context.Context) string {
	return c.Value(hashcode{}).(string)
}

func AttachHashcodeToContext(c context.Context, h string) context.Context {
	return context.WithValue(c, hashcode{}, h)
}

var (
	once   sync.Once
	logger *zerolog.Logger
)

func InitLogger() *zerolog.Logger {
	once.Do(func() {
		zerolog.DurationFieldUnit = time.Microsecond
		zerolog.ErrorFieldName = "error"
		zerolog.ErrorStackFieldName = "stack-trace"
		zerolog.LevelFieldName = "level"
		zerolog.MessageFieldName = "message"
		zerolog.TimestampFieldName = "timestamp"

		fileWriter := &lumberjack.Logger{
			Filename:   "url_shortener.jsonl",
			MaxSize:    10,
			MaxBackups: 10,
			MaxAge:     10,
			Compress:   true,
		}
		var logOutput io.Writer = os.Stdout
		if os.Getenv("env") == "dev" {
			logOutput = zerolog.ConsoleWriter{
				Out:          os.Stdout,
				TimeFormat:   time.RFC3339Nano,
				NoColor:      false,
				TimeLocation: time.UTC,
			}
		}
		output := zerolog.MultiLevelWriter(logOutput, fileWriter)
		log := zerolog.New(output).
			Level(zerolog.TraceLevel).
			With().
			Timestamp().
			Caller().
			Int("pid", os.Getpid()).
			Int("gid", os.Getegid()).
			Int("uid", os.Getuid()).
			Logger()
		logger = &log

		logger.Info().
			Str(KeyProcess, "InitLogger").
			Msg("finish initiating logging")
	})
	return logger
}
