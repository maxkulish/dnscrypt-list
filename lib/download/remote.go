package download

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/maxkulish/dnscrypt-list/lib/files"
	"github.com/maxkulish/dnscrypt-list/lib/logger"
	"github.com/maxkulish/dnscrypt-list/lib/target"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"time"
)

// LocalFile keeps information about file: path, temp, type
type LocalFile struct {
	Path   string
	Format target.Format
	Temp   bool
	Type   target.Type
}

// GetAndSaveTargets iterate targets and save response body to the temp files
func GetAndSaveTargets(tempDir string, targets *target.TargetsStore) ([]LocalFile, error) {

	err := files.MkdirAllIfNotExist(tempDir)
	if err != nil {
		logger.Error("temp dir creation error", zap.Error(err))
		return nil, err
	}

	var tempFiles []LocalFile
	var bytesDownloaded int64
	for i, targ := range targets.Targets {

		// Local target without URL but with Path
		if targ.URL.String() == "" {
			if targ.Path != "" {
				tempFiles = append(tempFiles, LocalFile{
					Path:   targ.Path,
					Format: targ.Format,
					Temp:   false,
					Type:   targ.TargetType,
				})
			}
			continue
		}

		// prepare local file name: tempDir + TargetType + Host
		// example: /tmp/dnscrypt/2_22_rescure.me
		fileName := fmt.Sprintf("%s/%d_%d_%s", tempDir, targ.TargetType, i, targ.URL.Host)

		// Download body of the file
		response := GetRemote(targ.URL)

		// Save response body to the temp file
		n, err := SaveToFile(fileName, response.Body)
		if err != nil {
			logger.Error("temp file saving error", zap.Error(err))
		}

		bytesDownloaded += n
		tempFiles = append(tempFiles, LocalFile{
			Path:   fileName,
			Format: targ.Format,
			Temp:   true,
			Type:   targ.TargetType,
		})

		response.Body.Close()
	}

	logger.Info("all remote targets downloaded", zap.String("size", humanize.Bytes(uint64(bytesDownloaded))))

	logger.Info("local files from config added")
	return tempFiles, nil
}

//GetRemote sends GET response and save resp.Body
func GetRemote(remoteURL *url.URL) *http.Response {

	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	// prepare request
	req, err := http.NewRequest(http.MethodGet, remoteURL.String(), nil)
	if err != nil {
		logger.Error("http.NewRequest error", zap.Error(err))
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("http.Get url error", zap.Error(err))
	}

	return resp
}

// SaveToFile copy response body to the file
// if local file not exist, file will be created
// string, *http.Response.Body -> err, int64
func SaveToFile(fileName string, body io.ReadCloser) (int64, error) {

	if fileName == "" {
		logger.Error("file name is empty", zap.String("file", fileName))
		return 0, nil
	}

	tempFile, err := files.CreateFileOrTruncate(fileName)
	if err != nil {
		logger.Debug("temp file creation", zap.Error(err))
	}

	n, err := io.Copy(tempFile, body)
	logger.Debug(
		"file saved",
		zap.String("file", fileName),
		zap.String("size", humanize.Bytes(uint64(n))))
	if err != nil {
		return 0, err
	}

	return n, nil
}
