#!/bin/sh
rm dist/*
yarn parcel build src/index.html --public-url ./
