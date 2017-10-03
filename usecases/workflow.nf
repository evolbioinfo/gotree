params.dataurl="http://evolution.gs.washington.edu/book/primates.dna"
params.nboot = 300
params.seed=2000
params.outpath="results"
params.itolconfig= "data/itol_image_config.txt"
params.mapfile="data/mapfile.txt"

dataurl=params.dataurl
nboot = params.nboot
seed = params.seed
outpath = file(params.outpath)
itolconfig=file(params.itolconfig)
mapfile=file(params.mapfile)

/**********************************/
/*     General tree inferences    */
/**********************************/

process downloadAlignment{
	input:
	val dataurl

	output:
	file "primates.phy" into refalign, refalign2, refaligncopy

	shell:
	'''
	wget -O primates_tmp.phy !{dataurl}
	goalign reformat phylip -p -i primates_tmp.phy --input-strict > primates.phy
	'''
}

process inferRefTree{
	input:
	file align from refalign
	val seed

	output:
	file "reftree.nw" into reftree, reftree2, reftreedraw, reftreecopy

	shell:
	outfile="reftree.nw"
	template 'phyml.sh'
}


process seqBoots {
	input:
	file align from refalign2
	val nboot
	val seed

	output:
	file "bootalign_*" into bootaligns mode flatten

	shell:
	'''
	#!/usr/bin/env bash

	# Will generate 1 outfile containing all alignments
	goalign build seqboot -n !{nboot} -i !{align} -p -o bootalign_ -S -s !{seed}
	'''
}

process inferBootstrapTrees{
	input:
	file align from bootaligns
	val seed
	
	output:
	file "boot.nw" into boottree

	shell:
	outfile="boot.nw"
	template 'phyml.sh'
}

boottree.collectFile(name: 'boottrees.nw').into{boottrees1; boottrees2; boottrees3}

/**********************************/
/*      Consensus computation     */
/**********************************/
process consensus {
	input:
	file boot from boottrees1

	output:
	file "consensus.nw" into consensuscopy,consensusdraw

	shell:
	'''
	#!/usr/bin/env bash
	gotree compute consensus -f 0.6 -i !{boot} -o consensus.nw
	'''
}

/**********************************/
/*      Bootstrap supports        */
/**********************************/
process supports {
	input:
	file ref from reftree
	file boot from boottrees2

	output:
	file "support.nw" into supportcopy, supportdraw, supportannot

	shell:
	'''
	#!/usr/bin/env bash
	gotree compute support classical -i !{ref} -b !{boot} -o support.nw
	'''
}

/**********************************/
/*  Comparison with NCBI taxonomy */
/**********************************/
process downloadNewickTaxonomy {
	output:
	file "ncbi.nw" into ncbitax

	shell:
	'''
	#!/usr/bin/env bash
	gotree download ncbitax -o ncbi.nw
	'''
}

process rerootSupport{
	input:
	file support from supportannot

	output:
	file "rerooted_support.nw" into rerootedsupportncbi

	shell:
	'''
	#!/usr/bin/env bash
	gotree reroot outgroup -i !{support} Mouse Bovine > rerooted_support.nw
	'''
}

process renameSupport {
	input:
	file tree from rerootedsupportncbi
	file mapfile from mapfile
	
	output:
	file "renamed_support.nw" into renamedsupport1, renamedsupport2

	shell:
	'''
	#!/usr/bin/env bash	
	gotree rename -i !{tree} -m !{mapfile} -o renamed_support.nw
	'''
}

process pruneNCBITax {

	input:
	file tree from renamedsupport1
	file ncbi from ncbitax

	output:
	file "ncbi_pruned.nw" into ncbipruned

	shell:
	'''
	#!/usr/bin/env bash
	gotree prune -i !{ncbi} -c !{tree} -o ncbi_pruned.nw
	'''
}

process annotateSupportTree{
	input:
	file renamed from renamedsupport2
	file ncbi from ncbipruned

	output:
	file "annotated_support.nw" into annotatedsupport, annotatedsupportcopy

	shell:
	'''
	#!/usr/bin/env bash
	gotree annotate -i !{renamed} -c !{ncbi} -o annotated_support.nw
	'''
}


/***********************************/
/* Comparison of bootstrap trees   */
/*       With reference tree       */
/***********************************/
process compareTrees {
	input:
	file ref from reftree2
	file boot from boottrees3

	output:
	file "common.txt" into compare

	shell:
	'''
	#!/usr/bin/env bash
	gotree compare trees -i !{ref} -c !{boot} > common.txt
	'''
}

process histCommonbranches {
	input:
	file compare 

	output:
	file "*.svg" into comparehist

	shell:
	'''
	#!/usr/bin/env Rscript

	comp=read.table("!{compare}",header=T)
	svg("common.svg",width=14,height=7)
	hist(100-comp$common*100/(comp$common+comp$reference),main="Distribution of distances",xlab="% Common branches")
	dev.off()
	'''
}


/**********************************************/
/*               Tree drawing                 */
/**********************************************/

// Reroot the trees to draw using an outgroup
process reroot{
	input:
	file tree from reftreedraw.mix(consensusdraw, supportdraw)

	output:
	file "${tree.baseName}_reroot.nw" into treestodraw, treestodrawitol

	shell:
	'''
	#/usr/bin/env bash
	gotree reroot outgroup -i !{tree} Mouse Bovine > !{tree.baseName}_reroot.nw
	'''
}

process drawTree {
	input:
	file tree from treestodraw.mix(annotatedsupport)

	output:
	file "*.svg" into treeimages

	shell:
	'''
	#!/usr/bin/env bash
	gotree draw svg -i !{tree} -w 1000 -H 1000 --with-branch-support --with-node-labels --support-cutoff 0.7 -o !{tree}.svg
	'''
}

process uploadiTOL{
	input:
	file tree from treestodrawitol
	file itolconfig

	output:
	file "*.txt" into iTOLurl
	file "*.svg" into iTOLimage

	shell:
	'''
	#!/usr/bin/env bash
	# Upload the tree
	gotree upload itol --name "consensustree" -i !{tree} > !{tree}_url.txt
	# We get the iTOL id
	ID=$(basename $(cat !{tree}_url.txt ))
	# We Download the image with options defined in data/itol_image_config.txt
	gotree download itol -c !{itolconfig} -f svg -o !{tree}_itol.svg -i $ID
	'''
}


/*********************************************/
/*                File  COPY                 */
/*********************************************/
reftreecopy.mix(consensuscopy, supportcopy, annotatedsupportcopy, refaligncopy).subscribe{
	f -> f.copyTo(outpath.resolve(f.name))
}
treeimages.subscribe{
	f -> f.copyTo(outpath.resolve(f.name))
}

iTOLurl.subscribe{
	f -> f.copyTo(outpath.resolve(f.name))
}

iTOLimage.subscribe{
	f -> f.copyTo(outpath.resolve(f.name))
}

comparehist.subscribe{
	f -> f.copyTo(outpath.resolve(f.name))
}
