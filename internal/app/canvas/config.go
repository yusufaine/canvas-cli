package canvas

import (
	"flag"
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

const tokenFileName = ".token"

// TODO: Add more config params (e.g max size, ignore extension, etc.)
type Config struct {
	AccessToken string
}

func NewConfig() *Config {
	var (
		debug, store bool
		c            Config
	)

	flag.StringVar(&c.AccessToken, "token", "", "Canvas access token. If none provided, the application will try to read the '.token' file")
	flag.BoolVar(&store, "store", false, "Stores token in a '.token' file in the same directory as the binary. (Default: false)")
	flag.BoolVar(&debug, "debug", false, "Log debug severity")
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}
	log.Debug("running in debug mode")

	c.tokenCheck()
	c.storeToken()

	return &c
}

func (c *Config) tokenCheck() {
	// If passed as a flag arg
	if c.AccessToken != "" {
		log.Info("access token supplied via flags, skip reading from file...")
		return
	}

	// Read from file, if exists
	log.Info("access token not provided, attempting to read from '.token'...")
	fileInfo, err := os.Stat(tokenFileName)
	if os.IsNotExist(err) {
		log.Error("'.token' file does not exist")
		os.Exit(1)
	}
	if fileInfo.IsDir() {
		log.Error("'.token' is not a file")
		os.Exit(1)
	}

	content, err := os.ReadFile(tokenFileName)
	if err != nil {
		log.Error("unable to read '.token'", "error", err)
		os.Exit(1)
	}

	c.AccessToken = strings.TrimSpace(string(content))
	log.Info("loaded access token from '.token' file")
}

func (c *Config) storeToken() {
	// If access token in struct was not set
	if c.AccessToken == "" {
		log.Error("unable to store token, none was provided")
		os.Exit(1)
	}

	f, err := os.Create(tokenFileName)
	if err != nil {
		log.Error("unable to create/truncate '.token' file", "error", err)
		os.Exit(1)
	}

	_, err = f.WriteString(c.AccessToken)
	if err != nil {
		log.Error("unable to (over)write content to '.token'", "error", err)
		os.Exit(1)
	}
}
