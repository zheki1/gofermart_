package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

// testLogger создаёт логгер с выводом в bytes.Buffer для тестирования
func testLogger(level Level) (*Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	l := &Logger{
		debug: log.New(buf, "[DEBUG] ", 0),
		info:  log.New(buf, "[INFO] ", 0),
		warn:  log.New(buf, "[WARN] ", 0),
		err:   log.New(buf, "[ERROR] ", 0),
		level: level,
	}
	return l, buf
}

func TestLoggerLevels(t *testing.T) {
	l, buf := testLogger(INFO)

	l.Debug("debug") // не должен выводиться
	l.Info("info")
	l.Warn("warn")
	l.Error("error")

	output := buf.String()
	if strings.Contains(output, "debug") {
		t.Error("DEBUG должен быть пропущен при уровне INFO")
	}
	if !strings.Contains(output, "info") {
		t.Error("INFO должно выводиться")
	}
	if !strings.Contains(output, "warn") {
		t.Error("WARN должно выводиться")
	}
	if !strings.Contains(output, "error") {
		t.Error("ERROR должно выводиться")
	}
}

func TestInitCreatesLogger(t *testing.T) {
	Init("/tmp/test.log", DEBUG)
	if Log == nil {
		t.Fatal("Init не инициализировал Log")
	}
	if Log.level != DEBUG {
		t.Fatalf("ожидали уровень DEBUG, получили %v", Log.level)
	}
}
