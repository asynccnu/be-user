package service

import (
	"context"
	"fmt"
	"github.com/anaskhan96/soup"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

type CCNUService interface {
	Login(ctx context.Context, studentId string, password string) (bool, error)
}

type ccnuService struct {
	timeout time.Duration
}

func NewCCNUService() CCNUService {
	return &ccnuService{timeout: time.Second * 5}
}

func (c *ccnuService) Login(ctx context.Context, studentId string, password string) (bool, error) {
	js, st := c.preload()
	url := fmt.Sprintf("https://account.ccnu.edu.cn/cas/login;jsessionid=%v?service=http%3A%2F%2Fone.ccnu.edu.cn%2Fcas%2Flogin_portal", js)
	text := fmt.Sprintf("username=%s&password=%s&lt=%s&execution=e1s1&_eventId=submit&submit=%E7%99%BB%E5%BD%95", studentId, password, st)
	body := strings.NewReader(text)
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Cookie", "JSESSIONID="+js)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://account.ccnu.edu.cn")
	req.Header.Set("Referer", "https://account.ccnu.edu.cn/cas/login?service=http%3A%2F%2Fone.ccnu.edu.cn%2Fcas%2Flogin_portal")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0")
	req.Header.Set("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120", "Microsoft Edge";v="120"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	resp, err := c.client().Do(req)
	if err != nil {
		return false, err
	}
	return resp.Header.Get("Pragma") == "", nil
}

// preload 返回值为 js 和 st
func (c *ccnuService) preload() (string, string) {
	htmlBody, _ := soup.Get("https://account.ccnu.edu.cn/cas/login?service=http%3A%2F%2Fone.ccnu.edu.cn%2Fcas%2Flogin_portal")
	doc := soup.HTMLParse(htmlBody)
	links1 := doc.Find("body", "id", "cas").FindAll("script")
	js := links1[2].Attrs()["src"][26:]
	links2 := doc.Find("div", "class", "logo").FindAll("input")
	st := links2[2].Attrs()["value"]
	return js, st
}

func (c *ccnuService) client() *http.Client {
	j, _ := cookiejar.New(&cookiejar.Options{})
	return &http.Client{
		Timeout: c.timeout,
		Jar:     j,
	}
}
