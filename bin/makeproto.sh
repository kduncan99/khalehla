#!/bin/bash

# Khalehla Project
# Copyright © 2023 by Kurt Duncan, BearSnake LLC
# All Rights Reserved
#
# Generates all the gRPC code for Khalehla
# Run this from the root Kahlehla directory

protoc exec/msg/console.proto --java_out ./kdte/src/main/generated
protoc exec/msg/console.proto --go_out .
