.. _clear-page:

gotree clear
============

**Informations**
----------------
Removes length or support information from an input tree.

**Usage**
---------

.. code:: bash

    Usage:
      gotree clear [command]
    Available Commands:
      lengths     Clear lengths from input trees
      supports    Clear supports from input trees
    Flags:
      -h, --help            help for clear
      -i, --input string    Input tree (default "stdin")
      -o, --output string   Cleared tree output file (default "stdout")
     
**Sub-commands**
----------------

gotree clear lengths
~~~~~~~~~~~~~~~~~~~~

- **Description** : Removes all branch lengths from an input tree.
  
- **Usage**:

.. code:: bash

   Usage:
      gotree clear lengths [flags]

      Global Flags:
        -i, --input string    Input tree (default "stdin")
        -o, --output string   Cleared tree output file (default "stdout")


- **Example**:

.. code-block:: bash

   gotree generate yuletree -l 20 -s 10 > tree.nw
   gotree generate yuletree -l 20 -s 10 | gotree clear lengths > tree2.nw


Should give a tree without branch lengths:

* tree.nw:

.. code::

   (((Tip16:0.235466,(Tip18:0.308786,Tip12:0.112378):0.043607):0.164845,Tip11:0.118784):0.010524,Tip0:0.044944,((Tip4:0.020616,(Tip7:0.097402,((Tip17:0.017246,Tip15:0.026931):0.002470,Tip2:0.134979):0.005444):0.129396):0.091234,((Tip8:0.027846,(Tip9:0.134926,Tip3:0.103093):0.005133):0.096048,((((Tip19:0.130964,Tip10:0.098771):0.017851,Tip6:0.061805):0.143344,(Tip14:0.004046,Tip5:0.077359):0.035620):0.029088,(Tip13:0.043410,Tip1:0.085651):0.064129):0.150752):0.022969):0.168367);

* tree2.nw:

.. code::

   (((Tip16,(Tip18,Tip12)),Tip11),Tip0,((Tip4,(Tip7,((Tip17,Tip15),Tip2))),((Tip8,(Tip9,Tip3)),((((Tip19,Tip10),Tip6),(Tip14,Tip5)),(Tip13,Tip1)))));

gotree clear supports
~~~~~~~~~~~~~~~~~~~~~~

- **Description** : Removes all branch supports from an input tree.
  
- **Usage**:

.. code:: bash

   Usage:
     gotree clear supports [flags]

     Global Flags:
       -i, --input string    Input tree (default "stdin")
       -o, --output string   Cleared tree output file (default "stdout")
