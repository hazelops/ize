package examples

import "embed"

var (
	//go:embed */*
	Examples embed.FS
)
