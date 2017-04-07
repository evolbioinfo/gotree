nextflow run gendata.nf -with-dag gendata.dot -with-timeline gendata.html -resume
nextflow run bootstrap.nf -with-dag bootstrap.dot -with-timeline bootstrap.html -resume
nextflow run consensus.nf -with-dag consensus.dot -with-timeline consensus.html -resume
nextflow run compare.nf -with-dag compare.dot -with-timeline compare.html
