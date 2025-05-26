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
		if answerStatus.StatusCode != http.StatusOK {
			t.Errorf("count=%d: expected status 200, got %d", r.count, answerStatus.StatusCode)
			continue
		}

		body := response.Body.String()
		got := 0
		if len(strings.TrimSpace(body)) > 0 {
			got = len(strings.Split(body, ","))
		}

		if got != r.want {
			t.Errorf("count=%d: expected %d cafes, got %d", r.count, r.want, got)
		}
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
		if answerStatus.StatusCode != http.StatusOK {
			t.Errorf("search=%s: expected 200, got %d", r.search, answerStatus.StatusCode)
			continue
		}

		body := resp.Body.String()
		var cafes []string
		if body != "" {
			cafes = strings.Split(body, ",")
		}

		if len(cafes) != r.want {
			t.Errorf("search=%s: expected %d, got %d", r.search, r.want, len(cafes))
			continue
		}

		lowerSearch := strings.ToLower(r.search)
		for _, cafe := range cafes {
			if !strings.Contains(strings.ToLower(cafe), lowerSearch) {
				t.Errorf("search=%s: string contains wrong cafe %s", r.search, cafe)
			}
		}
	}
}
