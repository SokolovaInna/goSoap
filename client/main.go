package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type EquationRequest struct {
	XMLName xml.Name `xml:"request"`
	A       float64  `xml:"a"`
	B       float64  `xml:"b"`
	C       float64  `xml:"c"`
}

type EquationResponse struct {
	XMLName xml.Name `xml:"response"`
	Formula string   `xml:"formula"`
	D       float64  `xml:"D"`
	X1      *float64 `xml:"x1,omitempty"`
	X2      *float64 `xml:"x2,omitempty"`
	Error   string   `xml:"error,omitempty"`
}
type EquationResponseJson struct {
	Formula string   `json:"formula"`
	D       float64  `json:"D"`
	X1      *float64 `json:"x1,omitempty"`
	X2      *float64 `json:"x2,omitempty"`
	Error   string   `json:"error,omitempty"`
}
type Response struct {
	Message string `json:"message"`
}

// Функция для отправки SOAP-запроса и получения ответа
func solveEquation(a, b, c float64) (*EquationResponse, error) {

	req := EquationRequest{A: a, B: b, C: c}
	reqXML, err := xml.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при маршалинге XML: %v", err)
	}

	// Оборачиваем XML в SOAP-конверт
	//soapEnvelope := fmt.Sprintf(`
	//    <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
	//        <soapenv:Body>
	//            %s
	//        </soapenv:Body>
	//    </soapenv:Envelope>`, reqXML)
	soapEnvelope := fmt.Sprintf(`%s`, reqXML)

	// Отправляем запрос на сервер
	resp, err := http.Post("http://localhost:8080/solve", "application/xml", bytes.NewBufferString(soapEnvelope))
	if err != nil {
		return nil, fmt.Errorf("Ошибка отправки запроса: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Ошибка чтения ответа: %v", err)
	}
	fmt.Println("Ответ:", string(body))

	// Парсим ответ XML
	//var soapResp struct {
	//	Body struct {
	//		EquationResponse EquationResponse `xml:"response"`
	//	} `xml:"Body"`
	//}
	//if err := xml.Unmarshal(body, &soapResp); err != nil {
	//	return nil, fmt.Errorf("Ошибка парсинга XML: %v", err)
	//}
	var resp1 EquationResponse
	if err := xml.Unmarshal(body, &resp1); err != nil {
		return nil, fmt.Errorf("Ошибка парсинга XML: %v", err)
	}
	//return &soapResp.Body.EquationResponse, nil
	return &resp1, nil
}
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	a := getfloat(r.URL.Query().Get("a"))
	b := getfloat(r.URL.Query().Get("b"))
	c := getfloat(r.URL.Query().Get("c"))
	resp, err := solveEquation(a, b, c)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json") // Устанавливаем заголовок Content-Type
	response := EquationResponseJson{Formula: resp.Formula,
		D:  resp.D,
		X1: resp.X1, X2: resp.X2,
		Error: resp.Error}
	json.NewEncoder(w).Encode(response) // Кодируем структуру в JSON и отправляем в ответ
}

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		a := getfloat(r.URL.Query().Get("a"))
		b := getfloat(r.URL.Query().Get("b"))
		c := getfloat(r.URL.Query().Get("c"))
		resp, err := solveEquation(a, b, c)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			return
		}

		// Выводим ответ
		if resp.Error != "" {
			fmt.Printf("Ошибка: %s\n", resp.Error)
		} else {
			fmt.Printf("Уравнение: %s\n", resp.Formula)
			fmt.Printf("Дискриминант: %g\n", resp.D)
			fmt.Printf("Корень x1: %g\n", *resp.X1)
			if resp.X2 != nil {
				fmt.Printf("Корень x2: %g\n", *resp.X2)
			}
		}
	})
	http.HandleFunc("/json", jsonHandler)
	http.ListenAndServe(":8081", nil)
}

func getfloat(bStr string) float64 {
	b, _ := strconv.ParseFloat(bStr, 64)

	return b
}
