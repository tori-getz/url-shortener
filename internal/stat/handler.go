package stat

import (
	"net/http"
	"slices"
	"time"
	"url-shortener/configs"
	"url-shortener/pkg/middleware"
	"url-shortener/pkg/res"
)

type StatHandlerDeps struct {
	*StatRepository
	*configs.Config
}

type StatHandler struct {
	*StatRepository
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{
		StatRepository: deps.StatRepository,
	}

	router.Handle("GET /stat", middleware.IsAuthed(handler.GetStat(), deps.Config))
}

// @Summary Получить статистику
// @Description Возвращает статистику за указанный период, сгруппированную по дням или месяцам.
// @Tags Stats
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param from query string true "Начальная дата в формате YYYY-MM-DD"
// @Param to query string true "Конечная дата в формате YYYY-MM-DD"
// @Param by query string true "Группировка: day (по дням) или month (по месяцам)"
// @Success 200 {object} map[string]interface{} "Статистика"
// @Failure 400 {object} res.ErrorResponse "Неверные параметры запроса"
// @Failure 500 {object} res.ErrorResponse "Внутренняя ошибка сервера"
// @Router /stat [get]
func (handler *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from, err := time.Parse("2006-01-02", r.URL.Query().Get("from"))
		if err != nil {
			res.Error(w, ErrInvalidFrom, http.StatusBadRequest)
			return
		}

		to, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
		if err != nil {
			res.Error(w, ErrInvalidTo, http.StatusBadRequest)
			return
		}

		by := r.URL.Query().Get("by")
		if !slices.Contains([]string{GroupByMonth, GroupByDay}, by) {
			res.Error(w, ErrInvalidBy, http.StatusBadRequest)
			return
		}

		stats := handler.StatRepository.GetStats(by, from, to)

		res.Json(w, stats, http.StatusOK)
	}
}
