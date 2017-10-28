#!/usr/bin/env bash

go test -v $@ $(go list sample/... | grep -v migrations)
