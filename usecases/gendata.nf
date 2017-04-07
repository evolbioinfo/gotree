params.outpath="data"
params.seed=1000
outpath=file(params.outpath)
outpath.with{mkdirs()}

seed=params.seed

process gentruetree {

	input:
	val seed

	output:
	file("true_tree.nw") into truetree

	shell:
	'''
	#!/usr/bin/env bash
	gotree generate yuletree -l 30 -s !{seed} | sed 's/Tip/Seq/g' > true_tree.nw
	'''
}

process genalign {
	input:
	file(truetree)
	val seed
	
	output:
	file("align.phy") into align

	shell:
	'''
	#!/usr/bin/env bash
	seq-gen -l 150 -mLG !{truetree} -z !{seed} | goalign reformat phylip -p -s > align.phy
	'''
}

align.subscribe{
	file -> file.copyTo(outpath.resolve(file.name))
}
