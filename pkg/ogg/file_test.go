package ogg

import (
	"errors"
	"testing"

	"github.com/tarvitz/dnq/pkg/tests"
)

func TestIsOpusOggFile(t *testing.T) {
	for _, entry := range []struct {
		name     string
		file     string
		expected bool
	}{
		{"is", testFileOggOpusHeadersOnly, true},
		{"is-not/vorbis", testFileOggVorbisHeadersOnly, false},
		{"is-not/bad-ogg", testFileBadOgg, false},
	} {
		t.Run(entry.name, func(in *testing.T) {
			fd, rollback := tests.OpenFile(entry.file)
			defer rollback()

			if result := IsOggOpusFile(fd); result != entry.expected {
				in.Errorf("expected: %v, got: %v", entry.expected, result)
			}
		})
	}
}

func TestReadPage(t *testing.T) {
	for _, entry := range []struct {
		name string
		file string
		err  error
	}{
		{"ok", testFileOggOpusHeadersOnly, nil},
		{"failure/not-ogg-file", testFilePlainText, errors.New("not an ogg file: this")},
		{"failure/bad-file", testFileBadOgg, errors.New("unexpected EOF")},
		{"failure/eof-in-segment-sizes", testFileOggOpusBrokenSegmentSizes, errors.New("EOF")},
		{"failure/eof-in-segments", testFileOggOpusBrokenSegments, errors.New("EOF")},
	} {
		t.Run(entry.name, func(in *testing.T) {
			fd, rollback := tests.OpenFile(entry.file)
			defer rollback()

			_, err := ReadPage(fd)
			if err != nil && entry.err != nil {
				if err.Error() != entry.err.Error() {
					in.Errorf("\nexp: `%v`\ngot: `%v`", entry.err, err)
				}
			}
			if entry.err == nil && err != nil {
				in.Errorf("got error: %v", err)
			}
		})
	}
}
