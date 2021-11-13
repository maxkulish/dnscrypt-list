package download

import (
	"fmt"
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

		// Remote target with URL
		fileName := fmt.Sprintf("%s/%d_%d_%s", tempDir, targ.TargetType, i, targ.URL.Host)
		fmt.Printf("%d: %s\n", i, targ.URL.String())

		response := GetRemote(targ.URL)
		tempFile, err := files.CreateFileOrTruncate(fileName)
		if err != nil {
			logger.Debug("temp file creation", zap.Error(err))
		}

		n, err := io.Copy(tempFile, response.Body)
		logger.Debug("file saved", zap.String("file", fileName), zap.Int64("size", n))
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

	logger.Info("all remote targets downloaded", zap.Int64("size", bytesDownloaded))

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
