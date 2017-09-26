package draw

import (
	"bufio"
	"fmt"

	"github.com/fredericlemoine/gotree/tree"
)

type cytoscapeLayout struct {
	supportCutoff float64
	hasSupport    bool
	writer        *bufio.Writer
}

func NewCytoscapeLayout(writer *bufio.Writer, hasSupport bool) TreeLayout {
	return &cytoscapeLayout{0.7, hasSupport, writer}
}

func (layout *cytoscapeLayout) SetSupportCutoff(c float64) {
	layout.supportCutoff = c
}

/*
Draw the tree on the specific drawer. Does not close the file. The caller must do it.
*/
func (layout *cytoscapeLayout) DrawTree(t *tree.Tree) error {
	var err error = nil
	_, err = layout.writer.WriteString(`
<html>
<head>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/cytoscape/3.2.3/cytoscape.min.js"></script>
  <script src="https://cdn.rawgit.com/cpettitt/dagre/v0.7.4/dist/dagre.min.js"></script>
  <script src="https://cdn.rawgit.com/cytoscape/cytoscape.js-dagre/1.5.0/cytoscape-dagre.js"></script>
  <style media="screen" type="text/css">
    #cy {
    width: 100%;
    height: 100%;
    display: block;
    }
  </style>
</head>
<body>
  <div id="cy">
  </div>
    <script>
    var cy = cytoscape({
    container: document.getElementById("cy"),
    layout: {
    name: 'dagre',
    },
    style: [{
      selector: 'node',
      style: {
        'content': 'data(label)',
        'text-opacity': 0.5,
        'text-valign': 'center',
        'text-halign': 'right',
        'background-color': '#11479e'
      }
    },{
      selector: 'edge',
      style: {
`)
	if layout.hasSupport {
		layout.writer.WriteString("'width': 'mapData(support, 0, 100, 4, 50)',\n")
	} else {
		layout.writer.WriteString("'width': 4,\n")
	}
	layout.writer.WriteString("'curve-style': 'bezier',\n")
	if t.Rooted() {
		_, err = layout.writer.WriteString("   'target-arrow-shape': 'triangle',\n")
	}
	_, err = layout.writer.WriteString(`
        'line-color': '#9dbaea',
        'target-arrow-color': '#9dbaea'
      }
    }],
    elements: {
`)
	layout.drawNodes(t)
	layout.drawEdges(t)
	if err != nil {
		return err
	}

	_, err = layout.writer.WriteString(`
    }});
  </script>
</body>
</html>
`)
	return err
}

func (layout *cytoscapeLayout) drawNodes(t *tree.Tree) {
	layout.writer.WriteString("    nodes: [\n")
	for i, n := range t.Nodes() {
		label := n.Name()
		if label == "" {
			label = fmt.Sprintf("n%d", i)
		}
		n.SetId(i)
		layout.writer.WriteString(fmt.Sprintf("{ data: { id: 'n%d' , label: '%s'}},\n", i, label))
	}
	layout.writer.WriteString("    ],\n")

}

func (layout *cytoscapeLayout) drawEdges(t *tree.Tree) {
	layout.writer.WriteString("    edges: [\n")
	for _, e := range t.Edges() {
		support := 0
		if e.Support() != tree.NIL_SUPPORT {
			support = int(e.Support() * 100)
		}
		layout.writer.WriteString(fmt.Sprintf("{data: {source: 'n%d', target: 'n%d',  support: %d}},\n", e.Left().Id(), e.Right().Id(), support))
	}
	layout.writer.WriteString("    ],\n")

}
