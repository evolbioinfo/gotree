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
((Tip4:0.020616211789029896,(Tip7:0.09740195047110385,Tip2:0.015450672710905129):0.12939642466438622):0.0912341925030609,Tip0:0.12959932895259058,((Tip8:0.027845992087631298,(Tip9:0.13492605122032592,Tip3:0.10309294031874587):0.005132906169455565):0.09604804621401375,((Tip6:0.3779897840448691,Tip5:0.1120177846434196):0.029087690784364996,Tip1:0.239082088939295)internal1:0.15075207292513051):0.022969404523534506);
EOF
gotree generate yuletree -s 10 | gotree annotate -m mapfile > result
diff result expected
rm -f expected result mapfile


# gotree clear lengths
echo "->gotree clear lengths"
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
gotree generate yuletree -s 10 -n 10 | gotree clear lengths > result
diff result expected
rm -f expected result


# gotree clear supports
echo "->gotree clear supports"
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
gotree generate yuletree -s 10 -n 10 | gotree randsupport | gotree clear supports | gotree clear lengths > result
diff result expected
rm -f expected result


# gotree collapse length
echo "->gotree collapse length"
cat > expected <<EOF
((Tip4:0.020616211789029896,(Tip7:0.09740195047110385,Tip2:0.015450672710905129):0.12939642466438622):0.0912341925030609,Tip0:0.12959932895259058,(Tip8:0.027845992087631298,Tip9:0.13492605122032592,Tip3:0.10309294031874587):0.09604804621401375,(Tip1:0.239082088939295,Tip6:0.3779897840448691,Tip5:0.1120177846434196):0.15075207292513051);
EOF
gotree generate yuletree -s 10 | gotree collapse length -l 0.05 > result
diff result expected
rm -f expected result


# gotree collapse support
echo "->gotree collapse support"
cat > expected <<EOF
(Tip0:0.12959932895259058,((Tip1:0.239082088939295,Tip6:0.3779897840448691,Tip5:0.1120177846434196)0.9167074899036827:0.15075207292513051,Tip8:0.027845992087631298,Tip9:0.13492605122032592,Tip3:0.10309294031874587)0.925128845219594:0.022969404523534506,Tip4:0.020616211789029896,Tip7:0.09740195047110385,Tip2:0.015450672710905129);
EOF
gotree generate yuletree -s 10 | gotree randsupport -s 10 | gotree collapse support -s 0.7 > result
diff result expected
rm -f expected result


