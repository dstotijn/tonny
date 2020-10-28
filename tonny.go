package tonny

import (
	"bytes"
	"io"
	"net"
	"time"
)

// TeeListener implements net.Listener.
type TeeListener struct {
	ln net.Listener
}

// TeeConn implements net.Conn.
type TeeConn struct {
	// Buffer containing all data read from the underlying net.Conn.
	ReadBuffer *bytes.Buffer
	// Buffer containing all data written to the underlying net.Conn.
	WriteBuffer *bytes.Buffer

	conn   net.Conn
	reader io.Reader
	writer io.Writer
}

// Accept calls Accept on the underlying net.Listener.
func (teeLn TeeListener) Accept() (net.Conn, error) {
	conn, err := teeLn.ln.Accept()
	if err != nil {
		return nil, err
	}

	readBuf := &bytes.Buffer{}
	writeBuf := &bytes.Buffer{}

	TeeConn := TeeConn{
		ReadBuffer:  readBuf,
		WriteBuffer: writeBuf,
		conn:        conn,
		reader:      io.TeeReader(conn, readBuf),
		writer:      io.MultiWriter(conn, writeBuf),
	}

	return TeeConn, nil
}

// Close calls Close on the underlying net.Listener.
func (teeLn TeeListener) Close() error {
	return teeLn.ln.Close()
}

// Addr calls Addr on the underlying net.Listener.
func (teeLn TeeListener) Addr() net.Addr {
	return teeLn.ln.Addr()
}

// Read reads data from the underlying net.Conn into b, and also writes it into
// the read buffer.
func (tc TeeConn) Read(b []byte) (int, error) {
	return tc.reader.Read(b)
}

// Write does a write on the underlying net.Conn and the write buffer.
func (tc TeeConn) Write(b []byte) (int, error) {
	return tc.writer.Write(b)
}

// Close calls Close on the underlying net.Conn.
func (tc TeeConn) Close() error {
	return tc.conn.Close()
}

// LocalAddr calls LocalAddr on the underlying net.Conn.
func (tc TeeConn) LocalAddr() net.Addr {
	return tc.conn.LocalAddr()
}

// RemoteAddr calls RemoteAddr on the underlying net.Conn.
func (tc TeeConn) RemoteAddr() net.Addr {
	return tc.conn.RemoteAddr()
}

// SetDeadline calls SetDeadline on the underlying net.Conn.
func (tc TeeConn) SetDeadline(t time.Time) error {
	return tc.conn.SetDeadline(t)
}

// SetReadDeadline calls SetReadDeadline on the underlying net.Conn.
func (tc TeeConn) SetReadDeadline(t time.Time) error {
	return tc.conn.SetReadDeadline(t)
}

// SetWriteDeadline calls SetWriteDeadline on the underlying net.Conn.
func (tc TeeConn) SetWriteDeadline(t time.Time) error {
	return tc.conn.SetWriteDeadline(t)
}

// Listen wraps net.Listen, and returns a net.Listener to be used for teeing
// net.Conn traffic.
func Listen(network, address string) (TeeListener, error) {
	ln, err := net.Listen(network, address)
	if err != nil {
		return TeeListener{}, err
	}
	return TeeListener{ln: ln}, nil
}
