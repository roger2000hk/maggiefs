package dataserver

import (
	"fmt"
	"io"
	"net"
	"syscall"
	"github.com/jbooth/maggiefs/maggiefs"
)

// wraps ReadWriteCloser with Stringer for debugging info
type Endpoint interface {
	io.ReadWriteCloser
	fmt.Stringer
	Rfd() int // returns an FD that can be used to read from this Endpoint.  Negative FDs indicate you should use the Read() method instead.
	Wfd() int // returns an FD that can be used write to this interface.  Negative FDs indicate you should use the Write() method instead
}

type endpt struct {
	r      io.Reader
	w      io.Writer
	desc   string
	sameFD bool // whether they point to the same fd or not
	rfd int
	wfd int
}

func (e *endpt) Read(p []byte) (n int, err error)  { return e.r.Read(p) }
func (e *endpt) Write(p []byte) (n int, err error) { return e.w.Write(p) }
func (e *endpt) Rfd() int { return e.rfd }
func (e *endpt) Wfd() int { return e.wfd }
func (e *endpt) String() string                    { return e.desc }
func (e *endpt) Close() error {
	closer, ok := e.r.(io.Closer)
	var err error
	if ok {
		err = closer.Close()
	}
	if !e.sameFD {
		closer, ok = e.w.(io.Closer)
		if ok {
			err = closer.Close()
		}
	}
	return err
}

// wrap a socket with nice desc
func SockEndpoint(c *net.TCPConn) (Endpoint) {
	desc := fmt.Sprintf("Socket(nonblock): local %s to remote %s", c.LocalAddr(), c.RemoteAddr())
	c.SetNoDelay(true)
	c.SetReadBuffer(128 * 1024)
	c.SetWriteBuffer(128*1024)
	return &endpt{c, c, desc, true, -1, -1}
}

func BlockingSockEndPoint(c *net.TCPConn) (Endpoint, error) {
	desc := fmt.Sprintf("Socket(blocking): local %s to remote %s", c.LocalAddr(), c.RemoteAddr())
	err := c.SetNoDelay(true)
	if err != nil { return nil,err }
	err = c.SetReadBuffer(128 * 1024)
	if err != nil { return nil,err }
	err = c.SetWriteBuffer(128 * 1024)
	if err != nil { return nil,err }
	f, err := c.File()
	defer c.Close()
	if err != nil {
		return nil, err
	}
	syscall.SetNonblock(int(f.Fd()), false)
	fd := int(f.Fd())

	return &endpt{f, f, desc, true, fd, fd}, nil
}

// matching pipe endpoints
func PipeEndpoints() (Endpoint, Endpoint) {
	leftRead, rightWrite := io.Pipe()
	rightRead, leftWrite := io.Pipe()
	left := &endpt{leftRead, leftWrite, "Pipe", false, -1, -1 }
	right := &endpt{rightRead, rightWrite, "Pipe", false, -1, -1 }
	return left, right
}

func Copy(dst io.Writer, src io.Reader, n int64) (int64, error) {
	// If the writer has a ReadFrom method, use it to do the copy.
	// Avoids an allocation and a copy.
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}
	// Similarly, if the reader has a WriteTo method, use it to do the copy.
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}

	buff := maggiefs.GetBuff()
	defer maggiefs.ReturnBuff(buff)
	nWritten := int64(0)
	for nWritten < n {
		r, err := src.Read(buff)
		if r > 0 {
			w, e2 := dst.Write(buff[0:r])
			if w > 0 {
				nWritten += int64(w)
			}
			if e2 == io.EOF {
				fmt.Printf("Copy: EOF, read %d write %d nWritten %d out of %d\n", r, w, nWritten, n)

			}
			if e2 != nil && e2 != io.EOF {
				return nWritten, e2
			}
			if r != w {
				fmt.Printf("Copy: ErrShortWrite, read %d write %d nWritten %d out of %d\n", r, w, nWritten, n)
				return nWritten, io.ErrShortWrite
			}
		}
		if err == io.EOF {
			return nWritten, nil
		}
		if err != nil {
			return nWritten, err
		}
	}
	return nWritten, nil
}

func NewSectionWriter(w io.WriterAt, off int64, length int64) io.Writer {
	return &SectionWriter{w, off, off + length}
}

// SectionReader implements Write on a section
// of an underlying WriteAt.
type SectionWriter struct {
	w     io.WriterAt
	off   int64
	limit int64
}

func (s *SectionWriter) Write(p []byte) (n int, err error) {
	if s.off >= s.limit {
		return 0, io.EOF
	}
	if max := s.limit - s.off; int64(len(p)) > max {
		p = p[0:max]
	}
	n, err = s.w.WriteAt(p, s.off)
	s.off += int64(n)
	return
}
