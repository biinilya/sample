#!/usr/bin/env bash
bee generate docs
docker run --rm -v \
    $(pwd)/swagger:/opt \
    swagger2markup/swagger2markup \
    convert \
    -i /opt/swagger.yml \
    -f /opt/swagger \
    -c /opt/config.properties
sed -e 's|<a name.*||g' -e 's/<br>//g' $(pwd)/swagger/swagger.md > README.md