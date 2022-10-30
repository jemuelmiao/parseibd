#!/bin/sh
go mod tidy
cd cmd
go build -o parseibd
