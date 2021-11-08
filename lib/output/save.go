package output

import (
	"github.com/maxkulish/byteconv"
	"github.com/maxkulish/dnscrypt-list/lib/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
)

var (
	// ErrorEmptyPath if file path is empty
	ErrorEmptyPath = errors.New("file path is empty")
)

//SaveDomainToFile iterate keys and append them to the file
func SaveDomainToFile(path string, keys []string) error {

	if path == "" {
		return ErrorEmptyPath
	}

	f, err := os.Create(path)
	if err != nil {
		logger.Error("output file creation error", zap.Error(err), zap.String("file", path))
		return nil
	}
	defer f.Close()

	var totalSize float64
	for _, k := range keys {
		n, err := f.WriteString(k + "\n")
		if err != nil {
			logger.Debug("append lines to file error", zap.Error(err))
		}
		totalSize += float64(n)
	}

	logger.Debug(
		"domain added to the file",
		zap.String("file", path),
		zap.String("size", byteconv.BytesSize(totalSize, "decimal", 0)))

	return nil
}
