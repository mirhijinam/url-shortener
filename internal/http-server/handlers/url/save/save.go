package save

import (
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/mirhijinam/url-shortener/internal/lib/api/response"
	"github.com/mirhijinam/url-shortener/internal/lib/logger/sl"
	"github.com/mirhijinam/url-shortener/internal/lib/random"
	"github.com/mirhijinam/url-shortener/internal/storage"
	"golang.org/x/exp/slog"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: move to config
const aliasLen = 6

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = slog.With(
			slog.String("op", op),
			slog.String("handler_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		// check for an empty response body
		//   and handle the issue separately
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, resp.Error("empty request"))
			return
		}

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			errValidator := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(errValidator))
			return
		}

		alias := req.Alias
		if alias == "" {
			//TODO: check for duplicate creation
			alias, err = random.NewRandomAlias(aliasLen)
			if err != nil {
				log.Error("error ", sl.Err(err))
				render.JSON(w, r, sl.Err(err))
				return
			}
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Error("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}

		if err != nil {
			log.Error("failed to save url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to save url"))
			return
		}

		log.Info("url was saved", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
