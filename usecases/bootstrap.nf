params.outpath="results"
params.align="data/align.phy"
params.nboot=100
params.itolconfig = "data/itol_image_config.txt"

align=file(params.align)
itolconfig=file(params.itolconfig)
outpath=file(params.outpath)
outpath.with{mkdirs()}

process buildtree {
	input:
	file(align)
	val nboot from params.nboot

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

process supports {
	input:
	file(ref) from treeref
	file(boot) from boottrees

	output:
	file "support.nw" into support, support2, support3

	shell:
	'''
	#!/usr/bin/env bash
	gotree compute support classical -i !{ref} -b !{boot} -o support.nw
	'''
}

support.subscribe{
	f -> f.copyTo(outpath.resolve(f.name))
}

process drawTree{
	input:
	file tree from support2

	output:
	file "*.svg" into supportimage

	shell:
	'''
	#!/usr/bin/env bash
	gotree draw svg -i !{tree} -r -w 1000 -H 1000 --with-branch-support --support-cutoff 0.8 -o support.svg
	'''
}

supportimage.subscribe{
	f -> f.copyTo(outpath.resolve(f.name))
}

process uploadiTOL{
	input:
	file tree from support3
	file itolconfig

	output:
	file "*.txt" into iTOLurl
	file "*.svg" into iTOLimage

	shell:
	'''
	#!/usr/bin/env bash
	# Upload the tree
	gotree upload itol --name "supporttree" -i !{tree} > support_url.txt
	# We get the iTOL id
	ID=$(basename $(cat support_url.txt ))
	# We Download the image with options defined in data/itol_image_config.txt
	gotree dlimage itol -c !{itolconfig} -f svg -o support_itol.svg -i $ID
	'''
}

iTOLurl.subscribe{
	f -> f.copyTo(outpath.resolve(f.name))
}

iTOLimage.subscribe{
	f -> f.copyTo(outpath.resolve(f.name))
}
