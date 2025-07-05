package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	SOURCE_URL = "https://raw.githubusercontent.com/barry-far/V2ray-Config/refs/heads/main/All_Configs_Sub.txt"
	CONFIG_DIR = "./.config"
)

func main() {

	// create directory if not exists
	if err := os.MkdirAll(CONFIG_DIR, os.ModePerm); err != nil {
		log.Fatalf("failed to create directory for configfiles: %v", err)
	}

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

	// printing all configs
	configs := extractConfigs(res.Body)
	counter := 1

	for i, config := range configs {
		if strings.HasPrefix(config, "vmess://") {
			cfg, err := parseVMess(config)
			if err != nil {
				fmt.Printf("error parsing VMess link on line %d: %v\n", i, err)
				continue
			}
			cfg.Inbounds = []Inbound{
				{
					Port:     10808 + counter, // Increment for each config
					Listen:   "127.0.0.1",
					Protocol: "socks",
					Settings: struct {
						Clients []map[string]string `json:"clients"`
					}{},
				},
			}

			// Save the config
			outputFile := fmt.Sprintf("%s/xray_vmess_%05d.json", CONFIG_DIR, i)
			err = saveConfigToFile(cfg, outputFile)
			if err != nil {
				fmt.Println("failed to write config: ", err)

			} else {
				fmt.Println("Saved: ", outputFile)
			}

			counter++
		}
	}
	log.Printf("the configs length: %d", len(configs))

	log.Println("it works")
}

func extractConfigs(r io.Reader) []string {
	scanner := bufio.NewScanner(r)
	var configs []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if checkProtocolValid(line) {
			configs = append(configs, line)
		}
	}
	return configs
}

func checkProtocolValid(s string) bool {
	return strings.HasPrefix(s, "vmess://") || strings.HasPrefix(s, "vless://") || strings.HasPrefix(s, "ss://") || strings.HasPrefix(s, "trojan://")
}

func saveConfigToFile(cfg *VMessConfig, filename string) error {
	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Clean(filename), out, 0644)
}
