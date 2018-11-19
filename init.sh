#! /bin/bash

if [ ! -e ./vendor ]; then
	$GOPATH/bin/govendor init
	$GOPATH/bin/govendor fetch github.com/gin-gonic/gin@v1.3
	curl https://raw.githubusercontent.com/gin-gonic/gin/master/examples/basic/main.go > main.go
fi
