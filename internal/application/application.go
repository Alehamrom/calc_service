package application

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Alehamrom/calc_service/pkg/calculation"
)

type Application struct {
}

func New() *Application {
	return &Application{}
}

// Запускает приложение в режиме командной строки.
// Считывает ввод пользователя и, после нажатия ENTER, выводит результат вычисления.
// Завершает работу при вводе команды 'exit'.
func (a *Application) Run() error {
	for {
		// Считываем выражение для вычисления из стандартного ввода.
		log.Println("Введите выражение (или 'exit' для выхода):")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Ошибка чтения выражения из консоли:", err)
			continue
		}
		// Удаляем лишние пробелы и символы переноса строки.
		text = strings.TrimSpace(text)
		// Завершаем работу, если введена команда 'exit'.
		if strings.ToLower(text) == "exit" {
			log.Println("Приложение успешно завершено.")
			return nil
		}
		// Выполняем вычисление введённого выражения.
		result, err := calculation.Calc(text)
		if err != nil {
			log.Println("Ошибка вычисления выражения:", err)
		} else {
			log.Printf("%s = %f\n", text, result)
		}
	}
}

type Request struct {
	Expression string `json:"expression"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\n  \"error\": \"%s\"\n}", err)
		return
	}
	result, err := calculation.Calc(request.Expression)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "{\n  \"error\": \"%s\"\n}", err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\n  \"result\": %f\n}", result)
	}
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	log.Println("Запуск сервера на порту :8080")
	return http.ListenAndServe(":8080", nil)
}
