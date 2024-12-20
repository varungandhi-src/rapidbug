package rapidbug

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	genslices "github.com/life4/genesis/slices"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

type FileContents struct {
	contents string
	lines    []string
}

func fileGen() *rapid.Generator[FileContents] {
	return rapid.Custom(func(t *rapid.T) FileContents {
		charGen := rapid.RuneFrom([]rune("abcdef\n"))
		contents := rapid.StringOfN(charGen, 0, 20, 20).Draw(t, "contents")
		if len(contents) > 0 && contents[len(contents)-1] != '\n' && rapid.Bool().Draw(t, "add final new line") {
			contents += "\n"
		}
		lines := strings.Split(contents, "\n")
		if len(contents) > 0 && contents[len(contents)-1] == '\n' {
			lines = lines[:len(lines)-1]
		}
		return FileContents{contents, lines}
	})
}

func TestCheckers(t *testing.T) {
	t.Parallel()

	t.Run("property-based tests", func(t *testing.T) {
		t.Parallel()
		rapid.Check(t, func(t *rapid.T) {
			//require.NotPanics(t, func() {
			checkerPBT(t, fileGen())
			//})
		})
	})
}

func checkerPBT(t *rapid.T, fg *rapid.Generator[FileContents]) {
	file := fg.Draw(t, "file")

	lo := max(1, len(file.lines)-2)
	hi := len(file.lines) + 2
	t.Logf("lo=%d, hi=%d", lo, hi)
	lineCountLimit := LineCount(rapid.IntRange(lo, hi).Draw(t, "line count limit"))

	lo = max(1, len(file.contents)-10)
	hi = len(file.contents) + 10
	sizeLimit := Size(rapid.IntRange(lo, hi).Draw(t, "file size limit"))

	sizeOverLimit := len(file.contents) > int(sizeLimit)
	linesOverLimit := len(file.lines) > int(lineCountLimit)
	_, _ = sizeOverLimit, linesOverLimit

	events := []ReadEvent{}
	haveErr := false
	readFailErr := errors.New("Failed to read data")
	_, _ = haveErr, readFailErr
	if len(file.contents) != 0 {
		idxGen := rapid.IntRange(0, len(file.contents)-1)
		idxs := rapid.SliceOfNDistinct(idxGen, 1, min(5, len(file.contents)), func(i int) int { return i }).Draw(t, "idxs")
		idxs = append(idxs, 0, len(file.contents))
		genslices.Sort(idxs)
		genslices.Dedup(idxs)
		haveErr = rapid.Float64Range(0, 1.0).Draw(t, "error chance") < 0.1
		errEventIdx := -1
		if haveErr {
			errEventIdx = rapid.IntRange(0, len(idxs)-2).Draw(t, "error event idx")
		}
		for i := range len(idxs) - 1 {
			data := file.contents[idxs[i]:idxs[i+1]]
			var err error
			if haveErr && i == errEventIdx {
				err = readFailErr
			}
			events = append(events, ReadEvent{data, err})
		}
		var allData strings.Builder
		for _, ev := range events {
			allData.WriteString(ev.Data)
		}
	}

	{
		reader := NewFakeReader(events)
		var buf bytes.Buffer
		teeReader := io.TeeReader(&reader, &buf)
		checker := NewTextFileChecker(teeReader, sizeLimit, lineCountLimit)
		gotStats, gotErr := checker.TryReadAll()

		if gotErr == nil {
			require.Equal(t, len(file.lines), int(gotStats.LineCount), "line count mismatch")
		}
	}
}
