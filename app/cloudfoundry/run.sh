#!/bin/bash

exec ./bot \
  --config=<( ./ytt template --recursive --file config --file config/lib ) \
  --config=config/default.delegatebot \
  run
