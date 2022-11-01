package strutil

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var regex = regexp.MustCompile(`[a-z0-9]`)

// source: https://github.com/jaegertracing/jaeger-operator/blob/91e3b69ee5c8761bbda9d3cf431400a73fc1112a/pkg/util/dns_name.go#L15
func DNSName(name string) string {
	var d []rune

	for i, x := range strings.ToLower(name) {
		if regex.Match([]byte(string(x))) {
			d = append(d, x)
		} else {
			if i == 0 || i == utf8.RuneCountInString(name)-1 {
				d = append(d, 'a')
			} else {
				d = append(d, '-')
			}
		}
	}

	return string(d)
}