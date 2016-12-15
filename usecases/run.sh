nextflow run gendata.nf -with-dag gendata.dot
nextflow run bootstrap.nf -with-dag bootstrap.dot
nextflow run consensus.nf -with-dag consensus.dot
nextflow run compare.nf -with-dag compare.dot