# gotree collapse depth
echo "->gotree collapse depth"
cat > expected <<EOF
((Tip4:0.020616211789029896,Tip7:0.09740195047110385,Tip2:0.015450672710905129):0.0912341925030609,Tip0:0.12959932895259058,((Tip8:0.027845992087631298,Tip9:0.13492605122032592,Tip3:0.10309294031874587):0.09604804621401375,(Tip1:0.239082088939295,Tip6:0.3779897840448691,Tip5:0.1120177846434196):0.15075207292513051):0.022969404523534506);
EOF
gotree generate yuletree -s 10 | gotree collapse depth -m 2 -M 2  > result
diff result expected
rm -f expected result


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
0	0	0.0912341925030609	N/A	false	1	3		false	2	-Tip2,-Tip7	Tip4
0	1	0.020616211789029896	N/A	true	0	1	Tip4	true	0		Tip4
0	2	0.12939642466438622	N/A	false	1	2		false	1	-Tip7	Tip2
0	3	0.09740195047110385	N/A	true	0	1	Tip7	true	0		Tip7
0	4	0.015450672710905129	N/A	true	0	1	Tip2	true	0		Tip2
0	5	0.12959932895259058	N/A	true	0	1	Tip0	true	0		Tip0
0	6	0.022969404523534506	N/A	false	1	4		false	3	+Tip0,+Tip2,+Tip7	Tip4
0	7	0.09604804621401375	N/A	false	1	3		false	1	-Tip3	
0	8	0.027845992087631298	N/A	true	0	1	Tip8	true	0		Tip8
0	9	0.005132906169455565	N/A	false	1	2		false	1	-Tip9	Tip3
0	10	0.13492605122032592	N/A	true	0	1	Tip9	true	0		Tip9
0	11	0.10309294031874587	N/A	true	0	1	Tip3	true	0		Tip3
0	12	0.15075207292513051	N/A	false	1	3		false	2	-Tip1,-Tip5	Tip6
0	13	0.029087690784364996	N/A	false	1	2		false	1	-Tip5	Tip6
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
0	0	0	3	0.009431375327575379	N/A	2	+Tip6,-Tip2,-Tip7
0	0	4	4	0.05806229063227526	N/A	3	+Tip0,+Tip6,-Tip2,-Tip7
0	0	6	5	0.006167968511774678	N/A	4	+Tip0,+Tip3,+Tip6,-Tip2,-Tip7
0	0	7	4	0.01928476038232059	N/A	3	+Tip8,+Tip9,-Tip4,-Tip7
0	0	8	5	0.0025657553764472595	N/A	2	+Tip0,+Tip1,+Tip3,+Tip5,+Tip6
0	0	12	4	0.06468741289168706	N/A	3	+Tip1,+Tip5,-Tip2,-Tip4
0	0	13	3	0.002591557233084285	N/A	2	+Tip5,-Tip2,-Tip4
0	2	0	4	0.009431375327575379	N/A	2	+Tip4,+Tip6,-Tip2,-Tip7
0	2	4	5	0.05806229063227526	N/A	3	+Tip0,+Tip4,+Tip6,-Tip2,-Tip7
0	2	6	4	0.006167968511774678	N/A	4	+Tip1,+Tip5,+Tip8,+Tip9
0	2	7	3	0.01928476038232059	N/A	3	+Tip8,+Tip9,-Tip7
0	2	8	4	0.0025657553764472595	N/A	2	+Tip8,+Tip9,-Tip2,-Tip7
0	2	12	3	0.06468741289168706	N/A	3	+Tip1,+Tip5,-Tip2
0	2	13	2	0.002591557233084285	N/A	2	+Tip5,-Tip2
0	6	0	4	0.009431375327575379	N/A	2	+Tip0,+Tip2,+Tip7,-Tip6
0	6	4	3	0.05806229063227526	N/A	3	+Tip2,+Tip7,-Tip6
0	6	6	4	0.006167968511774678	N/A	4	+Tip2,+Tip7,-Tip3,-Tip6
0	6	7	5	0.01928476038232059	N/A	3	+Tip0,+Tip4,+Tip7,-Tip8,-Tip9
0	6	8	4	0.0025657553764472595	N/A	2	-Tip1,-Tip3,-Tip5,-Tip6
0	6	12	5	0.06468741289168706	N/A	3	+Tip0,+Tip2,+Tip4,-Tip1,-Tip5
0	6	13	4	0.002591557233084285	N/A	2	+Tip0,+Tip2,+Tip4,-Tip5
0	7	0	5	0.009431375327575379	N/A	2	+Tip0,+Tip1,+Tip2,+Tip5,+Tip7
0	7	4	4	0.05806229063227526	N/A	3	+Tip1,+Tip2,+Tip5,+Tip7
0	7	6	5	0.006167968511774678	N/A	4	+Tip0,+Tip4,+Tip6,-Tip8,-Tip9
0	7	7	2	0.01928476038232059	N/A	3	+Tip2,-Tip3
0	7	8	1	0.0025657553764472595	N/A	2	-Tip3
0	7	12	4	0.06468741289168706	N/A	3	+Tip0,+Tip2,+Tip4,+Tip6
0	7	13	5	0.002591557233084285	N/A	2	+Tip0,+Tip1,+Tip2,+Tip4,+Tip6
0	9	0	4	0.009431375327575379	N/A	2	+Tip4,+Tip6,-Tip3,-Tip9
0	9	4	5	0.05806229063227526	N/A	3	+Tip0,+Tip4,+Tip6,-Tip3,-Tip9
0	9	6	4	0.006167968511774678	N/A	4	+Tip0,+Tip4,+Tip6,-Tip9
0	9	7	3	0.01928476038232059	N/A	3	+Tip2,+Tip8,-Tip3
0	9	8	2	0.0025657553764472595	N/A	2	+Tip8,-Tip3
0	9	12	5	0.06468741289168706	N/A	3	+Tip0,+Tip2,+Tip4,+Tip6,+Tip8
0	9	13	4	0.002591557233084285	N/A	2	+Tip5,+Tip7,-Tip3,-Tip9
0	12	0	3	0.009431375327575379	N/A	2	+Tip4,-Tip1,-Tip5
0	12	4	4	0.05806229063227526	N/A	3	+Tip0,+Tip4,-Tip1,-Tip5
0	12	6	5	0.006167968511774678	N/A	4	+Tip0,+Tip3,+Tip4,-Tip1,-Tip5
0	12	7	4	0.01928476038232059	N/A	3	+Tip0,+Tip3,+Tip4,+Tip7
0	12	8	5	0.0025657553764472595	N/A	2	+Tip0,+Tip2,+Tip3,+Tip4,+Tip7
0	12	12	2	0.06468741289168706	N/A	3	+Tip7,-Tip6
0	12	13	3	0.002591557233084285	N/A	2	+Tip7,-Tip1,-Tip6
0	13	0	2	0.009431375327575379	N/A	2	+Tip4,-Tip5
0	13	4	3	0.05806229063227526	N/A	3	+Tip0,+Tip4,-Tip5
0	13	6	4	0.006167968511774678	N/A	4	+Tip0,+Tip3,+Tip4,-Tip5
0	13	7	5	0.01928476038232059	N/A	3	+Tip0,+Tip1,+Tip3,+Tip4,+Tip7
0	13	8	4	0.0025657553764472595	N/A	2	+Tip8,+Tip9,-Tip5,-Tip6
0	13	12	3	0.06468741289168706	N/A	3	+Tip1,+Tip7,-Tip6
0	13	13	2	0.002591557233084285	N/A	2	+Tip7,-Tip6
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
(Tip0:0.12761,(Tip4:0.03009,(Tip7:0.09122,Tip2:0.01529)1:0.10923)1:0.09645,((Tip9:0.12833,(Tip8:0.01973,Tip3:0.10375)0.87:0.01234)1:0.08833,(Tip1:0.17627,(Tip6:0.30902,Tip5:0.0926)0.65:0.05594)0.97:0.07479)0.67:0.02348);
EOF
gotree compute support classical -i tests/data/bootstap_inferred_test.nw.gz -b tests/data/bootstap_test.nw.gz -o result 2>/dev/null
diff expected result
rm -f expected result


