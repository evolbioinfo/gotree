.. _getstarted-page:

*******************
Getting started
*******************

.. _getstarted-requirements:

Requirements
============
To install GoTree executables, you do not need to install any dependency. On the other hand, to install GoTree from sources, the only requirement is that you need to download_ and install_ Go on your system.

.. _download: https://golang.org/dl/
.. _install: https://golang.org/doc/install

.. _getstarted-install:

Installation
============

From binaries
~~~~~~~~~~~~~
You can download already compiled binaries for the latest release in the release-section_ on github. Binaries are available for MacOS, Linux, and Windows (32 and 64 bits).

Once downloaded, you can just run the executable without any other downloads.

.. _release-section: https://github.com/fredericlemoine/gotree/releases

From sources
~~~~~~~~~~~~

If Go is already installed on your system, then you just have to type :

.. code:: bash

    go get github.com/fredericlemoine/gotree/

This will download GoTree sources from github, and all its dependencies.

You can then build it with:

.. code:: bash

   cd $GOPATH/src/github.com/fredericlemoine/gotree/
   make

The :code:`gotree` executable should be located in the :code:`$GOPATH/bin` folder.

.. _getstarted-firsttry:

First try
==========

gotree implements several tree manipulation commands. Here are some short examples:

* Generate random unrooted uniform binary trees

.. code:: bash
	  
   $ gotree generate uniformtree -l 100 -n 10 | gotree stats


* Unrooting a tree

.. code:: bash

   $ gotree unroot -i tree.tre -o unrooted.tre

* Collapsing short branches

.. code:: bash

   $ gotree collapse length -i tree.tre -l 0.001 -o collapsed.tre

* Collapsing lowly supported branches

.. code:: bash

   $ gotree collapse support -i tree.tre -s 0.8 -o collapsed.tre

* Clearing length information

.. code:: bash

   $ gotree clear lengths -i tree.nw -o nolength.nw

* Clearing support information

.. code:: bash

   $ gotree clear supports -i tree.nw -o nosupport.nw

Note that you can pipe the two previous commands:

.. code:: bash

   $ gotree clear supports -i tree.nw | gotree clear lengths -o nosupport.nw

* Printing tree statistics

.. code:: bash

   $ gotree stats -i tree.tre

* Printing edge statistics

.. code:: bash

   $ gotree stats edges -i tree.tre

Example of result:

====== ======== ============ =========== ============ ========= ============= =============
tree     brid     length       support     terminal     depth     topodepth     rightname
====== ======== ============ =========== ============ ========= ============= =============
0        0        0.107614     N/A         false        1         6             
0        1        0.149560     N/A         true         0         1             Tip51
0        2        0.051126     N/A         false        1         5             
0        3        0.003992     N/A         false        1         4             
0        4        0.030974     N/A         false        1         3             
0        5        0.270017     N/A         true         0         1             Tip84
0        6        0.029931     N/A         false        1         2             
0        7        0.001136     N/A         true         0         1             Tip70
0        8        0.011658     N/A         true         0         1             Tip45
0        9        0.104188     N/A         true         0         1             Tip34
0        10       0.003361     N/A         true         0         1             Tip16
0        11       0.021988     N/A         true         0         1             Node0
====== ======== ============ =========== ============ ========= ============= =============

* Printing tips

.. code:: bash

   $ gotree stats tips -i tree.tre

Example of result:

====== ======  =======  ======
tree 	id 	nneigh 	name
====== ======  =======  ======
0 	1 	1 	Tip8
0 	2 	1 	Node0
0 	5 	1 	Tip4
0 	8 	1 	Tip9
0 	9 	1 	Tip7
0 	11 	1 	Tip6
0 	13 	1 	Tip5
0 	14 	1 	Tip3
0 	16 	1 	Tip2
0 	17 	1 	Tip1
====== ======  =======  ======

* Comparing tips of two trees

.. code:: bash

   $ gotree difftips -i tree.tre -c tree2.tre

This will compare the two sets of tips.

Example:

.. code:: bash

   $ gotree difftips -i <(gotree generate uniformtree -l 10 -n 1) \
	  -c <(gotree generate uniformtree -l 11 -n 1)
   > Tip10
   = 10

10 tips are equal, and "Tip10" is present only in the second tree.

* Removing tips that are absent from another tree

.. code:: bash

   $ gotree prune -i tree.tre -c other.tre -o pruned.tre

You can test with

.. code:: bash

   $ gotree prune -i <(gotree generate uniformtree -l 1000 -n 1) \
	  -c <(gotree generate uniformtree -l 100 -n 1) \
          | gotree stats

It should print 100 tips.

* Comparing bipartitions

Count the number of common/specific bipartitions between two trees.

.. code:: bash

   $ gotree compare -i tree.tre -c other.tre

You can test with random trees (there should be very few common bipartitions)

.. code:: bash

   $ gotree compare -i <(gotree generate uniformtree -l 100 -n 1) \
	  -c <(gotree generate uniformtree -l 100 -n 1)

==== ======= ======
Tree specref common
==== ======= ======
0         97      0
==== ======= ======

* Renaming tips of the tree

If you have a file containing the mapping between current names and new names of the tips, you can rename the tips:

.. code:: bash

   $ gotree rename -i tree.tre -m mapfile.txt -o newtree.tre

You can try by doing:

.. code:: bash

   $ gotree generate uniformtree -l 100 -n 1 -o tree.tre
   $ gotree stats tips -i tree.tre | awk '{if(NR>1){print $4 "\tNEWNAME" $4}}' > mapfile.txt
   $ gotree rename -i tree.tre -m mapfile.txt | gotree stats tips
