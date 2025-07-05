package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type VMessConfig struct {
	V    json.Number `json:"v"`
	Ps   string      `json:"ps"`
	Add  string      `json:"add"`
	Port json.Number `json:"port"`
	ID   string      `json:"id"`
	Aid  json.Number `json:"aid"`
	Net  string      `json:"net"`
	Type string      `json:"type"`
	Host string      `json:"host"`
	Path string      `json:"path"`
	TLS  string      `json:"tls"`
}

func decodeVMessLink(link string) (*VMessConfig, error) {
	if len(link) < 8 || !strings.HasPrefix(link, PROTOCOL_VMESS) {
		return nil, fmt.Errorf("invalid vmess link")
	}
	raw := link[8:]
	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, err
	}
	var cfg VMessConfig
	if err := json.Unmarshal(decoded, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *VMessConfig) generateV2rayConfig() *V2RayConfig {
	v2rayCfg := &V2RayConfig{
		Log:      Log{},
		Inbounds: []Inbound{},
		Outbounds: []Outbound{
			{
				Protocol: "vmess",
				Settings: map[string]any{
					"vnext": []map[string]any{
						{
							"address": cfg.Add,
							"port":    toInt(cfg.Port),
							"users": []map[string]any{
								{
									"id":       cfg.ID,
									"alterID":  toInt(cfg.Aid),
									"security": "auto",
								},
							},
						},
					},
				},
				StreamSettings: &StreamSettings{
					Network: cfg.Net,
					Security: func() string {
						if cfg.TLS == "tls" {
							return "tls"
						}
						return "none"
					}(),
					HTTPSettings: map[string]any{
						"path": cfg.Path,
						"host": []string{cfg.Host},
					},
				},
			},
		},
	}
	return v2rayCfg
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
