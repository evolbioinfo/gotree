params.dataurl="http://evolution.gs.washington.edu/book/primates.dna"
params.nboot = 100
params.seed=10000
params.outpath="results"
params.itolconfig= "data/itol_image_config.txt"

dataurl=params.dataurl
nboot = params.nboot
seed = params.seed
outpath = file(params.outpath)
itolconfig=file(params.itolconfig)

/**********************************/
/*     General tree inferences    */
/**********************************/

process downloadAlignment{
	input:
	val dataurl

	output:
	file "primates.phy" into truealign

	shell:
	'''
	wget -O primates_tmp.phy !{dataurl}
	goalign reformat phylip -p -i primates_tmp.phy --input-strict > primates.phy
	'''
}

process inferTrueTree{
	input:
	file align from truealign
	val seed

	output:
	file "truetree.nw" into truetree, truetree2, truetreedraw, truetreecopy

	shell:
	outfile="truetree.nw"
	template 'phyml.sh'
}

process simulAlign {
	input:
	file tree from truetree
	val seed

	output:
	file "align.phy" into simualign

	shell:
	'''
	#!/usr/bin/env bash
	seq-gen -l 50 -mLG -z !{seed} !{tree}  > align.phy
	'''
}

process reformatPhylip {
	input:
	file align from simualign

	output:
	file "align_clean.phy" into simualignphylip

	shell:
	'''
	#!/usr/bin/env bash
	goalign reformat phylip -p --input-strict -i !{align} > align_clean.phy
	'''
}

simualignphylip.into{refalign1; refalign2}

process inferReferenceTree{
	input:
	file align from refalign1
	val seed

	output:
	file "reftree.nw" into reftree

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
	gotree compute consensus -f 0.5 -i !{boot} -o consensus.nw
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
	file "support.nw" into supportcopy,supportdraw

	shell:
	'''
	#!/usr/bin/env bash
	gotree compute support classical -i !{ref} -b !{boot} -o support.nw
	'''
}

/***********************************/
/* Comparison of bootstrap trees   */
/*       With reference tree       */
/***********************************/
process compareTrees {
	input:
	file ref from truetree2
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
	file "*.png" into comparehist

	shell:
	'''
	#!/usr/bin/env Rscript

	comp=read.table("!{compare}",header=T)
	png("common.png")
	hist(comp$common*100/(comp$common+comp$reference),xlim=c(0,100),main="Distribution of distances",xlab="% Common branches")
	dev.off()
	'''
}


/**********************************************/
/*               Tree drawing                 */
/**********************************************/

// Reroot the trees to draw using an outgroup
process reroot{
	input:
	file tree from truetreedraw.mix(consensusdraw, supportdraw)

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
	file tree from treestodraw

	output:
	file "*.svg" into treeimages

	shell:
	'''
	#!/usr/bin/env bash
	gotree draw svg -i !{tree} -w 1000 -H 1000 --with-branch-support --support-cutoff 0.7 -o !{tree}.svg
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
	gotree dlimage itol -c !{itolconfig} -f svg -o !{tree}_itol.svg -i $ID
	'''
}

/*********************************************/
/*                File  COPY                 */
/*********************************************/
truetreecopy.mix(consensuscopy, supportcopy).subscribe{
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
