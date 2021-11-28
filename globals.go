package main

import "sync"

var config config_format
var writeLock sync.Mutex
