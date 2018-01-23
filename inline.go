package csp

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"io"
	"strings"

	"github.com/MJKWoolnough/memio"
)

type tag struct {
	Type                                    byte
	StartTagStart, StartTagEnd, EndTagStart int64
}

func findTags(r io.Reader) ([]tag, error) {
	var (
		currTag tag
		tags    []tag
	)
	d := xml.NewDecoder(r)
	d.AutoClose = xml.HTMLAutoClose
	d.Entity = xml.HTMLEntity
	d.Strict = false
	for {
		pos := d.InputOffset()
		t, err := d.RawToken()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		switch t := t.(type); t {
		case xml.StartElement:
			var tagType byte
			switch s := strings.ToLower(t.Name.Local); s {
			case "style":
				tagType = 1
				fallthrough
			case "script":
				tag.Type = tagType
				tag.StartTagStart = pos
				tag.StartTagEnd = d.InputOffset()
			}
		case xml.EndElement:
			if s := strings.ToLower(t.Name.Local); s == tag.Name {
				tag.EndTagStart = pos
				tags = append(tags, currTag)
			}
		}
	}
	return tags, nil
}

func Inline(r io.Reader) ([]byte, error) {
	var buf memio.Buffer
	tags, err := findTags(io.TeeReade(r, &buf))
	if err != nil {
		return nil, err
	}
	var offset int64
	for _, t := range tags {
		hash := base64.StdEncoding.EncodeToString((sha256.Sum256(buf[t.StartTagEnd+offset : t.EndTagStart+offset]))[:])

	}
	return buf, nil
}
