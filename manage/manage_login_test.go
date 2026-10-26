package manage

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/feelingsray/ray-utils-go/v2/rotp"
)

func TestLoginAuthLockAndExpiry(t *testing.T) {
	p := &AppManage{SuperAuth: map[string]string{"test-admin": "test-password"}}
	now := time.Now()
	for i := 0; i < loginFailureLimit; i++ {
		ok, err, _ := p.loginAuth("test-admin", "wrong", now)
		if ok || (i == loginFailureLimit-1 && !errors.Is(err, errLoginLocked)) {
			t.Fatalf("attempt %d: ok=%v err=%v", i+1, ok, err)
		}
	}
	for _, elapsed := range []time.Duration{0, loginLockDuration - time.Nanosecond} {
		ok, err, _ := p.loginAuth("test-admin", "test-password", now.Add(elapsed))
		if ok || !errors.Is(err, errLoginLocked) {
			t.Fatalf("locked login at %v: ok=%v err=%v", elapsed, ok, err)
		}
	}
	ok, err, _ := p.loginAuth("test-admin", "test-password", now.Add(loginLockDuration))
	if !ok || err != nil {
		t.Fatalf("expired lock: ok=%v err=%v", ok, err)
	}
}

func TestLoginAuthSuccessResetsFailures(t *testing.T) {
	p := &AppManage{SuperAuth: map[string]string{"test-admin": "test-password"}}
	now := time.Now()
	for cycle := 0; cycle < 2; cycle++ {
		for i := 0; i < loginFailureLimit-1; i++ {
			p.loginAuth("test-admin", "wrong", now)
		}
		if ok, err, _ := p.loginAuth("test-admin", "test-password", now); !ok || err != nil {
			t.Fatalf("cycle %d: ok=%v err=%v", cycle, ok, err)
		}
	}
}

func TestLoginAuthConcurrentOTPAliases(t *testing.T) {
	p := &AppManage{}
	now := time.Now()
	var unlocked atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			username := []string{"otp", "OTP"}[i%2]
			ok, err, _ := p.loginAuth(username, "invalid-code", now)
			if ok {
				t.Error("invalid OTP accepted")
			}
			if !errors.Is(err, errLoginLocked) {
				unlocked.Add(1)
			}
		}(i)
	}
	wg.Wait()
	if unlocked.Load() != loginFailureLimit-1 {
		t.Fatalf("unlocked responses = %d", unlocked.Load())
	}
	if len(p.loginFailures) != 1 || p.loginFailures["otp"].count != loginFailureLimit {
		t.Fatal("OTP aliases did not share a bounded failure counter")
	}
}

func TestLoginAuthUnknownUsersBounded(t *testing.T) {
	p := &AppManage{}
	for i := 0; i < 100; i++ {
		p.loginAuth(fmt.Sprintf("unknown-%d", i), "wrong", time.Now())
	}
	if len(p.loginFailures) != 1 {
		t.Fatalf("failure counters = %d", len(p.loginFailures))
	}
}

func TestLoginHTTP(t *testing.T) {
	p := &AppManage{SuperAuth: map[string]string{"test-admin": "test-password"}}
	router := gin.New()
	router.POST("/mapi/login", p.login)
	request := func(password string) *httptest.ResponseRecorder {
		body := fmt.Sprintf(`{"username":"test-admin","password":%q,"timestamp":"1"}`, password)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/mapi/login", strings.NewReader(body)))
		return w
	}
	w := request("test-password")
	var response struct {
		Code int `json:"code"`
		Data struct {
			Role string `json:"role"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || response.Code != 200 || response.Data.Role != "super" {
		t.Fatalf("successful login response changed: %s", w.Body.String())
	}
	for i := 0; i < loginFailureLimit; i++ {
		w = request("wrong")
	}
	if w.Code != http.StatusTooManyRequests || request("test-password").Code != http.StatusTooManyRequests {
		t.Fatal("HTTP login did not enforce the lock")
	}
}

func TestLoginRejectsClientSecret(t *testing.T) {
	p := &AppManage{SuperAuth: map[string]string{"test-admin": "test-password"}}
	stamp := time.Unix(1500000000, 0)
	secret := "TESTONLY"
	code := rotp.RTOTPCodeWithTime(secret, stamp)
	for {
		ok, _, _ := p.basicAuth("OTP", code, time.Now().Unix(), nil)
		if !ok {
			break
		}
		stamp = stamp.Add(30 * time.Second)
		code = rotp.RTOTPCodeWithTime(secret, stamp)
	}
	router := gin.New()
	router.POST("/mapi/login", p.login)
	body := fmt.Sprintf(`{"username":"OTP","password":%q,"timestamp":%q,"secret":%q}`, code, fmt.Sprint(stamp.Unix()), secret)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/mapi/login", strings.NewReader(body)))
	var response struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || response.Code == http.StatusOK {
		t.Fatalf("client-controlled OTP accepted: %s", w.Body.String())
	}
}
