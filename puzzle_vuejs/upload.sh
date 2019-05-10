#!/bin/bash

rm -rf ../puzzle_static/*
npm run build
cp -rf dist/* ../puzzle_static/

git statu

git add -u
git add ../puzzle_static/static

git statu
