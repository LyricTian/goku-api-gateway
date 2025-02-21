package admin

import (
	"net/http"
	"strconv"
	"strings"
)

//GetIPPort 获取客户端IP和端口
func GetIPPort(r *http.Request) (string, int, error) {
	ip := r.RemoteAddr
	ip = ip[:strings.Index(ip, ":")]
	if realIP := strings.TrimSpace(r.Header.Get("X-Real-Ip")); realIP != "" {
		ip = realIP
	}
	r.ParseForm()
	p := r.FormValue("port")
	port, err := strconv.Atoi(p)
	if err != nil {
		return ip, port, err
	}
	return ip, port, nil
}
