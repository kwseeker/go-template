package buf

import (
	"io"
	"log"
	"os"
	"runtime"
	"syscall"
)

// 是否使用 readv 系统调用
var useReadv = false

func init() {
	const v2rayReadvEnv = "V2RAY_BUF_READV"
	envVal := "auto"

	v, found := os.LookupEnv(v2rayReadvEnv)
	if found {
		envVal = v
	}

	switch envVal {
	case "enable":
		useReadv = true
	case "auto":
		if (runtime.GOARCH == "386" || runtime.GOARCH == "amd64" || runtime.GOARCH == "s390x") && (runtime.GOOS == "linux" || runtime.GOOS == "darwin" || runtime.GOOS == "windows") {
			useReadv = true
		}
	}
}

func ReadBuffer(r io.Reader) (*Buffer, error) {
	b := New()
	n, err := b.ReadFrom(r)
	if n > 0 {
		return b, err
	}
	b.Release()
	return nil, err
}

type Reader interface {
	ReadMultiBuffer() (MultiBuffer, error)
}

func NewReader(reader io.Reader) Reader {
	if mr, ok := reader.(Reader); ok {
		return mr
	}

	if _, isFile := reader.(*os.File); isFile && useReadv {
		if sc, ok := reader.(syscall.Conn); ok {
			rawConn, err := sc.SyscallConn()
			if err != nil {
				log.Println("failed to get syscallConn")
			} else {
				return NewReadVReader(reader, rawConn)
			}
		}
	}

	return &SingleReader{
		Reader: reader,
	}
}

type SingleReader struct {
	io.Reader
}

// ReadMultiBuffer implements Reader.
func (r *SingleReader) ReadMultiBuffer() (MultiBuffer, error) {
	b, err := ReadBuffer(r.Reader)
	return MultiBuffer{b}, err
}

type multiReader interface {
	Init([]*Buffer)
	Read(fd uintptr) int32
	Clear()
}

type allocStrategy struct {
	current uint32
}

func (s *allocStrategy) Current() uint32 {
	return s.current
}

func (s *allocStrategy) Adjust(n uint32) {
	if n >= s.current {
		s.current *= 4
	} else {
		s.current = n
	}

	if s.current > 32 {
		s.current = 32
	}

	if s.current == 0 {
		s.current = 1
	}
}

func (s *allocStrategy) Alloc() []*Buffer {
	bs := make([]*Buffer, s.current)
	for i := range bs {
		bs[i] = New()
	}
	return bs
}

type ReadVReader struct {
	io.Reader
	rawConn syscall.RawConn
	mr      multiReader
	alloc   allocStrategy
}

func NewReadVReader(reader io.Reader, rawConn syscall.RawConn) *ReadVReader {
	return &ReadVReader{
		Reader:  reader,
		rawConn: rawConn,
		alloc: allocStrategy{
			current: 1,
		},
		mr: newMultiReader(),
	}
}

func (r *ReadVReader) readMulti() (MultiBuffer, error) {
	bs := r.alloc.Alloc()

	r.mr.Init(bs)
	var nBytes int32
	err := r.rawConn.Read(func(fd uintptr) bool {
		n := r.mr.Read(fd)
		if n < 0 {
			return false
		}

		nBytes = n
		return true
	})
	r.mr.Clear()

	if err != nil {
		ReleaseMulti(MultiBuffer(bs))
		return nil, err
	}

	if nBytes == 0 {
		ReleaseMulti(MultiBuffer(bs))
		return nil, io.EOF
	}

	nBuf := 0
	for nBuf < len(bs) {
		if nBytes <= 0 {
			break
		}
		end := nBytes
		if end > Size {
			end = Size
		}
		bs[nBuf].end = end
		nBytes -= end
		nBuf++
	}

	for i := nBuf; i < len(bs); i++ {
		bs[i].Release()
		bs[i] = nil
	}

	return MultiBuffer(bs[:nBuf]), nil
}

