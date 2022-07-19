package game

import log "github.com/sirupsen/logrus"

// LogWriter implements io.WriterCloser in top of a string channel
type LogWriter struct {
}

func NewChannelWriter() *LogWriter {
	return &LogWriter{}
}

func (w *LogWriter) Write(p []byte) (int, error) {
	log.Info(string(p))

	return len(p), nil
}
