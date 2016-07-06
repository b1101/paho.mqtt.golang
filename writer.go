package mqtt

type Writer struct {
	publish PublishFunc
}

type PublishFunc func(interface{}) error

func NewWriter(pub Publisher, qos byte, retained bool, topic string) (w *Writer) {
	return &Writer{
		publish: PublishFunc(func(data interface{}) error {
			token := pub.Publish(topic, qos, retained, data)
			token.Wait()
			return token.Error()
		}),
	}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	if err := w.publish(p); err != nil {
		return 0, err
	}
	return len(p), nil
}
