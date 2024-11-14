package requests

import (
	"maps"
	"mime"
	"strings"
	"sync"

	"github.com/whitekid/goxp/mapx"
)

var mimeTypes = map[string]string{}

func mimeByExt(ext string) string {
	var mimeType string
	var ok bool

	if mimeType, ok = mimeTypes[ext]; !ok {
		mimeType = mime.TypeByExtension(ext)
	}

	if strings.HasPrefix(mimeType, "application/json") {
		mimeType, param, _ := mime.ParseMediaType(mimeType)
		mapx.SetNX[map[string]string](param, "charset", "utf-8")

		return mime.FormatMediaType(mimeType, param)
	}

	return mimeType
}

var muSetMimeTypes sync.Mutex

func setMimeTypes(types map[string]string) {
	muSetMimeTypes.Lock()
	defer muSetMimeTypes.Unlock()

	maps.Copy(mimeTypes, types)
}
