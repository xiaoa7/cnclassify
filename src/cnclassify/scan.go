//自己重写的buio.Scanner
//支持返回断词
package cnclassify

import (
	"errors"
	"io"
	"unicode/utf8"
)

type Scanner struct {
	r            io.Reader // The reader provided by the client.
	split        SplitFunc // The function to split the tokens.
	maxTokenSize int       // Maximum size of a token; modified by tests.
	token        []byte    // Last token returned by split.
	buf          []byte    // Buffer used as argument to split.
	start        int       // First non-processed byte in buf.
	end          int       // End of data in buf.
	err          error     // Sticky error.
	empties      int       // Count of successive empty tokens.
	stopchar     rune      //
}

type SplitFunc func(data []byte, atEOF bool) (advance int, r rune, token []byte, err error)

var (
	ErrTooLong         = errors.New("bufio.Scanner: token too long")
	ErrNegativeAdvance = errors.New("bufio.Scanner: SplitFunc returns negative advance count")
	ErrAdvanceTooFar   = errors.New("bufio.Scanner: SplitFunc returns advance count beyond input")
)

const (
	MaxScanTokenSize         = 64 * 1024
	minReadBufferSize        = 16
	maxConsecutiveEmptyReads = 100
)

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r:            r,
		maxTokenSize: MaxScanTokenSize,
		buf:          make([]byte, 4096), // Plausible starting size; needn't be large.
	}
}

func (s *Scanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}
func (s *Scanner) Bytes() []byte {
	return s.token
}

func (s *Scanner) Text() string {
	return string(s.token)
}
func (s *Scanner) Stopchar() rune {
	return s.stopchar
}
func (s *Scanner) Scan() bool {
	for {
		if s.end > s.start || s.err != nil {
			advance, r, token, err := s.split(s.buf[s.start:s.end], s.err != nil)
			if err != nil {
				s.setErr(err)
				return false
			}
			if !s.advance(advance) {
				return false
			}
			s.token = token
			s.stopchar = r
			if token != nil {
				if s.err == nil || advance > 0 {
					s.empties = 0
				} else {
					// Returning tokens not advancing input at EOF.
					s.empties++
					if s.empties > 100 {
						panic("bufio.Scan: 100 empty tokens without progressing")
					}
				}
				return true
			}
		}
		if s.err != nil {
			s.start = 0
			s.end = 0
			return false
		}
		if s.start > 0 && (s.end == len(s.buf) || s.start > len(s.buf)/2) {
			copy(s.buf, s.buf[s.start:s.end])
			s.end -= s.start
			s.start = 0
		}
		if s.end == len(s.buf) {
			if len(s.buf) >= s.maxTokenSize {
				s.setErr(ErrTooLong)
				return false
			}
			newSize := len(s.buf) * 2
			if newSize > s.maxTokenSize {
				newSize = s.maxTokenSize
			}
			newBuf := make([]byte, newSize)
			copy(newBuf, s.buf[s.start:s.end])
			s.buf = newBuf
			s.end -= s.start
			s.start = 0
			continue
		}
		for loop := 0; ; {
			n, err := s.r.Read(s.buf[s.end:len(s.buf)])
			s.end += n
			if err != nil {
				s.setErr(err)
				break
			}
			if n > 0 {
				s.empties = 0
				break
			}
			loop++
			if loop > maxConsecutiveEmptyReads {
				s.setErr(io.ErrNoProgress)
				break
			}
		}
	}
}

// advance consumes n bytes of the buffer. It reports whether the advance was legal.
func (s *Scanner) advance(n int) bool {
	if n < 0 {
		s.setErr(ErrNegativeAdvance)
		return false
	}
	if n > s.end-s.start {
		s.setErr(ErrAdvanceTooFar)
		return false
	}
	s.start += n
	return true
}

// setErr records the first error encountered.
func (s *Scanner) setErr(err error) {
	if s.err == nil || s.err == io.EOF {
		s.err = err
	}
}

// Split sets the split function for the Scanner. If called, it must be
// called before Scan. The default split function is ScanLines.
func (s *Scanner) Split(split SplitFunc) {
	s.split = split
}

//停词
func isStopChar(r rune) bool {
	switch r {
	case 0x2b /*+*/, 0x2d /*-*/, 0x5e /*^*/, 0x28 /*(*/, 0x29 /*)*/, 0x7c /*|*/ :
		return true
	default:
		return false
	}
}

//分割函数
func MyStopWord(data []byte, atEOF bool) (advance int, r rune, token []byte, err error) {
	start := 0
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if isStopChar(r) {
			return i + width, r, data[start:i], nil
		}
	}
	if atEOF && len(data) > start {
		return len(data), 0, data[start:], nil
	}
	return start, 0, nil, nil
}
