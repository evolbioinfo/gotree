.. _clear-page:

gotree reroot
=============

**Informations**
----------------
This command reroots an input tree. Two methods are implemented:
- Rooting using an outgroup;
- Rooting at the midpoint.


**Usage**
---------

.. code:: bash

   Usage:
     gotree reroot [command]
   
   Available Commands:
     midpoint    Reroot trees at midpoint
     outgroup    Reroot trees using an outgroup
   
   Flags:
     -h, --help            help for reroot
     -i, --input string    Input Tree (default "stdin")
     -o, --output string   Rerooted output tree file (default "stdout")

**Sub-commands**
----------------

gotree reroot outgroup
~~~~~~~~~~~~~~~~~~~~~~

- **Description** : Reroots an input tree using an outgoup. It first removes the current root if it exists, take the last common ancestor of all given tips, and reroot the tree on the incomming branch.
  
- **Usage**:

.. code:: bash

    Usage:
      gotree reroot outgroup [flags]
    
    Flags:
      -l, --tip-file string   File containing names of tips of the outgroup (default "stdin")
    
    Global Flags:
      -i, --input string    Input Tree (default "stdin")
      -o, --output string   Rerooted output tree file (default "stdout")

- **Example**:

.. code-block:: bash

    gotree generate balancedtree -d 4 > tree.nw
    echo "Tip6" | gotree reroot outgroup -i tree.nw -l - > reroot.nw

Should give:

.. image:: ../images/reroot_outgroup.png
    

gotree reroot midpoint
~~~~~~~~~~~~~~~~~~~~~~

- **Description** : Reroots an input tree using at the midpoint position (middle of the longest path between tips).It first removes the current root if it exists, take the longest path between two tips, and place the root at the middle.
  
- **Usage**:

.. code:: bash

    Usage:
      gotree reroot midpoint [flags]
    
    Global Flags:
      -i, --input string    Input Tree (default "stdin")
      -o, --output string   Rerooted output tree file (default "stdout")

- **Example**:

.. code-block:: bash

    gotree generate yuletreee -l 10 > tree.nw
    gotree reroot midpoint -i tree.nw > reroot.nw

Can give:

.. image:: ../images/reroot_midpoint.png
    
