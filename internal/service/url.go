package service

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"

	"github.com/Alturino/url-shortener/internal/cache"
	"github.com/Alturino/url-shortener/internal/log"
	"github.com/Alturino/url-shortener/internal/repository"
)

const name = "github.com/Alturino/url-shortener"

var tracer = otel.Tracer(name)

type UrlService struct {
	cache   *redis.Client
	db      *sql.DB
	encoder *base64.Encoding
	queries *repository.Queries
}

func NewUrlService(
	cache *redis.Client,
	db *sql.DB,
	encoder *base64.Encoding,
	queries *repository.Queries,
) *UrlService {
	return &UrlService{cache: cache, db: db, queries: queries, encoder: encoder}
}

func (s *UrlService) InsertUrl(
	c context.Context,
	param url.URL,
) (repository.Url, error) {
	c, span := tracer.Start(c, "UrlService InsertUrl")
	defer span.End()

	logger := zerolog.Ctx(c).With().Logger()

	logger.Info().Msg("generating uuid")
	id, err := uuid.NewRandom()
	if err != nil {
		logger.Error().Err(err).Msgf("failed to generate uuid for url with error=%s", err.Error())
		return repository.Url{}, err
	}
	logger.Info().Msgf("generated uuid=%s", id.String())

	logger.Info().Msgf("starting to encode url=%s id=%s to shortUrl", param.String(), id.String())
	encoded := s.encoder.EncodeToString([]byte(id.String()))
	shortUrl := encoded[:5]
	logger.Info().Msgf("encoded url=%s id=%s to shortUrl=%s", param.String(), id.String(), shortUrl)

	logger.Info().Msgf("inserting url=%s id=%s shortUrl=%s", param.String(), id.String(), shortUrl)
	inserted, err := s.queries.InsertUrl(c, repository.InsertUrlParams{
		ID:       id,
		Url:      param.String(),
		ShortUrl: shortUrl,
	})
	if err != nil {
		err = fmt.Errorf(
			"failed when inserting url=%s with id=%s to database with error=%w",
			param.String(),
			id.String(),
			err,
		)
		logger.Error().Err(err).Msg(err.Error())
		return repository.Url{}, err
	}
	logger.Info().
		Str(log.KeyUrlID, inserted.ID.String()).
		Msgf("inserted url=%s id=%s shortUrl=%s", param.String(), id, shortUrl)

	logger.Info().
		Msgf("inserting shortUrl=%s url=%s id=%s to cache", shortUrl, param.String(), id.String())
	err = s.cache.JSONSet(c, fmt.Sprintf(cache.KeyUrl, shortUrl), "$", inserted).Err()
	if err != nil {
		err = fmt.Errorf(
			"inserting shortUrl=%s url=%s id=%s to cache with error=%w",
			shortUrl,
			param.String(),
			id.String(),
			err,
		)
		logger.Error().Err(err).Msg(err.Error())
		return inserted, err
	}
	logger.Info().
		Msgf("inserting shortUrl=%s url=%s id=%s to cache", shortUrl, param.String(), id.String())

	return inserted, nil
}

func (s *UrlService) UpdateUrl(
	c context.Context,
	url url.URL,
	shortUrl string,
) (repository.Url, error) {
	c, span := tracer.Start(c, "UrlService UpdateUrl")
	defer span.End()

	logger := zerolog.Ctx(c).With().Logger()

	logger.Info().Msgf("finding shortUrl=%s", shortUrl)
	existing, err := s.queries.FindUrlByShortUrl(c, shortUrl)
	if err != nil {
		logger.Error().
			Err(err).
			Msgf("shortUrl=%s not found", shortUrl)
		return repository.Url{}, err
	}
	logger = logger.With().
		Str(log.KeyOldUrl, existing.Url).
		Str(log.KeyUrlID, existing.ID.String()).
		Logger()
	logger.Info().Msgf("found shortUrl=%s", shortUrl)

	logger.Info().
		Msgf("updating url=%s id=%s to url=%s", existing.Url, existing.ID.String(), url.String())
	updated, err := s.queries.UpdateUrl(
		c,
		repository.UpdateUrlParams{ShortUrl: shortUrl, Url: url.String()},
	)
	if err != nil {
		logger.Error().
			Err(err).
			Msgf("failed updating url=%s id=%s with error=%s", existing.Url, existing.ID.String(), err.Error())
		return repository.Url{}, err
	}
	logger.Info().
		Msgf("updated url=%s id=%s to url=%s", existing.Url, existing.ID.String(), url.String())

	logger.Info().
		Msgf("updating shortUrl=%s url=%s id=%s to cache", shortUrl, url.String(), existing.ID.String())
	err = s.cache.JSONSet(c, fmt.Sprintf(cache.KeyUrl, shortUrl), "$", updated).Err()
	if err != nil {
		err = fmt.Errorf(
			"failed updating shortUrl=%s url=%s id=%s to cache with error=%w",
			shortUrl,
			url.String(),
			existing.ID.String(),
			err,
		)
		logger.Error().Err(err).Msg(err.Error())
		return updated, err
	}
	logger.Info().
		Msgf("updated shortUrl=%s url=%s id=%s to cache", shortUrl, url.String(), existing.ID.String())

	return updated, nil
}

