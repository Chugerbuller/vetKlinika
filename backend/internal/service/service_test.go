package service

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestConvertDate(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{"ISO-дата из браузера", "2024-03-15T00:00:00.000Z", "20240315", false},
		{"время не влияет на дату", "2023-12-31T23:59:59.999Z", "20231231", false},
		{"пустая строка", "", "", true},
		{"уже в формате БД", "20240315", "", true},
		{"мусор", "вчера", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertDate(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ConvertDate(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ConvertDate(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestCalcAge(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		date    string
		age     int
		want    int
		wantErr bool
	}{
		{"записан в этом году", now.Format(layout), 3, 3, false},
		{"записан два года назад", now.AddDate(-2, 0, 0).Format(layout), 3, 5, false},
		{"пустая дата", "", 3, 0, true},
		{"неверный формат", "2024-01-01", 3, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalcAge(tt.date, tt.age)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CalcAge(%q, %d) error = %v, wantErr %v", tt.date, tt.age, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("CalcAge(%q, %d) = %d, want %d", tt.date, tt.age, got, tt.want)
			}
		})
	}
}

func TestMustLoad(t *testing.T) {
	t.Run("переменная CONFIG_PATH не задана", func(t *testing.T) {
		t.Setenv("CONFIG_PATH", "")
		if _, err := MustLoad(); err == nil {
			t.Fatal("MustLoad() error = nil, want error")
		}
	})

	t.Run("файла нет", func(t *testing.T) {
		t.Setenv("CONFIG_PATH", filepath.Join(t.TempDir(), "missing.yaml"))
		if _, err := MustLoad(); err == nil {
			t.Fatal("MustLoad() error = nil, want error")
		}
	})

	t.Run("корректный файл", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config.yaml")
		data := "db_path: \"./test.db\"\nport: 9000\nlogger_path: \"./app.log\"\n"
		if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
		t.Setenv("CONFIG_PATH", path)

		cfg, err := MustLoad()
		if err != nil {
			t.Fatalf("MustLoad() error = %v", err)
		}
		if cfg.DbPath != "./test.db" || cfg.Port != "9000" || cfg.LoggerPath != "./app.log" {
			t.Errorf("MustLoad() = %+v", cfg)
		}
	})
}

func TestNewLogger_WritesToFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.log")
	logger := NewLogger(path)
	logger.Info("hello")

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("log file not created: %v", err)
	}
	if len(data) == 0 {
		t.Error("log file is empty")
	}
}
