#!/bin/bash -e

go install github.com/ramya-rao-a/go-outline@latest && \
  go install golang.org/x/tools/gopls@latest && \
  go install golang.org/x/tools/cmd/goimports@latest && \
  go install honnef.co/go/tools/cmd/staticcheck@latest && \
  go install github.com/cosmtrek/air@latest && \
  go install github.com/go-delve/delve/cmd/dlv@latest && \
  go install github.com/volatiletech/sqlboiler/v4@latest && \
  go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest
