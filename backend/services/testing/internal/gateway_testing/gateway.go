package gateway_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/http/cookiejar"
	"social-network/services/testing/internal/configs"
	"social-network/services/testing/internal/utils"
	"sync"
	"time"
)

var baseURL = "http://api-gateway:8081"

func StartTest(ctx context.Context, cfgs configs.Configs) error {
	var wg sync.WaitGroup
	wg.Go(func() { utils.HandleErr("api-gateway", ctx, testAuthFlow) })
	time.Sleep(time.Second) //sleeping so that ratelimiting caused by the next tests doesn't affect the above test
	wg.Go(func() { utils.HandleErr("api-gateway", ctx, randomRegister) })
	wg.Go(func() { utils.HandleErr("api-gateway", ctx, randomLogin) })
	wg.Wait()
	return nil
}

func testAuthFlow(ctx context.Context) error {
	fmt.Println("api-gateway Start test auth flow")
	// Create client with cookie jar to persist cookies between requests
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("failed to create cookie jar: %w", err)
	}
	client := &http.Client{Jar: jar}

	// 1. Register
	registerData := newRegisterReq()

	resp, err := postJSON(client, baseURL+"/register", registerData)
	if err != nil {
		return fmt.Errorf("register failed: %w, body: %s", err, bodyToString(resp))
	}
	// fmt.Printf("Register: %d\n", resp.StatusCode)

	email, _ := (*registerData)["email"].(string)
	pass, _ := (*registerData)["password"].(string)

	// 2. Login
	loginData := map[string]any{
		"identifier": email,
		"password":   pass,
	}

	resp, err = postJSON(client, baseURL+"/login", loginData)
	if err != nil {
		return fmt.Errorf("login failed: %w, body: %s", err, bodyToString(resp))
	}
	// fmt.Printf("Login: %d\n", resp.StatusCode)

	// 3. Auth status
	resp, err = client.Get(baseURL + "/auth-status")
	if err != nil {
		return fmt.Errorf("auth status failed: %w, body: %s", err, bodyToString(resp))
	}
	// fmt.Printf("Auth Status: %d\n", resp.StatusCode)

	// 4. Logout
	resp, err = postJSON(client, baseURL+"/logout", nil)
	if err != nil {
		return fmt.Errorf("logout failed: %w, body: %s", err, bodyToString(resp))
	}
	// fmt.Printf("Logout: %d\n", resp.StatusCode)

	// 5. Auth status
	resp, err = client.Get(baseURL + "/auth-status")
	if err != nil {
		return fmt.Errorf("second auth status failed: %w, body: %s", err, bodyToString(resp))
	}
	// fmt.Printf("Auth Status: %d\n", resp.StatusCode)
	fmt.Println("api-gateway Finished test auth flow")
	return nil
}

func postJSON(client *http.Client, url string, data any) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

func randomRegister(ctx context.Context) error {
	fmt.Println("api-gateway starting register test")
	client := &http.Client{}
	gotRateLimited := false
	for range 100 {
		registerData := newRegisterReq()
		resp, err := postJSON(client, baseURL+"/register", registerData)
		if err != nil {
			return fmt.Errorf("spam register failed: %w", err)
		}

		if resp.StatusCode/200 != 1 && resp.StatusCode != 429 {
			return fmt.Errorf("bad status when spam registering: %d, body: %s", resp.StatusCode, bodyToString(resp))
		}

		if resp.StatusCode == 429 {
			gotRateLimited = true
		}
		time.Sleep(time.Millisecond * 50)
	}

	if gotRateLimited == false {
		return fmt.Errorf("register spam didn't get ratelimited when spamming?!")
	}

	fmt.Println("api-gateway spam api-gateway register test passed")
	return nil
}

func randomLogin(ctx context.Context) error {
	fmt.Println("api-gateway starting Login test")
	client := &http.Client{}
	gotRateLimited := false
	for range 100 {
		loginReq := newLoginReq()
		resp, err := postJSON(client, baseURL+"/login", loginReq)
		if err != nil {
			return fmt.Errorf("spam login failed: %w", err)
		}
		if resp.StatusCode == 200 {
			return fmt.Errorf("somehow managed to login while spamming bad logins??: %d, body: %s", resp.StatusCode, bodyToString(resp))
		}
		if resp.StatusCode == 429 {
			gotRateLimited = true
		}
		time.Sleep(time.Millisecond * 50)
	}

	if gotRateLimited == false {
		return fmt.Errorf("login spam didn't get ratelimited when spamming?!")
	}

	fmt.Println("api-gateway spam login test passed")
	return nil
}

func newRegisterReq() *map[string]any {
	registerData := map[string]any{
		"username":      utils.Title(utils.RandomString(10, false)),
		"first_name":    utils.Title(utils.RandomString(10, false)),
		"last_name":     utils.Title(utils.RandomString(10, false)),
		"date_of_birth": fmt.Sprintf("%d-%02d-%02d", rand.IntN(50)+1950, rand.IntN(11)+1, rand.IntN(20)+1),
		"about":         utils.RandomString(300, true),
		"public":        rand.IntN(3) == 0,
		"email":         utils.RandomString(10, false) + "@hotmail.com",
		"password":      utils.RandomPassword(),
	}

	return &registerData
}

func newLoginReq() *map[string]any {
	login := map[string]any{
		"identifier": utils.Title(utils.RandomString(10, false)),
		"password":   utils.RandomString(10, true),
	}
	return &login
}

func bodyToString(r *http.Response) string {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	return string(body)
}
