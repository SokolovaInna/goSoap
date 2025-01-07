package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"math"
	"net/http"
)

// Запрос SOAP для уравнения
type EquationRequest struct {
	XMLName xml.Name `xml:"request"`
	A       float64  `xml:"a"`
	B       float64  `xml:"b"`
	C       float64  `xml:"c"`
}

// Ответ SOAP для уравнения
type EquationResponse struct {
	XMLName xml.Name `xml:"response"`
	Formula string   `xml:"formula"`
	D       float64  `xml:"D"`
	X1      *float64 `xml:"x1,omitempty"`
	X2      *float64 `xml:"x2,omitempty"`
	Error   string   `xml:"error,omitempty"`
}

// Обработчик SOAP-запроса
func solveEquationHandler(w http.ResponseWriter, r *http.Request) {
	var req EquationRequest

	// Декодируем XML-запрос
	if err := xml.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка декодирования запроса", http.StatusBadRequest)
		return
	}

	// Вычисляем дискриминант
	D := req.B*req.B - 4*req.A*req.C
	formula := fmt.Sprintf("%gx^2 + %gx + %g = 0", req.A, req.B, req.C)

	// Создаем ответ
	resp := EquationResponse{
		Formula: formula,
		D:       D,
	}

	if D < 0 {
		resp.Error = "Дискриминант меньше нуля, уравнение не имеет действительных корней."
	} else {
		sqrtD := math.Sqrt(D)
		x1 := (-req.B + sqrtD) / (2 * req.A)
		resp.X1 = &x1

		if D > 0 {
			x2 := (-req.B - sqrtD) / (2 * req.A)
			resp.X2 = &x2
		}
	}

	// Устанавливаем заголовок Content-Type для SOAP-ответа
	w.Header().Set("Content-Type", "application/xml")

	// Кодируем ответ в XML и отправляем
	if err := xml.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Ошибка кодирования ответа", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/solve", solveEquationHandler)

	fmt.Println("Сервер запущен на http://localhost:8080/solve")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
