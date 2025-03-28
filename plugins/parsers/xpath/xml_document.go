package xpath

import (
	"strings"

	"github.com/antchfx/xmlquery"
	path "github.com/antchfx/xpath"
)

type xmlDocument struct{}

func (*xmlDocument) Parse(buf []byte) (dataNode, error) {
	return xmlquery.Parse(strings.NewReader(string(buf)))
}

func (*xmlDocument) QueryAll(node dataNode, expr string) ([]dataNode, error) {
	// If this panics it's a programming error as we changed the document type while processing
	native, err := xmlquery.QueryAll(node.(*xmlquery.Node), expr)
	if err != nil {
		return nil, err
	}

	nodes := make([]dataNode, 0, len(native))
	for _, n := range native {
		nodes = append(nodes, n)
	}
	return nodes, nil
}

func (*xmlDocument) CreateXPathNavigator(node dataNode) path.NodeNavigator {
	// If this panics it's a programming error as we changed the document type while processing
	return xmlquery.CreateXPathNavigator(node.(*xmlquery.Node))
}

func (d *xmlDocument) GetNodePath(node, relativeTo dataNode, sep string) string {
	names := make([]string, 0)

	// If these panic it's a programming error as we changed the document type while processing
	nativeNode := node.(*xmlquery.Node)
	nativeRelativeTo := relativeTo.(*xmlquery.Node)

	// Climb up the tree and collect the node names
	n := nativeNode.Parent
	for n != nil && n != nativeRelativeTo {
		nodeName := d.GetNodeName(n, sep, false)
		names = append(names, nodeName)
		n = n.Parent
	}

	if len(names) < 1 {
		return ""
	}

	// Construct the nodes
	nodepath := ""
	for _, name := range names {
		nodepath = name + sep + nodepath
	}

	return nodepath[:len(nodepath)-1]
}

func (*xmlDocument) GetNodeName(node dataNode, _ string, _ bool) string {
	// If this panics it's a programming error as we changed the document type while processing
	nativeNode := node.(*xmlquery.Node)

	return nativeNode.Data
}

func (*xmlDocument) OutputXML(node dataNode) string {
	native := node.(*xmlquery.Node)
	return native.OutputXML(false)
}