// ReadMultiBuffer implements Reader.
func (r *ReadVReader) ReadMultiBuffer() (MultiBuffer, error) {
	if r.alloc.Current() == 1 {
		b, err := ReadBuffer(r.Reader)
		if b.IsFull() {
			r.alloc.Adjust(1)
		}
		return MultiBuffer{b}, err
	}

	mb, err := r.readMulti()
	if err != nil {
		return nil, err
	}
	r.alloc.Adjust(uint32(len(mb)))
	return mb, nil
}

// BufferedReader 自带 Buffer 缓冲的 Reader
type BufferedReader struct {
	Reader Reader
	Buffer MultiBuffer
	// Splitter is a function to read bytes from MultiBuffer
	Splitter func(MultiBuffer, []byte) (MultiBuffer, int)
}

// BufferedBytes returns the number of bytes that is cached in this reader.
func (r *BufferedReader) BufferedBytes() int32 {
	return r.Buffer.Len()
}

// ReadByte implements io.ByteReader.
func (r *BufferedReader) ReadByte() (byte, error) {
	var b [1]byte
	_, err := r.Read(b[:])
	return b[0], err
}

// Read implements io.Reader. It reads from internal buffer first (if available) and then reads from the underlying reader.
func (r *BufferedReader) Read(b []byte) (int, error) {
	splitter := r.Splitter
	if splitter == nil {
		splitter = SplitBytes
	}

	if !r.Buffer.IsEmpty() {
		buffer, nBytes := splitter(r.Buffer, b)
		r.Buffer = buffer
		if r.Buffer.IsEmpty() {
			r.Buffer = nil
		}
		return nBytes, nil
	}

	mb, err := r.Reader.ReadMultiBuffer()
	if err != nil {
		return 0, err
	}

	mb, nBytes := splitter(mb, b)
	if !mb.IsEmpty() {
		r.Buffer = mb
	}
	return nBytes, nil
}

// ReadMultiBuffer implements Reader.
func (r *BufferedReader) ReadMultiBuffer() (MultiBuffer, error) {
	if !r.Buffer.IsEmpty() {
		mb := r.Buffer
		r.Buffer = nil
		return mb, nil
	}

	return r.Reader.ReadMultiBuffer()
}

// ReadAtMost returns a MultiBuffer with at most size.
func (r *BufferedReader) ReadAtMost(size int32) (MultiBuffer, error) {
	if r.Buffer.IsEmpty() {
		mb, err := r.Reader.ReadMultiBuffer()
		if mb.IsEmpty() && err != nil {
			return nil, err
		}
		r.Buffer = mb
	}

	rb, mb := SplitSize(r.Buffer, size)
	r.Buffer = rb
	if r.Buffer.IsEmpty() {
		r.Buffer = nil
	}
	return mb, nil
}

// TODO
//func (r *BufferedReader) writeToInternal(writer io.Writer) (int64, error) {
//	mbWriter := NewWriter(writer)
//	var sc SizeCounter
//	if r.Buffer != nil {
//		sc.Size = int64(r.Buffer.Len())
//		if err := mbWriter.WriteMultiBuffer(r.Buffer); err != nil {
//			return 0, err
//		}
//		r.Buffer = nil
//	}
//
//	err := Copy(r.Reader, mbWriter, CountSize(&sc))
//	return sc.Size, err
//}
//
//// WriteTo implements io.WriterTo.
//func (r *BufferedReader) WriteTo(writer io.Writer) (int64, error) {
//	nBytes, err := r.writeToInternal(writer)
//	if errors.Cause(err) == io.EOF {
//		return nBytes, nil
//	}
//	return nBytes, err
//}
//
//// Interrupt implements common.Interruptible.
//func (r *BufferedReader) Interrupt() {
//	common.Interrupt(r.Reader)
//}
//
//// Close implements io.Closer.
//func (r *BufferedReader) Close() error {
//	return common.Close(r.Reader)
//}
