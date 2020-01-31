package storage

type nodeId = int
type propId = int
type relId = int

const noId = -1

// Node is a label with relationship and property linked lists
type Node struct {
	Label     string
	firstProp propId
	firstRel  relId
}

// Relationship is a type, a from and to node, and are inside linked lists of
// relationships for from and to nodes
type Relationship struct {
	Typ      string
	from     nodeId
	to       nodeId
	fromNext relId
	toNext   relId
}

// Property is a name, value, and forms a linked list
type Property struct {
	name string
	val  string
	next propId
}

// InsertNode creates a node with a given label
// and inserts it into the graph
func InsertNode(label string) nodeId {
	n := Node{
		Label:     label,
		firstProp: noId,
		firstRel:  noId,
	}
	_nodes = append(_nodes, n)
	return len(_nodes) - 1
}

// SetProperty creates a property with a name and val for a node
func SetProperty(n nodeId, name, val string) {
	prop := Property{
		name: name,
		val:  val,
		next: _nodes[n].firstProp,
	}
	_props = append(_props, prop)
	_nodes[n].firstProp = len(_props) - 1
}

// AddRelationship creats a relationship from node to to with a given type
func AddRelationship(from nodeId, to nodeId, typ string) {
	rel := Relationship{
		Typ:      typ,
		from:     from,
		to:       to,
		fromNext: _nodes[from].firstRel,
		toNext:   _nodes[to].firstRel,
	}
	_rels = append(_rels, rel)

	_nodes[from].firstRel = len(_rels) - 1
	_nodes[to].firstRel = len(_rels) - 1
}

type notFoundError struct {
}

func (e notFoundError) Error() string {
	return "not found"
}

func FindProp(n nodeId, name string) (string, error) {
	for prop := _nodes[n].firstProp; prop != noId; prop = _props[prop].next {
		if _props[prop].name == name {
			return _props[prop].val, nil
		}
	}

	return "", notFoundError{}
}

func FindFirstRelTypeTo(n nodeId, to nodeId) (string, error) {
	for rel := _nodes[n].firstRel; rel != noId; rel = _rels[rel].fromNext {
		if _rels[rel].to == to {
			return _rels[rel].Typ, nil
		}
	}

	return "", notFoundError{}
}
