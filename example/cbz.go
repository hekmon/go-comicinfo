package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hekmon/go-comicinfo"
)

func writeCBZChapter(chapter Chapter, outputDir string) (err error) {
	// Prepare file
	file, err := os.Create(filepath.Join(
		outputDir,
		fmt.Sprintf("%s.cbz", sanitizeFileName(chapter.FullTitle())),
	))
	if err != nil {
		return fmt.Errorf("failed to create CBZ file: %w", err)
	}
	defer file.Close()
	cbzWriter := zip.NewWriter(file)
	defer cbzWriter.Close()
	// Prepare ComicInfo.xml
	ci := comicinfo.ComicInfov2{
		Title:         chapter.FullTitle(),
		Series:        chapter.Serie.Title,
		Number:        chapter.Number,
		Count:         len(chapter.Serie.Chapters),
		Summary:       chapter.Serie.Summary,
		Year:          chapter.PublishDate.Year(),
		Month:         int(chapter.PublishDate.Month()),
		Day:           chapter.PublishDate.Day(),
		Publisher:     chapter.Serie.Publisher,
		Genre:         chapter.Serie.Genre,
		Web:           chapter.Serie.URL.String(),
		PageCount:     len(chapter.Pages),
		LanguageISO:   comicinfo.LanguageEnglish,
		Format:        "Web",
		BlackAndWhite: comicinfo.No,
		Manga:         comicinfo.MangaNo,
		Writer:        strings.Join(chapter.Serie.Creators, ","),
		Pages: comicinfo.PagesV2{
			Pages: make([]comicinfo.PageV2, len(chapter.Pages)+1), // +1 for cover
		},
	}
	// Write cover
	coverFilename := fmt.Sprintf("cover%s", chapter.Serie.Cover.Type.Extension())
	zipImgFile, err := cbzWriter.Create(coverFilename)
	if err != nil {
		return fmt.Errorf("failed to create image file in ZIP: %w", err)
	}
	if _, err = io.Copy(zipImgFile, bytes.NewReader(chapter.Serie.Cover.Data)); err != nil {
		return fmt.Errorf("failed to write image data to ZIP: %w", err)
	}
	coverImg, err := chapter.Serie.Cover.Decode()
	if err != nil {
		return fmt.Errorf("failed to decode cover: %w", err)
	}
	ci.Pages.Pages[0] = comicinfo.PageV2{
		Image:       0,
		Type:        comicinfo.PageTypeFrontCover,
		ImageSize:   len(chapter.Serie.Cover.Data),
		Key:         coverFilename,
		Bookmark:    "Cover",
		ImageWidth:  coverImg.Bounds().Dx(),
		ImageHeight: coverImg.Bounds().Dy(),
	}
	// Add images
	for i, page := range chapter.Pages {
		pageName := fmt.Sprintf("p%03d%s", i+1, page.Type.Extension())
		zipImgFile, err := cbzWriter.Create(pageName)
		if err != nil {
			return fmt.Errorf("failed to create image file in ZIP: %w", err)
		}
		if _, err = io.Copy(zipImgFile, bytes.NewReader(page.Data)); err != nil {
			return fmt.Errorf("failed to write image data to ZIP: %w", err)
		}
		img, err := page.Decode()
		if err != nil {
			return fmt.Errorf("can not decode image at page #%d: %w", i, err)
		}
		ci.Pages.Pages[i+1] = comicinfo.PageV2{
			Image:       i + 1,
			Type:        comicinfo.PageTypeStory,
			ImageSize:   len(page.Data),
			Key:         pageName,
			Bookmark:    fmt.Sprintf("Page %d", i+1),
			ImageWidth:  img.Bounds().Dx(),
			ImageHeight: img.Bounds().Dy(),
		}
	}
	// Write ComicInfo.xml within the zip
	ciWriter, err := cbzWriter.Create(comicinfo.ComicInfoFileName)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", comicinfo.ComicInfoFileName, err)
	}
	if err = ci.Encode(ciWriter); err != nil {
		return fmt.Errorf("failed to generate ComicInfo.xml: %w", err)
	}
	// Set ZIP comment before closing
	if err = cbzWriter.SetComment(chapter.FullTitle()); err != nil {
		return fmt.Errorf("failed to set ZIP file's comment: %w", err)
	}
	return nil
}

type Serie struct {
	Title     string
	Summary   string
	Genre     string
	Publisher string
	URL       *url.URL
	Creators  []string
	Cover     Image
	Chapters  []Chapter
}

type Chapter struct {
	Number      int
	Title       string
	PublishDate time.Time
	Author      string
	Serie       *Serie
	Pages       []Image
}

func (c *Chapter) FullTitle() string {
	return fmt.Sprintf("%s - Chapter %03d - %s", c.Serie.Title, c.Number, c.Title)
}

type Image struct {
	Data []byte
	Type ImageType
}

func (i *Image) Decode() (image image.Image, err error) {
	switch i.Type {
	case PNG:
		image, err = png.Decode(bytes.NewReader(i.Data))
	case JPEG:
		image, err = jpeg.Decode(bytes.NewReader(i.Data))
	default:
		err = fmt.Errorf("unsupported image type: %q", i.Type)
	}
	return
}

type ImageType int

const (
	PNG ImageType = iota
	JPEG
)

func (it ImageType) String() string {
	switch it {
	case PNG:
		return "PNG"
	case JPEG:
		return "JPEG"
	default:
		return "Unknown"
	}
}

func (it ImageType) Extension() string {
	switch it {
	case PNG:
		return ".png"
	case JPEG:
		return ".jpg"
	default:
		return ""
	}
}

func sanitizeFileName(fileName string) (sanitized string) {
	sanitized = fileName
	for _, c := range []rune{'\\', '/', ':', '*', '?', '"', '<', '>', '|'} {
		sanitized = strings.ReplaceAll(sanitized, string(c), "")
	}
	return
}