func (s *UrlService) DeleteUrl(
	c context.Context,
	shortUrl string,
) (repository.Url, error) {
	c, span := tracer.Start(c, "UrlService DeleteUrl")
	defer span.End()

	logger := zerolog.Ctx(c).With().Logger()

	logger.Info().Msgf("deleting shortUrl=%s", shortUrl)
	deleted, err := s.queries.DeleteUrlByShortUrl(c, shortUrl)
	if err != nil {
		logger.Error().
			Err(err).
			Msgf("failed deleting shortUrl=%s not found", shortUrl)
		return repository.Url{}, err
	}
	logger.Info().Msgf("deleted url=%s id=%s", deleted.Url, deleted.ID.String())

	logger.Info().Msgf("deleting shortUrl=%s from cache", shortUrl)
	err = s.cache.JSONDel(c, fmt.Sprintf(cache.KeyUrl, shortUrl), "$").Err()
	if err != nil {
		err = fmt.Errorf("failed deleting shortUrl=%s from cache with error=%w", shortUrl, err)
		logger.Error().Err(err).Msg(err.Error())
		return deleted, err
	}
	logger.Info().Msgf("deleted shortUrl=%s from cache", shortUrl)

	return deleted, nil
}

func (s *UrlService) GetUrlByShortUrl(
	c context.Context,
	shortUrl string,
) (repository.Url, error) {
	c, span := tracer.Start(c, "UrlService GetUrlByShortUrl")
	defer span.End()

	logger := zerolog.Ctx(c).With().Logger()

	logger.Info().Msgf("incrementing visited_count for shortUrl=%s", shortUrl)
	err := s.cache.JSONNumIncrBy(c, fmt.Sprintf(cache.KeyUrl, shortUrl), "$.visited_count", 1).Err()
	if err != nil {
		err = fmt.Errorf(
			"failed incrementing visited_count for shortUrl=%s with error=%w",
			shortUrl,
			err,
		)
		logger.Error().Err(err).Msg(err.Error())
		return repository.Url{}, err
	}
	logger.Info().Msgf("incremented visited_count for shortUrl=%s", shortUrl)

	logger.Info().Msgf("finding shortUrl=%s from cache", shortUrl)
	jsonCache, err := s.cache.JSONGet(c, fmt.Sprintf(cache.KeyUrl, shortUrl)).Result()
	if err != nil {
		err = fmt.Errorf("failed finding shortUrl=%s from cache with error=%w", shortUrl, err)
		logger.Error().Err(err).Msg(err.Error())
		return repository.Url{}, err
	}
	logger.Info().Msgf("found shortUrl=%s from cache", shortUrl)

	updated := repository.Url{}
	logger.Info().Msg("marshalling jsonCache to url struct")
	err = json.Unmarshal([]byte(jsonCache), &updated)
	if err != nil {
		err = fmt.Errorf("failed marshalling jsonCache to url struct with error=%w", err)
		logger.Error().Err(err).Msg(err.Error())
		return repository.Url{}, err
	}
	logger.Info().Msgf("marshalled jsonCache to url struct")

	logger.Info().Msgf("updating visited_count for shortUrl=%s", shortUrl)
	_, err = s.queries.UpdateVisitedCountUrl(
		c,
		repository.UpdateVisitedCountUrlParams{ID: updated.ID, VisitedCount: updated.VisitedCount},
	)
	logger.Info().
		Msgf("updated visited_count for shortUrl=%s to %d", shortUrl, updated.VisitedCount)
	if err != nil {
		err = fmt.Errorf(
			"failed updating visited_count for shortUrl=%s with error=%w",
			shortUrl,
			err,
		)
		logger.Error().Err(err).Msg(err.Error())
		return updated, nil
	}

	return updated, nil
}

func (s *UrlService) GetUrlByShortUrlDetail(
	c context.Context,
	shortUrl string,
) (repository.Url, error) {
	c, span := tracer.Start(c, "UrlService GetUrlByShortUrlDetail")
	defer span.End()

	logger := zerolog.Ctx(c).With().Logger()

	logger.Info().Msgf("finding shortUrl=%s from cache", shortUrl)
	jsonCache, err := s.cache.JSONGet(c, fmt.Sprintf(cache.KeyUrl, shortUrl)).Result()
	if err != nil {
		err = fmt.Errorf("failed finding shortUrl=%s from cache with error=%w", shortUrl, err)
		logger.Error().Err(err).Msg(err.Error())

		logger.Info().Msgf("finding shortUrl=%s", shortUrl)
		existing, err := s.queries.FindUrlByShortUrl(c, shortUrl)
		if err != nil {
			err = fmt.Errorf("failed finding shortUrl=%s with error=%w", shortUrl, err)
			logger.Error().Err(err).Msg(err.Error())
			return repository.Url{}, err
		}
		logger.Info().Any("existing", existing).Msgf("finding shortUrl=%s", shortUrl)
		return existing, err
	}
	logger.Info().Msgf("found shortUrl=%s from cache", shortUrl)

	logger.Info().Msg("marshalling jsonCache to url struct")
	url := repository.Url{}
	err = json.Unmarshal([]byte(jsonCache), &url)
	if err != nil {
		err = fmt.Errorf("failed marshalling jsonCache to url struct with error=%w", err)
		logger.Error().Err(err).Msg(err.Error())
		return repository.Url{}, err
	}
	logger.Info().Msgf("marshalled jsonCache to url struct")

	return url, nil
}
