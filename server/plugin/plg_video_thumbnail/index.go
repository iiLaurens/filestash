package plg_video_thumbnail

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"os"
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

	f, err := os.CreateTemp("/tmp/videos/", "vid")
	if err != nil {
		Log.Error("plg_video_thumbnail::tmpfile::create %s", err.Error())
		return nil, err
	}
	defer os.Remove(f.Name())

	_, err = io.Copy(f, reader)
	if err != nil {
		Log.Error("plg_video_thumbnail::tmpfile::copy %s", err.Error())
		return nil, err
	}

	bitrate, duration, err = getVideoDetails(f.Name())
	
	cmd := exec.Command("ffmpeg",
		"-ss", "10",
		"-i", f.Name(),
		"-vf", "scale='if(gt(a,250/250),-1,250)':'if(gt(a,250/250),250,-1)'",
		"-frames:v", "1",
		"-f", "image2pipe",
		"-vcodec", "png",
		"pipe:1")

	cmd.Stderr = &str
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		Log.Debug("plg_video_thumbnail::ffmpeg::stderr %s", str.String())
		Log.Error("plg_video_thumbnail::ffmpeg::run %s", err.Error())
		return nil, err
	}
	return NewReadCloserFromBytes(buf.Bytes()), nil
}


func getVideoDetails(inputName string) (bitrate int64, duration float64, err error) {
	var buf bytes.Buffer
	var str bytes.Buffer

	cmd := exec.Command("ffprobe", 
	"-v", "error",
	 "-select_streams", "v:0",
	  "-show_entries", "stream=bit_rate:format=duration",
	  "-of", "default=noprint_wrappers=1:nokey=1",
	  inputName)

	cmd.Stderr = &str
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		Log.Debug("plg_video_thumbnail::ffmpeg::probe %s", str.String())
		Log.Error("plg_video_thumbnail::ffmpeg::probe %s", err.Error())
		return nil, err
	}

	return parseFfprobeOutput(buf.String())
}

func parseFfprobeOutput(raw string) (bitrate int64, duration float64, err error) {
	for _, output := range strings.Split(strings.Trim(raw, "\n"), "\n") {
		if bitrate == 0 {
			bitrate, err = strconv.ParseInt(output, 10, 64)
			if err != nil {
				if duration != 0 {
					Log.Debug("plg_video_thumbnail::ffmpeg::probe::parse_bitrate %s", output.String())
					Log.Error("plg_video_thumbnail::ffmpeg::probe::parse_bitrate %s", err.Error())
					return
				}
			} else {
				continue
			}
		}

		if duration == 0 {
			duration, err = strconv.ParseFloat(output, 64)
			if err != nil && bitrate != 0 {
				Log.Debug("plg_video_thumbnail::ffmpeg::probe::parse_duration %s", output.String())
				Log.Error("plg_video_thumbnail::ffmpeg::probe::parse_duration %s", err.Error())
				return
			}
		}
	}

	return
}