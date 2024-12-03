package listener

import (
	"net"
	"sync"
	"syscall"
	"time"
)

// CustomConn represents a TCP connection
type CustomConn struct {
	fd         int
	localAddr  net.Addr
	remoteAddr net.Addr
	closed     bool
	closeMutex sync.RWMutex
}

// NewCustomConn creates a new custom connection
// @toDo: Add support for custom connection options
func NewCustomConn(fd int, remoteAddr, localAddr net.Addr) *CustomConn {
	return &CustomConn{
		fd:         fd,
		remoteAddr: remoteAddr,
		localAddr:  localAddr,
	}
}

// Read implements the net.Conn Read method
func (c *CustomConn) Read(b []byte) (n int, err error) {
	return syscall.Read(c.fd, b)
}

// Write implements the net.Conn Write method
func (c *CustomConn) Write(b []byte) (n int, err error) {
	return syscall.Write(c.fd, b)
}

// Close implements the net.Conn Close method
func (c *CustomConn) Close() error {
	c.closeMutex.Lock()
	defer c.closeMutex.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	return syscall.Close(c.fd)
}

// LocalAddr returns the local network address
func (c *CustomConn) LocalAddr() net.Addr {
	return c.localAddr
}

// RemoteAddr returns the remote network address
func (c *CustomConn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

// SetDeadline implements the net.Conn SetDeadline method
func (c *CustomConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline implements the net.Conn SetReadDeadline method
func (c *CustomConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline implements the net.Conn SetWriteDeadline method
func (c *CustomConn) SetWriteDeadline(t time.Time) error {
	return nil
}
