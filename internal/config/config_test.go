package config

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestLoadConfigRejectsMissingCredentials(t *testing.T) {
	for _, missing := range []string{"HOST", "TOKEN"} {
		t.Run(missing, func(t *testing.T) {
			t.Chdir(t.TempDir())
			t.Setenv("HOST", "misskey.example")
			t.Setenv("TOKEN", "test-token")
			t.Setenv(missing, "")
			cfg, err := LoadConfig(slog.New(slog.NewTextHandler(io.Discard, nil)))
			if err == nil || !strings.Contains(err.Error(), missing) || cfg != nil {
				t.Fatalf("LoadConfig() = %v, %v; want error identifying %s", cfg, err, missing)
			}
		})
	}
}

func TestLoadConfigReadsDotEnvWithoutOverridingEnvironment(t *testing.T) {
	for _, useEnvironment := range []bool{false, true} {
		name := "dotenv"
		if useEnvironment {
			name = "environment"
		}
		t.Run(name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			for _, key := range []string{"HOST", "TOKEN"} {
				// Register restoration before removing inherited values.
				t.Setenv(key, "")
				if err := os.Unsetenv(key); err != nil {
					t.Fatal(err)
				}
			}
			if err := os.WriteFile(".env", []byte("HOST=dotenv.example\nTOKEN=dotenv-token\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			want := Config{Host: "dotenv.example", Token: "dotenv-token"}
			if useEnvironment {
				t.Setenv("HOST", "environment.example")
				t.Setenv("TOKEN", "environment-token")
				want = Config{Host: "environment.example", Token: "environment-token"}
			}
			got, err := LoadConfig(slog.New(slog.NewTextHandler(io.Discard, nil)))
			if err != nil {
				t.Fatal(err)
			}
			if *got != want {
				t.Fatalf("LoadConfig() = %v, want %v", *got, want)
			}
		})
	}
}
