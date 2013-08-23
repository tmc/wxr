package wxr

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

var (
	testFileURL  = "http://wpcandy.s3.amazonaws.com/resources/postsxml.zip"
	testFileName = "posts.xml"
)

// dlZip extracts the content of the named file from the zip located at url
func dlZip(url, fileName string) (io.ReadCloser, error) {
	h, err := http.Get(url)
	defer h.Body.Close()
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(h.Body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(b)
	z, err := zip.NewReader(buf, h.ContentLength)
	if err != nil {
		return nil, err
	}
	for _, f := range z.File {
		if f.Name == fileName {
			return f.Open()
		}
	}
	return nil, fmt.Errorf("filename '%s' not found in zip", fileName)
}

func TestBasicUnmarshal(t *testing.T) {
	z, err := dlZip(testFileURL, testFileName)
	if err != nil {
		t.Error(err)
	}

	b, err := NewRSS(z)
	if err != nil {
		t.Error(err)
	}

	if len(b.Channel.Categories) != 8 {
		t.Error("expected 8 categories, got ",
			len(b.Channel.Categories))
	}

	if len(b.Channel.Items) != 54 {
		t.Error("expected 54 items, got ",
			len(b.Channel.Items))
	}
	e := Item{
		XMLName: xml.Name{Space: "", Local: "item"},
		ID:      4,
		Name:    "a-simple-post-with-text",
		Title:   "A Simple Post with Text",
		Link:    "http://dev.wpcoder.com/dan/wordpress/2008/08/a-simple-post-with-text/",
		Type:    "post",
		PubDate: b.Channel.Items[0].PubDate,
		Categories: []ItemCategory{
			{XMLName: xml.Name{Space: "", Local: "category"}, Slug: "", Name: "Child Category I"},
			{XMLName: xml.Name{Space: "", Local: "category"}, Slug: "child-category-i", Name: "Child Category I"},
			{XMLName: xml.Name{Space: "", Local: "category"}, Slug: "", Name: "Parent Category I"},
			{XMLName: xml.Name{Space: "", Local: "category"}, Slug: "parent-category-i", Name: "Parent Category I"},
			{XMLName: xml.Name{Space: "", Local: "category"}, Slug: "", Name: "tag1"},
			{XMLName: xml.Name{Space: "", Local: "category"}, Slug: "tag1", Name: "tag1"},
			{XMLName: xml.Name{Space: "", Local: "category"}, Slug: "", Name: "tag2"},
			{XMLName: xml.Name{Space: "", Local: "category"}, Slug: "tag2", Name: "tag2"},
			{XMLName: xml.Name{Space: "", Local: "category"}, Slug: "", Name: "tag5"},
			{XMLName: xml.Name{Space: "", Local: "category"}, Slug: "tag5", Name: "tag5"}},
	}

	if !reflect.DeepEqual(b.Channel.Items[0], e) {
		t.Errorf("example item did not parse correctly.\nexpected:\n%+v\n\nactual:\n%+v",
			e, b.Channel.Items[0])
	}
}
