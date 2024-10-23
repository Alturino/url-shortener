package service

import (
	"context"
	"database/sql"
	"encoding/base64"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/Alturino/url-shortener/internal/log"
	"github.com/Alturino/url-shortener/internal/repository"
)

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
	startTime := time.Now()

	logger := zerolog.Ctx(c)
	hashcode := c.Value(log.KeyHashcode).(string)

	logger.Info().
		Str(log.KeyProcess, "InsertUrl").
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str(log.KeyHashcode, hashcode).
		Any(log.KeyRequestBody, url).
		Time(log.KeyStartTime, startTime).
		Any(log.KeyRequestBody, url.String()).
		Msg("generating uuid")
	gen, err := uuid.NewRandom()
	if err != nil {
		logger.Error().
			Err(err).
			Str(log.KeyProcess, "InsertUrl").
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Str(log.KeyHashcode, hashcode).
			Time(log.KeyStartTime, startTime).
			Msgf("failed to generate uuid for url", err.Error())
		return repository.Url{}, err
	}
	id := gen.String()
	logger.Info().
		Str(log.KeyProcess, "InsertUrl").
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str(log.KeyHashcode, hashcode).
		Any(log.KeyRequestBody, url).
		Time(log.KeyStartTime, startTime).
		Any(log.KeyRequestBody, url.String()).
		Msgf("generated uuid=%s", gen.String())

	logger.Info().
		Str(log.KeyProcess, "InsertUrl").
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str(log.KeyHashcode, hashcode).
		Any(log.KeyRequestBody, url).
		Time(log.KeyStartTime, startTime).
		Msgf("starting to encode url=%s id=%d to shortUrl", url.String(), gen.String())
	encoded := s.encoder.EncodeToString([]byte(gen.String()))
	shortUrl := encoded[:5]
	logger.Info().
		Str(log.KeyProcess, "InsertUrl").
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str(log.KeyHashcode, hashcode).
		Any(log.KeyRequestBody, url).
		Time(log.KeyStartTime, startTime).
		Msgf("encoded url=%s id=%d to shortUrl=%s", url.String(), id, shortUrl)

	logger.Info().
		Str(log.KeyProcess, "InsertUrl").
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str(log.KeyHashcode, hashcode).
		Any(log.KeyRequestBody, url).
		Int(log.KeyID, int(gen.ID())).
		Time(log.KeyStartTime, startTime).
		Msgf("inserting url=%s id=%d shortUrl=%s", url.String(), id, shortUrl)
	res, err := s.queries.InsertUrl(c, repository.InsertUrlParams{
		ID:       gen,
		Url:      url.String(),
		ShortUrl: shortUrl,
	})
	if err != nil {
		logger.Error().
			Err(err).
			Str(log.KeyProcess, "InsertUrl").
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Str(log.KeyHashcode, hashcode).
			Int(log.KeyID, int(gen.ID())).
			Time(log.KeyStartTime, startTime).
			Msgf("failed when inserting url=%s with id=%d to database with error=%s", id, url.String(), err.Error())
		return repository.Url{}, err
	}
	logger.Info().
		Str(log.KeyProcess, "InsertUrl").
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str(log.KeyHashcode, hashcode).
		Any(log.KeyRequestBody, url).
		Int(log.KeyID, int(gen.ID())).
		Time(log.KeyStartTime, startTime).
		Msgf("inserted url=%s id=%d shortUrl=%s", url.String(), id, shortUrl)

	return res, nil
}

