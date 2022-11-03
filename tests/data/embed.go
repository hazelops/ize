package data

import "embed"

var (
	//go:embed */*
	TestData embed.FS
)
