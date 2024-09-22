package save

import (
	"errors"
	"github.com/evgenii-gsv/url-shortener/internal/config"
	resp "github.com/evgenii-gsv/url-shortener/internal/lib/api/response"
	"github.com/evgenii-gsv/url-shortener/internal/lib/logger/sl"
	"github.com/evgenii-gsv/url-shortener/internal/lib/random"
	"github.com/evgenii-gsv/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.46.0 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave, alias string) error
}

func New(log *slog.Logger, urlSaver URLSaver, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// logging
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// decoding
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			msg := "failed to decode request body"
			log.Error(msg, sl.Err(err))

			render.JSON(w, r, resp.Error(msg))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// validating
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		// TODO: handle existing random alias
		if alias == "" {
			alias = random.NewRandomString(cfg.AliasLength)
		}

		err = urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			msg := "url already exists"
			log.Info(msg, slog.String("url", req.URL))

			render.JSON(w, r, resp.Error(msg))

			return
		}
		if err != nil {
			msg := "failed to save url"
			log.Error(msg, sl.Err(err))

			render.JSON(w, r, resp.Error(msg))

			return
		}

		log.Info("url added", slog.String("alias", alias))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
