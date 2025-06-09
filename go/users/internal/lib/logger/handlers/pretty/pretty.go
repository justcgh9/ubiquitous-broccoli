package pretty

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"io"
	stdLog "log"
	"strings"

	"log/slog"

	"github.com/fatih/color"
)

// PrettyHandlerOptions configures PrettyHandler behavior.
type PrettyHandlerOptions struct {
	SlogOpts       *slog.HandlerOptions
	TimeFormat     string
	LevelColors    map[slog.Level]*color.Color
	TimeColor      *color.Color
	MessageColor   *color.Color
	AttrsColor     *color.Color
	ShowCaller     bool
	CallerColor    *color.Color
	IndentJSON     bool
	IndentStr      string
	DisableJSON    bool // If true, don't print fields JSON, print in key=val instead
	LevelTextFunc  func(level slog.Level) string
}

// PrettyHandler implements slog.Handler with pretty colored output.
type PrettyHandler struct {
	opts  PrettyHandlerOptions
	out   io.Writer
	l     *stdLog.Logger
	attrs []slog.Attr
	group string
}

func defaultLevelColors() map[slog.Level]*color.Color {
	return map[slog.Level]*color.Color{
		slog.LevelDebug: color.New(color.FgMagenta),
		slog.LevelInfo:  color.New(color.FgBlue),
		slog.LevelWarn:  color.New(color.FgYellow),
		slog.LevelError: color.New(color.FgRed),
	}
}

func pcToFileLine(pc uintptr) (file string, line int) {
	if pc == 0 {
		return "", 0
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "", 0
	}
	file, line = fn.FileLine(pc - 1)
	return
}

// NewPrettyHandler creates a new PrettyHandler writing to out with given options.
func NewPrettyHandler(out io.Writer, opts PrettyHandlerOptions) *PrettyHandler {
	if opts.TimeFormat == "" {
		opts.TimeFormat = "[15:04:05.000]"
	}
	if opts.LevelColors == nil {
		opts.LevelColors = defaultLevelColors()
	}
	if opts.TimeColor == nil {
		opts.TimeColor = color.New(color.FgHiBlack)
	}
	if opts.MessageColor == nil {
		opts.MessageColor = color.New(color.FgCyan)
	}
	if opts.AttrsColor == nil {
		opts.AttrsColor = color.New(color.FgWhite)
	}
	if opts.CallerColor == nil {
		opts.CallerColor = color.New(color.FgGreen)
	}
	if opts.IndentStr == "" {
		opts.IndentStr = "  "
	}
	if opts.LevelTextFunc == nil {
		opts.LevelTextFunc = func(l slog.Level) string {
			return strings.ToUpper(l.String()) + ":"
		}
	}

	return &PrettyHandler{
		opts: opts,
		out:  out,
		l:    stdLog.New(out, "", 0),
	}
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	levelText := h.opts.LevelTextFunc(r.Level)
	levelColor, ok := h.opts.LevelColors[r.Level]
	if !ok {
		levelColor = color.New()
	}
	levelStr := levelColor.Sprint(levelText)

	timeStr := h.opts.TimeColor.Sprintf("%s", r.Time.Format(h.opts.TimeFormat))
	msgStr := h.opts.MessageColor.Sprint(r.Message)

	// Collect attributes from record and handler persistent attrs
	fields := make(map[string]interface{}, r.NumAttrs()+len(h.attrs))
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})
	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	// Add group prefix to attribute keys if any
	if h.group != "" {
		grouped := make(map[string]interface{}, len(fields))
		for k, v := range fields {
			grouped[h.group+"."+k] = v
		}
		fields = grouped
	}

	var attrsStr string
	if len(fields) > 0 {
		if h.opts.DisableJSON {
			// key=val style
			var sb strings.Builder
			for k, v := range fields {
				fmt.Fprintf(&sb, "%s=%v ", k, v)
			}
			attrsStr = strings.TrimSpace(sb.String())
			attrsStr = h.opts.AttrsColor.Sprint(attrsStr)
		} else {
			var b []byte
			var err error
			if h.opts.IndentJSON {
				b, err = json.MarshalIndent(fields, "", h.opts.IndentStr)
			} else {
				b, err = json.Marshal(fields)
			}
			if err != nil {
				return err
			}
			attrsStr = h.opts.AttrsColor.Sprint(string(b))
		}
	}

	// Add caller info if enabled
	var callerStr string
	if h.opts.ShowCaller {
		if pc := r.PC; pc != 0 {
			file, line := pcToFileLine(pc)
			if file != "" && line != 0 {
				callerStr = h.opts.CallerColor.Sprintf("%s:%d", file, line)
			}
		}
	}

	// Compose final output line
	var parts []string
	parts = append(parts, timeStr, levelStr, msgStr)
	if callerStr != "" {
		parts = append(parts, callerStr)
	}
	if attrsStr != "" {
		parts = append(parts, attrsStr)
	}

	h.l.Println(strings.Join(parts, " "))

	return nil
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if h.opts.SlogOpts != nil && h.opts.SlogOpts.Level != nil {
		return level >= h.opts.SlogOpts.Level.Level()
	}
	// By default enable all levels
	return true
}


func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Append to existing attrs slice (make new slice to avoid mutation)
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &PrettyHandler{
		opts:  h.opts,
		out:   h.out,
		l:     h.l,
		attrs: newAttrs,
		group: h.group,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	newGroup := name
	if h.group != "" {
		newGroup = h.group + "." + name
	}
	return &PrettyHandler{
		opts:  h.opts,
		out:   h.out,
		l:     h.l,
		attrs: h.attrs,
		group: newGroup,
	}
}
