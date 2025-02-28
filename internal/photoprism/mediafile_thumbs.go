package photoprism

import (
	"fmt"
	"image"
	"time"

	"github.com/disintegration/imaging"
	"github.com/dustin/go-humanize/english"

	"github.com/photoprism/photoprism/internal/thumb"
	"github.com/photoprism/photoprism/pkg/capture"
	"github.com/photoprism/photoprism/pkg/clean"
	"github.com/photoprism/photoprism/pkg/fs"
)

// Bounds returns the media dimensions as image.Rectangle.
func (m *MediaFile) Bounds() image.Rectangle {
	return image.Rectangle{Min: image.Point{}, Max: image.Point{X: m.Width(), Y: m.Height()}}
}

// Thumbnail returns a thumbnail filename.
func (m *MediaFile) Thumbnail(path string, sizeName thumb.Name) (filename string, err error) {
	size, ok := thumb.Sizes[sizeName]

	if !ok {
		log.Errorf("media: invalid type %s", sizeName)
		return "", fmt.Errorf("media: invalid type %s", sizeName)
	}

	// Choose the smallest fitting size if the original image is smaller.
	if size.Fit && m.Bounds().In(size.Bounds()) {
		size = thumb.FitBounds(m.Bounds())
		log.Tracef("media: smallest fitting size for %s is %s (width %d, height %d)", clean.Log(m.RootRelName()), size.Name, size.Width, size.Height)
	}

	thumbName, err := size.FromFile(m.FileName(), m.Hash(), path, m.Orientation())

	if err != nil {
		err = fmt.Errorf("media: failed creating thumbnail for %s (%s)", clean.Log(m.BaseName()), err)
		log.Debug(err)
		return "", err
	}

	return thumbName, nil
}

// Resample returns a resampled image of the file.
func (m *MediaFile) Resample(path string, sizeName thumb.Name) (img image.Image, err error) {
	thumbName, err := m.Thumbnail(path, sizeName)

	if err != nil {
		return nil, err
	}

	return imaging.Open(thumbName)
}

// CreateThumbnails creates the default thumbnail sizes if the media file
// is a JPEG and they don't exist yet (except force is true).
func (m *MediaFile) CreateThumbnails(thumbPath string, force bool) (err error) {
	if !m.IsJpeg() {
		// Skip.
		return
	}

	count := 0
	start := time.Now()

	defer func() {
		switch count {
		case 0:
			log.Debug(capture.Time(start, fmt.Sprintf("media: created no new thumbnails for %s", clean.Log(m.RootRelName()))))
		default:
			log.Info(capture.Time(start, fmt.Sprintf("media: created %s for %s", english.Plural(count, "thumbnail", "thumbnails"), clean.Log(m.RootRelName()))))
		}
	}()

	hash := m.Hash()

	var original image.Image

	var srcImg image.Image
	var srcName thumb.Name

	for _, name := range thumb.Names {
		var size thumb.Size
		var fileName string

		if size = thumb.Sizes[name]; size.Uncached() {
			// Skip, exceeds pre-cached size limit.
			continue
		} else if fileName, err = size.FileName(hash, thumbPath); err != nil {
			log.Errorf("media: failed creating %s (%s)", clean.Log(string(name)), err)
			return err
		} else if force || !fs.FileExists(fileName) {
			// Open original if needed.
			if original == nil {
				img, err := thumb.Open(m.FileName(), m.Orientation())

				if err != nil {
					log.Debugf("media: %s in %s", err.Error(), clean.Log(m.RootRelName()))
					return err
				}

				original = img

				log.Debugf("media: opened %s [%s]", clean.Log(m.RootRelName()), thumb.MemSize(original).String())
			}

			// Thumb size too large
			// for the original image?
			if size.Skip(original) {
				continue
			}

			// Reuse existing thumb to improve performance
			// and reduce server load?
			if size.Source != "" {
				if size.Source == srcName && srcImg != nil {
					_, err = size.Create(srcImg, fileName)
				} else {
					_, err = size.Create(original, fileName)
				}
			} else {
				srcImg, err = size.Create(original, fileName)
				srcName = name
			}

			// RunFailed?
			if err != nil {
				log.Errorf("media: failed creating %s (%s)", name.String(), err)
				return err
			}

			count++
		}
	}

	return nil
}
