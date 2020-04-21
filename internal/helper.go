package internal

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
)

func isSameHost(a, b string) bool {
	aa, _ := url.Parse(a)
	bb, _ := url.Parse(b)
	return aa.Host == bb.Host
}

func urlWithScheme(s string) string {
	uri, _ := url.Parse(s)
	return fmt.Sprintf("%s://%s", uri.Scheme, uri.Host)
}

func Open(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
