package plg_video_thumbnail

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"os"
	"strings"
	"strconv"
	"math"
	. "github.com/mickael-kerjean/filestash/server/common"
)

//go:embed dist/placeholder.png
var placeholder []byte

func init() {
	err := os.MkdirAll("/tmp/videos/", os.ModePerm)
	if err != nil {
		Log.Error("plg_video_thumbnail::init %s", err.Error())
	}
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
	r, err := generateThumbnailFromVideo(reader, "mp4")
	if err != nil {
		h.Set("Content-Type", "image/png")
		return NewReadCloserFromBytes(placeholder), nil
	}
	h.Set("Content-Type", "image/webp")
	h.Set("Cache-Control", fmt.Sprintf("max-age=%d", 3600*12))
	return r, nil
}

func generateThumbnailFromVideo(reader io.ReadCloser, ext string) (io.ReadCloser, error) {
	var str bytes.Buffer

	f, err := os.CreateTemp("/tmp/videos/", "vid_*")
	if err != nil {
		Log.Error("plg_video_thumbnail::tmpfile::create %s", err.Error())
		return nil, err
	}
	defer os.Remove(f.Name())
	tmp_out := f.Name() + ".webp"
	tmp_img := f.Name() + "_%02d.jpeg"

	_, err = io.Copy(f, reader)
	if err != nil {
		Log.Error("plg_video_thumbnail::tmpfile::copy %s", err.Error())
		return nil, err
	}

	duration, err := getVideoDetails(f.Name())
	if err != nil {
		return nil, err
	}

	for i := 1; i <= 10; i++ {
		cmd := exec.Command("ffmpeg",
		"-ss", strconv.FormatFloat((float64(i) - 0.5) * duration / 10, 'g', 6, 64),
		"-f", ext,
		"-i", f.Name(),
		"-vf", "select='eq(pict_type,I)',scale='if(gt(a,250/250),-1,250)':'if(gt(a,250/250),250,-1)'",
		"-vframes", "1",
		fmt.Sprintf(tmp_img, i))

		Log.Debug("plg_video_thumbnail:ffmpeg::make_img %s", cmd.String())

		cmd.Stderr = &str
		if err := cmd.Run(); err != nil {
			Log.Debug("plg_video_thumbnail::ffmpeg::stderr %s", str.String())
			Log.Error("plg_video_thumbnail::ffmpeg::run %s", err.Error())
			return nil, err
		} 
	}
	
	cmd := exec.Command("ffmpeg",
		"-itsscale", strconv.FormatFloat(math.Min(5.0/duration, 1), 'g', 6, 64),
		"-f", ext,
		"-i", f.Name(),
		"-vf", "scale='if(gt(a,250/250),-1,250)':'if(gt(a,250/250),250,-1)',fps=2",
		"-f", "webp",
		"-lossless", "0",
		"-compression_level", "6",
		"-loop", "0",
		"-an",
		"-preset", "picture",
		"-vcodec", "libwebp",
		tmp_out)

	Log.Debug("plg_video_thumbnail:ffmpeg::cmd %s", cmd.String())

	cmd.Stderr = &str
	if err := cmd.Run(); err != nil {
		Log.Debug("plg_video_thumbnail::ffmpeg::stderr %s", str.String())
		Log.Error("plg_video_thumbnail::ffmpeg::run %s", err.Error())
		return nil, err
	} else {
		data, _ := os.ReadFile(tmp_out)
		os.Remove(tmp_out)
		return NewReadCloserFromBytes(data), nil
	}
}


func getVideoDetails(inputName string) (duration float64, err error) {
	var buf bytes.Buffer
	var str bytes.Buffer

	cmd := exec.Command("ffprobe", 
	"-v", "error",
	 "-select_streams", "v:0",
	  "-show_entries", "format=duration",
	  "-of", "default=noprint_wrappers=1:nokey=1",
	  inputName)

	cmd.Stderr = &str
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		Log.Debug("plg_video_thumbnail::ffmpeg::probe %s", str.String())
		Log.Error("plg_video_thumbnail::ffmpeg::probe %s", err.Error())
		return 0, err
	}

	return parseFfprobeOutput(buf.String())
}

func parseFfprobeOutput(raw string) (duration float64, err error) {
	duration, err = strconv.ParseFloat(strings.Trim(raw, "\n"), 64)
	if err != nil {
		Log.Error("plg_video_thumbnail::ffmpeg::probe::parse %s", err.Error())
		return
	}

	return
}