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

function usage {
	echo "Usage: $(basename $0) {linux|rpi|tmpl-dev|tmpl-prod}"
}

function build_for_linux {
	gen_tmpl_prod
	echo -n "Building for Linux... "
    go build -v -o $EXE  *.go
	echo done
}

function build_for_rpi {
	gen_tmpl_prod
	echo -n "Building for Raspberry PI... "
    GOOS=linux GOARCH=arm go build -v -o $EXE  *.go
	echo done
}

function gen_tmpl_dev {
	echo -n "Generating templates for development... "
	go-bindata -debug -pkg $PKG -prefix $ASSET -o $PKG/bindata.go $ASSET/...
	gentmpl -d -c $PKG/gentmpl.conf -p $TMPL -o $PKG/templates.go
	echo done
}

function gen_tmpl_prod {
	echo -n "Generating templates for production... "
	go-bindata -pkg $PKG -prefix $ASSET -o $PKG/bindata.go $ASSET/...
	gentmpl -c $PKG/gentmpl.conf -p $TMPL -o $PKG/templates.go
	echo done
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
