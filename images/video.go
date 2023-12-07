package images

import (
	"fmt"
	"image"
	"log/slog"
	"os"
	"os/exec"

	"github.com/jcocozza/go_fractals/utils"
)

// create a video from a set of images
// inputPattern should be something like:
// dirPath + "/image%01d.png"
func CreateVideo(inputPattern, outputVideo string, fps int) {
	cmd := exec.Command("ffmpeg",
		"-framerate", fmt.Sprint(fps),            // Frame rate
		"-i", inputPattern,           			  // Input image pattern
		"-c:v", "libx264",            			  // Video codec
		"-pix_fmt", "yuv420p",        			  // Pixel format
		outputVideo)							  // File output path

	out, err := cmd.CombinedOutput()

	if err != nil {
		slog.Error("Error running ffmpeg command:", err)
		slog.Error("Combined Output: " + string(out))
		return
	}
}

// turn a list of images into an mp4
func VideoFromImages(imgList []image.Image, outputPath string, fps int) {
	dir, _ := os.MkdirTemp("","video")
	defer os.RemoveAll(dir)

	for i,img := range imgList {
		fileName := fmt.Sprintf("image%d", i)
		path := dir + "/" + fileName + ".png"
		SavePNG(img, path)
		utils.ProgressBar(i,len(imgList))
	}

	inputPattern := dir+"/image%01d.png"
	CreateVideo(inputPattern, outputPath, fps)
}