#!/usr/bin/env bash
bee generate docs
docker run --rm -v \
    $(pwd)/swagger:/opt \
    swagger2markup/swagger2markup \
    convert \
    -i /opt/swagger.yml \
    -f /opt/swagger \
    -c /opt/config.properties
cp $(pwd)/swagger/swagger.md README.md