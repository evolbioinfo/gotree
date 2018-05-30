########### Test Suite for Gotree command line tools ##############

set -e
set -u
set -o pipefail

# gotree annotate
echo "->gotree annotate"
cat > mapfile <<EOF
internal1:Tip6,Tip5,Tip1
EOF
cat > expected <<EOF
((Tip4,(Tip7,Tip2)),Tip0,((Tip8,(Tip9,Tip3)),((Tip6,Tip5),Tip1)internal1));
EOF
gotree generate yuletree -s 10 | gotree brlen clear | gotree annotate -m mapfile > result
diff result expected
rm -f expected result mapfile

# gotree annotate
echo "->gotree annotate 2"
cat > mapfile <<EOF
ACGACTCATCTA:Tip6,Tip5,Tip1
ACGACTCATCTA:internal1
EOF
cat > intree <<EOF
((Tip4,(Tip7,Tip2)),Tip0,((Tip8,(Tip9,Tip3))internal1,((Tip6,Tip5),Tip1)));
EOF
cat > expected <<EOF
((Tip4,(Tip7,Tip2)),Tip0,((Tip8,(Tip9,Tip3))ACGACTCATCTA,((Tip6,Tip5),Tip1)ACGACTCATCTA));
EOF
gotree annotate -i intree -m mapfile > result
diff result expected
rm -f expected result intree mapfile

# gotree annotate
echo "->gotree annotate comment"
cat > mapfile <<EOF
internal1:Tip6,Tip5,Tip1
EOF
cat > expected <<EOF
((Tip4,(Tip7,Tip2)),Tip0,((Tip8,(Tip9,Tip3)),((Tip6,Tip5),Tip1)[internal1]));
EOF
gotree generate yuletree -s 10 | gotree brlen clear | gotree annotate --comment -m mapfile > result
diff result expected
rm -f expected result mapfile

# gotree annotate
echo "->gotree annotate comment 2"
cat > mapfile <<EOF
ACGACTCATCTA:Tip6,Tip5,Tip1
ACGACTCATCTA:internal1
EOF
cat > intree <<EOF
((Tip4,(Tip7,Tip2)),Tip0,((Tip8,(Tip9,Tip3))internal1,((Tip6,Tip5),Tip1)));
EOF
cat > expected <<EOF
((Tip4,(Tip7,Tip2)),Tip0,((Tip8,(Tip9,Tip3))internal1[ACGACTCATCTA],((Tip6,Tip5),Tip1)[ACGACTCATCTA]));
EOF
gotree annotate -i intree --comment -m mapfile > result
diff result expected
rm -f expected result intree mapfile

# gotree annotate with tree
echo "->gotree annotate with tree"
cat > intree <<EOF
((Tip4,(Tip7,Tip2)n1)n2,Tip0,((Tip8,(Tip9,Tip3)n3)n4,((Tip6,Tip1)n5,Tip5)n6)n7);
EOF
cat > expected <<EOF
((Tip4,(Tip7,Tip2)n1_0_2)n2_0_3,Tip0,((Tip8,(Tip9,Tip3)n3_0_2)n4_0_3,((Tip6,Tip5)Tip6_1_2,Tip1)n6_0_3)n7_0_6);
EOF
gotree generate yuletree -s 10 | gotree brlen clear | gotree annotate -c intree > result
diff result expected
rm -f expected result intree

# gotree annotate with tree
echo "->gotree annotate with tree comments"
cat > intree <<EOF
((Tip4,(Tip7,Tip2)n1)n2,Tip0,((Tip8,(Tip9,Tip3)n3)n4,((Tip6,Tip1)n5,Tip5)n6)n7);
EOF
cat > expected <<EOF
((Tip4,(Tip7,Tip2)[n1_0_2])[n2_0_3],Tip0,((Tip8,(Tip9,Tip3)[n3_0_2])[n4_0_3],((Tip6,Tip5)[Tip6_1_2],Tip1)[n6_0_3])[n7_0_6]);
EOF
gotree generate yuletree -s 10 | gotree brlen clear | gotree annotate --comment -c intree > result
diff result expected
rm -f expected result intree


# gotree brlen clear
echo "->gotree brlen clear"
cat > expected <<EOF
((Tip4,(Tip7,Tip2)),Tip0,((Tip8,(Tip9,Tip3)),((Tip6,Tip5),Tip1)));
(Tip5,Tip0,((Tip6,(Tip7,Tip4)),(Tip2,((Tip8,(Tip9,Tip3)),Tip1))));
(Tip6,Tip0,((((Tip5,Tip4),((Tip9,Tip8),Tip3)),(Tip7,Tip2)),Tip1));
(Tip3,Tip0,(((Tip6,Tip5),Tip2),(((Tip9,Tip8),Tip4),(Tip7,Tip1))));
(((Tip4,(Tip8,Tip3)),((Tip6,Tip5),Tip2)),Tip0,(Tip7,(Tip9,Tip1)));
(Tip4,Tip0,(Tip3,((((Tip8,Tip7),Tip5),((Tip9,Tip6),Tip2)),Tip1)));
(Tip5,Tip0,((Tip8,Tip2),((Tip7,Tip3),((Tip9,Tip4),(Tip6,Tip1)))));
((Tip9,Tip6),Tip0,((((Tip8,Tip4),Tip3),Tip2),((Tip7,Tip5),Tip1)));
(Tip9,Tip0,(Tip8,((Tip7,Tip4),((Tip6,Tip2),(Tip3,(Tip5,Tip1))))));
((((Tip7,(Tip8,(Tip9,Tip6))),Tip5),Tip4),Tip0,(Tip2,(Tip3,Tip1)));
EOF
gotree generate yuletree -s 10 -n 10 | gotree brlen clear > result
diff result expected
rm -f expected result

# gotree brlen scale
echo "->gotree brlen scale"
cat > input.tree <<EOF
((Tip4:1,(Tip7:1,Tip2:1):1):1,Tip0:1,((Tip8:1,(Tip9:1,Tip3:1):1):1,((Tip6:1,Tip5:1):1,Tip1:1):1):1);
(Tip5,Tip0,((Tip6,(Tip7,Tip4)),(Tip2,((Tip8,(Tip9,Tip3)),Tip1))));
(Tip6:0.5,Tip0:0.5,((((Tip5:0.5,Tip4:0.5):0.5,((Tip9:0.5,Tip8:0.5),Tip3:0.5):0.5):0.5,(Tip7:0.5,Tip2:0.5):0.5):0.5,Tip1:0.5):0.5);
EOF
cat > expected <<EOF
((Tip4:3,(Tip7:3,Tip2:3):3):3,Tip0:3,((Tip8:3,(Tip9:3,Tip3:3):3):3,((Tip6:3,Tip5:3):3,Tip1:3):3):3);
(Tip5,Tip0,((Tip6,(Tip7,Tip4)),(Tip2,((Tip8,(Tip9,Tip3)),Tip1))));
(Tip6:1.5,Tip0:1.5,((((Tip5:1.5,Tip4:1.5):1.5,((Tip9:1.5,Tip8:1.5),Tip3:1.5):1.5):1.5,(Tip7:1.5,Tip2:1.5):1.5):1.5,Tip1:1.5):1.5);
EOF
gotree brlen scale -i input.tree -f 3.0 > result
diff result expected
rm -f expected result input.tree

