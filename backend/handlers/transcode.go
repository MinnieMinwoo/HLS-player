package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

type Quality struct {
	Name       string
	Resolution string
	Bitrate    string
	Bandwidth  int
	Size       string
}

var qualities = []Quality{
    {
        Name:       "1080p",
        Resolution: "1920:1080",
        Bitrate:    "5000k",
        Bandwidth:  5000000,
				Size: 		  "1920x1080",
    },
    {
        Name:       "720p",
        Resolution: "1280:720",
        Bitrate:    "2800k",
        Bandwidth:  2800000,
				Size: 		  "1280x720",
    },
    {
        Name:       "480p",
        Resolution: "854:480",
        Bitrate:    "1400k",
        Bandwidth:  1400000,
				Size: 		  "854x480",
    },
}


func Transcode(inputPath string) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(qualities))

	for _, q := range qualities {
		wg.Add(1)
		go func (q Quality) {
			defer wg.Done()
			if err := transcodeToQuality(inputPath, q); err != nil {
				errCh <- err
			}
		}(q)
	}
	wg.Wait()
	close(errCh)

	if err := <- errCh; err != nil {
		return err
	}
	return generateMasterPlaylist()
}

func transcodeToQuality(inputPath string, q Quality) error {
	outputDir := "hls/" + q.Name
	os.MkdirAll(outputDir, os.ModePerm)

	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-vf", "scale="+q.Resolution,
		"-b:v", q.Bitrate,
		"-hls_time", "10",
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", outputDir + "/seg%03d.ts",
		outputDir + "/index.m3u8",
	)
	
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateMasterPlaylist() error {
	f, err := os.Create("hls/master.m3u8")
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, "#EXTM3U")

	for _, q := range qualities {
		fmt.Fprintf(f, "#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%s\n", q.Bandwidth, q.Size)
		fmt.Fprintf(f, "%s/index.m3u8\n", q.Name)
	}

	return nil
}