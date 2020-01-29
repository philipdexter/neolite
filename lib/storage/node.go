package storage

// Node is a label with relationship and property linked lists
type Node struct {
	Label     string
	firstRel  *Relationship
	firstProp *Property
}

// Relationship is a type, a from and to node, and are inside linked lists of
// relationships for from and to nodes
type Relationship struct {
	Typ      string
	from     *Node
	to       *Node
	fromNext *Relationship
	toNext   *Relationship
}

// Property is a name, value, and forms a linked list
type Property struct {
	name string
	val  string
	next *Property
}

// NewNode creates a node with a given label
func NewNode(label string) Node {
	return Node{
		Label: label,
	}
}

// SetProperty creates a property with a name and val for a node
func (n *Node) SetProperty(name, val string) {
	prop := Property{
		name: name,
		val:  val,
		next: n.firstProp,
	}

	n.firstProp = &prop
}

// AddRelationship creats a relationship from node to to with a given type
func (n *Node) AddRelationship(to *Node, typ string) {
	rel := Relationship{
		Typ:      typ,
		from:     n,
		to:       to,
		fromNext: n.firstRel,
		toNext:   to.firstRel,
	}

	n.firstRel = &rel
	to.firstRel = &rel
}

type notFoundError struct {
}

func (e notFoundError) Error() string {
	return "not found"
}

func (n *Node) FindProp(name string) (string, error) {
	for prop := n.firstProp; prop != nil; prop = prop.next {
		if prop.name == name {
			return prop.val, nil
		}
	}

	return "", notFoundError{}
}

func (n *Node) FindFirstRelTypeTo(to *Node) (string, error) {
	for rel := n.firstRel; rel != nil; rel = rel.fromNext {
		if rel.to == to {
			return rel.Typ, nil
		}
	}

	return "", notFoundError{}
}
