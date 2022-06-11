package client

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/keystop/yaDiploma/pkg/logger"
)

func MakeRequest(rtype, address, ctype, authorization string, b []byte) ([]byte, bool) {
	req, _ := http.NewRequest(rtype, address, bytes.NewReader(b))
	if len(ctype) != 0 {
		req.Header.Add("Content-Type", ctype)
	}

	if len(authorization) != 0 {
		req.Header.Add("Authorization", authorization)
	}

	tr := &http.Transport{
		MaxIdleConns:          10,
		IdleConnTimeout:       30 * time.Second,
		DisableCompression:    true,
		ResponseHeaderTimeout: 10 * time.Second,
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)

	if err != nil {
		logger.Info("Error", "Ошибка выполнения запроса", err)
		return nil, false
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		logger.Info(http.StatusTooManyRequests, "Приостановка запросов к серверу расчета начислений")
		time.Sleep(time.Second)
	}

	if resp.StatusCode == http.StatusInternalServerError {
		logger.Info(http.StatusInternalServerError, "Ошибка обработки запроса сервером")
		return nil, false
	}

	if resp.StatusCode < 200 || resp.StatusCode > 230 {
		logger.Info(resp.StatusCode, "Ответ сервера не 2**")
		return nil, false
	}
	rb, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Info("Error", "Ошибка чтения ответа на запрос", err)
		return nil, false
	}
	defer resp.Body.Close()
	return rb, true
}