func (s *UrlService) UpdateUrl(
	c context.Context,
	url url.URL,
	shortUrl string,
) (repository.Url, error) {
	startTime := time.Now()

	logger := zerolog.Ctx(c)
	hashcode := c.Value(log.KeyHashcode).(string)

	logger.Info().
		Any(log.KeyHashcode, hashcode).
		Any(log.KeyRequestBody, url).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str("new_url", url.String()).
		Str(log.KeyProcess, "UpdateUrl").
		Time(log.KeyStartTime, startTime).
		Msgf("finding shortUrl=%s", shortUrl)
	existing, err := s.queries.FindUrlByShortUrl(c, shortUrl)
	if err != nil {
		logger.Error().
			Err(err).
			Any(log.KeyHashcode, hashcode).
			Any(log.KeyRequestBody, url).
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Str("new_url", url.String()).
			Str(log.KeyProcess, "UpdateUrl").
			Time(log.KeyStartTime, startTime).
			Msgf("shortUrl=%s not found", shortUrl)
		return repository.Url{}, err
	}
	logger.Info().
		Any(log.KeyHashcode, hashcode).
		Any(log.KeyRequestBody, url).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str("new_url", url.String()).
		Str("old_url", existing.Url).
		Str(log.KeyID, existing.ID.String()).
		Str(log.KeyProcess, "UpdateUrl").
		Time(log.KeyStartTime, startTime).
		Msgf("found shortUrl=%s", shortUrl)

	logger.Info().
		Any(log.KeyHashcode, hashcode).
		Any(log.KeyRequestBody, url).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str("new_url", url.String()).
		Str("old_url", existing.Url).
		Str(log.KeyID, existing.ID.String()).
		Str(log.KeyProcess, "UpdateUrl").
		Time(log.KeyStartTime, startTime).
		Msgf("updating url=%s id=%d to url=%s", existing.Url, existing.ID.String(), url.String())
	updated, err := s.queries.UpdateUrl(
		c,
		repository.UpdateUrlParams{ShortUrl: shortUrl, Url: url.String()},
	)
	if err != nil {
		logger.Error().
			Err(err).
			Any(log.KeyHashcode, hashcode).
			Any(log.KeyRequestBody, url).
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Str("new_url", url.String()).
			Str("old_url", existing.Url).
			Str(log.KeyID, existing.ID.String()).
			Str(log.KeyProcess, "UpdateUrl").
			Time(log.KeyStartTime, startTime).
			Msgf("failed updating url=%s id=%s with error=%s", existing.Url, existing.ID.String(), err.Error())
		return repository.Url{}, err
	}
	logger.Info().
		Any(log.KeyHashcode, hashcode).
		Any(log.KeyRequestBody, url).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str("new_url", url.String()).
		Str("old_url", existing.Url).
		Str(log.KeyID, existing.ID.String()).
		Str(log.KeyProcess, "UpdateUrl").
		Time(log.KeyStartTime, startTime).
		Msgf("updated url=%s id=%d to url=%s", existing.Url, existing.ID.String(), url.String())

	return updated, nil
}

func (s *UrlService) DeleteUrl(
	c context.Context,
	shortUrl string,
) (repository.Url, error) {
	startTime := time.Now()

	logger := zerolog.Ctx(c)
	hashcode := c.Value(log.KeyHashcode).(string)

	logger.Info().
		Any(log.KeyHashcode, hashcode).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str(log.KeyProcess, "DeleteUrl").
		Str(log.KeyShortUrl, shortUrl).
		Time(log.KeyStartTime, startTime).
		Msgf("deleting shortUrl=%s", shortUrl)
	deleted, err := s.queries.DeleteUrlByShortUrl(c, shortUrl)
	if err != nil {
		logger.Error().
			Err(err).
			Any(log.KeyHashcode, hashcode).
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Str(log.KeyProcess, "DeleteUrl").
			Time(log.KeyStartTime, startTime).
			Msgf("failed deleting shortUrl=%d not found", shortUrl)
		return repository.Url{}, err
	}
	logger.Info().
		Any(log.KeyHashcode, hashcode).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str(log.KeyProcess, "DeleteUrl").
		Str(log.KeyShortUrl, shortUrl).
		Time(log.KeyStartTime, startTime).
		Msgf("deleted url=%s id=%d", deleted.Url, deleted.ID.String())

	return deleted, nil
}

