package game

// ChannelWriter implements io.WriterCloser in top of a string channel
type ChannelWriter struct {
	reader chan string
}

func NewChannelWriter() *ChannelWriter {
	return &ChannelWriter{make(chan string)}
}

func (w *ChannelWriter) Write(p []byte) (int, error) {
	w.reader <- string(p)

	return len(p), nil
}

func (w *ChannelWriter) Close() error {
	close(w.reader)
	return nil
}

func (w *ChannelWriter) Chan() <-chan string {
	return w.reader
}
