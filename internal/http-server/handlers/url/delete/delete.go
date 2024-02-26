package delete

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/mirhijinam/url-shortener/internal/lib/api/response"
	"github.com/mirhijinam/url-shortener/internal/lib/logger/sl"
	"github.com/mirhijinam/url-shortener/internal/storage"
	"golang.org/x/exp/slog"
)

type Response struct {
	resp.Response
}

// TODO: move to config
const aliasLen = 6

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = slog.With(
			slog.String("op", op),
			slog.String("handler_id", middleware.GetReqID(r.Context())),
		)

		// get parameter alias from the router.Get pattern
		// note that here in handler we are addicted to the router
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		if alias == "" {
			log.Error("empty delete request")
			render.JSON(w, r, resp.Error("empty delete request"))
			return
		}

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Error("alias does not exist")
			render.JSON(w, r, resp.Error("alias does not exists"))
			return
		}

		if err != nil {
			log.Error("failed to delete alias", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to delete alias"))
			return
		}

		log.Info("alias was deleted", slog.String("alias", alias))

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
