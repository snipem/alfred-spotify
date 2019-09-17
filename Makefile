SHELL := /bin/bash

install:
	go build -o workflow/alfred-spotify

test:
	source secrets.env && alfred_workflow_data=workflow alfred_workflow_cache=/tmp/alfred alfred_workflow_bundleid=mk_testing go run alfred-spotify.go album mötley