echo "->gotree compute booster supports"
cat > expected <<EOF
(Tip0:0.12761,(Tip4:0.03009,(Tip7:0.09122,Tip2:0.01529)1:0.10923)1:0.09645,((Tip9:0.12833,(Tip8:0.01973,Tip3:0.10375)0.87:0.01234)1:0.08833,(Tip1:0.17627,(Tip6:0.30902,Tip5:0.0926)0.65:0.05594)0.985:0.07479)0.89:0.02348);
EOF
gotree compute support booster -i tests/data/bootstap_inferred_test.nw.gz -b tests/data/bootstap_test.nw.gz --silent -l /dev/null -o result
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
((Tip4:0.020616211789029896,(Tip7:0.09740195047110385,Tip2:0.015450672710905129):0.12939642466438622):0.0912341925030609,Tip0:0.12959932895259058,((Tip8:0.027845992087631298,(Tip9:0.13492605122032592,Tip3:0.10309294031874587):0.005132906169455565):0.09604804621401375,((Tip6:0.3779897840448691,Tip5:0.1120177846434196):0.029087690784364996,Tip1:0.239082088939295):0.15075207292513051):0.022969404523534506);
EOF
cat > expected2 <<EOF
(Tip5:0.08565127428804534,Tip0:0.021705093532846203,((Tip6:0.07735928108468929,(Tip7:0.13497893602480787,Tip4:0.0867277155865104):0.0006245685056451973):0.002022779393968061,(Tip2:0.06180495431033288,((Tip8:0.04262994378368574,(Tip9:0.026931366387984084,Tip3:0.05072862497546402):0.00862321921222685):0.11773314378908921,Tip1:0.0059751058197291896):0.003071914457094241):0.16885535253756265):0.5612486967794823);
EOF
gotree generate yuletree -s 10 -n 2 | gotree divide -o div
diff expected1 div_000.nw
diff expected2 div_001.nw
rm -f expected1 expected2 div_000.nw div_001.nw


