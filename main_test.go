package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	city := "moscow"
	total := len(cafeList[city])

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, total},
	}

	for _, r := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?count=%d&city=%s", r.count, city), nil)

		handler.ServeHTTP(response, req)

		answerStatus := response.Result()

		assert.Equal(t, http.StatusOK, answerStatus.StatusCode)

		body := response.Body.String()
		got := 0
		if len(strings.TrimSpace(body)) > 0 {
			got = len(strings.Split(body, ","))
		}

		assert.Equal(t, r.want, got)
	}

}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	city := "moscow"

	requests := []struct {
		search string
		want   int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, r := range requests {
		resp := httptest.NewRecorder()

		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=%s&search=%s", city, r.search), nil)

		handler.ServeHTTP(resp, req)

		answerStatus := resp.Result()

		assert.Equal(t, http.StatusOK, answerStatus.StatusCode)

		body := resp.Body.String()
		var cafes []string
		if body != "" {
			cafes = strings.Split(body, ",")
		}

		assert.Equal(t, r.want, len(cafes))

		lowerSearch := strings.ToLower(r.search)
		for _, cafe := range cafes {
			if !strings.Contains(strings.ToLower(cafe), lowerSearch) {
				t.Errorf("search=%s: string contains wrong cafe %s", r.search, cafe)
			}
		}
	}
}
