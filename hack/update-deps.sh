#!/bin/bash

dep ensure
find vendor -name "BUILD" -delete
find vendor -name "BUILD.bazel" -delete
bazel run //:gazelle -- -proto=disable

