#!/bin/sh
GOOS=windows go build -o ./target/game.exe .
zip ./target/GameWin ./target/game.exe ./res/*.png
