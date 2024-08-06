package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/natefinch/lumberjack"
)

// LogConfig is the log config
type LogConfig struct {
	Level     string       `mapstructure:"level"`
	Format    string       `mapstructure:"format"`
	AddSource bool         `mapstructure:"add_source"`
	File      string       `mapstructure:"file"`
	Rotate    RotateConfig `mapstructure:"rotate"`
}

// Parse parses the log config
func (c *LogConfig) parse() {
	if c.Level == "" {
		c.Level = "debug"
	}
	if c.Format == "" {
		c.Format = "text"
	}
	c.Rotate.Parse()
}

// InitLog Config is the utils config
func InitLog(cfg LogConfig) {
	cfg.parse()

	opts := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     slog.LevelInfo,
	}
	// init logger
	switch cfg.Level {
	case "debug":
		opts.Level = slog.LevelDebug
	case "info":
		opts.Level = slog.LevelInfo
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	default:
		slog.Error("invalid log level", "level", cfg.Level)
	}

	var w io.WriteCloser
	if cfg.File == "" {
		w = os.Stdout
	} else {
		w = NewRotate(cfg.File, cfg.Rotate)
	}

	var logger *slog.Logger
	switch cfg.Format {
	case "json":
		h := slog.NewJSONHandler(w, opts)
		logger = slog.New(h)
	default:
		h := slog.NewTextHandler(w, opts)
		logger = slog.New(h)
	}

	slog.SetDefault(logger)
}

var (
	defaultMaxSize    = 100
	defaultMaxAge     = 7
	defaultMaxBackups = 5
)

// Config rotate config
type RotateConfig struct {
	MaxSize    int  `mapstructure:"max_size"`
	MaxAge     int  `mapstructure:"max_age"`
	MaxBackups int  `mapstructure:"max_backups"`
	LocalTime  bool `mapstructure:"local_time"`
	Compress   bool `mapstructure:"compress"`
}

// Parse parses the rotate config
func (c *RotateConfig) Parse() {
	if c.MaxSize == 0 {
		c.MaxSize = defaultMaxSize
	}
	if c.MaxAge == 0 {
		c.MaxAge = defaultMaxAge
	}
	if c.MaxBackups == 0 {
		c.MaxBackups = defaultMaxBackups
	}
}

// NewRotate creates a new rotate writer
func NewRotate(filename string, c RotateConfig) io.WriteCloser {
	w := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    c.MaxSize,
		MaxAge:     c.MaxAge,
		MaxBackups: c.MaxBackups,
		LocalTime:  c.LocalTime,
		Compress:   c.Compress,
	}
	return w
}
