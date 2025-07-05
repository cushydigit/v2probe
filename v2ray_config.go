package main

type V2RayConfig struct {
	Log       Log        `json:"log"`
	Inbounds  []Inbound  `json:"inbounds"`
	Outbounds []Outbound `json:"outbounds"`
}

type Inbound struct {
	Port     int            `json:"port"`
	Listen   string         `json:"listen,omitempty"`
	Protocol string         `json:"protocol"`
	Settings map[string]any `json:"settings"`
	Sniffing *Sniffing      `json:"sniffing,omitempty"`
	Tag      string         `json:"tag,omitempty"`
}

type Sniffing struct {
	Enabled      bool     `json:"enabled"`
	DestOverride []string `json:"destOverride"`
}

type Outbound struct {
	Protocol       string          `json:"protocol"` // vmess, vless, shadowsocks, trojan
	Settings       map[string]any  `json:"settings"`
	StreamSettings *StreamSettings `json:"streamSettings,omitempty"`
	Tag            string          `json:"tag,omitempty"`
	Mux            map[string]any  `json:"mux,omitempty"`
}

type StreamSettings struct {
	Network      string         `json:"network,omitempty"`  // tcp, ws, h2, grpc
	Security     string         `json:"security,omitempty"` // tls or none
	TLSSettings  map[string]any `json:"tlsSettings,omitempty"`
	WSSettings   map[string]any `json:"wsSettings,omitempty"`
	HTTPSettings map[string]any `json:"httpSettings,omitempty"`
	GRPCSettings map[string]any `json:"grpcSettings,omitempty"`
}

type Log struct {
	Loglevel string `json:"loglevel"`
}
