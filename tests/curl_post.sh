#!/bin/sh

curl -X POST "http://localhost:8001/?sl=it&tl=en" --data '{"0" : "ciao come stai", "1" : "va tutto bene?" }'
