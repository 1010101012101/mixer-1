package proxy

import (
	"fmt"
)

type schema struct {
	server   *Server
	cfg      *config
	db       string
	nodes    []*node
	rw_split bool
}

func newSchema(server *Server, cfgSchema *configSchema, nodes []*node) *schema {
	s := new(schema)

	s.server = server
	s.cfg = server.cfg
	s.db = cfgSchema.DB
	s.nodes = nodes
	s.rw_split = cfgSchema.RWSplit

	return s
}

type routeQuery struct {
	Query string
	Args  []interface{}
}

//return a map key is node and value is the routeQuery the node will run
func (s *schema) Route(l *lex) (map[*node]routeQuery, error) {
	//todo
	//rebuild query for different node
	//now we only return first datanode

	return map[*node]routeQuery{s.nodes[0]: routeQuery{l.Query, l.Args}}, nil
}

//return a node for prepare query
func (s *schema) PrepareNode(l *lex) (*node, error) {
	return s.nodes[0], nil
}

type schemas map[string]*schema

func (ss schemas) GetSchema(db string) *schema {
	if s, ok := ss[db]; ok {
		return s
	} else {
		return nil
	}
}

func newSchemas(server *Server, nodes nodes) schemas {
	cfg := server.cfg

	s := make(schemas, len(cfg.Schemas))

	for _, v := range cfg.Schemas {
		if len(v.Nodes) == 0 {
			panic(fmt.Sprintf("schema %s has no node", v.DB))
		}

		nds := make([]*node, 0, len(v.Nodes))
		for _, nodeName := range v.Nodes {
			if node := nodes.GetNode(nodeName); node == nil {
				panic(fmt.Sprintf("schema %s has invalid node name %s", v.DB, nodeName))
			} else {
				nds = append(nds, node)
			}
		}

		s[v.DB] = newSchema(server, &v, nds)
	}

	return s
}
