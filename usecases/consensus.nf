params.outpath="results"
params.align="data/align.phy"
params.nboot=100
params.itolconfig = "data/itol_image_config.txt"

align=file(params.align)
itolconfig=file(params.itolconfig)
outpath=file(params.outpath)
outpath.with{mkdirs()}

process buildboots {
	input:
	file(align)
	val nboot from params.nboot

	output:
	file("bootalign_*") into bootaligns mode flatten

	shell:
	'''
	#!/usr/bin/env bash

	# Will generate 1 outfile containing all alignments
	goalign build seqboot -n !{params.nboot} -i !{align} -p -o bootalign_ -S
	'''
}

process treeboot {
	input:
	file(align) from bootaligns

	output:
	file("boottree.nw") into boottree

	shell:
	'''
	#!/usr/bin/env bash
	phyml -i !{align} -m LG -o tlr -b 0 -d aa
	mv !{align}_phyml_tree.txt boottree.nw
	'''
}

boottrees = boottree.collectFile(name: 'boottrees.nw')

process consensus {
	input:
	file(boot) from boottrees

	output:
	file "consensus.nw" into consensus,consensus2,consensus3

	shell:
	'''
	#!/usr/bin/env bash
	gotree compute consensus -f 0.5 -i !{boot} -o consensus.nw
	'''
}

consensus.subscribe{
	f -> f.copyTo(outpath.resolve(f.name))
}

process drawConsensus {
	input:
	file consensus from consensus2

	output:
	file "*.svg" into consensusimage

	shell:
	'''
	#!/usr/bin/env bash
	gotree draw svg -i !{consensus} -r -w 1000 -H 1000 --with-branch-support --support-cutoff 0.8 -o consensus.svg
	'''
}

consensusimage.subscribe{
	f -> f.copyTo(outpath.resolve(f.name))
}

process uploadiTOL{
	input:
	file tree from consensus3
	file itolconfig

	output:
	file "*.txt" into iTOLurl
	file "*.svg" into iTOLimage

	shell:
	'''
	#!/usr/bin/env bash
	# Upload the tree
	gotree upload itol --name "consensustree" -i !{tree} > consensus_url.txt
	# We get the iTOL id
	ID=$(basename $(cat consensus_url.txt ))
	# We Download the image with options defined in data/itol_image_config.txt
	gotree dlimage itol -c !{itolconfig} -f svg -o consensus_itol.svg -i $ID
	'''
}

iTOLurl.subscribe{
	f -> f.copyTo(outpath.resolve(f.name))
}

iTOLimage.subscribe{
	f -> f.copyTo(outpath.resolve(f.name))
}
