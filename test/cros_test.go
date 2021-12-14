package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ahmetb/go-linq/v3"
	"github.com/turboyang-cn/middleware"
)

func TestCros(t *testing.T) {
	endpoint := "/test"
	origin := "https://localhost"
	headers := []string{"X-App-Id", "Authorization", "Content-Type"}

	router := middleware.NewRouter()
	router.UseCros()

	handler := router.Add(endpoint, testHandler)

	testMethodT(handler, origin, endpoint, http.MethodGet, strings.Join(headers, ", "), t)
	testMethodT(handler, origin, endpoint, http.MethodPost, strings.Join(headers, ", "), t)
	testMethodT(handler, origin, endpoint, http.MethodPut, strings.Join(headers, ", "), t)
	testMethodT(handler, origin, endpoint, http.MethodDelete, strings.Join(headers, ", "), t)
	testMethodT(handler, origin, endpoint, http.MethodHead, strings.Join(headers, ", "), t)
	testMethodT(handler, origin, endpoint, http.MethodPatch, strings.Join(headers, ", "), t)
	testMethodT(handler, origin, endpoint, http.MethodConnect, strings.Join(headers, ", "), t)
	testMethodT(handler, origin, endpoint, http.MethodTrace, strings.Join(headers, ", "), t)
	testMethodOptions(handler, origin, endpoint, http.MethodPost, strings.Join(headers, ", "), t)
}