# gotree support scale
echo "->gotree support scale"
cat > input.tree <<EOF
((a,b)99,(c,d)50,(e,f)0,g);
EOF
cat > expected.1 <<EOF
((a,b)0.99,(c,d)0.5,(e,f)0,g);
EOF
cat > expected.2 <<EOF
((a,b)198,(c,d)100,(e,f)0,g);
EOF
gotree support scale -i input.tree -f 0.01 > result.1
gotree support scale -i input.tree -f 2 > result.2
diff result.1 expected.1
diff result.2 expected.2
rm -f expected.1 expected.2 result.1 result.2 input.tree


# gotree support clear
echo "->gotree support clear"
cat > expected <<EOF
((Tip4,(Tip7,Tip2)),Tip0,((Tip8,(Tip9,Tip3)),((Tip6,Tip5),Tip1)));
(Tip5,Tip0,((Tip6,(Tip7,Tip4)),(Tip2,((Tip8,(Tip9,Tip3)),Tip1))));
(Tip6,Tip0,((((Tip5,Tip4),((Tip9,Tip8),Tip3)),(Tip7,Tip2)),Tip1));
(Tip3,Tip0,(((Tip6,Tip5),Tip2),(((Tip9,Tip8),Tip4),(Tip7,Tip1))));
(((Tip4,(Tip8,Tip3)),((Tip6,Tip5),Tip2)),Tip0,(Tip7,(Tip9,Tip1)));
(Tip4,Tip0,(Tip3,((((Tip8,Tip7),Tip5),((Tip9,Tip6),Tip2)),Tip1)));
(Tip5,Tip0,((Tip8,Tip2),((Tip7,Tip3),((Tip9,Tip4),(Tip6,Tip1)))));
((Tip9,Tip6),Tip0,((((Tip8,Tip4),Tip3),Tip2),((Tip7,Tip5),Tip1)));
(Tip9,Tip0,(Tip8,((Tip7,Tip4),((Tip6,Tip2),(Tip3,(Tip5,Tip1))))));
((((Tip7,(Tip8,(Tip9,Tip6))),Tip5),Tip4),Tip0,(Tip2,(Tip3,Tip1)));
EOF
gotree generate yuletree -s 10 -n 10 | gotree support setrand | gotree support clear | gotree brlen clear > result
diff result expected
rm -f expected result

# gotree comment clear
echo "->gotree comment clear"
cat > input <<EOF
(t1[c1],t2[c2],(t3[c3],t4[c4])[c5]);
EOF
cat > expected <<EOF
(t1,t2,(t3,t4));
EOF
gotree comment clear -i input > result
diff result expected
rm -f expected result input

# gotree comment clear nodes+edges
echo "->gotree comment clear (nodes+edges)"
cat > input <<EOF
(t1[n1]:1[e1],t2[n2]:1[e2],(t3[n3]:1[e3],t4[n4]:1[e4])[n5]:1[e5]);
EOF
cat > expected <<EOF
(t1:1,t2:1,(t3:1,t4:1):1);
EOF
gotree comment clear -i input > result
diff result expected
rm -f expected result input

# gotree comment clear nodes only
echo "->gotree comment clear (nodes)"
cat > input <<EOF
(t1[n1]:1[e1],t2[n2]:1[e2],(t3[n3]:1[e3],t4[n4]:1[e4])[n5]:1[e5]);
EOF
cat > expected <<EOF
(t1:1[e1],t2:1[e2],(t3:1[e3],t4:1[e4]):1[e5]);
EOF
gotree comment clear --nodes-only -i input > result
diff result expected
rm -f expected result input

# gotree comment clear edges only
echo "->gotree comment clear (edges)"
cat > input <<EOF
(t1[n1]:1[e1],t2[n2]:1[e2],(t3[n3]:1[e3],t4[n4]:1[e4])[n5]:1[e5]);
EOF
cat > expected <<EOF
(t1[n1]:1,t2[n2]:1,(t3[n3]:1,t4[n4]:1)[n5]:1);
EOF
gotree comment clear --edges-only -i input > result
diff result expected
rm -f expected result input

# gotree collapse length
echo "->gotree collapse length"
cat > expected <<EOF
((Tip4,(Tip7,Tip2)),Tip0,(Tip8,Tip9,Tip3),((Tip6,Tip5),Tip1));
EOF
gotree generate yuletree -s 10 | gotree collapse length -l 0.05 | gotree brlen clear > result
diff result expected
rm -f expected result


# gotree collapse support
echo "->gotree collapse support"
cat > expected <<EOF
(Tip0,((Tip1,Tip6,Tip5)0.9167074899036827,Tip8,Tip9,Tip3)0.925128845219594,Tip4,Tip7,Tip2);
EOF
gotree generate yuletree -s 10 | gotree support setrand -s 10 | gotree collapse support -s 0.7 | gotree brlen clear > result
diff result expected
rm -f expected result


# gotree collapse depth
echo "->gotree collapse depth"
cat > expected <<EOF
((Tip4,Tip7,Tip2),Tip0,((Tip8,Tip9,Tip3),(Tip1,Tip6,Tip5)));
EOF
gotree generate yuletree -s 10 | gotree collapse depth -m 2 -M 2 | gotree brlen clear > result
diff result expected
rm -f expected result

# gotree collapse single
echo "->gotree collapse single"
cat > test_input <<EOF
((((A,B)),((C))),(D,(E)));
EOF
cat > expected <<EOF
(((A,B),C),(D,E));
EOF
gotree collapse single -i test_input > result
diff result expected
rm -f expected result test_input

# gotree compare trees
echo "->gotree compare trees"
cat > expected <<EOF
tree	reference	common	compared
0	7	0	7
EOF
gotree compare trees -i <(gotree generate yuletree -s 10) -c <(gotree generate yuletree -s 12 -n 1) > result
diff expected result
rm -f expected result


