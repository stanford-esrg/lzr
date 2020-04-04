package ipmi

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

// A ReadWriter that counts how much data is read / written on the underlying io.Reader/io.Writer.
type readWriteCounter struct {
	Writer     io.Writer
	Reader     io.Reader
	numWritten int
	numRead    int
}

// Write data to the Writer, and add the number of bytes written to numWritten
func (c *readWriteCounter) Write(data []byte) (int, error) {
	n, err := c.Writer.Write(data)
	c.numWritten += n
	return n, err
}

// Read from the Reader, adding the number of bytes actually read to numRead
func (c *readWriteCounter) Read(data []byte) (int, error) {
	n, err := c.Reader.Read(data)
	c.numRead += n
	return n, err
}

// Helper to do a single binary.Write and return the number of bytes written
func writeAndCount(dst io.Writer, data interface{}) (int, error) {
	temp := &readWriteCounter{Writer: dst, numRead: 0, numWritten: 0}
	err := binary.Write(temp, binary.BigEndian, data)
	return temp.numWritten, err
}

// Helper too write data to a writer while checking for short writes
func write(writer io.Writer, v []byte) (int, error) {
	n, err := writer.Write(v)
	if err != nil {
		return n, err
	}
	if n != len(v) {
		return n, io.ErrShortWrite
	}
	return n, nil
}

// Helper to create a readWriteCounter for a Reader
func countingReader(src io.Reader) *readWriteCounter {
	return &readWriteCounter{Reader: src}
}

// Helper to do a single binary.Read and count the number of bytes actually written
func readAndCount(src io.Reader, data interface{}) (int, error) {
	temp := countingReader(src)
	err := binary.Read(temp, binary.BigEndian, data)
	return temp.numRead, err
}

// Helper to read an object with binary.Read, and optionally confirm that the data read is the
// expected size. Use a negative size to skip checking.
func readObject(src io.Reader, v interface{}, expectedSize int) (int, error) {
	n, err := readAndCount(src, v)
	if err != nil {
		return n, err
	}
	if expectedSize >= 0 && n != expectedSize {
		return n, fmt.Errorf("short read: got %d, expected %d", n, expectedSize)
	}
	return n, err
}

// Helper to read exactly len(v) bytes
func read(src io.Reader, v []byte) (int, error) {
	n, err := io.ReadFull(src, v)
	if err != nil {
		return n, err
	}
	return n, nil
}

// Helper to get a JSON string suitable for logging
func toJSON(v interface{}) string {
	ret, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Errorf("JSON Marshal failed: %v", err)
		return ""
	}
	return string(ret)
}
