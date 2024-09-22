package remove

import (
	"errors"
	"fmt"
	resp "github.com/evgenii-gsv/url-shortener/internal/lib/api/response"
	"github.com/evgenii-gsv/url-shortener/internal/lib/logger/sl"
	"github.com/evgenii-gsv/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Response struct {
	resp.Response
	DeletedAlias string `json:"deleted_alias,omitempty"`
}

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.remove.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))

			render.JSON(w, r, resp.Error(fmt.Sprintf("No url with alias %s found", alias)))

			return
		}
		if err != nil {
			log.Error("failed to remove url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("deleted url", slog.String("alias", alias))

		render.JSON(w, r, Response{
			Response:     resp.OK(),
			DeletedAlias: alias,
		})
	}
}
