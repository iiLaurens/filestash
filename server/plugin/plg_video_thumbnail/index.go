package plg_video_thumbnail

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	. "github.com/iilaurens/filestash/server/common"
)

//go:embed dist/placeholder.png
var placeholder []byte

func init() {
	Hooks.Register.Thumbnailer("video/mp4", thumbnailBuilder{thumbnailMp4})
}

type thumbnailBuilder struct {
	fn func(reader io.ReadCloser, ctx *App, res *http.ResponseWriter, req *http.Request) (io.ReadCloser, error)
}

func (this thumbnailBuilder) Generate(reader io.ReadCloser, ctx *App, res *http.ResponseWriter, req *http.Request) (io.ReadCloser, error) {
	return this.fn(reader, ctx, res, req)
}

func thumbnailMp4(reader io.ReadCloser, ctx *App, res *http.ResponseWriter, req *http.Request) (io.ReadCloser, error) {
	h := (*res).Header()
	r, err := generateThumbnailFromVideo(reader)
	if err != nil {
		h.Set("Content-Type", "image/png")
		return NewReadCloserFromBytes(placeholder), nil
	}
	h.Set("Content-Type", "image/png")
	h.Set("Cache-Control", fmt.Sprintf("max-age=%d", 3600*12))
	return r, nil
}

func generateThumbnailFromVideo(reader io.ReadCloser) (io.ReadCloser, error) {
	var buf bytes.Buffer
	var str bytes.Buffer

	cmd := exec.Command("ffmpeg",
		"-ss", "10",
		"-i", "pipe:0",
		"-vf", "scale='if(gt(a,250/250),-1,250)':'if(gt(a,250/250),250,-1)'",
		"-frames:v", "1",
		"-f", "image2pipe",
		"-vcodec", "png",
		"pipe:1")

	cmd.Stdin = reader
	cmd.Stderr = &str
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		Log.Debug("plg_video_thumbnail::ffmpeg::stderr %s", str.String())
		Log.Error("plg_video_thumbnail::ffmpeg::run %s", err.Error())
		return nil, err
	}
	return NewReadCloserFromBytes(buf.Bytes()), nil
}
