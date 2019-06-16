# cacheapi

An alternative [BigCache HTTP Server](https://github.com/allegro/bigcache/tree/master/server)

## Web API documents

todo

## Installing

```bash
go get -u github.com/joehjo/cacheapi/cmd/cacheapi
```

## CLI

```
Usage of ./cacheapi:
  -clearWindow duration
        Interval between removing expired entries (clean up).
  -hardMaxCacheSize int
        HardMaxCacheSize is a limit for cache size in MB. Cache will not allocate more memory than this limit.
  -host string
        The hostname to listen on.
  -lifeWindow duration
        Time after which entry can be evicted. (default 10m0s)
  -maxEntrySize int
        Max size of entry in bytes. (default 500)
  -maxInWindow int
        Max number of entries in life window. (default 600000)
  -port int
        The port to listen on. (default 8888)
  -shards int
        Number of cache shards, value must be a power of two. (default 1024)
  -verbose
        Verbose mode prints information about new memory allocation.
  -version
        Print server version.
```

## License

MIT
