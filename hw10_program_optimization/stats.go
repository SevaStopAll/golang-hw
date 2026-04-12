package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	json := jsoniter.ConfigFastest
	suffix := "." + domain

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		// Извлекаем только поле Email без полного парсинга структуры
		email := json.Get(line, "Email").ToString()
		if email == "" {
			continue
		}

		atIdx := strings.IndexByte(email, '@')
		if atIdx == -1 || atIdx == len(email)-1 {
			continue
		}

		emailDomain := email[atIdx+1:]
		lowerDomain := strings.ToLower(emailDomain)

		if strings.HasSuffix(lowerDomain, suffix) {
			result[lowerDomain]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return result, nil
}
