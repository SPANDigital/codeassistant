#!/bin/sh
rm dist/*
nvm use
yarn parcel build src/index.html --public-url ./