echo "->gotree generate yuletree"
cat > expected <<EOF
((Tip4:0.020616211789029896,(Tip7:0.09740195047110385,Tip2:0.015450672710905129):0.12939642466438622):0.0912341925030609,Tip0:0.12959932895259058,((Tip8:0.027845992087631298,(Tip9:0.13492605122032592,Tip3:0.10309294031874587):0.005132906169455565):0.09604804621401375,((Tip6:0.3779897840448691,Tip5:0.1120177846434196):0.029087690784364996,Tip1:0.239082088939295):0.15075207292513051):0.022969404523534506);
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
((((((((Tip9:0.09740195047110385,Tip8:0.015450672710905129):0.12939642466438622,Tip7:0.18347097513974125):0.18899489202243455,Tip6:0.05817538156872999):0.1195410444696475,Tip5:0.08184535681853511):0.05016605374897058,Tip4:0.12467449897149811):0.0912341925030609,Tip3:0.002749016032849596):0.06802497368877697,Tip2:0.04593880904706901):0.033593025352688384,Tip0:0.0270343539164149,Tip1:0.054743875470795914);
EOF
gotree generate caterpillartree -s 10  > result
diff expected result
rm -f expected result

echo "->gotree generate uniformtree"
cat > expected <<EOF
(Tip5:0.08184535681853511,Tip0:0.15075207292513051,((Tip9:0.13492605122032592,Tip6:0.10309294031874587):0.005132906169455565,((Tip7:0.09740195047110385,(Tip8:0.027845992087631298,((Tip4:0.020616211789029896,Tip3:0.12467449897149811):0.0912341925030609,Tip2:0.19852695409349608):0.0440885662122905):0.09604804621401375):0.12939642466438622,Tip1:0.06718605070537677):0.1120177846434196):0.029087690784364996);
EOF
gotree generate uniformtree -s 10  > result
diff expected result
rm -f expected result


echo "->gotree matrix"
cat > expected <<EOF
5
Tip4	0.000000000000	0.145290710761	0.241449733245	0.270869756193	0.333346762909
Tip2	0.145290710761	0.000000000000	0.345508020427	0.374928043376	0.437405050092
Tip0	0.241449733245	0.345508020427	0.000000000000	0.288618680854	0.351095687570
Tip3	0.270869756193	0.374928043376	0.288618680854	0.000000000000	0.334576901471
Tip1	0.333346762909	0.437405050092	0.351095687570	0.334576901471	0.000000000000
EOF
gotree generate yuletree -s 10 -l 5 | gotree matrix > result
diff expected result
rm -f expected result


echo "->gotree minbrlen"
cat > expected <<EOF
((Tip4:1,(Tip7:1,Tip2:1):1):1,Tip0:1,((Tip8:1,(Tip9:1,Tip3:1):1):1,((Tip6:1,Tip5:1):1,Tip1:1):1):1);
EOF
gotree generate yuletree -s 10 -l 10 | gotree minbrlen -l 1 > result
diff expected result
rm -f expected result

echo "->gotree prune"
cat > expected <<EOF
((Tip4:0.020616211789029896,(Tip7:0.09740195047110385,Tip2:0.14042286306572152):0.12939642466438622):0.0912341925030609,((Tip8:0.027845992087631298,(Tip9:0.13492605122032592,Tip3:0.10309294031874587):0.005132906169455565):0.09604804621401375,((Tip6:0.20514935183767607,Tip5:0.11297927682761447):0.029087690784364996,Tip1:0.1497800057234515):0.15075207292513051):0.022969404523534506,Tip0:0.21331104181132168);
EOF
gotree generate yuletree -s 10 -l 20 | gotree prune -i - -c <(gotree generate yuletree -s 12 -l 10) > result
diff expected result
rm -f expected result


echo "->gotree randbrlen"
cat > expected <<EOF
((Tip4:0.11181011331618643,(Tip7:0.21688356961855743,Tip2:0.21695890315161873):0.007486847792469759):0.02262551762264341,Tip0:0.07447903650558614,((Tip8:0.05414175839023796,(Tip9:0.34924246250387486,Tip3:0.023925115233614132):0.1890483249199916):0.03146499978313507,((Tip6:0.31897358778004786,Tip5:0.29071259678750266):0.04826059603128351,Tip1:0.02031669307269784):0.052025373286913534):0.011401847253594477);
EOF
gotree generate yuletree -s 10 | gotree randbrlen -s 13 > result
diff expected result
rm -f expected result


