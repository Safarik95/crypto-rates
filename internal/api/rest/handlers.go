package rest

import (
	"crypto-rates/internal/service"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type Handlers struct {
	service *service.RateService
}

func NewHandlers(service *service.RateService) *Handlers {
	return &Handlers{
		service: service,
	}
}

// GetRates возвращает все курсы
func (h *Handlers) GetRates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	rates := h.service.GetAllRateInfo()
	if len(rates) == 0 {
		h.sendError(w, "Нет данных о курсах", http.StatusNotFound)
		return
	}

	response := RatesResponse{
		Rates: make([]RateResponse, 0, len(rates)),
	}

	for _, rate := range rates {
		response.Rates = append(response.Rates, ConvertRateInfo(rate))
	}

	h.sendJSON(w, response, http.StatusOK)

	zap.L().Debug("REST API: возвращены все курсы",
		zap.Int("количество", len(rates)),
	)
}

// GetRate возвращает курс конкретной валюты
func (h *Handlers) GetRate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем currency из URL пути /rates/{currency}
	path := strings.TrimPrefix(r.URL.Path, "/rates/")
	currency := strings.ToUpper(strings.TrimSpace(path))

	if currency == "" {
		h.sendError(w, "Необходимо указать валюту", http.StatusBadRequest)
		return
	}

	// Получаем данные для одной валюты
	rateInfo, err := h.service.GetRateInfo(currency)
	if err != nil {
		h.sendError(w, fmt.Sprintf("Валюта %s не найдена", currency), http.StatusNotFound)
		return
	}

	// Для одного курса возвращаем объект, а не массив
	response := ConvertRateInfo(rateInfo)
	h.sendJSON(w, response, http.StatusOK)

	zap.L().Debug("REST API: возвращен курс для валюты",
		zap.String("валюта", currency),
	)
}

// sendJSON отправляет JSON ответ
func (h *Handlers) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		zap.L().Error("Ошибка кодирования JSON ответа", zap.Error(err))
	}
}

// sendError отправляет ошибку
func (h *Handlers) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := ErrorResponse{Error: message}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		zap.L().Error("Ошибка кодирования ответа с ошибкой", zap.Error(err))
	}

	zap.L().Warn("Ошибка REST API",
		zap.String("ошибка", message),
		zap.Int("код_статуса", statusCode),
	)
}
