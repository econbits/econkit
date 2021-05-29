#!/bin/sh

PRJPATH=$(pwd)
PWD="$PRJPATH/third_party/netrc"
TMP="$PWD/.tmp"

find $PWD -name LICENSE -exec rm {} \;
find $PWD -name *.go -exec rm {} \;
find $PWD -name examples -exec rm -Rf {} \;

if [ -d $TMP ]; then
    rm -Rf $TMP
fi

git clone https://github.com/git-lfs/go-netrc.git $TMP

cp $TMP/LICENSE $PWD/
cp $TMP/netrc/*.go $PWD/
cp -R $TMP/netrc/examples $PWD/

rm -Rf $TMP