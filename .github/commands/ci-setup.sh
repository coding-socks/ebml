#!/usr/bin/env bash

echo "=====> Install Go"
sudo add-apt-repository -y ppa:longsleep/golang-backports
sudo apt update -y
sudo apt-get install -y --no-install-recommends golang-1.15-go

echo "=====> Install git"
sudo apt-get install -y --no-install-recommends git
