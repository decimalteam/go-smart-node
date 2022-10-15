package upgrade

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"syscall"

	"gopkg.in/ini.v1"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/upgrade/types"

	cmdcfg "bitbucket.org/decimalteam/go-smart-node/cmd/config"
)

// DownloadSecured downloads file from provided URL to the specified path and checks it's hash.
func DownloadSecured(ctx sdk.Context, url string, path string, hash string) {
	// download
	ctx.Logger().Info(fmt.Sprintf("start download binary \"%s\" to file \"%s\"", url, path))
	err := downloadFile(url, path)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("error while downloading binary \"%s\" to file \"%s\" with error '%s'", url, path, err.Error()))
		return
	}
	ctx.Logger().Info(fmt.Sprintf("successful download binary \"%s\" to file \"%s\"", url, path))
	// check hash
	if !checkFile(path, hash) {
		err = os.Remove(path)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("error while remove wrong file \"%s\": '%s'", path, err.Error()))
		}
		ctx.Logger().Error(fmt.Sprintf("error check hash for file \"%s\"", path))
		return
	}
	ctx.Logger().Info(fmt.Sprintf("check hash successful for file \"%s\"", path))
}

// downloadFile downloads file from provided URL to the specified path.
func downloadFile(url string, path string) error {
	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download '%s' reply code is %d", url, resp.StatusCode)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// checkFile gets the hash of the download file, then checks what was in the transaction.
func checkFile(path, hash string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return false
	}
	return hash == hex.EncodeToString(h.Sum(nil))
}

// getDownloadFileName generates name of a download file.
func getDownloadFileName(name string) string {
	ex, _ := os.Executable()
	return filepath.Join(filepath.Dir(ex), fmt.Sprintf("%s.nv", name))
}

// resolveDownloadURL returns exact URL to download correct binary.
func resolveDownloadURL(s string) string {
	// example: "linux/ubuntu/20.04"
	u, err := url.Parse(osArchForURL())
	if err != nil {
		return ""
	}

	// example: "http://127.0.0.1/90500/decd"
	myUrl, err := url.Parse(s)
	if err != nil {
		return ""
	}

	// result: "http://127.0.0.1/90500/linux/ubuntu/20.04/decd"
	return fmt.Sprintf("%s/%s", myUrl.ResolveReference(u), path.Base(myUrl.Path))
}

// doesPageExist checks if the page at provided URL exists.
func doesPageExist(s string) bool {
	resp, err := http.Head(s)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

// changeBinary changes the old binary to a new one.
func changeBinary(plan types.Plan) error {
	mapping := planMapping(plan)
	if mapping == nil {
		return fmt.Errorf("error: mapping decode")
	}

	ex, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error: get current dir")
	}

	currPath := filepath.Dir(ex)

	downloadName := getDownloadFileName(cmdcfg.AppBinName)
	if _, err := os.Stat(downloadName); os.IsNotExist(err) {
		return err
	}

	hashes, ok := mapping[osArchForURL()]
	if !ok {
		return fmt.Errorf("error: mapping[os] undefined")
	}

	if !checkFile(downloadName, hashes[0]) {
		os.Remove(downloadName)
		return fmt.Errorf("error: hash does not match")
	}

	currBin := filepath.Join(currPath, cmdcfg.AppBinName)
	mode, err := getMode(currBin)
	if err != nil {
		os.Remove(downloadName)
		return err
	}

	err = markExecutableWithMode(downloadName, mode)
	if err != nil {
		os.Remove(downloadName)
		return err
	}

	ok = isRunSuccess(downloadName)
	if !ok {
		os.Remove(downloadName)
		return fmt.Errorf("error: file not running")
	}

	syscall.Unlink(currBin)
	err = os.Rename(downloadName, currBin)
	if err != nil {
		os.Remove(downloadName)
		return err
	}

	return nil
}

// getMode returns file mode.
func getMode(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("stating binary: %w", err)
	}
	return info.Mode().Perm(), nil
}

// markExecutableWithMode sets executable flag for specified file.
func markExecutableWithMode(path string, mode os.FileMode) error {
	return os.Chmod(path, mode|0111)
}

// isRunSuccess runs specified file with argument `version` to check if it runs.
func isRunSuccess(path string) bool {
	cmd := exec.Command(path, "version")
	err := cmd.Run()
	return err == nil
}

// osArchForURL detects and returns OS to create an URL.
func osArchForURL() string {
	switch runtime.GOOS {
	case "windows", "darwin":
		return runtime.GOOS
	case "linux":
		distr := readOSRelease("ID")
		if distr == "" {
			distr = "<unknown>"
		}
		version := readOSRelease("VERSION_ID")
		if version == "" {
			version = "<unknown>"
		}
		return fmt.Sprintf("linux/%s/%s", distr, version)
	default:
		return runtime.GOOS
	}
}

// readOSRelease reads the file under /etc/os-release to get the distribution name or version.
func readOSRelease(key string) string {
	const cfgfile = "/etc/os-release"
	cfg, err := ini.Load(cfgfile)
	if err != nil {
		return ""
	}
	return cfg.Section("").Key(key).String()
}

// planMapping returns plans info as map.
func planMapping(plan types.Plan) map[string][]string {
	var mapping map[string][]string
	err := json.Unmarshal([]byte(plan.Info), &mapping)
	if err != nil {
		return nil
	}
	return mapping
}
