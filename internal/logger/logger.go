package logger

import (
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Level задает уровень логирования.
type Level int

const (
	DEBUG Level = iota // отладочная информация
	INFO               // информационные сообщения
	WARN               // предупреждения
	ERROR              // ошибки
)

// Logger — структура логгера с разными уровнями.
type Logger struct {
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
	level Level
}

// Log — глобальный экземпляр логгера.
var Log *Logger

// Init инициализирует глобальный логгер Log.
// filepath — путь к файлу логов.
// level — минимальный уровень логирования (DEBUG, INFO, WARN, ERROR).
// Логи пишутся одновременно в stdout и в файл с ротацией через lumberjack.
func Init(filepath string, level Level) {
	rotator := &lumberjack.Logger{
		Filename:   filepath,
		MaxSize:    5,    // MB
		MaxBackups: 3,    // количество старых файлов
		MaxAge:     28,   // дни
		Compress:   true, // сжатие gzip
	}

	writer := io.MultiWriter(os.Stdout, rotator)

	Log = &Logger{
		debug: log.New(writer, "[DEBUG] ", log.LstdFlags|log.Lmicroseconds),
		info:  log.New(writer, "[INFO] ", log.LstdFlags|log.Lmicroseconds),
		warn:  log.New(writer, "[WARN] ", log.LstdFlags|log.Lmicroseconds),
		err:   log.New(writer, "[ERROR] ", log.LstdFlags|log.Lmicroseconds),
		level: level,
	}
}

// Debug выводит сообщение уровня DEBUG, если текущий уровень <= DEBUG.
func (l *Logger) Debug(v ...interface{}) {
	if l.level <= DEBUG {
		l.debug.Println(v...)
	}
}

// Info выводит сообщение уровня INFO, если текущий уровень <= INFO.
func (l *Logger) Info(v ...interface{}) {
	if l.level <= INFO {
		l.info.Println(v...)
	}
}

// Warn выводит сообщение уровня WARN, если текущий уровень <= WARN.
func (l *Logger) Warn(v ...interface{}) {
	if l.level <= WARN {
		l.warn.Println(v...)
	}
}

// Error выводит сообщение уровня ERROR.
func (l *Logger) Error(v ...interface{}) {
	l.err.Println(v...)
}