func (s *UrlService) GetUrlByShortUrl(
	c context.Context,
	shortUrl string,
) (repository.Url, error) {
	startTime := time.Now()

	logger := zerolog.Ctx(c)
	hashcode := c.Value(log.KeyHashcode).(string)

	logger.Info().
		Any(log.KeyHashcode, hashcode).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str(log.KeyProcess, "GetUrlByShortUrl").
		Time(log.KeyStartTime, startTime).
		Msgf("initializing transaction", shortUrl)
	tx, err := s.db.BeginTx(c, &sql.TxOptions{})
	if err != nil {
		logger.Info().
			Err(err).
			Any(log.KeyHashcode, hashcode).
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Str(log.KeyProcess, "GetUrlByShortUrl").
			Time(log.KeyStartTime, startTime).
			Msgf("failed initializing transaction with error=%s", err.Error())
		return repository.Url{}, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			logger.Info().
				Err(err).
				Any(log.KeyHashcode, hashcode).
				Dur(log.KeyProcessingTime, time.Since(startTime)).
				Str(log.KeyProcess, "GetUrlByShortUrl").
				Time(log.KeyStartTime, startTime).
				Msgf("failed rollback transaction with error=%s", err.Error())
			return
		}
		logger.Info().
			Any(log.KeyHashcode, hashcode).
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Str(log.KeyProcess, "GetUrlByShortUrl").
			Time(log.KeyStartTime, startTime).
			Msgf("rollback transaction")
	}()
	logger.Info().
		Any(log.KeyHashcode, hashcode).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str(log.KeyProcess, "GetUrlByShortUrl").
		Time(log.KeyStartTime, startTime).
		Msgf("initialized transaction", shortUrl)

	logger.Info().
		Any(log.KeyHashcode, hashcode).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str(log.KeyProcess, "GetUrlByShortUrl").
		Str(log.KeyShortUrl, shortUrl).
		Time(log.KeyStartTime, startTime).
		Msgf("finding shortUrl=%s", shortUrl)
	existing, err := s.queries.WithTx(tx).FindUrlByShortUrl(c, shortUrl)
	if err != nil {
		logger.Error().
			Err(err).
			Any(log.KeyHashcode, hashcode).
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Str(log.KeyProcess, "GetUrlByShortUrl").
			Str(log.KeyShortUrl, shortUrl).
			Time(log.KeyStartTime, startTime).
			Msgf("shortUrl=%s not found", shortUrl)
		return repository.Url{}, err
	}
	logger.Info().
		Any(log.KeyHashcode, hashcode).
		Any("url", existing).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str(log.KeyProcess, "GetUrlByShortUrl").
		Str(log.KeyShortUrl, shortUrl).
		Time(log.KeyStartTime, startTime).
		Msgf("found shortUrl=%s", shortUrl)

	logger.Info().
		Any("url", existing).
		Any(log.KeyHashcode, hashcode).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Int32("NewVisitedCount", existing.VisitedCount+1).
		Int32("OldVisitedCount", existing.VisitedCount).
		Str(log.KeyProcess, "GetUrlByShortUrl").
		Str(log.KeyShortUrl, shortUrl).
		Time(log.KeyStartTime, startTime).
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
			Any("url", existing).
			Any(log.KeyHashcode, hashcode).
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Int32("NewVisitedCount", existing.VisitedCount+1).
			Int32("OldVisitedCount", existing.VisitedCount).
			Str(log.KeyProcess, "GetUrlByShortUrl").
			Str(log.KeyShortUrl, shortUrl).
			Time(log.KeyStartTime, startTime).
			Msgf("failed updating id=%s shortUrl=%s VisitedCount", existing.ID.String(), shortUrl)
		return repository.Url{}, err
	}
	logger.Info().
		Any("url", existing).
		Any(log.KeyHashcode, hashcode).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Int32("NewVisitedCount", existing.VisitedCount+1).
		Int32("OldVisitedCount", existing.VisitedCount).
		Str(log.KeyProcess, "GetUrlByShortUrl").
		Str(log.KeyShortUrl, shortUrl).
		Time(log.KeyStartTime, startTime).
		Msgf("updated VisitedCount=%d to VisitedCount=%d", existing.VisitedCount, existing.VisitedCount+1)

	logger.Info().
		Any("url", existing).
		Any(log.KeyHashcode, hashcode).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Int32("NewVisitedCount", existing.VisitedCount+1).
		Int32("OldVisitedCount", existing.VisitedCount).
		Str(log.KeyProcess, "GetUrlByShortUrl").
		Str(log.KeyShortUrl, shortUrl).
		Time(log.KeyStartTime, startTime).
		Msgf("committing transaction")
	if err := tx.Commit(); err != nil {
		logger.Error().
			Err(err).
			Any("url", existing).
			Any(log.KeyHashcode, hashcode).
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Int32("NewVisitedCount", existing.VisitedCount+1).
			Int32("OldVisitedCount", existing.VisitedCount).
			Str(log.KeyProcess, "GetUrlByShortUrl").
			Str(log.KeyShortUrl, shortUrl).
			Time(log.KeyStartTime, startTime).
			Msgf("failed committing transaction with error=%s", err.Error())
		return repository.Url{}, err
	}
	logger.Info().
		Any("url", existing).
		Any(log.KeyHashcode, hashcode).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Int32("NewVisitedCount", existing.VisitedCount+1).
		Int32("OldVisitedCount", existing.VisitedCount).
		Str(log.KeyProcess, "GetUrlByShortUrl").
		Str(log.KeyShortUrl, shortUrl).
		Time(log.KeyStartTime, startTime).
		Msgf("committed transaction")

	return updated, nil
}

func (s *UrlService) GetUrlByShortUrlDetail(
	c context.Context,
	shortUrl string,
) (repository.Url, error) {
	startTime := time.Now()

	logger := zerolog.Ctx(c)
	hashcode := c.Value(log.KeyHashcode).(string)

	logger.Info().
		Any(log.KeyHashcode, hashcode).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str(log.KeyProcess, "GetUrlByShortUrl").
		Str(log.KeyShortUrl, shortUrl).
		Time(log.KeyStartTime, startTime).
		Msgf("finding shortUrl=%s", shortUrl)
	existing, err := s.queries.FindUrlByShortUrl(c, shortUrl)
	if err != nil {
		logger.Error().
			Err(err).
			Any(log.KeyHashcode, hashcode).
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Str(log.KeyProcess, "GetUrlByShortUrl").
			Str(log.KeyShortUrl, shortUrl).
			Time(log.KeyStartTime, startTime).
			Msgf("shortUrl=%s not found", shortUrl)
		return repository.Url{}, err
	}
	logger.Info().
		Any(log.KeyHashcode, hashcode).
		Any("existing", existing).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Str(log.KeyProcess, "GetUrlByShortUrl").
		Str(log.KeyShortUrl, shortUrl).
		Time(log.KeyStartTime, startTime).
		Msgf("finding shortUrl=%s", shortUrl)

	return existing, nil
}