# gotree compare edges
echo "->gotree compare edges"
cat > expected <<EOF
tree	brid	length	support	terminal	depth	topodepth	rightname	found	transfer	taxatomove	comparednodename
0	0	0.1824683850061218	N/A	false	1	3		false	2	-Tip2,-Tip7	Tip4
0	1	0.020616211789029896	N/A	true	0	1	Tip4	true	0		Tip4
0	2	0.25879284932877245	N/A	false	1	2		false	1	-Tip7	Tip2
0	3	0.09740195047110385	N/A	true	0	1	Tip7	true	0		Tip7
0	4	0.015450672710905129	N/A	true	0	1	Tip2	true	0		Tip2
0	5	0.25919865790518115	N/A	true	0	1	Tip0	true	0		Tip0
0	6	0.04593880904706901	N/A	false	1	4		false	3	+Tip0,+Tip2,+Tip7	Tip4
0	7	0.1920960924280275	N/A	false	1	3		false	1	-Tip3	
0	8	0.027845992087631298	N/A	true	0	1	Tip8	true	0		Tip8
0	9	0.01026581233891113	N/A	false	1	2		false	1	-Tip9	Tip3
0	10	0.13492605122032592	N/A	true	0	1	Tip9	true	0		Tip9
0	11	0.10309294031874587	N/A	true	0	1	Tip3	true	0		Tip3
0	12	0.30150414585026103	N/A	false	1	3		false	2	-Tip1,-Tip5	Tip6
0	13	0.05817538156872999	N/A	false	1	2		false	1	-Tip5	Tip6
0	14	0.3779897840448691	N/A	true	0	1	Tip6	true	0		Tip6
0	15	0.1120177846434196	N/A	true	0	1	Tip5	true	0		Tip5
0	16	0.239082088939295	N/A	true	0	1	Tip1	true	0		Tip1
EOF
gotree compare edges -i <(gotree generate yuletree -s 10) -c <(gotree generate yuletree -s 12 -n 1) -m --moved-taxa > result 2>/dev/null
diff expected result
rm -f expected result


# gotree compare tips
echo "->gotree compare tips"
cat > expected <<EOF
(Tree 0) > Tip11
(Tree 0) > Tip10
(Tree 0) = 10
EOF
gotree compare tips -i <(gotree generate yuletree -s 10) -c <(gotree generate yuletree -s 12 -n 1 -l 12) > result
diff expected result
rm -f expected result


# gotree compare distances
echo "->gotree compare distances"
cat > expected <<EOF
tree_id	er_id	ec_id	tdist	ec_length	ec_support	ec_topodepth	moving_taxa
0	0	0	3	0.018862750655150758	N/A	2	+Tip6,-Tip2,-Tip7
0	0	4	4	0.05806229063227526	N/A	3	+Tip0,+Tip6,-Tip2,-Tip7
0	0	6	5	0.006167968511774678	N/A	4	+Tip0,+Tip3,+Tip6,-Tip2,-Tip7
0	0	7	4	0.03856952076464118	N/A	3	+Tip8,+Tip9,-Tip4,-Tip7
0	0	8	5	0.005131510752894519	N/A	2	+Tip0,+Tip1,+Tip3,+Tip5,+Tip6
0	0	12	4	0.12937482578337411	N/A	3	+Tip1,+Tip5,-Tip2,-Tip4
0	0	13	3	0.00518311446616857	N/A	2	+Tip5,-Tip2,-Tip4
0	2	0	4	0.018862750655150758	N/A	2	+Tip4,+Tip6,-Tip2,-Tip7
0	2	4	5	0.05806229063227526	N/A	3	+Tip0,+Tip4,+Tip6,-Tip2,-Tip7
0	2	6	4	0.006167968511774678	N/A	4	+Tip1,+Tip5,+Tip8,+Tip9
0	2	7	3	0.03856952076464118	N/A	3	+Tip8,+Tip9,-Tip7
0	2	8	4	0.005131510752894519	N/A	2	+Tip8,+Tip9,-Tip2,-Tip7
0	2	12	3	0.12937482578337411	N/A	3	+Tip1,+Tip5,-Tip2
0	2	13	2	0.00518311446616857	N/A	2	+Tip5,-Tip2
0	6	0	4	0.018862750655150758	N/A	2	+Tip0,+Tip2,+Tip7,-Tip6
0	6	4	3	0.05806229063227526	N/A	3	+Tip2,+Tip7,-Tip6
0	6	6	4	0.006167968511774678	N/A	4	+Tip2,+Tip7,-Tip3,-Tip6
0	6	7	5	0.03856952076464118	N/A	3	+Tip0,+Tip4,+Tip7,-Tip8,-Tip9
0	6	8	4	0.005131510752894519	N/A	2	-Tip1,-Tip3,-Tip5,-Tip6
0	6	12	5	0.12937482578337411	N/A	3	+Tip0,+Tip2,+Tip4,-Tip1,-Tip5
0	6	13	4	0.00518311446616857	N/A	2	+Tip0,+Tip2,+Tip4,-Tip5
0	7	0	5	0.018862750655150758	N/A	2	+Tip0,+Tip1,+Tip2,+Tip5,+Tip7
0	7	4	4	0.05806229063227526	N/A	3	+Tip1,+Tip2,+Tip5,+Tip7
0	7	6	5	0.006167968511774678	N/A	4	+Tip0,+Tip4,+Tip6,-Tip8,-Tip9
0	7	7	2	0.03856952076464118	N/A	3	+Tip2,-Tip3
0	7	8	1	0.005131510752894519	N/A	2	-Tip3
0	7	12	4	0.12937482578337411	N/A	3	+Tip0,+Tip2,+Tip4,+Tip6
0	7	13	5	0.00518311446616857	N/A	2	+Tip0,+Tip1,+Tip2,+Tip4,+Tip6
0	9	0	4	0.018862750655150758	N/A	2	+Tip4,+Tip6,-Tip3,-Tip9
0	9	4	5	0.05806229063227526	N/A	3	+Tip0,+Tip4,+Tip6,-Tip3,-Tip9
0	9	6	4	0.006167968511774678	N/A	4	+Tip0,+Tip4,+Tip6,-Tip9
0	9	7	3	0.03856952076464118	N/A	3	+Tip2,+Tip8,-Tip3
0	9	8	2	0.005131510752894519	N/A	2	+Tip8,-Tip3
0	9	12	5	0.12937482578337411	N/A	3	+Tip0,+Tip2,+Tip4,+Tip6,+Tip8
0	9	13	4	0.00518311446616857	N/A	2	+Tip5,+Tip7,-Tip3,-Tip9
0	12	0	3	0.018862750655150758	N/A	2	+Tip4,-Tip1,-Tip5
0	12	4	4	0.05806229063227526	N/A	3	+Tip0,+Tip4,-Tip1,-Tip5
0	12	6	5	0.006167968511774678	N/A	4	+Tip0,+Tip3,+Tip4,-Tip1,-Tip5
0	12	7	4	0.03856952076464118	N/A	3	+Tip0,+Tip3,+Tip4,+Tip7
0	12	8	5	0.005131510752894519	N/A	2	+Tip0,+Tip2,+Tip3,+Tip4,+Tip7
0	12	12	2	0.12937482578337411	N/A	3	+Tip7,-Tip6
0	12	13	3	0.00518311446616857	N/A	2	+Tip7,-Tip1,-Tip6
0	13	0	2	0.018862750655150758	N/A	2	+Tip4,-Tip5
0	13	4	3	0.05806229063227526	N/A	3	+Tip0,+Tip4,-Tip5
0	13	6	4	0.006167968511774678	N/A	4	+Tip0,+Tip3,+Tip4,-Tip5
0	13	7	5	0.03856952076464118	N/A	3	+Tip0,+Tip1,+Tip3,+Tip4,+Tip7
0	13	8	4	0.005131510752894519	N/A	2	+Tip8,+Tip9,-Tip5,-Tip6
0	13	12	3	0.12937482578337411	N/A	3	+Tip1,+Tip7,-Tip6
0	13	13	2	0.00518311446616857	N/A	2	+Tip7,-Tip6
EOF
gotree compare distances -i <(gotree generate yuletree -s 10) -c <(gotree generate yuletree -s 12 -n 1) > result 2>/dev/null
diff expected result
rm -f expected result


