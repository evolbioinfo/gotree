params.nboot = 100
params.seed=2000
params.outpath="results"
params.itolconfig= "data/itol_image_config.txt"
//params.mapfile="data/mapfile.txt"

nboot = params.nboot
seed = params.seed
outpath = file(params.outpath)
itolconfig=file(params.itolconfig)
//mapfile=file(params.mapfile)

/************************************/
/*       Get OrthoDB sequences      */
/* That are present in all primates */
/*      and unique (no paralogs)    */
/*   with "ribosomal" keyword       */
/*         199 Proteins             */
/************************************/
process getOrthoDBIds{
	publishDir "${outpath}", mode: 'copy'

	output:
	file "ids.txt" into allorthoids
	
	shell:
	'''
	wget -O orthoids.json "https://v100.orthodb.org/search?query=ribosomal&level=9443&species=9443&universal=1&singlecopy=1"
	jq '.data | join(",")' orthoids.json | sed 's/"//g' | sed 's/,/\\n/g' | sort > ids.txt
	'''
}

/* Split line by line */
protids = allorthoids.splitText().map{ it -> it.trim() }

/***********************************/
/* Download sequences and metadata */
/***********************************/
process downloadSequences{
	maxForks 1
	tag "${id}"

	input:
	val id from protids
	
	output:
	set val(id), file("sequences.fasta") into sequences

	shell:
	template "dl_seq.sh"
}

process getMapTable{
	maxForks 1
	tag "${id}"

	input:
	set val(id), file(seq) from sequences

	output:
	set val(id), file(seq), file("map.txt") into mapfile
	file "gene.txt" into genefile

	shell:
	'''
	wget -O align.tab https://v100.orthodb.org/tab?query=!{id}
	cut -f 5,6 align.tab > map.txt
	cut -f 1,2 align.tab | tail -n+2 | sort -u > gene.txt
	'''
}

genefile.collectFile(name: 'genes.txt').subscribe{file -> file.copyTo(outpath.resolve(file.name))}

/***************************/
/*  Cleaning, reformating  */
/* and aligning sequences  */
/***************************/
process renameSequences{
	tag "${id}"

	input:
	set val(id), file(sequences), file(mapfile) from mapfile

	output:
	file "renamed.fasta" into renamed

	shell:
	'''
	goalign rename -r -m !{mapfile} -i !{sequences} --unaligned | goalign rename --regexp " " --replace "_" --unaligned  > renamed.fasta
	'''
}

process cleanSequences{
	input:
	file(sequences) from renamed

	output:
	file "cleaned.fasta" into cleaned

	shell:
	'''
	goalign replace -s U -n X -i !{sequences} -o cleaned.fasta --unaligned
	'''
}

process alignSequences{
	input:
	file cleaned

	output:
	file "aligned.fasta" into alignment

	shell:
	'''
	mafft --quiet !{cleaned} > aligned.fasta
	'''
}

process concatSequences {
	input:
	file 'align_fasta' from alignment.toList()

	output:
	file "concat.fasta" into concat

	shell:
	'''
	goalign concat -o concat.fasta -i none align_fasta*
	'''
}

process cleanAlign {
	input:
	file align from concat

	output:
	file "cleanalign.fasta" into cleanalign

	shell:
	'''
	BMGE -i !{align} -t AA -m BLOSUM62 -w 3 -g 0.2 -h 0.5  -b 5 -of cleanalign.fasta
	'''
}

process reformatAlignment{
	input:
	file cleanalign

	output:
	file "aligned.phylip" into alignmentphylip

	shell:
	'''
	goalign reformat phylip -i !{cleanalign} -o aligned.phylip
	'''
}

/***************************/
/*     Inferring tree      */
/***************************/
process inferTrueTree{
	publishDir "${outpath}", mode: 'copy'

	input:
	file align from alignmentphylip
	val seed

	output:
	file "tree.nw" into tree, tree2

	shell:
	'''
	iqtree -s !{align} -seed !{seed} -m MFP -b 100 -nt !{task.cpus}
	mv *.treefile tree.nw
	'''
}

process rerootTree{
	publishDir "${outpath}", mode: 'copy'

	input:
	file tree from tree2

	output:
	file "rerooted.nw" into reroottree1, reroottree2

	shell:
	'''
	gotree reroot outgroup -i !{tree} -o rerooted.nw Otolemur_garnettii Microcebus_murinus Propithecus_coquereli 
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

process pruneNCBITax {

	input:
	file tree from tree
	file ncbi from ncbitax

	output:
	file "ncbi_pruned.nw" into ncbipruned

	shell:
	'''
	#!/usr/bin/env bash
	gotree prune -i !{ncbi} -c !{tree} -o ncbi_pruned.nw
	'''
}

process rerootNCBITax {

	input:
	file tree from ncbipruned

	output:
	file "ncbi_rerooted.nw" into ncbirerooted1, ncbirerooted2

	shell:
	'''
	#!/usr/bin/env bash
	gotree reroot outgroup -i !{tree} -o ncbi_rerooted.nw Otolemur_garnettii Microcebus_murinus Propithecus_coquereli 
	'''
}

process annotateTree{
	publishDir "${outpath}", mode: 'copy'

	input:
	file tree from reroottree1
	file ncbi from ncbirerooted1

	output:
	file "annotated.nw" into annotated

	shell:
	'''
	#!/usr/bin/env bash
	gotree annotate -i !{tree} -c !{ncbi} -o annotated.nw
	'''
}

process compareTrees{
	publishDir "${outpath}", mode: 'copy'

	input:
	file tree from reroottree2
	file ncbi from ncbirerooted2

	output:
	file "tree_comparison.txt" into comparison

	shell:
	'''
	#!/usr/bin/env bash
	gotree compare trees -i !{tree} -c !{ncbi} > tree_comparison.txt
	'''
}


/**********************************/
/*      Uploading tree to ITol    */
/*       Downloading the image    */
/**********************************/
process uploadTree{
	publishDir "${outpath}", mode: 'copy'

	input:
	file tree from annotated
	file itolconfig

	output:
	file "tree_url.txt" into iTOLurl
	file "tree_image.svg" into iTOLimage

	shell:
	'''
	#!/usr/bin/env bash
	# Upload the tree
	gotree upload itol --name "AnnotatedTree" -i !{tree} > tree_url.txt
	# We get the iTOL id
	ID=$(basename $(cat tree_url.txt ))
	# We Download the image with options defined in data/itol_image_config.txt
	gotree download itol -c !{itolconfig} -f svg -o tree_image.svg -i $ID
	'''
}
