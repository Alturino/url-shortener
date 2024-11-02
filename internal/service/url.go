package service

import (
	"context"
	"database/sql"
	"encoding/base64"
	"net/url"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"

	"github.com/Alturino/url-shortener/internal/log"
	"github.com/Alturino/url-shortener/internal/repository"
)

const name = "github.com/Alturino/url-shortener"

var tracer = otel.Tracer(name)

type UrlService struct {
	db      *sql.DB
	queries *repository.Queries
	encoder *base64.Encoding
}

func NewUrlService(
	db *sql.DB,
	queries *repository.Queries,
	encoder *base64.Encoding,
) *UrlService {
	return &UrlService{db: db, queries: queries, encoder: encoder}
}

func (s *UrlService) InsertUrl(
	c context.Context,
	url url.URL,
) (repository.Url, error) {
	c, span := tracer.Start(c, "UrlService InsertUrl")
	defer span.End()

	logger := zerolog.Ctx(c)

	logger.Info().Msg("generating uuid")
	gen, err := uuid.NewRandom()
	if err != nil {
		logger.Error().
			Err(err).
			Msgf("failed to generate uuid for url with error=%s", err.Error())
		return repository.Url{}, err
	}
	id := gen.String()
	logger.Info().Msgf("generated uuid=%s", gen.String())

	logger.Info().Msgf("starting to encode url=%s id=%s to shortUrl", url.String(), gen.String())
	encoded := s.encoder.EncodeToString([]byte(gen.String()))
	shortUrl := encoded[:5]
	logger.Info().Msgf("encoded url=%s id=%s to shortUrl=%s", url.String(), id, shortUrl)

	logger.Info().Msgf("inserting url=%s id=%s shortUrl=%s", url.String(), id, shortUrl)
	inserted, err := s.queries.InsertUrl(c, repository.InsertUrlParams{
		ID:       gen,
		Url:      url.String(),
		ShortUrl: shortUrl,
	})
	if err != nil {
		logger.Error().
			Err(err).
			Msgf("failed when inserting url=%s with id=%s to database with error=%s", id, url.String(), err.Error())
		return repository.Url{}, err
	}
	logger.Info().
		Str(log.KeyUrlID, inserted.ID.String()).
		Msgf("inserted url=%s id=%s shortUrl=%s", url.String(), id, shortUrl)

	return inserted, nil
}

func (s *UrlService) UpdateUrl(
	c context.Context,
	url url.URL,
	shortUrl string,
) (repository.Url, error) {
	c, span := tracer.Start(c, "UrlService UpdateUrl")
	defer span.End()

	logger := zerolog.Ctx(c)

	logger.Info().Msgf("finding shortUrl=%s", shortUrl)
	existing, err := s.queries.FindUrlByShortUrl(c, shortUrl)
	if err != nil {
		logger.Error().
			Err(err).
			Msgf("shortUrl=%s not found", shortUrl)
		return repository.Url{}, err
	}
	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(log.KeyOldUrl, existing.Url).
			Str(log.KeyUrlID, existing.ID.String())
	})
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

	return updated, nil
}

func (s *UrlService) DeleteUrl(
	c context.Context,
	shortUrl string,
) (repository.Url, error) {
	c, span := tracer.Start(c, "UrlService DeleteUrl")
	defer span.End()

	logger := zerolog.Ctx(c)

	logger.Info().Msgf("deleting shortUrl=%s", shortUrl)
	deleted, err := s.queries.DeleteUrlByShortUrl(c, shortUrl)
	if err != nil {
		logger.Error().
			Err(err).
			Msgf("failed deleting shortUrl=%s not found", shortUrl)
		return repository.Url{}, err
	}
	logger.Info().Msgf("deleted url=%s id=%s", deleted.Url, deleted.ID.String())

	return deleted, nil
}

func (s *UrlService) GetUrlByShortUrl(
	c context.Context,
	shortUrl string,
) (repository.Url, error) {
	c, span := tracer.Start(c, "UrlService GetUrlByShortUrl")
	defer span.End()

	logger := zerolog.Ctx(c)

	logger.Info().Msg("initializing transaction")
	tx, err := s.db.BeginTx(c, &sql.TxOptions{})
	if err != nil {
		logger.Info().
			Err(err).
			Msgf("failed initializing transaction with error=%s", err.Error())
		return repository.Url{}, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			logger.Info().
				Err(err).
				Msgf("failed rollback transaction with error=%s", err.Error())
			return
		}
		logger.Info().Msg("rollback transaction")
	}()
	logger.Info().Msg("initialized transaction")

	logger.Info().Msgf("finding shortUrl=%s", shortUrl)
	existing, err := s.queries.WithTx(tx).FindUrlByShortUrl(c, shortUrl)
	if err != nil {
		logger.Error().
			Err(err).
			Msgf("shortUrl=%s not found", shortUrl)
		return repository.Url{}, err
	}
	logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Any(log.KeyUrl, existing).
			Int32("newVisitedCount", existing.VisitedCount+1).
			Int32("oldVisitedCount", existing.VisitedCount)
	})
	logger.Info().Msgf("found shortUrl=%s", shortUrl)

	logger.Info().
		Msgf("updating VisitedCount=%d to VisitedCount=%d", existing.VisitedCount, existing.VisitedCount+1)
	updated, err := s.queries.WithTx(tx).UpdateVisitedCountUrl(
		c, repository.UpdateVisitedCountUrlParams{
			ID:           existing.ID,
			VisitedCount: existing.VisitedCount + 1,
		},
	)
	if err != nil {
		logger.Error().
			Err(err).
			Msgf("failed updating id=%s shortUrl=%s VisitedCount", existing.ID.String(), shortUrl)
		return repository.Url{}, err
	}
	logger.Info().
		Msgf("updated VisitedCount=%d to VisitedCount=%d", existing.VisitedCount, existing.VisitedCount+1)

	logger.Info().Msgf("committing transaction")
	if err := tx.Commit(); err != nil {
		logger.Error().
			Err(err).
			Msgf("failed committing transaction with error=%s", err.Error())
		return repository.Url{}, err
	}
	logger.Info().Msgf("committed transaction")

	return updated, nil
}

func (s *UrlService) GetUrlByShortUrlDetail(
	c context.Context,
	shortUrl string,
) (repository.Url, error) {
	c, span := tracer.Start(c, "UrlService GetUrlByShortUrlDetail")
	defer span.End()

	logger := zerolog.Ctx(c)

	logger.Info().Msgf("finding shortUrl=%s", shortUrl)
	existing, err := s.queries.FindUrlByShortUrl(c, shortUrl)
	if err != nil {
		logger.Error().
			Err(err).
			Msgf("shortUrl=%s not found", shortUrl)
		return repository.Url{}, err
	}
	logger.Info().
		Any("existing", existing).
		Msgf("finding shortUrl=%s", shortUrl)

	return existing, nil
}