# goree compute bipartitiontree
echo "->gotree compute bipartitiontree"
cat > expected <<EOF
((Tip4:1,Tip7:1,Tip0:1,Tip8:1,Tip9:1,Tip6:1,Tip5:1):1,Tip1:1,Tip2:1,Tip3:1);
EOF
gotree generate yuletree -s 10 | gotree compute bipartitiontree Tip1 Tip2 Tip3 > result
diff expected result
rm -f result expected


# gotree compute consensus
echo "->gotree compute consensus"
cat > expected <<EOF
(Tip0:0.12347870000000004,(Tip9:0.12811019999999992,(Tip8:0.018413999999999993,Tip3:0.10146340000000001)0.87:0.012998850574712643)1:0.08962310000000002,(Tip1:0.17112029999999998,Tip6:0.30890189999999984,Tip5:0.0929547)0.97:0.07447020618556698,(Tip4:0.030012499999999987,(Tip7:0.09010099999999997,Tip2:0.015264200000000006)1:0.10734890000000002)1:0.09298299999999998);
EOF
gotree compute consensus -i tests/data/bootstap_test.nw.gz -f 0.7 -o result
diff expected result
rm -f expected result


echo "->gotree compute classical bootstrap"
cat > expected <<EOF
(Tip0,(Tip4,(Tip7,Tip2)1)1,((Tip9,(Tip8,Tip3)0.87)1,(Tip1,(Tip6,Tip5)0.65)0.97)0.67);
EOF
gotree compute support classical -i tests/data/bootstap_inferred_test.nw.gz -b tests/data/bootstap_test.nw.gz  2>/dev/null |gotree brlen clear > result
diff expected result
rm -f expected result


echo "->gotree compute booster supports"
cat > expected <<EOF
(Tip0,(Tip4,(Tip7,Tip2)1)1,((Tip9,(Tip8,Tip3)0.87)1,(Tip1,(Tip6,Tip5)0.65)0.985)0.89);
EOF
gotree compute support booster -i tests/data/bootstap_inferred_test.nw.gz -b tests/data/bootstap_test.nw.gz --silent -l /dev/null | gotree brlen clear > result
diff expected result
rm -f expected result


echo "->gotree compute edgetrees"
cat > expected <<EOF
((Tip4:1,Tip7:1,Tip2:1):1,Tip0:1,Tip8:1,Tip9:1,Tip3:1,Tip6:1,Tip5:1,Tip1:1);
((Tip7:1,Tip2:1):1,Tip4:1,Tip0:1,Tip8:1,Tip9:1,Tip3:1,Tip6:1,Tip5:1,Tip1:1);
((Tip8:1,Tip9:1,Tip3:1,Tip6:1,Tip5:1,Tip1:1):1,Tip4:1,Tip7:1,Tip2:1,Tip0:1);
((Tip8:1,Tip9:1,Tip3:1):1,Tip4:1,Tip7:1,Tip2:1,Tip0:1,Tip6:1,Tip5:1,Tip1:1);
((Tip9:1,Tip3:1):1,Tip4:1,Tip7:1,Tip2:1,Tip0:1,Tip8:1,Tip6:1,Tip5:1,Tip1:1);
((Tip6:1,Tip5:1,Tip1:1):1,Tip4:1,Tip7:1,Tip2:1,Tip0:1,Tip8:1,Tip9:1,Tip3:1);
((Tip6:1,Tip5:1):1,Tip4:1,Tip7:1,Tip2:1,Tip0:1,Tip8:1,Tip9:1,Tip3:1,Tip1:1);
EOF
gotree generate yuletree -s 10  | gotree compute edgetrees > result
diff expected result
rm -f expected result

echo "->gotree compute edgetrees"
cat > expected <<EOF
((Tip4:1,Tip7:1,Tip2:1):1,Tip0:1,Tip8:1,Tip9:1,Tip3:1,Tip6:1,Tip5:1,Tip1:1);
((Tip7:1,Tip2:1):1,Tip4:1,Tip0:1,Tip8:1,Tip9:1,Tip3:1,Tip6:1,Tip5:1,Tip1:1);
((Tip8:1,Tip9:1,Tip3:1,Tip6:1,Tip5:1,Tip1:1):1,Tip4:1,Tip7:1,Tip2:1,Tip0:1);
((Tip8:1,Tip9:1,Tip3:1):1,Tip4:1,Tip7:1,Tip2:1,Tip0:1,Tip6:1,Tip5:1,Tip1:1);
((Tip9:1,Tip3:1):1,Tip4:1,Tip7:1,Tip2:1,Tip0:1,Tip8:1,Tip6:1,Tip5:1,Tip1:1);
((Tip6:1,Tip5:1,Tip1:1):1,Tip4:1,Tip7:1,Tip2:1,Tip0:1,Tip8:1,Tip9:1,Tip3:1);
((Tip6:1,Tip5:1):1,Tip4:1,Tip7:1,Tip2:1,Tip0:1,Tip8:1,Tip9:1,Tip3:1,Tip1:1);
EOF
gotree generate yuletree -s 10  | gotree compute edgetrees > result
diff expected result
rm -f expected result

echo "->gotree divide"
cat > expected1 <<EOF
((Tip4,(Tip7,Tip2)),Tip0,((Tip8,(Tip9,Tip3)),((Tip6,Tip5),Tip1)));
EOF
cat > expected2 <<EOF
(Tip5,Tip0,((Tip6,(Tip7,Tip4)),(Tip2,((Tip8,(Tip9,Tip3)),Tip1))));
EOF
gotree generate yuletree -s 10 -n 2 |gotree brlen clear | gotree divide -o div
diff expected1 div_000.nw
diff expected2 div_001.nw
rm -f expected1 expected2 div_000.nw div_001.nw


echo "->gotree generate yuletree"
cat > expected <<EOF
((Tip4:0.020616211789029896,(Tip7:0.09740195047110385,Tip2:0.015450672710905129):0.25879284932877245):0.1824683850061218,Tip0:0.25919865790518115,((Tip8:0.027845992087631298,(Tip9:0.13492605122032592,Tip3:0.10309294031874587):0.01026581233891113):0.1920960924280275,((Tip6:0.3779897840448691,Tip5:0.1120177846434196):0.05817538156872999,Tip1:0.239082088939295):0.30150414585026103):0.04593880904706901);
EOF
gotree generate yuletree -s 10 -n 1 > result
diff expected result
rm -f expected result


