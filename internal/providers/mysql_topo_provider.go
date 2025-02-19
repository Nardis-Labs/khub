package providers

import (
	"errors"

	"github.com/sullivtr/k8s_platform/internal/modules"
	"github.com/sullivtr/k8s_platform/internal/types"
)

// IMySQLTopoProvider represents the interface for the MySQL topology capture provider
type MySQLTopoProvider struct {
	Session MySQLTopoSession
}

// Compile time proof of implementation
var _ IMySQLTopoProvider = (*MySQLTopoProvider)(nil)

// MySQLTopoSession represents the MySQLTopo SDK session
type MySQLTopoSession struct {
	SDK modules.MySQLTopoSDK
}

func (p *ModuleProviders) InitMySQLTopoProvider() {
	p.MySQLTopoProvider = &MySQLTopoProvider{
		Session: MySQLTopoSession{
			SDK: modules.MySQLTopoSDK{
				MySQLDBPassword: p.Config.MySQLCatalogDBPassword,
			},
		},
	}
}

func (p *MySQLTopoProvider) CaptureReplicationTopology() ([]types.ReplTopoTreeNode, []types.ReplTopoTreeEdge, error) {
	if len(p.Session.SDK.Databases) == 0 {
		return nil, nil, errors.New("No databases found")
	}
	p.Session.SDK.CaptureReplicationTopo()
	p.Session.SDK.SetDBReplicas()

	nodeMap := make(map[string]bool)
	edgeMap := make(map[string]bool)

	nodes := []types.ReplTopoTreeNode{}
	edges := []types.ReplTopoTreeEdge{}

	for _, db := range p.Session.SDK.Databases {
		if _, ok := nodeMap[db.Shortname]; !ok {
			nodes = append(nodes, types.ReplTopoTreeNode{
				ID:   db.Shortname,
				Data: *db,
				Position: types.ReplTopoTreeNodePosition{
					X: 0,
					Y: 0,
				},
			})

			// // If the database is a DMS database, add an edge to the DMS node from its source
			// if strings.Contains(db.Host, "-dms") {
			// 	edges = append(edges, types.ReplTopoTreeEdge{
			// 		ID:       db.Source + "-" + db.Shortname,
			// 		Source:   db.Source,
			// 		Target:   db.Shortname,
			// 		EdgeType: "dms",
			// 	})
			// }
			nodeMap[db.Shortname] = true
		}

		for _, replica := range db.Replicas {
			if _, ok := edgeMap[db.Shortname+"-"+replica]; !ok {
				if db.Source != "" && db.Source == replica {

					// The bidirectional edge already exists, dont add it again
					if _, ok := edgeMap[replica+"-"+db.Shortname]; ok {
						continue
					}

					edges = append(edges, types.ReplTopoTreeEdge{
						ID:       db.Shortname + "-" + replica,
						Source:   db.Shortname,
						Target:   replica,
						EdgeType: "bidirectional",
						Animated: false,
					})
				} else {
					edges = append(edges, types.ReplTopoTreeEdge{
						ID:       db.Shortname + "-" + replica,
						Source:   db.Shortname,
						Target:   replica,
						EdgeType: "unidirectional",
						Animated: false,
					})
				}
				edgeMap[db.Shortname+"-"+replica] = true
			}

		}
	}

	return nodes, edges, nil
}
