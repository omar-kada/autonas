package socket

import (
	"log/slog"
	"omar-kada/air-compose/api"
)

// MapLogLine maps a logs.Line to an api.LogLine
func MapLogLine(line slog.Record) api.LogLine {
	log := api.LogLine{
		Msg:   line.Message,
		Level: api.Level(line.Level.String()),
		Time:  line.Time,
		Meta:  make(map[string]string, line.NumAttrs()),
	}
	line.Attrs(func(a slog.Attr) bool {
		log.Meta[a.Key] = a.Value.String()
		return true
	})
	return log
}
