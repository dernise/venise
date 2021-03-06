package cache

import (
	"github.com/dernise/venise/cache/binary"
	"github.com/dernise/venise/structures"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// The coord cache is a cache of every nodes' coordinates.
type Coords struct {
	*leveldb.DB
}

func NewCoordsCache(path string, options *opt.Options) (*Coords, error) {
	db, err := leveldb.OpenFile(path, options)
	if err != nil {
		return nil, err
	}

	cache := Coords{db}

	return &cache, err
}

func (c *Coords) PutCoord(node structures.Node) error {
	bytes, err := binary.MarshalCoords(node)
	if err != nil {
		return err
	}

	c.Put(idToKeyBuf(node.ID), bytes, nil)
	return nil
}

func (c *Coords) FillWay(way *structures.Way) error {
	for _, nodeId := range way.NodeIDs {
		nodeBuf, err := c.Get(idToKeyBuf(nodeId), nil)
		if err != nil {
			return err
		}

		node := structures.Node{
			ID:   nodeId,
			Info: structures.Info{},
		}
		unMarshalledCoords := binary.UnmarshallCoord(nodeBuf, node)
		way.Nodes = append(way.Nodes, unMarshalledCoords)
	}
	return nil
}
