package types

import (
	"github.com/gotd/td/tg"
)

type Message struct {
	*tg.Message
}

//func (m *Message)

//func (m *Message) DownloadPhoto(ctx context.Context, client *tg.Client) {
//	photo := m.Media.(*tg.MessageMediaPhoto)
//	d, ok := photo.Photo.(*tg.Photo)
//	if ok {
//		file, err := client.UploadGetFile(ctx, &tg.UploadGetFileRequest{
//			Location: &tg.InputPhotoFileLocation{
//				ID:            d.ID,
//				AccessHash:    d.AccessHash,
//				FileReference: d.FileReference,
//			},
//			Limit:  1024 * 1024,
//			Offset: 0,
//		})
//		bytes := file.(*tg.UploadFile).Bytes
//
//		if err != nil {
//			fmt.Println("failed to download photo", err.Error())
//			return
//		}
//	}
//}
