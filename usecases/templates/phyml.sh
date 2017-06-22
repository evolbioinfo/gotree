#!/usr/bin/env bash

phyml -i !{align} -m LG -o tlr -b 0 -d aa
mv !{align}_phyml_tree* !{outfile}
