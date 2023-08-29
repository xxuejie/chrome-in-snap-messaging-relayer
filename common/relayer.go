package common

import (
	"io"
)

func RelayData(
	reader1 io.Reader,
	writer1 io.Writer,
	reader2 io.Reader,
	writer2 io.Writer,
) error {
	data1, error1 := channelFromReader(reader1)
	data2, error2 := channelFromReader(reader2)

	for {
		select {
		case d := <-data1:
			if err := writeAll(writer2, d); err != nil {
				return err
			}
		case d := <-data2:
			if err := writeAll(writer1, d); err != nil {
				return err
			}
		case err := <-error1:
			if err == io.EOF {
				return nil
			}
			return err
		case err := <-error2:
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func channelFromReader(reader io.Reader) (chan []byte, chan error) {
	d := make(chan []byte)
	e := make(chan error)

	go func() {
		for {
			buf := make([]byte, 1024)
			l, err := reader.Read(buf)
			if err != nil {
				e <- err
				return
			}
			d <- buf[0:l]
		}
	}()

	return d, e
}

func writeAll(writer io.Writer, data []byte) error {
	written := 0
	for written < len(data) {
		w, err := writer.Write(data[written:])
		if err != nil {
			return err
		}
		written += w
	}
	return nil
}
