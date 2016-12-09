.. _introduction-page:

*******************
Introduction
*******************

.. _introduction-general:

General
=======

GoTree integrates a set of utilities that aim at facilitating phylogenetic tree manipulation via commandline and especially unix shell. To facilitate integration in complex workflows, GoTree was written to be able to read from stdin and print to stdout. Moreover, GoTree is implemented in Go and is thus portable to most platforms (Linux, Windows, MacOs).

Currently, GoTree accepts and generates trees in Newick format.

.. _introduction-examples:

GoTree basic usage
==================



.. _introduction-toollist:

List of commands
================

================================ ======================================================================================
command                           Description
================================ ======================================================================================
**clear**                            Clear lengths or supports from input trees
**collapse**                         Collapse branches of input trees
**compare**                          Compare a reference tree with a set of trees
**compute**                          Computations such as consensus and supports
**difftips**                         Print diff between tip names of two trees
**divide**                           Divide an input tree file into several tree files
**generate**                         Generate random trees
**minbrlen**                         Set a min branch length to all branches with length < cutoff
**prune**                            Remove tips of the input tree that are not in the compared tree
**randbrlen**                        Assign a random length to edges of input trees
**rename**                           Rename tips of the input tree, given a map file
**reroot**                           Reroot trees using an outgroup
**shuffletips**                      Shuffle tip names of an input tree
**stats**                            Print statistics about the tree
**unroot**                           Unroot input trees
**version**                          Print version of gotree
================================ ======================================================================================