echo "->gotree generate balancedtree"
cat > expected <<EOF
(((Tip0:0.04593880904706901,Tip1:0.13604994737755394):0.06718605070537677,(Tip2:0.19852695409349608,Tip3:0.002749016032849596):0.2485396648662035):0.25919865790518115,((Tip4:0.12467449897149811,Tip5:0.10033210749794116):0.1824683850061218,(Tip6:0.30150414585026103,Tip7:0.08184535681853511):0.020616211789029896):0.054743875470795914,(((Tip8:0.1120177846434196,Tip9:0.18347097513974125):0.05817538156872999,(Tip10:0.25879284932877245,Tip11:0.09740195047110385):0.3779897840448691):0.239082088939295,((Tip12:0.1920960924280275,Tip13:0.027845992087631298):0.015450672710905129,(Tip14:0.0440885662122905,Tip15:0.14809735366802398):0.17182241382980687):0.03199874235185574):0.13756099791982077);
EOF
gotree generate balancedtree -s 10 -d 4 > result
diff expected result
rm -f expected result


echo "->gotree generate caterpillartree"
cat > expected <<EOF
((((((((Tip9:0.09740195047110385,Tip8:0.015450672710905129):0.25879284932877245,Tip7:0.18347097513974125):0.3779897840448691,Tip6:0.05817538156872999):0.239082088939295,Tip5:0.08184535681853511):0.10033210749794116,Tip4:0.12467449897149811):0.1824683850061218,Tip3:0.002749016032849596):0.13604994737755394,Tip2:0.04593880904706901):0.06718605070537677,Tip0:0.0540687078328298,Tip1:0.054743875470795914);
EOF
gotree generate caterpillartree -s 10  > result
diff expected result
rm -f expected result

echo "->gotree generate uniformtree"
cat > expected <<EOF
(Tip5:0.08184535681853511,Tip0:0.30150414585026103,((Tip9:0.13492605122032592,Tip6:0.10309294031874587):0.01026581233891113,((Tip7:0.09740195047110385,(Tip8:0.027845992087631298,((Tip4:0.020616211789029896,Tip3:0.12467449897149811):0.1824683850061218,Tip2:0.19852695409349608):0.0440885662122905):0.1920960924280275):0.25879284932877245,Tip1:0.06718605070537677):0.1120177846434196):0.05817538156872999);
EOF
gotree generate uniformtree -s 10  > result
diff expected result
rm -f expected result


echo "->gotree matrix"
cat > expected <<EOF
5
Tip4	0.000000000000	0.145290710761	0.462283254700	0.385073353220	0.447550359936
Tip2	0.145290710761	0.000000000000	0.566341541883	0.489131640402	0.551608647118
Tip0	0.462283254700	0.566341541883	0.000000000000	0.441187414330	0.503664421046
Tip3	0.385073353220	0.489131640402	0.441187414330	0.000000000000	0.334576901471
Tip1	0.447550359936	0.551608647118	0.503664421046	0.334576901471	0.000000000000
EOF
gotree generate yuletree -s 10 -l 5 | gotree matrix > result
diff expected result
rm -f expected result

echo "->gotree brlen setmin 1"
cat > expected <<EOF
((Tip4:1,(Tip7:1,Tip2:1):1):1,Tip0:1,((Tip8:1,(Tip9:1,Tip3:1):1):1,((Tip6:1,Tip5:1):1,Tip1:1):1):1);
EOF
gotree generate yuletree -s 10 -l 10 | gotree brlen setmin  -l 1 > result
diff expected result
rm -f expected result

echo "->gotree brlen setmin 10"
cat > expected <<EOF
((Tip4:10,(Tip7:10,Tip2:10):10):10,Tip0:10,((Tip8:10,(Tip9:10,Tip3:10):10):10,((Tip6:10,Tip5:10):10,Tip1:10):10):10);
EOF
gotree generate yuletree -s 10 -l 10 | gotree brlen clear | gotree brlen setmin  -l 10 > result
diff expected result
rm -f expected result

echo "->gotree prune"
cat > expected <<EOF
((Tip4,(Tip7,Tip2)),((Tip8,(Tip9,Tip3)),((Tip6,Tip5),Tip1)),Tip0);
EOF
gotree generate yuletree -s 10 -l 20 | gotree prune -i - -c <(gotree generate yuletree -s 12 -l 10) | gotree brlen clear > result
diff expected result
rm -f expected result


echo "->gotree brlen setrand"
cat > expected <<EOF
((Tip4:0.11181011331618643,(Tip7:0.21688356961855743,Tip2:0.21695890315161873):0.007486847792469759):0.02262551762264341,Tip0:0.07447903650558614,((Tip8:0.05414175839023796,(Tip9:0.34924246250387486,Tip3:0.023925115233614132):0.1890483249199916):0.03146499978313507,((Tip6:0.31897358778004786,Tip5:0.29071259678750266):0.04826059603128351,Tip1:0.02031669307269784):0.052025373286913534):0.011401847253594477);
EOF
gotree generate yuletree -s 10 | gotree brlen setrand -s 13 > result
diff expected result
rm -f expected result


echo "->gotree support setrand"
cat > expected <<EOF
((Tip4,(Tip7,Tip2)0.2550878763278657)0.6418716208549535,Tip0,((Tip8,(Tip9,Tip3)0.9581212767194948)0.24992593115716047,((Tip6,Tip5)0.2962112349523319,Tip1)0.2923644736644398)0.20284376043157062);
EOF
gotree generate yuletree -s 10 | gotree support setrand -s 12  | gotree brlen clear > result
diff expected result
rm -f expected result


echo "->gotree rename"
cat > mapfile <<EOF
Tip0	Tax0
Tip1	Tax1
Tip2	Tax2
Tip3	Tax3
Tip4	Tax4
Tip5	Tax5
Tip6	Tax6
Tip7	Tax7
Tip8	Tax8
Tip9	Tax9
EOF
cat > expected <<EOF
((Tax4,(Tax7,Tax2)),Tax0,((Tax8,(Tax9,Tax3)),((Tax6,Tax5),Tax1)));
EOF
gotree generate yuletree -s 10 | gotree rename -m mapfile | gotree brlen clear > result
diff expected result
rm -f expected result mapfile

echo "->gotree rename auto"
cat > mapfile <<EOF
Tip4	T0001
Tip7	T0002
Tip0	T0004
Tip9	T0006
Tip3	T0007
Tip1	T0010
Tip2	T0003
Tip8	T0005
Tip6	T0008
Tip5	T0009
EOF

cat > expected <<EOF
((T0001,(T0002,T0003)),T0004,((T0005,(T0006,T0007)),((T0008,T0009),T0010)));
EOF
gotree generate yuletree -s 10 | gotree rename -a -m mapfile2 -l 5  | gotree brlen clear > result
diff expected result
diff <(sort mapfile) <(sort mapfile2)
rm -f expected result mapfile mapfile2

echo "->gotree reroot outgroup"
cat > expected <<EOF
((((Tip4,(Tip7,Tip2)),Tip0),((Tip6,Tip5),Tip1)),(Tip8,(Tip9,Tip3)));
EOF
gotree generate yuletree -s 10 | gotree reroot outgroup Tip3 Tip8 Tip9 | gotree brlen clear > result
diff expected result
rm -f expected result

