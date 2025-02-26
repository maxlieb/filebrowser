//go:generate go-enum --sql --marshal --names --file $GOFILE
package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/img"
)

type PreviewSize int

type ImgService interface {
	FormatFromExtension(ext string) (img.Format, error)
	Resize(ctx context.Context, in io.Reader, width, height int, out io.Writer, options ...img.Option) error
}

type FileCache interface {
	Store(ctx context.Context, key string, value []byte) error
	Load(ctx context.Context, key string) ([]byte, bool, error)
	Delete(ctx context.Context, key string) error
}

const privateFilePerm os.FileMode = 0600
const ffmpegJobLimit = 4

var (
	jobTokens = make(chan struct{}, ffmpegJobLimit) // Limit to 4 concurrent ffmpeg jobs
)

func previewHandler(imgSvc ImgService, fileCache FileCache, enableThumbnails, resizePreview bool) handleFunc {
	return withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		if !d.user.Perm.Download {
			return http.StatusAccepted, nil
		}
		vars := mux.Vars(r)

		previewSize, err := ParsePreviewSize(vars["size"])
		if err != nil {
			return http.StatusBadRequest, err
		}

		file, err := files.NewFileInfo(&files.FileOptions{
			Fs:         d.user.Fs,
			Path:       "/" + vars["path"],
			Modify:     d.user.Perm.Modify,
			Expand:     true,
			ReadHeader: d.server.TypeDetectionByHeader,
			Checker:    d,
		})
		if err != nil {
			return errToStatus(err), err
		}

		setContentDisposition(w, r, file)

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		switch file.Type {
		case "image":
			return handleImagePreview(ctx, w, r, imgSvc, fileCache, file, previewSize, enableThumbnails, resizePreview)
		case "video":
			return handleVideoPreview(ctx, w, r, fileCache, file, previewSize)
		default:
			return http.StatusNotImplemented, fmt.Errorf("can't create preview for %s type", file.Type)
		}
	})
}


// detectGPU detects the GPU based on the OS.
func detectGPU() string {
	switch runtime.GOOS {
	case "windows":
		return detectWindowsGPU()
	case "linux":
		return detectLinuxGPU()
	case "darwin":
		return detectMacGPU()
	default:
		return "Unsupported OS"
	}
}

// detectWindowsGPU detects the GPU on Windows using dxdiag (faster than WMIC).
func detectWindowsGPU() string {
	cmd := exec.Command("cmd", "/C", "dxdiag /t dxdiag.txt && findstr /C:\"Card name\" dxdiag.txt")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "Unknown GPU"
	}

	lines := strings.Split(out.String(), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(strings.Split(lines[0], ":")[1])
	}
	return "Unknown GPU"
}

// detectLinuxGPU detects the GPU on Linux using lspci.
func detectLinuxGPU() string {
	cmd := exec.Command("sh", "-c", "lspci | grep VGA")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "Unknown GPU"
	}

	lines := strings.Split(out.String(), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return "Unknown GPU"
}

// detectMacGPU detects the GPU on macOS using system_profiler.
func detectMacGPU() string {
	cmd := exec.Command("system_profiler", "SPDisplaysDataType")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "Unknown GPU"
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Chipset Model") {
			return strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}
	return "Unknown GPU"
}

// selectFFmpegAccel selects the appropriate FFmpeg hardware acceleration based on the detected GPU.
func selectFFmpegAccel(gpu string) string {
	gpu = strings.ToLower(gpu)

	if strings.Contains(gpu, "nvidia") {
		return "cuda"
	} else if strings.Contains(gpu, "amd") || strings.Contains(gpu, "radeon") {
		if runtime.GOOS == "windows" {
			return "d3d11va" // AMD on Windows uses d3d11va
		} else {
			return "vaapi" // AMD on Linux uses VAAPI
		}
	} else if strings.Contains(gpu, "intel") {
		if runtime.GOOS == "darwin" {
			return "videotoolbox" // Intel-based Macs use VideoToolbox
		}
		return "qsv"
	} else if runtime.GOOS == "darwin" {
		return "videotoolbox"
	}
	return "software"
}

