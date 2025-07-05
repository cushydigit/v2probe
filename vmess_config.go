package main

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
)

type VMessConfig struct {
	Inbounds  []Inbound  `json:"inbound"`
	Outbounds []Outbound `json:"outbound"`
}

type Inbound struct {
	Port     int    `json:"port"`
	Listen   string `json:"listen,omitempty"`
	Protocol string `json:"protocol"`
	Settings struct {
		Clients []map[string]string `json:"clients"`
	} `json:"settings"`
}

type Outbound struct {
	Protocol       string         `json:"protocol"`
	Settings       map[string]any `json:"settings"`
	StreamSettings map[string]any `json:"streamSettings,omitempty"`
	Tag            string         `json:"tag,omitempty"`
}

func parseVMess(link string) (*VMessConfig, error) {
	encoded := strings.TrimPrefix(link, "vmess://")
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	var config map[string]any
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	outbound := Outbound{
		Protocol: "vmess",
		Settings: map[string]any{
			"vnext": []map[string]any{
				{
					"address": config["add"],
					"port":    toInt(config["port"]),
					"users": []map[string]any{
						{
							"id":       config["id"],
							"alterId":  toInt(config["aid"]),
							"security": "auto",
						},
					},
				},
			},
		},
	}

	return &VMessConfig{
		Inbounds:  []Inbound{},
		Outbounds: []Outbound{outbound},
	}, nil
}

func toInt(v any) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	case float32:
		return int(val)
	case string:
		i, _ := strconv.Atoi(val)
		return i
	case []byte:
		i, _ := strconv.Atoi(string(val))
		return i
	default:
		return 0 // fallback if unsupported type
	}
}