echo "->gotree reroot midpoint"
cat > expected <<EOF
(((Tip6,Tip5),Tip1),(((Tip4,(Tip7,Tip2)),Tip0),(Tip8,(Tip9,Tip3))));
EOF
gotree generate yuletree -s 10 | gotree reroot midpoint | gotree brlen clear> result
diff expected result
rm -f expected result

echo "->gotree resolve"
cat > expected <<EOF
((Tip4,(Tip7,Tip2)),(Tip3,(Tip9,Tip8)),(((Tip6,Tip5),Tip1),Tip0));
EOF
gotree generate yuletree -s 10 | gotree collapse length -l 0.05 | gotree resolve -s 10 | gotree brlen clear > result
diff expected result
rm -f expected result


echo "->gotree shuffletips"
cat > expected <<EOF
((Tip5,(Tip2,Tip3)),Tip7,((Tip8,(Tip4,Tip0)),((Tip6,Tip1),Tip9)));
EOF
gotree generate yuletree -s 10 | gotree shuffletips -s 12 | gotree brlen clear > result
diff expected result
rm -f expected result


echo "->gotree subtree"
cat > clade <<EOF
clade:Tip2,Tip4,Tip7
EOF
cat > expected <<EOF
(Tip4,(Tip7,Tip2))clade;
EOF
gotree generate yuletree -s 10 | gotree annotate -m clade | gotree subtree -n clade | gotree brlen clear > result
diff expected result
rm -f expected result clade


echo "->gotree stats"
cat > expected <<EOF
tree	nodes	tips	edges	meanbrlen	sumbrlen	meansupport	mediansupport	rooted	nbcherries	colless	sackin
0	18	10	17	0.14334492	2.43686361	NaN	NaN	unrooted	3	7	35
EOF
gotree generate yuletree -s 10 | gotree stats > result
diff expected result
rm -f expected result


echo "->gotree stats edges"
cat > expected <<EOF
tree	brid	length	support	terminal	depth	topodepth	rightname	comments	leftname	rightcomment	leftcomment
0	0	0.1824683850061218	N/A	false	1	3		[]		[]	[]
0	1	0.020616211789029896	N/A	true	0	1	Tip4	[]		[]	[]
0	2	0.25879284932877245	N/A	false	1	2		[]		[]	[]
0	3	0.09740195047110385	N/A	true	0	1	Tip7	[]		[]	[]
0	4	0.015450672710905129	N/A	true	0	1	Tip2	[]		[]	[]
0	5	0.25919865790518115	N/A	true	0	1	Tip0	[]		[]	[]
0	6	0.04593880904706901	N/A	false	1	4		[]		[]	[]
0	7	0.1920960924280275	N/A	false	1	3		[]		[]	[]
0	8	0.027845992087631298	N/A	true	0	1	Tip8	[]		[]	[]
0	9	0.01026581233891113	N/A	false	1	2		[]		[]	[]
0	10	0.13492605122032592	N/A	true	0	1	Tip9	[]		[]	[]
0	11	0.10309294031874587	N/A	true	0	1	Tip3	[]		[]	[]
0	12	0.30150414585026103	N/A	false	1	3		[]		[]	[]
0	13	0.05817538156872999	N/A	false	1	2		[]		[]	[]
0	14	0.3779897840448691	N/A	true	0	1	Tip6	[]		[]	[]
0	15	0.1120177846434196	N/A	true	0	1	Tip5	[]		[]	[]
0	16	0.239082088939295	N/A	true	0	1	Tip1	[]		[]	[]
EOF
gotree generate yuletree -s 10 | gotree stats edges > result
diff expected result
rm -f expected result


echo "->gotree stats nodes"
cat > expected <<EOF
tree	nid	nneigh	name	depth	comments
0	0	3		1	[]
0	1	3		1	[]
0	2	1	Tip4	0	[]
0	3	3		1	[]
0	4	1	Tip7	0	[]
0	5	1	Tip2	0	[]
0	6	1	Tip0	0	[]
0	7	3		2	[]
0	8	3		1	[]
0	9	1	Tip8	0	[]
0	10	3		1	[]
0	11	1	Tip9	0	[]
0	12	1	Tip3	0	[]
0	13	3		1	[]
0	14	3		1	[]
0	15	1	Tip6	0	[]
0	16	1	Tip5	0	[]
0	17	1	Tip1	0	[]
EOF
gotree generate yuletree -s 10 | gotree stats nodes > result
diff expected result
rm -f expected result


echo "->gotree stats rooted"
cat > expected <<EOF
tree	rooted
0	unrooted
EOF
gotree generate yuletree -s 10 | gotree stats rooted > result
diff expected result
rm -f expected result


echo "->gotree stats splits"
cat > expected <<EOF
Tree	Tip9|Tip8|Tip7|Tip6|Tip5|Tip4|Tip3|Tip2|Tip1|Tip0
0	0010010100.
0	0000010000.
0	0010000100.
0	0010000000.
0	0000000100.
0	0000000001.
0	1101101010.
0	1100001000.
0	0100000000.
0	1000001000.
0	1000000000.
0	0000001000.
0	0001100010.
0	0001100000.
0	0001000000.
0	0000100000.
0	0000000010.
EOF
gotree generate yuletree -s 10 | gotree stats splits > result
diff expected result
rm -f expected result


echo "->gotree stats tips"
cat > expected <<EOF
tree	id	nneigh	name
0	2	1	Tip4
0	4	1	Tip7
0	5	1	Tip2
0	6	1	Tip0
0	9	1	Tip8
0	11	1	Tip9
0	12	1	Tip3
0	15	1	Tip6
0	16	1	Tip5
0	17	1	Tip1
EOF
gotree generate yuletree -s 10 | gotree stats tips > result
diff expected result
rm -f expected result


echo "->gotree unroot"
cat > expected <<EOF
((Tip9,Tip2),(Tip3,((((Tip8,Tip6),Tip5),Tip4),Tip1)),(Tip7,Tip0));
EOF
gotree generate yuletree -r -s 10 | gotree brlen clear | gotree unroot > result
diff expected result
rm -f expected result

echo "->gotree draw text"
cat > expected <<EOF
           + Tip4                                           
+----------|                                                
|          |                +----- Tip7                     
|          +----------------|                               
|                           + Tip2                          
|                                                           
|--------------- Tip0                                       
|                                                           
|              + Tip8                                       
| +------------|                                            
| |            |-------- Tip9                               
| |            |                                            
| |            +------ Tip3                                 
+-|                                                         
  |                      +----------------------- Tip6      
  |                   +--|                                  
  +-------------------|  +------- Tip5                      
                      |                                     
                      +-------------- Tip1                  
                                                            
EOF
gotree generate yuletree -s 10 | gotree draw text -w 50 > result
diff expected result
rm -f expected result

