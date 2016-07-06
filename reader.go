package mqtt

import "io"

type Reader struct {
	unsubscribe UnsubscribeFunc
	reader      *io.PipeReader
}

type UnsubscribeFunc func() error

func NewReader(sub Subscriber, qos byte, topics ...string) (r *Reader, err error) {
	reader, writer := io.Pipe()

	r = &Reader{
		reader:      reader,
		unsubscribe: unsubscribeFunc(sub, topics...),
	}

	if err := subscribe(sub, qos, payloadHandler(writer), topics...); err != nil {
		r.unsubscribe()
		return nil, err
	}

	return r, nil
}

func subscribe(sub Subscriber, qos byte, handler MessageHandler, topics ...string) (err error) {
	subscribed := make([]string, 0, len(topics))

	for _, topic := range topics {
		token := sub.Subscribe(topic, qos, handler)

		token.Wait()
		if err = token.Error(); err != nil {
			return err
		}

		subscribed = append(subscribed, topic)
	}

	return nil
}

func payloadHandler(w io.Writer) MessageHandler {
	return func(cl *Client, msg Message) {
		// TODO error handling?
		if _, err := w.Write(msg.Payload()); err != nil {
			panic(err)
		}
	}

}

func unsubscribeFunc(sub Subscriber, topics ...string) UnsubscribeFunc {
	return func() error {
		token := sub.Unsubscribe(topics...)
		token.Wait()
		return token.Error()
	}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

func (r *Reader) Close() (err error) {
	defer r.Close()
	return r.unsubscribe()
}
