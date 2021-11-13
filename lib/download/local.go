package download

import (
	"bufio"
	"github.com/maxkulish/dnscrypt-list/lib/db"
	"github.com/maxkulish/dnscrypt-list/lib/logger"
	"github.com/maxkulish/dnscrypt-list/lib/target"
	"go.uber.org/zap"
	"os"
	"path"
	"strconv"
	"strings"
)

//ReadFilesAndSaveToDB reads local files and save them to the DB
func ReadFilesAndSaveToDB(tempFiles []LocalFile, conn *db.Conn, targetType target.Type) error {

	var total int64
	var err error

	for _, tmpFile := range tempFiles {

		// skip wrong type
		if tmpFile.Type != targetType {
			continue
		}

		foundDomain := make(map[string]string)

		var fileLines int
		logger.Debug(
			"start reading file",
			zap.String("fileName", tmpFile.Path))
		f, err := os.Open(tmpFile.Path)
		if err != nil {
			logger.Debug("open file error", zap.Error(err))
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			domain := scanner.Text()
			foundDomain[strings.TrimSpace(domain)] = strconv.Itoa(int(targetType))
			total++
			fileLines++
		}

		if err := scanner.Err(); err != nil {
			logger.Debug("file scanner error", zap.Error(err))
			return err
		}

		logger.Debug("scanned", zap.Int("lines", fileLines))

		err = conn.AddBunch(foundDomain)
		if err != nil {
			logger.Debug("")
		}
	}

	logger.Debug("files scanning finished", zap.Int64("total", total))

	return err
}

// FilterByTargetType filter file names by prefix
// []string{"1_2_example.com"}, 1 -> ["1_2_example.com"]
// []string{"1_2_example.com"}, 2 -> []
func FilterByTargetType(targetType target.Type, targets ...string) []string {

	var toRead []string
	for _, tmp := range targets {

		fileName := path.Base(tmp)

		if strings.HasPrefix(fileName, strconv.Itoa(int(targetType))) {
			toRead = append(toRead, tmp)
		}
	}

	return toRead
}
