#!/usr/bin/env bash
NETWORKID="89"

echo "Update config"
go run utils/seed-data/main.go seeds $NETWORKID
