package add

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"mood_tracker/internal/lib/api/response"
	"mood_tracker/internal/storage"
	"net/http"
)

type Request struct {
	Score  int8  `json:"score" validate:"required,min=0,max=10"`
	UserId int64 `json:"user_id" validate:"required"`
}

type Response struct {
	response.Response
}

type ScoreAdder interface {
	AddMoodScore(dto storage.AddMoodScoreDto) (int64, error)
}

func New(log *slog.Logger, scoreAdder ScoreAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.score.add.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, response.Error("request is empty"))
			return
		}
		if err != nil {
			log.Error("failed to parse request", "error", err)
			render.JSON(w, r, response.Error("failed to parse request"))
			return
		}

		log.Info("request body decoded", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", "error", err)
			render.JSON(w, r, response.Error(validateErr.Error()))
			return
		}

		id, err := scoreAdder.AddMoodScore(storage.AddMoodScoreDto{
			Score:  req.Score,
			UserId: req.UserId,
		})
		if err != nil {
			log.Error("failed to add mood score", "error", err)
			render.JSON(w, r, response.Error("failed to add mood score"))
			return
		}
		log.Info("mood score added", slog.Int64("id", id))
		render.JSON(w, r, Response{
			Response: response.Ok(),
		})
	}
}
