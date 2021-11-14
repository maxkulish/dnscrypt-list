package download

import (
	"bufio"
	"github.com/maxkulish/dnscrypt-list/lib/db"
	"github.com/maxkulish/dnscrypt-list/lib/logger"
	"github.com/maxkulish/dnscrypt-list/lib/target"
	"github.com/maxkulish/dnscrypt-list/lib/validator"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
)

var (
	// ErrInvalidDomain an error shows that domain has wrong format or forbidden symbols
	ErrInvalidDomain = errors.New("Invalid domain")
	// ErrCommentedLine sould be returned if line has # as a first symbol
	ErrCommentedLine = errors.New("Commented line")
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

		foundDomains := make(map[string]string)

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
			line := scanner.Text()

			// TODO: collect ips
			// TODO: work with host files: ip domain
			domain, err := NormalizedDomain(line)
			if err != nil {
				continue
			}

			foundDomains[domain] = strconv.Itoa(int(targetType))
			total++
			fileLines++
		}

		if err := scanner.Err(); err != nil {
			logger.Debug("file scanner error", zap.Error(err))
			return err
		}

		logger.Debug("scanned", zap.Int("lines", fileLines))

		err = conn.AddBunch(foundDomains)
		if err != nil {
			logger.Debug("")
		}
	}

	logger.Debug("files scanning finished", zap.Int64("total", total))

	return err
}

// NormalizedDomain cleans symbols and returns valid domain
func NormalizedDomain(line string) (string, error) {

	domain := strings.TrimSpace(line)

	switch {
	case len(domain) == 0:
		return "", ErrInvalidDomain
	case domain[0:1] == "#":
		return "", ErrCommentedLine
	case !validator.IsLetter(domain[0:1]):
		return "", ErrInvalidDomain
	case !validator.IsValidHost(domain):
		return "", ErrInvalidDomain
	default:
		return domain, nil
	}
}
