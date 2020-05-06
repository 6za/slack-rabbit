#!/bin/bash
docker build -t slack_listener  -f $BASE_DIR/go/Dockerfile.slack_listener $BASE_DIR/go/
docker build -t  slack_writer -f $BASE_DIR/go/Dockerfile.slack_writer $BASE_DIR/go/