echo "->gotree annotate"
cat > inferred <<EOF
(((((Hylobates_pileatus:0.23988592,(Pongo_pygmaeus_abelii:0.11809071,(Gorilla_gorilla_gorilla:0.13596645,(Homo_sapiens:0.11344407,Pan_troglodytes:0.11665038)0.62:0.02364476)0.78:0.04257513)0.93:0.15711475)0.56:0.03966791,(Macaca_sylvanus:0.06332916,(Macaca_fascicularis_fascicularis:0.07605049,(Macaca_mulatta:0.06998962,Macaca_fuscata:0)0.98:0.08492791)0.47:0.02236558)0.89:0.11208218)0.43:0.0477543,Saimiri_sciureus:0.25824985)0.71:0.14311537,(Tarsius_tarsier:0.62272677,Lemur_sp.:0.40249393)0.35:0)0.62:0.077084225,(Mus_musculus:0.4057381,Bos_taurus:0.65776307)0.62:0.077084225);
EOF
cat > annotated <<EOF
((((((((Gorilla_gorilla_gorilla[subspecies],Pan_troglodytes[species],Homo_sapiens[species])Homo/Pan/Gorilla_group[subfamily],Pongo_pygmaeus_abelii[species])Pongidae[family],Hylobates_pileatus[species])Hominoidea[superfamily],(Macaca_sylvanus[species],Macaca_fascicularis_fascicularis[subspecies],Macaca_fuscata[species],Macaca_mulatta[species])Macaca[genus])Catarrhini[parvorder],Saimiri_sciureus[species])Simiiformes[infraorder],Tarsius_tarsier[species])Haplorrhini[suborder],Lemur_sp.[species])Primates[order],Mus_musculus[species],Bos_taurus[species])Euarchontoglires[superorder];
EOF
cat > expected <<EOF
(((((Hylobates_pileatus:0.23988592,(Pongo_pygmaeus_abelii:0.11809071,(Gorilla_gorilla_gorilla:0.13596645,(Homo_sapiens:0.11344407,Pan_troglodytes:0.11665038)Pan_troglodytes_1_2:0.02364476)Homo/Pan/Gorilla_group_0_3:0.04257513)Pongidae_0_4:0.15711475)Hominoidea_0_5:0.03966791,(Macaca_sylvanus:0.06332916,(Macaca_fascicularis_fascicularis:0.07605049,(Macaca_mulatta:0.06998962,Macaca_fuscata:0)Macaca_fuscata_1_2:0.08492791)Macaca_1_3:0.02236558)Macaca_0_4:0.11208218)Catarrhini_0_9:0.0477543,Saimiri_sciureus:0.25824985)Simiiformes_0_10:0.14311537,(Tarsius_tarsier:0.62272677,Lemur_sp.:0.40249393)Tarsius_tarsier_1_2:0)Primates_0_12:0.077084225,(Mus_musculus:0.4057381,Bos_taurus:0.65776307)Primates_0_2:0.077084225);
EOF
gotree annotate -i inferred -c annotated -o result
diff expected result
rm -f expected result annotated inferred


echo "->gotree reformat nexus"
cat > newick <<EOF
(fish, (frog, (snake, mouse)));
(fish, (snake, (frog, mouse)));
(fish, (mouse, (snake, frog)));
(mouse, (frog, (snake, fish)));
EOF
cat > expected <<EOF
#NEXUS
BEGIN TAXA;
 DIMENSIONS NTAX=4;
 TAXLABELS fish frog snake mouse;
END;
BEGIN TREES;
  TREE tree0 = (fish,(frog,(snake,mouse)));
  TREE tree1 = (fish,(snake,(frog,mouse)));
  TREE tree2 = (fish,(mouse,(snake,frog)));
  TREE tree3 = (mouse,(frog,(snake,fish)));
END;
EOF
gotree reformat nexus -i newick -f newick -o result
diff expected result
rm -f expected result newick

echo "->gotree reformat newick 0"
cat > newick <<EOF
(fish , (frog , (snake , mouse)));
(fish , (snake , (frog , mouse)));
(fish , (mouse , (snake , frog)));
(mouse , (frog , (snake , fish)));
EOF
cat > expected <<EOF
(fish,(frog,(snake,mouse)));
(fish,(snake,(frog,mouse)));
(fish,(mouse,(snake,frog)));
(mouse,(frog,(snake,fish)));
EOF
gotree reformat newick -i newick -f newick -o result
diff expected result
rm -f expected result newick


echo "->gotree reformat newick 1"
cat > nexus <<EOF
#NEXUS
BEGIN TAXA;
      DIMENSIONS NTAX=4;
      TaxLabels fish frog snake mouse;
END;

BEGIN CHARACTERS;
      Dimensions NChar=40;
      Format DataType=DNA;
      Matrix
        fish   ACATA GAGGG TACCT CTAAA
        frog   ACATA GAGGG TACCT CTAAC
        snake  ACATA GAGGG TACCT CTAAG
        mouse  ACATA GAGGG TACCT CTAAT

        fish   ACATA GAGGG TACCT CTAAG
        frog   CCATA GAGGG TACCT CTAAG
        snake  GCATA GAGGG TACCT CTAAG
        mouse  TCATA GAGGG TACCT CTAAG
;
END;

BEGIN TREES;
      Tree best=(fish, (frog, (snake, mouse)));
END;
EOF
cat > expected <<EOF
(fish,(frog,(snake,mouse)));
EOF
gotree reformat newick -i nexus -f nexus -o result
diff expected result
rm -f expected result nexus

echo "->gotree reformat newick 2"
cat > nexus <<EOF
#NEXUS
[NEXUS COMMENT]
BEGIN TAXA;
      DIMENSIONS NTAX=4;[NEXUS COMMENT]
      [NEXUS COMMENT]
      TaxLabels fish frog snake mouse;[NEXUS COMMENT]
END;
[NEXUS COMMENT]
BEGIN CHARACTERS;
      [NEXUS COMMENT]
      Dimensions NChar=40;
      [NEXUS COMMENT]
      Format DataType=DNA;
      Matrix
        fish   ACATA GAGGG TACCT CTAAA
        frog   ACATA GAGGG TACCT CTAAC
        snake  ACATA GAGGG TACCT CTAAG
        mouse  ACATA GAGGG TACCT CTAAT

        fish   ACATA GAGGG TACCT CTAAG
        frog   CCATA GAGGG TACCT CTAAG
        snake  GCATA GAGGG TACCT CTAAG
        mouse  TCATA GAGGG TACCT CTAAG
;
[NEXUS COMMENT]
END;

BEGIN TREES;
      Tree best= [&R] (fish[COMMENT], (frog, (snake, mouse)));
END;
EOF
cat > expected <<EOF
(fish[COMMENT],(frog,(snake,mouse)));
EOF
gotree reformat newick -i nexus -f nexus -o result
diff expected result
rm -f expected result nexus

echo "->gotree reformat nexus->newick with translate"
cat > nexus <<EOF
#NEXUS
[NEXUS COMMENT]
BEGIN TREES;
      translate 1 fish
      , 2 frog
      , 3 snake,
      4 mouse;
      Tree best= [&R] (1[COMMENT], (2, (3, 4)));
