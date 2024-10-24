package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog"

	"github.com/Alturino/url-shortener/internal/log"
	"github.com/Alturino/url-shortener/internal/request"
	"github.com/Alturino/url-shortener/internal/response"
	"github.com/Alturino/url-shortener/internal/service"
)

type UrlController struct {
	service *service.UrlService
}

func AttachUrlController(mux *http.ServeMux, service *service.UrlService) {
	controller := UrlController{service: service}
	mux.HandleFunc("GET /urls/{shortUrl}", controller.GetUrlByShortUrl)
	mux.HandleFunc("GET /urls/{shortUrl}/stats", controller.GetUrlByShortUrlDetail)
	mux.HandleFunc("PUT /urls/{shortUrl}", controller.UpdateUrl)
	mux.HandleFunc("DELETE /urls/{shortUrl}", controller.DeleteUrl)
	mux.HandleFunc("POST /urls", controller.InsertUrl)
}

func (u *UrlController) InsertUrl(w http.ResponseWriter, r *http.Request) {
	requestStartTime := r.Context().Value(log.KeyRequestReceivedAt).(time.Time)

	logger := zerolog.Ctx(r.Context())

	logger.Info().Msg("decoding requestBody")
	req := request.UrlRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error().
			Err(err).
			Str(log.KeyProcess, "InsertUrl").
			Msg("failed decoding requestBody")
		response.WriteJsonResponse(
			r.Context(),
			w,
			map[string]string{},
			map[string]interface{}{},
			http.StatusBadRequest,
		)
		return
	}
	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Any(log.KeyRequestBody, req).
			Str(log.KeyProcess, "InsertUrl").
			Str(log.KeyUrl, req.Url)
	})
	c := logger.WithContext(r.Context())
	logger.Info().Msg("decoded requestBody")

	logger.Info().Msgf("validating url=%s", req.Url)
	validatedUrl, err := url.Parse(req.Url)
	if err != nil {
		logger.Error().
			Err(err).
			Dur(log.KeyRequestProcessingTime, time.Since(requestStartTime)).
			Msgf("failed validating url=%s with error=%s", req.Url, err.Error())
		response.WriteJsonResponse(
			c,
			w,
			map[string]string{},
			map[string]interface{}{},
			http.StatusBadRequest,
		)
		return
	}
	logger.Info().Msgf("validated url=%s", req.Url)

	logger.Info().Msgf("inserting url=%s", req.Url)
	inserted, err := u.service.InsertUrl(c, *validatedUrl)
	if err != nil {
		logger.Error().
			Err(err).
			Dur(log.KeyRequestProcessingTime, time.Since(requestStartTime)).
			Msgf("failed inserting url=%s with error=%s", req.Url, err.Error())
		response.WriteJsonResponse(
			c,
			w,
			map[string]string{},
			map[string]interface{}{},
			http.StatusBadRequest,
		)
		return
	}
	logger.Info().
		Dur(log.KeyRequestProcessingTime, time.Since(requestStartTime)).
		Str(log.KeyShortUrl, inserted.ShortUrl).
		Msgf("inserted url=%s shortUrl=%s", req.Url, inserted.ShortUrl)

	response.WriteJsonResponse(
		c,
		w,
		map[string]string{},
		map[string]interface{}{
			"status":  "success",
			"message": fmt.Sprintf("inserted url=%s to shortUrl=%s", req.Url, inserted.ShortUrl),
			"data":    inserted,
		},
		http.StatusOK,
	)
}

func (u *UrlController) UpdateUrl(w http.ResponseWriter, r *http.Request) {
	requestStartTime := r.Context().Value(log.KeyRequestReceivedAt).(time.Time)

	shortUrl := r.PathValue("shortUrl")
	logger := zerolog.Ctx(r.Context())

	req, decodingStartTime := request.UrlRequest{}, time.Now()
	logger.Info().Time(log.KeyStartTime, decodingStartTime).Msg("decoding requestBody")
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Error().Err(err).
			Dur(log.KeyProcessingTime, time.Since(decodingStartTime)).
			Dur(log.KeyRequestProcessingTime, time.Since(requestStartTime)).
			Time(log.KeyStartTime, decodingStartTime).
			Msg("failed decoding requestBody")

		response.WriteJsonResponse(
			r.Context(),
			w,
			map[string]string{},
			map[string]interface{}{},
			http.StatusBadRequest,
		)
		return
	}
	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(log.KeyProcess, "UpdateUrl").
			Any(log.KeyRequestBody, req).
			Str(log.KeyShortUrl, shortUrl).
			Str(log.KeyNewUrl, req.Url)
	})
	c := logger.WithContext(r.Context())
	logger.Info().
		Dur(log.KeyProcessingTime, time.Since(decodingStartTime)).
		Time(log.KeyStartTime, decodingStartTime).
		Msg("decoded requestBody")

	logger.Info().Msgf("validating url=%s", req.Url)
	validatingUrlStartTime := time.Now()
	validatedUrl, err := url.Parse(req.Url)
	if err != nil {
		logger.Error().
			Err(err).
			Dur(log.KeyProcessingTime, time.Since(validatingUrlStartTime)).
			Dur(log.KeyRequestProcessingTime, time.Since(requestStartTime)).
			Time(log.KeyStartTime, validatingUrlStartTime).
			Msgf("failed validating url=%s with error=%s", req.Url, err.Error())
		response.WriteJsonResponse(
			c,
			w,
			map[string]string{},
			map[string]interface{}{},
			http.StatusBadRequest,
		)
		return
	}
	logger.Info().
		Dur(log.KeyProcessingTime, time.Since(validatingUrlStartTime)).
		Time(log.KeyStartTime, validatingUrlStartTime).
		Msgf("validated url=%s", req.Url)

	logger.Info().Msgf("updating url=%s", req.Url)
	updateUrlStartTime := time.Now()
	updated, err := u.service.UpdateUrl(c, *validatedUrl, shortUrl)
	if err != nil {
		logger.Error().
			Err(err).
			Dur(log.KeyProcessingTime, time.Since(updateUrlStartTime)).
			Dur(log.KeyRequestProcessingTime, time.Since(requestStartTime)).
			Time(log.KeyStartTime, updateUrlStartTime).
			Msgf("failed updating url=%s with error=%s", req.Url, err.Error())
		response.WriteJsonResponse(
			c,
			w,
			map[string]string{},
			map[string]interface{}{},
			http.StatusBadRequest,
		)
		return
	}
	logger.Info().
		Dur(log.KeyProcessingTime, time.Since(updateUrlStartTime)).
		Dur(log.KeyRequestProcessingTime, time.Since(requestStartTime)).
		Time(log.KeyStartTime, updateUrlStartTime).
		Msgf("updated url=%s shortUrl=%s", req.Url, updated.ShortUrl)

	response.WriteJsonResponse(
		c,
		w,
		map[string]string{},
		map[string]interface{}{
			"status":  "success",
			"message": fmt.Sprintf("updated url=%s to shortUrl=%s", req.Url, updated.ShortUrl),
			"data":    updated,
		},
		http.StatusOK,
	)
}

