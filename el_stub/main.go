package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"math"
	"net/http"
)



func okHandler(w http.ResponseWriter, r *http.Request) {
	//xmlResponse := "<UpdateEntryResponseType><SystemInfo><From>PCVP</From><To>10.15.44.71</To><MessageId>eeaddb147b61591b8d91545ebbc4f5c8</MessageId></SystemInfo><Response><Response_Code>0</Response_Code><Response_Description>OK</Response_Description></Response></UpdateEntryResponseType>"
	xmlResponse := "<Envelope xmlns=\"http://schemas.xmlsoap.org/soap/envelope/\"><Body xmlns=\"http://schemas.xmlsoap.org/soap/envelope/\"><UpdateEntryResponse xmlns=\"http://xmlns.dit.mos.ru/sudir/itb/connector\"><SystemInfo><From>PCVP</From><To>10.15.44.71</To><MessageId>d4f3947b468adbe9835f71143e207a81</MessageId><SentDateTime>2024-12-16T11:54:02.647388117+03:00</SentDateTime></SystemInfo><Response><Response_Code>0</Response_Code><Response_Description>OK</Response_Description></Response><EntryItem><EntryName>Users</EntryName><Object><Name>HOUSE</Name><Attribute><Name>EpdId</Name><Value>B3DD365CFE2D47FD8B089F345C5443AA</Value></Attribute></Object></EntryItem></UpdateEntryResponse></Body></Envelope> "

	// Устанавливаем заголовок Content-Type для SOAP-ответа
	w.Header().Set("Content-Type", "application/xml")

	//// Кодируем ответ в XML и отправляем
	//if err := xml.NewEncoder(w).Encode(s); err != nil {
	//	http.Error(w, "Ошибка кодирования ответа", http.StatusInternalServerError)
	//}
	w.Write([]byte(xmlResponse))
}

func main() {
	http.HandleFunc("/ok", okHandler)

	fmt.Println("Сервер запущен на http://localhost:9998/ok")
	log.Fatal(http.ListenAndServe(":9998", nil))
}
