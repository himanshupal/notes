package main

import "os"

func init() {
	if err := os.Setenv("HOST", ""); err != nil {
		panic(err)
	}

	if err := os.Setenv("PQ_HOST", ""); err != nil {
		panic(err)
	}

	if err := os.Setenv("PQ_PORT", ""); err != nil {
		panic(err)
	}

	if err := os.Setenv("PQ_USER", ""); err != nil {
		panic(err)
	}

	if err := os.Setenv("PQ_PASSWORD", ""); err != nil {
		panic(err)
	}

	if err := os.Setenv("PQ_DATABASE", ""); err != nil {
		panic(err)
	}

	if err := os.Setenv("JWT_SECRET", ""); err != nil {
		panic(err)
	}

	if err := os.Setenv("EMAIL_USERNAME", ""); err != nil {
		panic(err)
	}

	if err := os.Setenv("EMAIL_PASSWORD", ""); err != nil {
		panic(err)
	}

	if err := os.Setenv("EMAIL_FROM", ""); err != nil {
		panic(err)
	}
}
