#!/bin/sh
GOOS=linux go build -o ./target/game .
zip ./target/GameLin ./target/game ./res/*.png
