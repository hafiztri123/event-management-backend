package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/google/uuid"
)

type LogLevel string

const (
	DebugLevel LogLevel = "DEBUG"
	InfoLevel LogLevel = "INFO"
	WarnLevel LogLevel = "WARN"
	ErrorLevel LogLevel = "ERROR"
	FatalLevel LogLevel = "FATAL"
)

type Config struct {
	AppName 		string
	Environment 	string
	MinLogLevel 	LogLevel
	Output 			io.Writer
	EnableConsole 	bool
	EnableFile 		bool
	LogFilePath 	string
}

type Logger struct {
	config Config
	fileWriter io.Writer
}

type LogEntry struct {
	Timestamp   string                 `json:"timestamp"`
	Level       LogLevel               `json:"level"`
	Message     string                 `json:"message"`
	App         string                 `json:"app"`
	Environment string                 `json:"environment"`
	RequestID   string                 `json:"request_id,omitempty"`
	UserID      string                 `json:"user_id,omitempty"`
	TraceID     string                 `json:"trace_id,omitempty"`
	File        string                 `json:"file,omitempty"`
	Line        int                    `json:"line,omitempty"`
	Function    string                 `json:"function,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

type contextKey string

const (
	requestIDKey contextKey = "request_id"
	userIDKey contextKey = "user_id"
	traceIDKey contextKey = "trace_id"
)

func New(config Config) (*Logger, error) {
	if config.MinLogLevel == "" {
		config.MinLogLevel = InfoLevel
	}

	if config.Output == nil {
		config.Output = os.Stdout
	}

	var filewriter io.Writer
	if config.EnableFile {
		if config.LogFilePath == "" {
			config.LogFilePath = "application.log"
		}
		file, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("[FAIL] failed to open log file: %w", err)
		}
		filewriter = file
	}

	return &Logger{
		config: config,
		fileWriter: filewriter,
	}, nil
}


func WithRequestID(ctx context.Context, requestID string) context.Context {
	if requestID == "" {
		requestID = uuid.New().String()
	}
	return context.WithValue(ctx, requestIDKey, requestID)
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	if traceID == "" {
		traceID = uuid.New().String()
	}
	return context.WithValue(ctx, traceIDKey, traceID)
}

func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	requestID, _ := ctx.Value(requestIDKey).(string)
	return requestID
}

func GetUserID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	userID, _ := ctx.Value(userIDKey).(string)
	return userID
}

func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	traceID, _ := ctx.Value(traceIDKey).(string)
	return traceID
}


func (l *Logger) log(ctx context.Context, level LogLevel, msg string, data map[string]interface{}) {
	if !l.shouldLog(level) {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level: level,
		Message: msg,
		App: l.config.AppName,
		Environment: l.config.Environment,
		Data: data,
	}

	if ctx != nil {
		entry.RequestID = GetRequestID(ctx)
		entry.UserID = GetUserID(ctx)
		entry.TraceID = GetTraceID(ctx)
	}

	if pc, file, line, ok := runtime.Caller(2); ok {
		entry.File = file
		entry.Line = line
		if fn := runtime.FuncForPC(pc); fn != nil {
			entry.Function = fn.Name()
		}
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal log entry : %v\n", err)
		return
	}

	if l.config.EnableConsole {
		fmt.Fprintln(l.config.Output, string(jsonData))
	}

	if l.config.EnableFile && l.fileWriter != nil {
		fmt.Fprintln(l.fileWriter, string(jsonData))
	}

	if level == FatalLevel {
		os.Exit(1)
	}
	
	
}

func (l *Logger) shouldLog(level LogLevel) bool {
	switch l.config.MinLogLevel {
	case DebugLevel:
		return true
	case InfoLevel:
		return level != DebugLevel
	case WarnLevel:
		return level != DebugLevel && level != InfoLevel
	case ErrorLevel:
		return level == ErrorLevel || level == FatalLevel
	case FatalLevel:
		return level == FatalLevel
	default:
		return true
	}
}

func (l *Logger) Debug(ctx context.Context, msg string, data map[string]interface{}) {
	l.log(ctx, DebugLevel, msg, data)
}

func (l *Logger) Info(ctx context.Context, msg string, data map[string]interface{}) {
	l.log(ctx, InfoLevel, msg, data)
}

func (l *Logger) Warn(ctx context.Context, msg string, data map[string]interface{}) {
	l.log(ctx, WarnLevel, msg, data)
}

func (l *Logger) Error(ctx context.Context, msg string, err error, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	
	if err != nil {
		data["error"] = err.Error()
	}
	
	l.log(ctx, ErrorLevel, msg, data)
}

func (l *Logger) Fatal(ctx context.Context, msg string, err error, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	
	if err != nil {
		data["error"] = err.Error()
	}
	
	l.log(ctx, FatalLevel, msg, data)
}

func (l *Logger) Close() error {
	if closer, ok := l.fileWriter.(io.Closer); ok && closer != nil {
		return closer.Close()
	}
	return nil
}