func (u *UrlController) DeleteUrl(w http.ResponseWriter, r *http.Request) {
	requestStartTime := r.Context().Value(log.KeyRequestReceivedAt).(time.Time)
	logger := zerolog.Ctx(r.Context())
	shortUrl := r.PathValue("shortUrl")

	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(log.KeyProcess, "DeleteUrl").
			Str(log.KeyShortUrl, shortUrl)
	})
	c := logger.WithContext(r.Context())

	logger.Info().Msgf("deleting shortUrl=%s", shortUrl)
	deleted, err := u.service.DeleteUrl(r.Context(), shortUrl)
	if err != nil {
		logger.Error().
			Err(err).
			Dur(log.KeyRequestProcessingTime, time.Since(requestStartTime)).
			Msgf("failed deleting shortUrl=%s with error=%s", shortUrl, err.Error())
		response.WriteJsonResponse(
			c,
			w,
			map[string]string{},
			map[string]interface{}{},
			http.StatusBadRequest,
		)
		return
	}
	logger.Info().
		Dur(log.KeyRequestProcessingTime, time.Since(requestStartTime)).
		Msgf("deleted url=%s shortUrl=%s", deleted.Url, deleted.ShortUrl)

	response.WriteJsonResponse(
		c,
		w,
		map[string]string{},
		map[string]interface{}{
			"status":  "success",
			"message": fmt.Sprintf("deleted url=%s to shortUrl=%s", deleted.Url, deleted.ShortUrl),
			"data":    deleted,
		},
		http.StatusOK,
	)
}

func (u *UrlController) GetUrlByShortUrl(w http.ResponseWriter, r *http.Request) {
	requestStartTime := r.Context().Value(log.KeyRequestReceivedAt).(time.Time)
	logger := zerolog.Ctx(r.Context())
	shortUrl := r.PathValue("shortUrl")
	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(log.KeyProcess, "GetUrlByShortUrlDetail").
			Str(log.KeyShortUrl, shortUrl)
	})
	c := logger.WithContext(r.Context())

	logger.Info().Msgf("finding shortUrl=%s", shortUrl)
	existed, err := u.service.GetUrlByShortUrl(c, shortUrl)
	if err != nil {
		logger.Error().
			Err(err).
			Dur(log.KeyRequestProcessingTime, time.Since(requestStartTime)).
			Msgf("failed finding shortUrl=%s with error=%s", shortUrl, err.Error())
		response.WriteJsonResponse(
			c,
			w,
			map[string]string{},
			map[string]interface{}{},
			http.StatusBadRequest,
		)
		return
	}
	logger.Info().
		Dur(log.KeyRequestProcessingTime, time.Since(requestStartTime)).
		Msgf("found url=%s shortUrl=%s", existed.Url, existed.ShortUrl)

	response.WriteJsonResponse(
		c,
		w,
		map[string]string{},
		map[string]interface{}{
			"status":  "success",
			"message": fmt.Sprintf("found url=%s to shortUrl=%s", existed.Url, existed.ShortUrl),
			"data":    existed,
		},
		http.StatusOK,
	)
}

func (u *UrlController) GetUrlByShortUrlDetail(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())
	shortUrl := r.PathValue("shortUrl")

	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(log.KeyProcess, "GetUrlByShortUrlDetail").
			Str(log.KeyShortUrl, shortUrl)
	})
	c := logger.WithContext(r.Context())

	logger.Info().Msgf("finding shortUrl=%s", shortUrl)
	existed, err := u.service.GetUrlByShortUrlDetail(c, shortUrl)
	if err != nil {
		logger.Error().
			Err(err).
			Msgf("failed finding shortUrl=%s with error=%s", shortUrl, err.Error())
		response.WriteJsonResponse(
			c,
			w,
			map[string]string{},
			map[string]interface{}{},
			http.StatusBadRequest,
		)
		return
	}
	logger.Info().Msgf("found url=%s shortUrl=%s", existed.Url, existed.ShortUrl)

	response.WriteJsonResponse(
		c,
		w,
		map[string]string{},
		map[string]interface{}{
			"status":  "success",
			"message": fmt.Sprintf("found url=%s to shortUrl=%s", existed.Url, existed.ShortUrl),
			"data":    existed,
		},
		http.StatusOK,
	)
}
