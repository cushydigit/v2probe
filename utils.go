package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func checkProtocolValid(s string) bool {
	return (strings.HasPrefix(s, PROTOCOL_VLESS) ||
		strings.HasPrefix(s, PROTOCOL_VMESS) ||
		strings.HasPrefix(s, PROTOCOL_TROJAN) ||
		strings.HasPrefix(s, PROTOCOL_SS))
}

func saveConfigToFile(cfg *V2RayConfig, filename string) error {
	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Clean(filename), out, 0644)
}

func fetchAndExtractLinks() []string {
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, SOURCE_URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("unexpected statusCode: %d", res.StatusCode)
	}

	links := extractLinks(res.Body)
	return links

}

func extractLinks(r io.Reader) []string {
	scanner := bufio.NewScanner(r)
	var links []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if checkProtocolValid(line) {
			links = append(links, line)
		}
	}
	return links
}
