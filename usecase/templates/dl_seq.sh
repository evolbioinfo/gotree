#!/usr/bin/env bash

goalign rename -i "https://v100.orthodb.org/fasta?id=!{id}" \
	--regexp '([^\s]+).*' --replace '$1' \
	--unaligned \
	> sequences.fasta
sleep 2
