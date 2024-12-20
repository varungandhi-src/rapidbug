package rapidbug

import (
	"io"
	"math"
)

type TextFileChecker struct {
	checker checker
}

// NO_LINE_LIMIT is a potential parameter for NewTextFileChecker.
const NO_LINE_LIMIT = LineCount(math.MaxInt)

// NewTextFileChecker creates a new TextFileChecker
// which will read from the given reader in a streaming
// fashion and enforce the provided size and line count
// limits.
//
// If you want the output to be copied to a buffer,
// you can use the following idiom:
//
//	var buf bytes.Buffer
//	teeReader := io.TeeReader(baseReader, &buf)
//	checker := limits.NewTextFileChecker(teeReader, ...)
//
// If you don't want to enforce a line limit, pass NO_LINE_LIMIT;
// don't use BinaryFileChecker.
//
// Pre-condition: sizeLimit and lineLimit must be positive.
func NewTextFileChecker(reader io.Reader, sizeLimit Size, lineLimit LineCount) TextFileChecker {
	return TextFileChecker{newChecker(reader, sizeLimit, lineLimit, true)}
}

type TextFileErrorKind string

const (
	TextFileExceededByteSizeLimit  TextFileErrorKind = "exceeded byte size limit"
	TextFileExceededLineCountLimit TextFileErrorKind = "exceeded line count limit"
	TextFileReadError              TextFileErrorKind = "read error"
)

type TextFileLimitError struct {
	Kind TextFileErrorKind
	// BytesRead may be lower than the total size of the file,
	// but will be higher than the limit passed during initialization.
	BytesRead Size
	// LinesRead may be lower than the total number of lines in the file,
	// but will be higher than the limit passed during initialization.
	LinesRead LineCount
	IOError   error
}

func (t *TextFileChecker) TryReadAll() (TextFileSize, *TextFileLimitError) {
	return t.checker.tryReadAll()
}

func newChecker(r io.Reader, sizeLimit Size, lineLimit LineCount, needLines bool) checker {
	return checker{make([]byte, 16*1024), r, sizeLimit, lineLimit, needLines, 0}
}

type checker struct {
	scratchBuf []byte
	reader     io.Reader
	sizeLimit  Size
	lineLimit  LineCount
	needLines  bool
	lastByte   byte
}

func (c checker) tryReadAll() (TextFileSize, *TextFileLimitError) {
	soFar := TextFileSize{Size(0), LineCount(0)}
	makeError := func(k TextFileErrorKind, err error) *TextFileLimitError {
		return &TextFileLimitError{k, soFar.Size, soFar.LineCount, err}
	}
	for {
		n, readErr := c.reader.Read(c.scratchBuf)
		var err error
		if readErr != io.EOF {
			err = readErr
		}
		// We handle n first because of this in the docs for io.Reader.Read:
		// > Callers should always process the n > 0 bytes returned before
		// > considering the error err.
		if n > 0 {
			c.lastByte = c.scratchBuf[n-1]
			// It is sufficient to check writeErr instead of needing some
			// call like io.WriteFull (non-existent API) because of this in
			// the io.Writer.Write docs:
			// > Write must return a non-nil error if it returns n < len(p).
			soFar.Size += Size(n)
			if soFar.Size > c.sizeLimit {
				return soFar, makeError(TextFileExceededByteSizeLimit, err)
			}
			if c.needLines {
				for i := range n { // Not using i, c := range syntax as we don't need full runes
					if c.scratchBuf[i] == '\n' {
						soFar.LineCount++
					}
				}
				if soFar.LineCount > c.lineLimit {
					return soFar, makeError(TextFileExceededLineCountLimit, err)
				}
			}
			if err != nil {
				return soFar, makeError(TextFileReadError, err)
			}
		}
		// NOTE: n >= 0 here.
		if readErr == nil || readErr == io.EOF {
			if readErr == io.EOF {
				if c.lastByte != '\n' && soFar.Size > 0 {
					soFar.LineCount++
					if soFar.LineCount > c.lineLimit {
						return soFar, makeError(TextFileExceededLineCountLimit, nil)
					}
				}
				return soFar, nil
			}
			// The io.Reader.Read docs state:
			//
			// > Callers should treat a return of 0 and nil as indicating that
			// > nothing happened; in particular it does not indicate EOF.
			//
			// So treat n == 0 and n > 0 uniformly by continuing to the next iteration.
			continue
		}
		return soFar, makeError(TextFileReadError, err)
	}
}
