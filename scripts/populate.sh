#!/usr/bin/env bash

./kvz kv set bg-color blue
./kvz hook add bg-update-file-path --link ./scripts/hookTest.sh
./kvz hook add bg-update-saved-file --file ./scripts/hookTest.sh

# shellcheck disable=SC2016
./kvz hook add bg-update-script 'echo "The value was updated to $NEW_VAL"'

./kvz hook attach bg-color bg-update-file-path
./kvz hook attach bg-color bg-update-saved-file
./kvz hook attach bg-color bg-update-script