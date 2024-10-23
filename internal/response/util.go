package response

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/Alturino/url-shortener/internal/log"
)

func WriteJsonResponse(
	c context.Context,
	w http.ResponseWriter,
	jsonEncoder *json.Encoder,
	header map[string]string,
	body map[string]string,
	statusCode int,
) {
	logger := zerolog.Ctx(c)
	hashcode := c.Value(log.KeyHashcode).(string)

	for k, v := range header {
		w.Header().Add(k, v)
	}

	res := map[string]interface{}{}
	for k, v := range body {
		res[k] = v
	}

	w.WriteHeader(statusCode)
	err := jsonEncoder.Encode(&res)
	if err != nil {
		logger.Error().
			Err(err).
			Str(log.KeyProcess, "WriteJsonResponse").
			Str(log.KeyHashcode, hashcode).
			Msgf("failed to write json response with error=%s", err.Error())
	}
}
