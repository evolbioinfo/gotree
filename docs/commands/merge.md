# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### merge
This command merges two rooted trees.

If one of the tree is not rooted, returns an error
Tip names must be different between the two trees, otherwise returns an error

Edges connecting new root with old roots have length of 1.0.

#### Usage

```
Usage:
  gotree merge [flags]

Flags:
  -c, --compared string   Compared tree input file (default "stdin")
  -o, --output string     Merged tree output file (default "stdout")
  -i, --reftree string    Reference tree input file (default "stdin")
```

#### Example

* Merging two trees from newick format

```
# bash short syntax
	gotree merge -i <(echo "(Tip0,(Tip3,(Tip2,Tip1)0.2)0.9);") \
                 -c <(echo "(Tip0_2,(Tip3_2,(Tip2_2,Tip1_2)0.2)0.9);") \
	| gotree draw text -w 20
# Or all steps
echo "(Tip0,(Tip3,(Tip2,Tip1)0.2)0.9);" > t1
echo "(Tip0_2,(Tip3_2,(Tip2_2,Tip1_2)0.2)0.9);" > t2
gotree merge -i t1 -c t2 -o merged
gotree draw text -i merged -w 20
```

It should give the following tree:
```
     +--- Tip0                
+--- |                        
|    |    +--- Tip3           
|    +--- |                   
|         |    +--- Tip2      
|         +--- |              
|              +--- Tip1      
|                             
|    +--- Tip0_2              
+--- |                        
     |    +--- Tip3_2         
     +--- |                   
          |    +--- Tip2_2    
          +--- |              
               +--- Tip1_2    

```
