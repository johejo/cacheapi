package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/allegro/bigcache"

	"github.com/johejo/cacheapi"
)

var (
	config = bigcache.Config{}
	ver    bool
	host   string
	port   int
)

const version = "0.0.1"

func init() {
	// Flags override values set from env vars.
	getVarsFromEnv()
	getVarsFromFlag()
}

func main() {
	flag.Parse()

	if ver {
		fmt.Printf("cacheapi server v%s", version)
		os.Exit(0)
	}

	logger := cacheapi.GetLogger()
	config.Logger = logger

	cs := cacheapi.NewCacheService(config)
	ch := cacheapi.NewCacheHandler(cs)
	api := cacheapi.NewCacheAPI(ch)
	s := cacheapi.NewServer(api, host, port)

	logger.Fatal(s.Run())
}

func getVarsFromEnv() {
	var err error

	host = cacheapi.GetEnv("HOST", "")

	port, err = strconv.Atoi(cacheapi.GetEnv("PORT", "8888"))
	if err != nil {
		port = 8888
	}

	config.Shards, err = strconv.Atoi(cacheapi.GetEnv("SHARDS", "1024"))
	if err != nil {
		config.Shards = 1024
	}

	lifeWindow, err := strconv.Atoi(cacheapi.GetEnv("LIFE_WINDOW", "-1"))
	if err != nil {
		lifeWindow = -1
	}
	config.LifeWindow = time.Duration(lifeWindow)

	cleanWindow, err := strconv.Atoi(cacheapi.GetEnv("CLEAN_WINDOW", "0"))
	if err != nil {
		cleanWindow = 0
	}
	config.CleanWindow = time.Duration(cleanWindow)

	maxEntriesInWindow, err := strconv.Atoi(cacheapi.GetEnv("MAX_ENTRIES_IN_WINDOW", "600000"))
	if err != nil {
		maxEntriesInWindow = 1000 * 10 * 60
	}
	config.MaxEntriesInWindow = maxEntriesInWindow

	maxEntrySize, err := strconv.Atoi(cacheapi.GetEnv("MAX_ENTRY_SIZE", "500"))
	if err != nil {
		maxEntrySize = 500
	}
	config.MaxEntrySize = maxEntrySize

	verbose, err := strconv.ParseBool(cacheapi.GetEnv("VERBOSE", "false"))
	if err != nil {
		verbose = false
	}
	config.Verbose = verbose

	hardMaxCacheSize, err := strconv.Atoi(cacheapi.GetEnv("HARD_MAX_CACHE_SIZE", "0"))
	if err != nil {
		hardMaxCacheSize = 0
	}
	config.HardMaxCacheSize = hardMaxCacheSize
}

func getVarsFromFlag() {
	flag.IntVar(
		&config.Shards,
		"shards",
		1024,
		"Number of cache shards, value must be a power of two.",
	)
	flag.DurationVar(
		&config.LifeWindow,
		"lifeWindow",
		100000*100000*60,
		"Time after which entry can be evicted.",
	)
	flag.DurationVar(
		&config.CleanWindow,
		"clearWindow",
		0,
		"Interval between removing expired entries (clean up).",
	)
	flag.IntVar(
		&config.MaxEntriesInWindow,
		"maxInWindow",
		1000*10*60,
		"Max number of entries in life window.",
	)
	flag.IntVar(
		&config.MaxEntrySize,
		"maxEntrySize",
		500,
		"Max size of entry in bytes.",
	)
	flag.BoolVar(
		&config.Verbose,
		"verbose",
		false,
		"Verbose mode prints information about new memory allocation.",
	)
	flag.IntVar(
		&config.HardMaxCacheSize,
		"hardMaxCacheSize",
		0,
		"HardMaxCacheSize is a limit for cache size in MB. Cache will not allocate more memory than this limit.",
	)
	flag.StringVar(
		&host,
		"host",
		"",
		"The hostname to listen on.",
	)
	flag.IntVar(
		&port,
		"port",
		8888,
		"The port to listen on.",
	)
	flag.BoolVar(
		&ver,
		"version",
		false,
		"Print server version.",
	)
}
