package download

import (
	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	"github.com/maxkulish/dnscrypt-list/lib/logger"
	"github.com/maxkulish/dnscrypt-list/lib/target"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

const (
	getTimeout   = 5 * time.Second
	tempFilePerm = os.FileMode(0600)
)

// LocalFile keeps information about file: path, temp, type
type LocalFile struct {
	Path   string
	Format target.Format
	Temp   bool
	Type   target.Type
}

// Temp structure which keeps list of local files and map of responses
type Temp struct {
	Storage map[string][]byte
	Files   []LocalFile
	Errors  []error
	mu      sync.Mutex
}

// NewTemp create Temp instance and returns pointer to the new struct
func NewTemp() *Temp {
	return &Temp{
		Storage: make(map[string][]byte),
	}
}

// AddToStorage adds concurrently []byte to the Storage map
func (t *Temp) AddToStorage(targetID string, data []byte) {

	// create new uuid if targetID is empty
	if targetID == "" {
		targetID = uuid.New().String()
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	t.Storage[targetID] = append(t.Storage[targetID], data...)
}

// GetAndSaveTargets iterate targets and save response body to the temp files
func GetAndSaveTargets(targets *target.TargetsStore) ([]LocalFile, error) {

	wg := &sync.WaitGroup{}

	temp := NewTemp()

	// Download all remote files to memory
	for _, targ := range targets.Targets {

		// Local target without URL but with Path
		if targ.URL.String() == "" {
			if targ.Path != "" {
				temp.Files = append(temp.Files, LocalFile{
					Path:   targ.Path,
					Format: targ.Format,
					Temp:   false,
					Type:   targ.TargetType,
				})
			}
			continue
		}

		// Download body of the file
		wg.Add(1)
		go temp.GetRemote(wg, targ.TargetID, targ.URL)

		temp.Files = append(temp.Files, LocalFile{
			Path:   targ.TempFile,
			Format: targ.Format,
			Temp:   true,
			Type:   targ.TargetType,
		})

	}

	wg.Wait()

	// Save response body to the temp file
	for _, targ := range targets.Targets {

		n, err := SaveToFile(targ.TempFile, temp.Storage[targ.TargetID])
		if err != nil {
			logger.Error("temp file saving error", zap.Error(err))
		}
		logger.Debug(
			"saved data to the tempFile",
			zap.Int("bytes", n),
			zap.String("file", targ.TempFile))
	}

	logger.Info("local files from config added")
	return temp.Files, nil
}

//GetRemote sends GET response and save resp.Body
func (t *Temp) GetRemote(wg *sync.WaitGroup, targID string, remoteURL *url.URL) {

	if wg != nil {
		defer wg.Done()
	}

	client := &http.Client{
		Timeout: getTimeout,
	}

	// prepare request
	req, err := http.NewRequest(http.MethodGet, remoteURL.String(), nil)
	if err != nil {
		logger.Error("http.NewRequest error", zap.Error(err))
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errors = append(t.Errors, err)
		logger.Debug("http.Get url error", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	// response body is copied to the buffer (memory)
	// possible problem if the file is too large to keep in memory

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("copying response body error", zap.Error(err))
		return
	}

	// AddToStorage append response bodies to the temprorary storage
	t.AddToStorage(targID, body)

	return
}

// SaveToFile copy response body to the file
// if local file not exist, file will be created
// string, *http.Response.Body -> int64, err
func SaveToFile(fileName string, data []byte) (int, error) {

	if fileName == "" {
		logger.Error("file name is empty", zap.String("file", fileName))
		return 0, nil
	}

	err := ioutil.WriteFile(fileName, data, tempFilePerm)
	logger.Debug(
		"file saved",
		zap.String("file", fileName),
		zap.String("size", humanize.Bytes(uint64(len(data)))))
	if err != nil {
		return 0, err
	}

	return len(data), nil
}