END;
EOF
cat > expected <<EOF
(fish[COMMENT],(frog,(snake,mouse)));
EOF
gotree reformat newick -i nexus -f nexus -o result
diff expected result
rm -f expected result nexus

echo "->gotree reformat newick 3"
cat > nexus <<EOF
#NEXUS
BEGIN TAXA;
      TaxLabels fish frog snake mouse;
END;

BEGIN CHARACTERS;
      Dimensions NChar=40;
      Format DataType=DNA;
      Matrix
        fish   ACATA GAGGG TACCT CTAAA
        fish   ACATA GAGGG TACCT CTAAG

        frog   ACATA GAGGG TACCT CTAAC
        frog   CCATA GAGGG TACCT CTAAG

        snake  ACATA GAGGG TACCT CTAAG
        snake  GCATA GAGGG TACCT CTAAG

        mouse  ACATA GAGGG TACCT CTAAT
        mouse  TCATA GAGGG TACCT CTAAG
;
END;

BEGIN TREES;
      Tree best=(fish, (frog, (snake, mouse)));
END;
EOF
cat > expected <<EOF
(fish,(frog,(snake,mouse)));
EOF
gotree reformat newick -i nexus -f nexus -o result
diff expected result
rm -f expected result nexus

echo "->gotree acr acctran"
cat > tmp_states.txt <<EOF
1,A
2,A
3,B
4,B
5,A
6,B
7,B
8,A
9,A
10,A
11,A
EOF
cat > tmp_tree.txt <<EOF
(1,(2,((3,(4,5)),(6,((7,8),((9,10),11))))));
EOF
cat > expected <<EOF
(1[A],(2[A],((3[B],(4[B],5[A])[B])[B],(6[B],((7[B],8[A])[A],((9[A],10[A])[A],11[A])[A])[A])[B])[B])[A])[A];
EOF
gotree acr -i tmp_tree.txt --states tmp_states.txt --algo acctran -o result
diff expected result
rm -f expected result tmp_tree.txt tmp_states.txt

echo "->gotree acr downpass"
cat > tmp_states.txt <<EOF
t1,A
t2,A
t3,B
t4,B
t5,A
t6,B
t7,B
t8,A
t9,A
t10,A
t11,A
EOF
cat > tmp_tree.txt <<EOF
(t1,(t2,((t3,(t4,t5)),(t6,((t7,t8),((t9,t10),t11))))));
EOF
cat > expected <<EOF
(t1[A],(t2[A],((t3[B],(t4[B],t5[A])[A|B])[A|B],(t6[B],((t7[B],t8[A])[A|B],((t9[A],t10[A])[A],t11[A])[A])[A|B])[A|B])[A|B])[A])[A];
EOF
gotree acr -i tmp_tree.txt --states tmp_states.txt --algo downpass -o result
diff expected result
rm -f expected result tmp_tree.txt tmp_states.txt

echo "->gotree acr deltran"
cat > tmp_states.txt <<EOF
t1,A
t2,A
t3,B
t4,B
t5,A
t6,B
t7,B
t8,A
t9,A
t10,A
t11,A
EOF
cat > tmp_tree.txt <<EOF
(t1,(t2,((t3,(t4,t5)),(t6,((t7,t8),((t9,t10),t11))))));
EOF
cat > expected <<EOF
(t1[A],(t2[A],((t3[B],(t4[B],t5[A])[A])[A],(t6[B],((t7[B],t8[A])[A],((t9[A],t10[A])[A],t11[A])[A])[A])[A])[A])[A])[A];
EOF
gotree acr -i tmp_tree.txt --states tmp_states.txt --algo deltran -o result
diff expected result
rm -f expected result tmp_tree.txt tmp_states.txt


echo "->gotree asr acctran"
cat > tmp_states.txt <<EOF
11 2
1 AA
2 AA
3 CC
4 CC
5 AA
6 CC
7 CC
8 AA
9 AA
10 AA
11 AA
EOF
cat > tmp_tree.txt <<EOF
(1,(2,((3,(4,5)),(6,((7,8),((9,10),11))))));
EOF
cat > expected <<EOF
(1[AA],(2[AA],((3[CC],(4[CC],5[AA])[CC])[CC],(6[CC],((7[CC],8[AA])[AA],((9[AA],10[AA])[AA],11[AA])[AA])[AA])[CC])[CC])[AA])[AA];
EOF
gotree asr -i tmp_tree.txt -p -a tmp_states.txt --algo acctran -o result
diff expected result
rm -f expected result tmp_tree.txt tmp_states.txt

echo "->gotree asr downpass"
cat > tmp_states.txt <<EOF
11 4
t1 AAAT
t2 AAAT
t3 CCCT
t4 CCCT
t5 AAAT
t6 CCCT
t7 CCCT
t8 AAAT
t9 AAAT
t10 AAAT
t11 AAAT
EOF
cat > tmp_tree.txt <<EOF
(t1,(t2,((t3,(t4,t5)),(t6,((t7,t8),((t9,t10),t11))))));
EOF
cat > expected <<EOF
(t1[AAAT],(t2[AAAT],((t3[CCCT],(t4[CCCT],t5[AAAT])[{AC}{AC}{AC}T])[{AC}{AC}{AC}T],(t6[CCCT],((t7[CCCT],t8[AAAT])[{AC}{AC}{AC}T],((t9[AAAT],t10[AAAT])[AAAT],t11[AAAT])[AAAT])[{AC}{AC}{AC}T])[{AC}{AC}{AC}T])[{AC}{AC}{AC}T])[AAAT])[AAAT];
EOF
gotree asr -i tmp_tree.txt -p -a tmp_states.txt --algo downpass -o result
diff expected result
rm -f expected result tmp_tree.txt tmp_states.txt

echo "->gotree rotate sort"
cat > expected <<EOF
(6,(1,2),(5,(3,4)));
EOF
echo "((1,2),((3,4),5),6);" | gotree rotate sort > result
diff expected result
rm -f expected result

echo "->gotree rotate sort 2"
cat > expected <<EOF
((6,7),(8,(9,10)),(5,((3,4),(1,(2,8)))));
EOF
echo "(((9,10),8),(((1,(2,8)),(3,4)),5),(6,7));" | gotree rotate sort > result
diff expected result
rm -f expected result

#gotree generate all topologies
cat > expected <<EOF
(B,D,(E,(C,A)));
(B,D,(A,(E,C)));
(B,D,(C,(E,A)));
(D,(C,A),(E,B));
(B,(C,A),(E,D));
(D,(E,A),(C,B));
(A,D,(E,(C,B)));
(A,D,(B,(E,C)));
(A,D,(C,(E,B)));
(A,(C,B),(E,D));
(B,(E,A),(C,D));
(A,(E,B),(C,D));
(A,B,(E,(C,D)));
(A,B,(D,(E,C)));
(A,B,(C,(E,D)));
EOF
echo "(A,(B,D),(C,E));" | gotree generate topologies -i - | gotree rotate sort > result
diff expected result
rm -f expected result
