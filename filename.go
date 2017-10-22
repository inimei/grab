package grab

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// guessFilename returns a filename for the given http.Response. If none can be
// determined ErrNoFilename is returned.
func guessFilename(resp *http.Response) (string, error) {
	filename := resp.Request.URL.Path
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		if _, params, err := mime.ParseMediaType(cd); err == nil {
			filename = params["filename"]
		}
	}

	// sanitize
	if filename == "" || strings.HasSuffix(filename, "/") || strings.Contains(filename, "\x00") {
		return "", ErrNoFilename
	}

	filename = filepath.Base(path.Clean("/" + filename))
	if filename == "" || filename == "." || filename == "/" {
		return "", ErrNoFilename
	}

	for {
		fi, err := os.Stat(filename)
		if !os.IsNotExist(err) && fi != nil {
			ext := filepath.Ext(filename)
			if ext == "" {
				filename = filename + ".1"
			} else {
				ext = ext[1:]
				if ext == "" {
					filename = filename + "1"
				} else {
					n, err := strconv.Atoi(ext)
					if err != nil {
						filename = filename + ".1"
					} else {
						filename = filename[:len(filename)-len(ext)] + fmt.Sprintf("%v", n+1)
					}
				}
			}
		} else {
			break
		}
	}

	return filename, nil
}