echo "->gotree randsupport"
cat > expected <<EOF
((Tip4:0.020616211789029896,(Tip7:0.09740195047110385,Tip2:0.015450672710905129)0.2550878763278657:0.12939642466438622)0.6418716208549535:0.0912341925030609,Tip0:0.12959932895259058,((Tip8:0.027845992087631298,(Tip9:0.13492605122032592,Tip3:0.10309294031874587)0.9581212767194948:0.005132906169455565)0.24992593115716047:0.09604804621401375,((Tip6:0.3779897840448691,Tip5:0.1120177846434196)0.2962112349523319:0.029087690784364996,Tip1:0.239082088939295)0.2923644736644398:0.15075207292513051)0.20284376043157062:0.022969404523534506);
EOF
gotree generate yuletree -s 10 | gotree randsupport -s 12 > result
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
((Tax4:0.020616211789029896,(Tax7:0.09740195047110385,Tax2:0.015450672710905129):0.12939642466438622):0.0912341925030609,Tax0:0.12959932895259058,((Tax8:0.027845992087631298,(Tax9:0.13492605122032592,Tax3:0.10309294031874587):0.005132906169455565):0.09604804621401375,((Tax6:0.3779897840448691,Tax5:0.1120177846434196):0.029087690784364996,Tax1:0.239082088939295):0.15075207292513051):0.022969404523534506);
EOF
gotree generate yuletree -s 10 | gotree rename -m mapfile > result
diff expected result
rm -f expected result mapfile


echo "->gotree reroot outgroup"
cat > expected <<EOF
((((Tip4:0.020616211789029896,(Tip7:0.09740195047110385,Tip2:0.015450672710905129):0.12939642466438622):0.0912341925030609,Tip0:0.12959932895259058):0.022969404523534506,((Tip6:0.3779897840448691,Tip5:0.1120177846434196):0.029087690784364996,Tip1:0.239082088939295):0.15075207292513051):0.048024023107006875,(Tip8:0.027845992087631298,(Tip9:0.13492605122032592,Tip3:0.10309294031874587):0.005132906169455565):0.048024023107006875);
EOF
gotree generate yuletree -s 10 | gotree reroot outgroup Tip3 Tip8 Tip9 > result
diff expected result
rm -f expected result

echo "->gotree reroot midpoint"
cat > expected <<EOF
(((Tip6:0.3779897840448691,Tip5:0.1120177846434196):0.029087690784364996,Tip1:0.239082088939295):0.04233828512899088,(((Tip4:0.020616211789029896,(Tip7:0.09740195047110385,Tip2:0.015450672710905129):0.12939642466438622):0.0912341925030609,Tip0:0.12959932895259058):0.022969404523534506,(Tip8:0.027845992087631298,(Tip9:0.13492605122032592,Tip3:0.10309294031874587):0.005132906169455565):0.09604804621401375):0.10841378779613964);
EOF
gotree generate yuletree -s 10 | gotree reroot midpoint> result
diff expected result
rm -f expected result

echo "->gotree resolve"
cat > expected <<EOF
((Tip4:0.020616211789029896,(Tip7:0.09740195047110385,Tip2:0.015450672710905129):0.12939642466438622):0.0912341925030609,Tip0:0.12959932895259058,((Tip3:0.10309294031874587,(Tip9:0.13492605122032592,Tip8:0.027845992087631298):0):0.09604804621401375,(Tip5:0.1120177846434196,(Tip6:0.3779897840448691,Tip1:0.239082088939295):0):0.15075207292513051):0);
EOF
gotree generate yuletree -s 10 | gotree collapse length -l 0.05 | gotree resolve -s 10 > result
diff expected result
rm -f expected result


echo "->gotree shuffletips"
cat > expected <<EOF
((Tip5:0.020616211789029896,(Tip2:0.09740195047110385,Tip3:0.015450672710905129):0.12939642466438622):0.0912341925030609,Tip7:0.12959932895259058,((Tip8:0.027845992087631298,(Tip4:0.13492605122032592,Tip0:0.10309294031874587):0.005132906169455565):0.09604804621401375,((Tip6:0.3779897840448691,Tip1:0.1120177846434196):0.029087690784364996,Tip9:0.239082088939295):0.15075207292513051):0.022969404523534506);
EOF
gotree generate yuletree -s 10 | gotree shuffletips -s 12 > result
diff expected result
rm -f expected result


