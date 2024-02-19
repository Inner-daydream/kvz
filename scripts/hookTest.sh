#!/usr/bin/env bash


interpreter=$(ps h -p $$ -o args='' | cut -f1 -d' ')

echo "I am running under $interpreter, the val I got was $1 and from the env: $NEW_VAL"