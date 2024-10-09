package get

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"mood_tracker/internal/lib/api/response"
	"mood_tracker/internal/lib/api/validate"
	"mood_tracker/internal/storage"
	"net/http"
)

type Request struct {
	UserId int64 `json:"user_id" validate:"required"`
}

type Response struct {
	response.Response
	Scores []storage.MoodScore `json:"scores"`
}

type ScoreGetter interface {
	GetMoodScoresByUserId(id int64) (scores []storage.MoodScore, err error)
}

func New(log *slog.Logger, scoreGetter ScoreGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.score.get.New"
		log.With(
			slog.String("op", op),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("empty request body")
			render.JSON(w, r, response.Error("empty request body"))
			return
		}
		if err != nil {
			log.Error("error while decode request body", "error", err)
			render.JSON(w, r, response.Error("error while decode request body"))
			return
		}
		log.Info("decode body completed", "req", req)

		msg, err := validate.Struct(&req)
		if err != nil {
			log.Error(msg, "error", err)
			render.JSON(w, r, response.Error(msg))
			return
		}

		scores, err := scoreGetter.GetMoodScoresByUserId(req.UserId)
		if err != nil {
			log.Error("error while getting mood scores", "error", err)
			render.JSON(w, r, response.Error("error while getting mood scores"))
			return
		}
		log.Info("mood scores", "count", len(scores))
		render.JSON(w, r, Response{
			Response: response.Ok(),
			Scores:   scores,
		})
	}
}