echo "->gotree subtree"
cat > clade <<EOF
clade:Tip2,Tip4,Tip7
EOF
cat > expected <<EOF
(Tip4:0.020616211789029896,(Tip7:0.09740195047110385,Tip2:0.015450672710905129):0.12939642466438622)clade;
EOF
gotree generate yuletree -s 10 | gotree annotate -m clade | gotree subtree -n clade > result
diff expected result
rm -f expected result clade


echo "->gotree stats"
cat > expected <<EOF
tree	nodes	tips	edges	meanbrlen	sumbrlen	meansupport	mediansupport	rooted
0	18	10	17	0.10486138	1.78264354	NaN	NaN	unrooted
EOF
gotree generate yuletree -s 10 | gotree stats > result
diff expected result
rm -f expected result


echo "->gotree stats edges"
cat > expected <<EOF
tree	brid	length	support	terminal	depth	topodepth	rightname
0	0	0.0912341925030609	N/A	false	1	3	
0	1	0.020616211789029896	N/A	true	0	1	Tip4
0	2	0.12939642466438622	N/A	false	1	2	
0	3	0.09740195047110385	N/A	true	0	1	Tip7
0	4	0.015450672710905129	N/A	true	0	1	Tip2
0	5	0.12959932895259058	N/A	true	0	1	Tip0
0	6	0.022969404523534506	N/A	false	1	4	
0	7	0.09604804621401375	N/A	false	1	3	
0	8	0.027845992087631298	N/A	true	0	1	Tip8
0	9	0.005132906169455565	N/A	false	1	2	
0	10	0.13492605122032592	N/A	true	0	1	Tip9
0	11	0.10309294031874587	N/A	true	0	1	Tip3
0	12	0.15075207292513051	N/A	false	1	3	
0	13	0.029087690784364996	N/A	false	1	2	
0	14	0.3779897840448691	N/A	true	0	1	Tip6
0	15	0.1120177846434196	N/A	true	0	1	Tip5
0	16	0.239082088939295	N/A	true	0	1	Tip1
EOF
gotree generate yuletree -s 10 | gotree stats edges > result
diff expected result
rm -f expected result


echo "->gotree stats nodes"
cat > expected <<EOF
tree	nid	nneigh	name	depth
0	0	3		1
0	1	3		1
0	2	1	Tip4	0
0	3	3		1
0	4	1	Tip7	0
0	5	1	Tip2	0
0	6	1	Tip0	0
0	7	3		2
0	8	3		1
0	9	1	Tip8	0
0	10	3		1
0	11	1	Tip9	0
0	12	1	Tip3	0
0	13	3		1
0	14	3		1
0	15	1	Tip6	0
0	16	1	Tip5	0
0	17	1	Tip1	0
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
((Tip9:0.10309294031874587,Tip2:0.03707446096764306):0.06746302561016296,(Tip3:0.19852695409349608,((((Tip8:0.0440885662122905,Tip6:0.14809735366802398):0.013922996043815649,Tip5:0.18347097513974125):0.18899489202243455,Tip4:0.03199874235185574):0.040922678409267554,Tip1:0.10033210749794116):0.010308105894514948):0.06802497368877697,(Tip7:0.015450672710905129,Tip0:0.17182241382980687):0.07607291297094988);
EOF
gotree generate yuletree -r -s 10 | gotree unroot > result
diff expected result
rm -f expected result

echo "->gotree draw text"
cat > expected <<EOF
       +- Tip4                                              
+------|                                                    
|      |          +-------- Tip7                            
|      +----------|                                         
|                 +- Tip2                                   
|                                                           
|---------- Tip0                                            
|                                                           
|         +- Tip8                                           
|+--------|                                                 
||        |----------- Tip9                                 
||        |                                                 
||        +-------- Tip3                                    
+|                                                          
 |               +------------------------------- Tip6      
 |            +--|                                          
 +------------|  +--------- Tip5                            
              |                                             
              +-------------------- Tip1                    
                                                            
EOF
gotree generate yuletree -s 10 | gotree draw text -w 50 > result
diff expected result
rm -f expected result

