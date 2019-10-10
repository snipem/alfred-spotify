SHELL := /bin/bash

install:
	go build -o workflow/alfred-spotify

watch:
	ls alfred-spotify.go | entr go build -o ${HOME}/cloud/data/alfred/Alfred.alfredpreferences/workflows/user.workflow.D965F327-9F53-4BE1-BE8E-AD4E955CB629/alfred-spotify

test:
	source secrets.env && alfred_workflow_data=workflow alfred_workflow_cache=/tmp/alfred alfred_workflow_bundleid=mk_testing go run alfred-spotify.go track mötley