func handleVideoPreview(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	fileCache FileCache,
	file *files.FileInfo,
	previewSize PreviewSize,
) (int, error) {
	path := afero.FullBaseFsPath(file.Fs.(*afero.BasePathFs), file.Path)
	thumbDir := filepath.Join(filepath.Dir(path), ".thumbnails")
	thumbPath := filepath.Join(thumbDir, previewCacheKey(file, previewSize)+".webp")

	// Check if the thumbnail file exists
	if _, err := os.Stat(thumbPath); err == nil {
		http.ServeFile(w, r, thumbPath)
		return 0, nil
	}

	cacheKey := previewCacheKey(file, previewSize)
	resizedImage, ok, err := fileCache.Load(ctx, cacheKey)
	if err != nil {
		return errToStatus(err), err
	}
	if !ok {
		jobTokens <- struct{}{}
		defer func() { <-jobTokens }()

		gpu := detectGPU()
		accel := selectFFmpegAccel(gpu)
		fmt.Println("Detected GPU:", gpu)
		fmt.Println("Selected FFmpeg Acceleration:", accel)

		// Example integration into preview.go (replace "vaapi" with detected accel method)
		ffmpegArgs := []string{
			"-y", "-hwaccel", accel, "-i", "input.mp4",
			"-vf", "thumbnail,crop=w='min(iw,ih)':h='min(iw,ih)',scale=128:128",
			"-quality", "40", "-frames:v", "1", "-c:v", "webp", "-f", "image2pipe", "-",
		}
		fmt.Println("FFmpeg command:", strings.Join(ffmpegArgs, " "))

		resizedImage, err = runFFmpeg(ctx, ffmpegArgs)
		if err != nil {
			fmt.Printf("Hardware acceleration failed: %v, retrying without it...\n", err)

			// Retry without hardware acceleration
			ffmpegArgs = []string{
				"-y", "-i", path,
				"-vf", "thumbnail,crop=w='min(iw,ih)':h='min(iw,ih)',scale=128:128",
				"-quality", "40", "-frames:v", "1", "-c:v", "webp", "-f", "image2pipe", "-",
			}
			resizedImage, err = runFFmpeg(ctx, ffmpegArgs)
			if err != nil {
				return errToStatus(err), err
			}
		}

		// Save the generated thumbnail
		if err := os.MkdirAll(thumbDir, os.ModePerm); err != nil {
			return http.StatusInternalServerError, err
		}
		if err := os.WriteFile(thumbPath, resizedImage, privateFilePerm); err != nil {
			return http.StatusInternalServerError, err
		}

		go func() {
			if err := fileCache.Store(context.Background(), cacheKey, resizedImage); err != nil {
				fmt.Printf("Failed to cache resized image: %v\n", err)
			}
		}()
	}

	w.Header().Set("Cache-Control", "private")
	w.Header().Set("Content-Type", "image/webp")
	http.ServeContent(w, r, "", file.ModTime, bytes.NewReader(resizedImage))
	return 0, nil
}

func runFFmpeg(ctx context.Context, args []string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	done := make(chan error)
	go func() {
		_, err := io.Copy(&buf, stdout)
		done <- err
	}()

	select {
	case <-ctx.Done():
		_ = cmd.Process.Kill()
		<-done
		return nil, ctx.Err()
	case err := <-done:
		if err != nil {
			return nil, err
		}
		if err := cmd.Wait(); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func handleImagePreview(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
	imgSvc ImgService,
	fileCache FileCache,
	file *files.FileInfo,
	previewSize PreviewSize,
	enableThumbnails, resizePreview bool,
) (int, error) {
	if (previewSize == PreviewSizeBig && !resizePreview) ||
		(previewSize == PreviewSizeThumb && !enableThumbnails) {
		return rawFileHandler(w, r, file)
	}

	format, err := imgSvc.FormatFromExtension(file.Extension)
	// Unsupported extensions directly return the raw data
	if errors.Is(err, img.ErrUnsupportedFormat) || format == img.FormatGif {
		return rawFileHandler(w, r, file)
	}
	if err != nil {
		return errToStatus(err), err
	}

	cacheKey := previewCacheKey(file, previewSize)
	resizedImage, ok, err := fileCache.Load(ctx, cacheKey)
	if err != nil {
		return errToStatus(err), err
	}
	if !ok {
		resizedImage, err = createPreview(ctx, imgSvc, fileCache, file, previewSize)
		if err != nil {
			return errToStatus(err), err
		}
	}

	w.Header().Set("Cache-Control", "private")
	http.ServeContent(w, r, file.Name, file.ModTime, bytes.NewReader(resizedImage))

	return 0, nil
}

func createPreview(ctx context.Context, imgSvc ImgService, fileCache FileCache,
	file *files.FileInfo, previewSize PreviewSize) ([]byte, error) {
	fd, err := file.Fs.Open(file.Path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	var (
		width   int
		height  int
		options []img.Option
	)

	switch {
	case previewSize == PreviewSizeBig:
		width = 1080
		height = 1080
		options = append(options, img.WithMode(img.ResizeModeFit), img.WithQuality(img.QualityMedium))
	case previewSize == PreviewSizeThumb:
		width = 256
		height = 256
		options = append(options, img.WithMode(img.ResizeModeFill), img.WithQuality(img.QualityLow), img.WithFormat(img.FormatJpeg))
	default:
		return nil, img.ErrUnsupportedFormat
	}

	buf := &bytes.Buffer{}
	if err := imgSvc.Resize(ctx, fd, width, height, buf, options...); err != nil {
		return nil, err
	}

	go func() {
		cacheKey := previewCacheKey(file, previewSize)
		if err := fileCache.Store(context.Background(), cacheKey, buf.Bytes()); err != nil {
			fmt.Printf("failed to cache resized image: %v", err)
		}
	}()

	return buf.Bytes(), nil
}

// func previewCacheKey(f *files.FileInfo, previewSize PreviewSize) string {
// return fmt.Sprintf("%x%x%x", f.RealPath(), f.ModTime.Unix(), previewSize)
// }
func previewCacheKey(f *files.FileInfo, previewSize PreviewSize) string {
	const maxLength = 100 // maximum length for the filename
	realPath := f.RealPath()
	modTime := f.ModTime.Unix()
	previewSizeStr := fmt.Sprintf("%d", previewSize)

	// Construct the cache key
	cacheKey := fmt.Sprintf("%x%x%x", realPath, modTime, previewSizeStr)
	if len(cacheKey) > maxLength {
		cacheKey = fmt.Sprintf("%x%x", cacheKey[:maxLength-10], cacheKey[len(cacheKey)-10:])
	}
	return cacheKey
}
