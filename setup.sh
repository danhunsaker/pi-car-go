#!/usr/bin/env bash

set -e

start_dir=$(pwd)

# Dependencies
sudo apt install -y build-essential automake autoconf libtool bison swig golang portaudio19-dev libsoxr-dev

# Setup
mkdir -p ~/git
cd ~/git

if [ ! -d sphinxbase ]
then
	git clone https://github.com/cmusphinx/sphinxbase
	cd sphinxbase
else
	cd sphinxbase
	git pull
fi
./autogen.sh
make
sudo make install
cd ..

if [ ! -d sphinxbase ]
then
	git clone https://github.com/cmusphinx/pocketsphinx
	cd pocketsphinx
else
	cd pocketsphinx
	git pull
fi
./autogen.sh
make
sudo make install

mkdir -p ~/go/{src/github.com/danhunsaker,bin}
ln -sf ${start_dir} ~/go/src/github.com/danhunsaker/
export GOPATH=~/go
export PATH=~/go/bin:${PATH}
cd ~/go/src/github.com/danhunsaker/$(basename ${start_dir})
go get -v ./...
