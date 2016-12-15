params.outpath="results"
params.align="data/align.phy"
params.nboot=10

align=file(params.align)
outpath=file(params.outpath)
outpath.with{mkdirs()}

process buildtree {
	input:
	file(align)

	output:
	file("tree.nw") into treeref

	shell:
	'''
	#!/usr/bin/env bash

	phyml -i !{align} -m LG -o tlr -b 0 -d aa
	mv !{align}_phyml_tree.txt tree.nw
	'''
}

process buildboots {
	input:
	file(align)

	output:
	file("bootalign_*") into bootaligns mode flatten

	shell:
	'''
	#!/usr/bin/env bash

	# Will generate bootstrap alignments
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

boottrees=boottree.collectFile(name: 'boottrees.nw')

process compare {
	input:
	file(ref) from treeref
	file(boot) from boottrees

	output:
	file("common.txt") into compare

	shell:
	'''
	#!/usr/bin/env bash
	gotree compare -i !{ref} -c !{boot} > common.txt
	'''
}

process hist {
	input:
	file(compare)

	output:
	file("*.png") into comparehist

	shell:
	'''
	#!/usr/bin/env Rscript

	comp=read.table("!{compare}",header=T)
	png("common.png")
	hist(comp$common/(comp$common+comp$specref),xlim=c(0,1),main="Distribution of distances",xlab="Distance")
	dev.off()
	'''
}

comparehist.subscribe{
	f -> f.copyTo(outpath.resolve(f.name))
}