func testMethodT(handler http.Handler, origin string, endpoint string, method string, headers string, t *testing.T) {
	isSuccess := true

	req := httptest.NewRequest(method, endpoint, bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	req.Header.Set("Origin", origin)
	req.Header.Set("Access-Control-Request-Method", method)
	req.Header.Set("Access-Control-Request-Headers", headers)

	handler.ServeHTTP(w, req)

	accessControlAllowOrigin, ok := w.Result().Header["Access-Control-Allow-Origin"]
	if !ok {
		t.Errorf("\033[41mERROR\033[0m Not found header [Access-Control-Allow-Origin]")
		isSuccess = false
	} else {
		if !linq.From(accessControlAllowOrigin).AnyWith(func(x interface{}) bool {
			return strings.EqualFold(x.(string), "*") || strings.EqualFold(x.(string), origin)
		}) {
			t.Errorf("\033[41mERROR\033[0m Header [Access-Control-Allow-Origin = %v] is not correct", strings.Join(accessControlAllowOrigin, ", "))
			isSuccess = false
		}
	}

	accessControlAllowMethods, ok := w.Result().Header["Access-Control-Allow-Methods"]
	if !ok {
		t.Errorf("\033[41mERROR\033[0m Not found header [Access-Control-Allow-Methods]")
		isSuccess = false
	} else {
		if !linq.From(strings.Split(accessControlAllowMethods[0], ",")).Select(func(x interface{}) interface{} {
			return strings.TrimSpace(x.(string))
		}).AnyWith(func(x interface{}) bool {
			return strings.EqualFold(x.(string), "*") || strings.EqualFold(x.(string), method)
		}) {
			t.Errorf("\033[41mERROR\033[0m Header [Access-Control-Allow-Methods = %v] is not correct", strings.Join(accessControlAllowMethods, ", "))
			isSuccess = false
		}
	}

	accessControlAllowHeaders, ok := w.Result().Header["Access-Control-Allow-Headers"]
	if !ok {
		t.Errorf("\033[41mERROR\033[0m Not found header [Access-Control-Allow-Headers]")
		isSuccess = false
	} else {
		if !linq.From(strings.Split(headers, ",")).Select(func(x interface{}) interface{} {
			return strings.TrimSpace(x.(string))
		}).All(func(x interface{}) bool {
			return accessControlAllowHeaders[0] == "*" || linq.From(strings.Split(accessControlAllowHeaders[0], ",")).Select(func(x interface{}) interface{} {
				return strings.TrimSpace(x.(string))
			}).AnyWith(func(y interface{}) bool {
				return strings.EqualFold(x.(string), y.(string))
			})
		}) {
			t.Errorf("\033[41mERROR\033[0m Header [Access-Control-Allow-Headers = %v] is not correct", strings.Join(accessControlAllowHeaders, ", "))
			isSuccess = false
		}
	}

	accessControlExposeHeaders, ok := w.Result().Header["Access-Control-Expose-Headers"]
	if !ok {
		t.Errorf("\033[41mERROR\033[0m Not found header [Access-Control-Expose-Headers]")
		isSuccess = false
	}
	if !linq.From(strings.Split(headers, ",")).Select(func(x interface{}) interface{} {
		return strings.TrimSpace(x.(string))
	}).All(func(x interface{}) bool {
		return accessControlExposeHeaders[0] == "*" || linq.From(strings.Split(accessControlExposeHeaders[0], ",")).Select(func(x interface{}) interface{} {
			return strings.TrimSpace(x.(string))
		}).AnyWith(func(y interface{}) bool {
			return strings.EqualFold(x.(string), y.(string))
		})
	}) {
		t.Errorf("\033[41mERROR\033[0m Header [Access-Control-Expose-Headers = %v] is not correct", strings.Join(accessControlExposeHeaders, ", "))
		isSuccess = false
	}

	if isSuccess {
		t.Logf("\033[42mPASS\033[0m Method: %v", method)
	} else {
		t.Logf("\033[41mERROR\033[0m Method: %v", method)
	}
}

func testMethodOptions(handler http.Handler, origin string, endpoint string, method string, headers string, t *testing.T) {
	isSuccess := true

	req := httptest.NewRequest(http.MethodOptions, endpoint, bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	req.Header.Set("Origin", origin)
	req.Header.Set("Access-Control-Request-Method", method)
	req.Header.Set("Access-Control-Request-Headers", headers)

	handler.ServeHTTP(w, req)

	accessControlAllowOrigin, ok := w.Result().Header["Access-Control-Allow-Origin"]
	if !ok {
		t.Errorf("\033[41mERROR\033[0m Not found header [Access-Control-Allow-Origin]")
		isSuccess = false
	} else {
		if !linq.From(accessControlAllowOrigin).AnyWith(func(x interface{}) bool {
			return strings.EqualFold(x.(string), "*") || strings.EqualFold(x.(string), origin)
		}) {
			t.Errorf("\033[41mERROR\033[0m Header [Access-Control-Allow-Origin = %v] is not correct", strings.Join(accessControlAllowOrigin, ", "))
			isSuccess = false
		}
	}

	accessControlAllowMethods, ok := w.Result().Header["Access-Control-Allow-Methods"]
	if !ok {
		t.Errorf("\033[41mERROR\033[0m Not found header [Access-Control-Allow-Methods]")
		isSuccess = false
	} else {
		if !linq.From(strings.Split(accessControlAllowMethods[0], ",")).Select(func(x interface{}) interface{} {
			return strings.TrimSpace(x.(string))
		}).AnyWith(func(x interface{}) bool {
			return strings.EqualFold(x.(string), "*") || strings.EqualFold(x.(string), method)
		}) {
			t.Errorf("\033[41mERROR\033[0m Header [Access-Control-Allow-Methods = %v] is not correct", strings.Join(accessControlAllowMethods, ", "))
			isSuccess = false
		}
	}

	accessControlAllowHeaders, ok := w.Result().Header["Access-Control-Allow-Headers"]
	if !ok {
		t.Errorf("\033[41mERROR\033[0m Not found header [Access-Control-Allow-Headers]")
		isSuccess = false
	} else {
		if !linq.From(strings.Split(headers, ",")).Select(func(x interface{}) interface{} {
			return strings.TrimSpace(x.(string))
		}).All(func(x interface{}) bool {
			return accessControlAllowHeaders[0] == "*" || linq.From(strings.Split(accessControlAllowHeaders[0], ",")).Select(func(x interface{}) interface{} {
				return strings.TrimSpace(x.(string))
			}).AnyWith(func(y interface{}) bool {
				return strings.EqualFold(x.(string), y.(string))
			})
		}) {
			t.Errorf("\033[41mERROR\033[0m Header [Access-Control-Allow-Headers = %v] is not correct", strings.Join(accessControlAllowHeaders, ", "))
			isSuccess = false
		}
	}

	accessControlExposeHeaders, ok := w.Result().Header["Access-Control-Expose-Headers"]
	if !ok {
		t.Errorf("\033[41mERROR\033[0m Not found header [Access-Control-Expose-Headers]")
		isSuccess = false
	} else {
		if !linq.From(strings.Split(headers, ",")).Select(func(x interface{}) interface{} {
			return strings.TrimSpace(x.(string))
		}).All(func(x interface{}) bool {
			return accessControlExposeHeaders[0] == "*" || linq.From(strings.Split(accessControlExposeHeaders[0], ",")).Select(func(x interface{}) interface{} {
				return strings.TrimSpace(x.(string))
			}).AnyWith(func(y interface{}) bool {
				return strings.EqualFold(x.(string), y.(string))
			})
		}) {
			t.Errorf("\033[41mERROR\033[0m Header [Access-Control-Expose-Headers = %v] is not correct", strings.Join(accessControlExposeHeaders, ", "))
			isSuccess = false
		}
	}

	if isSuccess {
		t.Logf("\033[42mPASS\033[0m Method: %v", http.MethodOptions)
	} else {
		t.Logf("\033[41mERROR\033[0m Method: %v", http.MethodOptions)
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
