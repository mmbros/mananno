#!/bin/bash
#set -x # echo on

# //go:generate go-bindata -nometadata -pkg templates -prefix templates/tmpl -o templates/bindata.go templates/tmpl/...
#//go:generate gentmpl -c templates/gentmpl.conf -o templates/templates.go
##GOOS=linux GOARCH=arm GOARM=7 go build -v -o mananno2 main.go
#GOOS=linux GOARCH=arm go build -v -o mananno2 main.go

PKG="templates"
EXE="mananno2"

function build_for_rpi {
	gen_tmpl_prod
	echo -n "Building for Raspberry PI... "
    GOOS=linux GOARCH=arm go build -v -o $EXE  main.go
	echo done
}

function usage {
	echo "Usage: $(basename $0) {rpi|tmpl-dev|tmpl-prod}"
}

function gen_tmpl_dev {
	echo -n "Generating templates for development... "
	go-bindata -debug -pkg $PKG -prefix $PKG/tmpl -o $PKG/bindata.go $PKG/tmpl/...
	gentmpl -d -c $PKG/gentmpl.conf -o $PKG/templates.go
	echo done
}

function gen_tmpl_prod {
	echo -n "Generating templates for production... "
	go-bindata -nometadata -pkg $PKG -prefix $PKG/tmpl -o $PKG/bindata.go $PKG/tmpl/...
	gentmpl -c $PKG/gentmpl.conf -o $PKG/templates.go
	echo done
}

case "$1" in
	rpi)
		build_for_rpi
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
