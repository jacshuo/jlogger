# A Simple Logger in Golang.

## Install

```bash
go get github.com/jacshuo/jlogger
```

## Usage

```go
package main

import (
	"os"
)

func main() {
	myLogger := GetMultiWriteLogger("./", "testlog.log", os.Stdout) // 2020/12/07 15:40:02.002401 JLogger/main.go:8	DEBUG	Hello world!
	myLogger.Debug("Hello world!")
	myLogger.Debugf("%s, world!", "Hello")
	myLogger.Info("Hello world!")
	myLogger.Infof("%s, world!", "Hello")
	myLogger.Warn("Hello world!")
	myLogger.Warnf("%s, world!", "Hello")
	myLogger.Error("Hello world!")
	myLogger.Errorf("%s, world!", "Hello")
	myLogger.Critical("Hello world!")
	myLogger.Criticalf("%s, world!", "Hello") // 2020/12/07 15:40:02.002401 JLogger/main.go:18	CRITICAL	Hello world!
	myLogger.Fatal("Hello world!")
	myLogger.Fatalf("%s, world!", "Hello")
}
```

Please refer to source code for more info.