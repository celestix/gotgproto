package functions

import (
	"fmt"

	mtperrors "github.com/celestix/gotgproto/errors"
	"github.com/gotd/td/tg"
)

// GetMediaFileNameWithId
// Return media's filename in format "{id}-{name}.{extension}"
func GetMediaFileNameWithId(media tg.MessageMediaClass) (string, error) {
	switch v := media.(type) {
	case *tg.MessageMediaPhoto: // messageMediaPhoto#695150d7
		f, ok := v.Photo.AsNotEmpty()
		if !ok {
			return "", mtperrors.ErrUnknownTypeMedia
		}

		return fmt.Sprintf("%d.png", f.ID), nil
	case *tg.MessageMediaDocument: // messageMediaDocument#4cf4d72d
		var (
			attr             tg.DocumentAttributeClass
			ok               bool
			filenameFromAttr *tg.DocumentAttributeFilename
			f                *tg.Document
			filename         = "undefined"
		)

		f, ok = v.Document.AsNotEmpty()
		if !ok {
			return "", mtperrors.ErrUnknownTypeMedia
		}

		for _, attr = range f.Attributes {
			filenameFromAttr, ok = attr.(*tg.DocumentAttributeFilename)
			if ok {
				filename = filenameFromAttr.FileName
			}
		}

		return fmt.Sprintf("%d-%s", f.ID, filename), nil
	case *tg.MessageMediaStory: // messageMediaStory#68cb6283
		f, ok := v.Story.(*tg.StoryItem)
		if !ok {
			return "", mtperrors.ErrUnknownTypeMedia
		}
		return GetMediaFileNameWithId(f.Media)
	}
	return "", mtperrors.ErrUnknownTypeMedia
}

// GetMediaFileName
// Return media's filename in format "{name}.{extension}"
// Warning, stickers will always have name "sticker.webp", if you need distinction, use GetMediaFileNameWithId
func GetMediaFileName(media tg.MessageMediaClass) (string, error) {
	switch v := media.(type) {
	case *tg.MessageMediaPhoto: // messageMediaPhoto#695150d7
		f, ok := v.Photo.AsNotEmpty()
		if !ok {
			return "", mtperrors.ErrUnknownTypeMedia
		}

		return fmt.Sprintf("%d.png", f.ID), nil
	case *tg.MessageMediaDocument: // messageMediaDocument#4cf4d72d
		var (
			attr             tg.DocumentAttributeClass
			ok               bool
			filenameFromAttr *tg.DocumentAttributeFilename
			f                *tg.Document
			filename         = "undefined"
		)

		f, ok = v.Document.AsNotEmpty()
		if !ok {
			return "", mtperrors.ErrUnknownTypeMedia
		}

		for _, attr = range f.Attributes {
			filenameFromAttr, ok = attr.(*tg.DocumentAttributeFilename)
			if ok {
				filename = filenameFromAttr.FileName
			}
		}

		return filename, nil
	case *tg.MessageMediaStory: // messageMediaStory#68cb6283
		f, ok := v.Story.(*tg.StoryItem)
		if !ok {
			return "", mtperrors.ErrUnknownTypeMedia
		}
		return GetMediaFileName(f.Media)
	}
	return "", mtperrors.ErrUnknownTypeMedia
}

// GetInputFileLocation
// Returns tg.InputFileLocationClass, which can be used to download media
// used by ext.DownloadMedia()
func GetInputFileLocation(media tg.MessageMediaClass) (tg.InputFileLocationClass, error) {
	switch v := media.(type) {
	case *tg.MessageMediaPhoto: // messageMediaPhoto#695150d7
		f, ok := v.Photo.AsNotEmpty()
		if !ok {
			return nil, mtperrors.ErrUnknownTypeMedia
		}
		thumbSize := ""
		if len(f.Sizes) > 1 {
			// Lowest (f.Sizes[0]) size has the lowest resolution
			// Highest (f.Sizes[len(f.Sizes)-1]) has the highest resolution
			thumbSize = f.Sizes[len(f.Sizes)-1].GetType()
		}
		return &tg.InputPhotoFileLocation{
			ID:            f.ID,
			AccessHash:    f.AccessHash,
			FileReference: f.FileReference,
			ThumbSize:     thumbSize,
		}, nil
	case *tg.MessageMediaDocument: // messageMediaDocument#4cf4d72d
		f, ok := v.Document.AsNotEmpty()
		if !ok {
			return nil, mtperrors.ErrUnknownTypeMedia
		}
		return f.AsInputDocumentFileLocation(""), nil
	case *tg.MessageMediaStory: // messageMediaStory#68cb6283
		f, ok := v.Story.(*tg.StoryItem)
		if !ok {
			return nil, mtperrors.ErrUnknownTypeMedia
		}
		return GetInputFileLocation(f.Media)
	}
	return nil, mtperrors.ErrUnknownTypeMedia
}
