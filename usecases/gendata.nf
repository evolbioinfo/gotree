params.outpath="data"
outpath=file(params.outpath)
outpath.with{mkdirs()}

process gentruetree {

	output:
	file("true_tree.nw") into truetree

	shell:
	'''
	#!/usr/bin/env bash
	gotree generate yuletree -l 30 | sed 's/Tip/Seq/g' > true_tree.nw
	'''
}

process genalign {
	input:
	file(truetree)
	
	output:
	file("align.phy") into align

	shell:
	'''
	#!/usr/bin/env bash
	seq-gen -l 150 -mLG !{truetree} | goalign reformat phylip -p -s > align.phy
	'''
}

align.subscribe{
	file -> file.copyTo(outpath.resolve(file.name))
}
