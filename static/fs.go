package static

import "embed"

//go:embed out/*
var Assets embed.FS
