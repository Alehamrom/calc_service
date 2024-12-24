package application

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCalcHandler проверяет корректность работы функции CalcHandler.
// Он тестирует различные сценарии, включая валидные и невалидные выражения,
// а также особые случаи, чтобы убедиться, что обработчик возвращает правильный ответ и статус код.
func TestCalcHandler(t *testing.T) {
	tests := []struct {
		name           string
		request        Request
		expected       string
		expectedStatus int
	}{
		{
			name:           "Valid Expression",
			request:        Request{Expression: "2+2*2"},
			expected:       "{\n  \"result\": 6.000000\n}",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Expression",
			request:        Request{Expression: "invalid"},
			expected:       "{\n  \"error\": \"invalid expression\"\n}",
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Преобразуем запрос в JSON
			requestBody, err := json.Marshal(tc.request)
			if err != nil {
				t.Fatalf("Ошибка при сериализации запроса: %v", err)
			}

			// Создаём новый HTTP-запрос с заданным телом
			req, err := http.NewRequest("GET", "/api/v1/calculate", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatalf("Ошибка при создании нового запроса: %v", err)
			}

			// Создаём ResponseRecorder для получения ответа
			w := httptest.NewRecorder()

			// Вызываем обрабатывающую функцию напрямую
			CalcHandler(w, req)

			// Получаем результат из ResponseRecorder
			resp := w.Result()
			defer resp.Body.Close()

			// Проверяем статус код ответа
			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Ожидаемый статус код %d, получен %d", tc.expectedStatus, resp.StatusCode)
			}

			// Читаем тело ответа
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Ошибка при чтении тела ответа: %v", err)
			}

			// Сравниваем полученный ответ с ожидаемым
			if string(body) != tc.expected {
				t.Errorf("Ожидаемый ответ:\n%s\nПолученный ответ:\n%s", tc.expected, string(body))
			}
		})
	}
}
