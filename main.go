package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	SOURCE_URL      = "https://raw.githubusercontent.com/barry-far/V2ray-Config/refs/heads/main/All_Configs_Sub.txt"
	CONFIG_DIR      = "./.config"
	PROTOCOL_VMESS  = "vmess://"
	PROTOCOL_VLESS  = "lmess://"
	PROTOCOL_SS     = "ss://"
	PROTOCOL_TROJAN = "trojan://"
)

func main() {

	// remove if any existing config exist
	if err := os.RemoveAll(CONFIG_DIR); err != nil {
		log.Fatalf("failed to clear the directory for config files: %v", err)
	}

	// create directory if not exists
	if err := os.MkdirAll(CONFIG_DIR, os.ModePerm); err != nil {
		log.Fatalf("failed to create directory for config files: %v", err)
	}

	// clear directory

	// printing all configs
	links := fetchAndExtractLinks()
	counter := 0

	for _, link := range links {
		var v2RayCfg *V2RayConfig

		switch {
		case strings.HasPrefix(link, PROTOCOL_VMESS):
			cfg, err := decodeVMessLink(link)
			if err != nil {
				log.Printf("failed to decode vmess link: %v", err)
				continue
			}
			v2RayCfg = cfg.generateV2rayConfig()

		case strings.HasPrefix(link, PROTOCOL_VLESS):
			continue

		case strings.HasPrefix(link, PROTOCOL_SS):
			continue

		case strings.HasPrefix(link, PROTOCOL_TROJAN):
			continue

		default:
			log.Printf("unknown protocol: %s", link[:8])
			continue

		}

		counter++
		v2RayCfg.Log = Log{Loglevel: "debug"}
		v2RayCfg.Inbounds = []Inbound{
			{
				// Port:     10808 + counter, // Increment for each config
				Port:     10800,
				Listen:   "127.0.0.1",
				Protocol: "socks",
				Settings: map[string]any{"auth": "noauth"},
			},
		}

		outputFile := fmt.Sprintf("%s/xray_vmess_%05d.json", CONFIG_DIR, counter)
		if err := saveConfigToFile(v2RayCfg, outputFile); err != nil {
			log.Printf("failed to write config: %v", err)
			counter--
			continue
		}

	}

	log.Printf("the total configs links: %d", counter)

}
