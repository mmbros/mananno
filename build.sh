#!/bin/bash
#set -x # echo on

# //go:generate go-bindata -nometadata -pkg templates -prefix templates/tmpl -o templates/bindata.go templates/tmpl/...
#//go:generate gentmpl -c templates/gentmpl.conf -o templates/templates.go
##GOOS=linux GOARCH=arm GOARM=7 go build -v -o mananno2 main.go
#GOOS=linux GOARCH=arm go build -v -o mananno2 main.go

PKG="templates"
ASSET="asset"
TMPL="tmpl"
EXE="mananno2"


BUILDFLAGS="-v"

# Shrink your Go binaries with this one weird trick
# https://blog.filippo.io/shrink-your-go-binaries-with-this-one-weird-trick/
LDFLAGS="-s -w"

LOG_PREFIX=">>> "

function log {
	echo "$LOG_PREFIX$1"
}

function usage {
	echo "Usage: $(basename $0) {linux|rpi|tmpl-dev|tmpl-prod}"
}

function compress_exe {
	if [ -x "$(command -v upx)" ]; then
		log "Compressing executeble with UPX... "
		upx $EXE
		log "Compressing executeble with UPX... Done"
	else
		log "Skip compression: UPX not found!"

	fi
}

function build_for_linux {
	gen_tmpl_prod
	log "Building for Linux... "
	go build $BUILDFLAGS -ldflags "$LDFLAGS" -o $EXE  *.go
	log "Building for Linux... Done"
	compress_exe
}

function build_for_rpi {
	gen_tmpl_prod
	log "Building for Raspberry PI... "
	GOOS=linux GOARCH=arm go build $BUILDFLAGS -ldflags "$LDFLAGS" -o $EXE  *.go
	log "Building for Raspberry PI... Done"
	compress_exe
}

function gen_tmpl_dev {
	log "Generating templates for development... "
	go-bindata -debug -pkg $PKG -prefix $ASSET -o $PKG/bindata.go $ASSET/...
	gentmpl -d -c $PKG/gentmpl.conf -p $TMPL -o $PKG/templates.go
	log "Generating templates for development... Done"
}

function gen_tmpl_prod {
	log "Generating templates for production... "
	go-bindata -pkg $PKG -prefix $ASSET -o $PKG/bindata.go $ASSET/...
	gentmpl -c $PKG/gentmpl.conf -p $TMPL -o $PKG/templates.go
	log "Generating templates for production... Done"
}

case "$1" in
	rpi)
		build_for_rpi
		;;
	linux)
		build_for_linux
		;;
	tmpl-dev)
		gen_tmpl_dev
		;;
	tmpl-prod)
		gen_tmpl_prod
		;;
	*)
		usage
		exit 1
esac
