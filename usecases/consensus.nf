params.outpath="results"
params.align="data/align.phy"
params.nboot=10

align=file(params.align)
outpath=file(params.outpath)
outpath.with{mkdirs()}

process buildboots {
	input:
	file(align)

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
	file("consensus.nw") into consensus

	shell:
	'''
	#!/usr/bin/env bash
	gotree compute consensus -f 0.5 -i !{boot} -o consensus.nw
	'''
}

consensus.subscribe{
	f -> f.copyTo(outpath.resolve(f.name))
}
