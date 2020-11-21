package ogg

import (
	"bytes"
	b "encoding/binary"
	"fmt"
	"io"
	"log"
)

var (
	oggHeader = [4]byte{'O', 'g', 'g', 'S'}

	// codecs
	codecHeaderOpus   = []byte("OpusHead")
	codecHeaderVorbis = []byte("vorbis")
)

// Segment is for ogg page segment
type Segment []byte

// Metadata is ogg page metadata item
type Metadata struct {
	CapturePattern        [4]byte
	Version               int8
	HeaderType            int8
	GranulePosition       int64
	BitStreamSerialNumber int32
	PageSequenceNumber    int32
	Checksum              int32
	PageSegments          uint8
}

// Page keeps ogg page info
type Page struct {
	Metadata *Metadata
	Segments []Segment
}

// Read reads binary chuck from the reader into metadata fields
func (m *Metadata) Read(reader io.Reader) (err error) {
	err = b.Read(reader, b.LittleEndian, m)
	return
}

func (p *Page) readSegments(reader io.Reader) (err error) {
	var (
		i     uint8
		in    []byte
		sizes []uint8
	)

	in = make([]byte, 1)
	for i = 0; i < p.Metadata.PageSegments; i++ {
		if _, err = reader.Read(in); err != nil {
			return
		}
		sizes = append(sizes, in[0])
	}

	for _, length := range sizes {
		buf := make([]byte, length)
		if _, err = reader.Read(buf); err != nil {
			return
		}
		p.Segments = append(p.Segments, buf)
	}
	return
}

// ReadPage reads a page of OGG file.
func ReadPage(reader io.Reader) (page *Page, err error) {
	metadata := Metadata{}
	page = new(Page)
	page.Metadata = &metadata

	if err = metadata.Read(reader); err != nil {
		return
	}
	if metadata.CapturePattern != oggHeader {
		return nil, fmt.Errorf("not an ogg file: %s", metadata.CapturePattern)
	}
	err = page.readSegments(reader)
	return
}

// IsOggOpusFile detects if given content:
//   - file is an ogg file
//   - ogg is encoded with opus codec
func IsOggOpusFile(reader io.Reader) bool {
	page, err := ReadPage(reader)
	if err == nil {
		for _, segment := range page.Segments {
			if bytes.Equal(segment[:8], codecHeaderOpus) {
				return true
			}
			// return false instantly without going forward
			if bytes.Equal(segment[1:7], codecHeaderVorbis) {
				log.Printf("vorbis codec detected, please recode file with opus\n")
				return false
			}
		}
	}
	return false
}
