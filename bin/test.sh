#!/usr/bin/env bash

go test -v $@ $(go list toptal/... | grep -v migrations)
