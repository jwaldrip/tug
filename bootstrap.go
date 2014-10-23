package main

import "os"

func BootstrapCommand() string {
	bootstrap := `
		#!/usr/bin/env bash

		# derived from http://protofusion.org/wordpress/2011/01/downloading-wget-with-bash/
		function __wget() {
			read proto server doc <<<$(echo ${1//// })

			local path=/${doc// //}
			local host=${server//:*}

			exec {HTTP_FD}<>/dev/tcp/${host}/80
			echo -ne "GET ${path} HTTP/1.1\r\nHost: ${host}\r\nConnection: close\r\n\r\n" >&${HTTP_FD}
			sed -e '1,/^.$/d' <&${HTTP_FD}
		}

		__wget http://s3.amazonaws.com/tug-binaries/tugd.tgz | tar -C /usr/sbin -xz
	`

	if os.Getenv("DDEBUG") == "true" {
		bootstrap += "env DDEBUG=true /usr/sbin/tugd"
	} else {
		bootstrap += "/usr/sbin/tugd"
	}

	return bootstrap
}
