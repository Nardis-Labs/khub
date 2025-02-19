
import React, { useCallback, useEffect, useMemo, useState } from "react";
import { useGetMySQLReplicationTopologyGraphQuery } from "../../../service/khub";

import 'reactflow/dist/style.css';

import ReactFlow, {
  addEdge,
  ConnectionLineType,
  useNodesState,
  useEdgesState,
  Background,
  MiniMap,
  isNode,
  isEdge,
  MarkerType,
  Controls,
  Panel
} from 'reactflow';

import { stratify, tree } from 'd3';

import { ReplTopoNode } from "./CustomNodes";
import { BsArrows } from "react-icons/bs";
import { TbArrowWaveRightDown } from "react-icons/tb";

import { TbDatabase, TbDatabaseHeart, TbDatabaseStar } from "react-icons/tb";

const nodeWidth = 285;
const nodeHeight = 260;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const getLayoutedElements = (nodes: any[], edges: any[], direction = 'TB') => {
  // Create a hierarchy from the nodes and edges
  const root = stratify()
    .id((d: any) => d.id)
    .parentId((d: any) => {
      const edge = edges.find(edge => edge.target === d.id);
      return edge ? edge.source : null;
    })(nodes);

  // Create a tree layout
  const treeLayout = tree().nodeSize([nodeWidth, nodeHeight]);

  // Apply the layout to the hierarchy
  const treeData = treeLayout(root);

  // Map the calculated positions back to the nodes
  const updatedNodes = nodes.map(node => {
    const treeNode = treeData.descendants().find(d => d.id === node.id);
    return {
      ...node,
      type: 'replTopoNode',
      position: {
        x: treeNode ? treeNode.x - nodeWidth / 2 : 0,
        y: treeNode ? treeNode.y - nodeHeight / 2 : 0
      }
    };
  });

  // Update edges with styles and handles
  const updatedEdges = edges.map(edge => {
    if (edge.edgeType === 'unidirectional' || edge.edgeType === 'dms') {
      return { 
        ...edge, 
        sourceHandle: 'right',
        targetHandle: 'left',
        type: ConnectionLineType.SimpleBezier, 
        animated: false,
        style: {
          strokeWidth: 1,
          stroke: '#fafafa',
        },
        markerEnd: {
          type: MarkerType.Arrow,
          width: 25,
          height: 25,
          color: '#fafafa'
        }
      };
    } else if (edge.edgeType === 'bidirectional') {
      return { 
        ...edge,
        sourceHandle: 'bottom',
        targetHandle: 'top',
        type: ConnectionLineType.SimpleBezier, 
        animated: true, 
        style: {
          strokeWidth: 2,
          stroke: '#FF0072',
        },
        markerEnd: {
          type: MarkerType.ArrowClosed, 
          width: 20,
          height: 20,
          color: '#FF0072'
        }, 
        markerStart:{
          type: MarkerType.ArrowClosed, 
          width: 20,
          height: 20,
          color: '#FF0072'
        }
      };
    }
  });

  return {
    nodes: updatedNodes,
    edges: updatedEdges
  };
};

export const MySQLReplTopo = () => {
  const {data: topologyGraph} = useGetMySQLReplicationTopologyGraphQuery({});

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  
  useEffect(() => {
    if (topologyGraph) {
      const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(topologyGraph!.nodes, topologyGraph!.edges, 'LR');
      setNodes(layoutedNodes);
      setEdges(layoutedEdges);
    }
  }, [topologyGraph, setNodes, setEdges]);
 
  const onConnect = useCallback(
    (params: any) =>
      setEdges((eds) =>
        addEdge({ ...params, type: ConnectionLineType.SmoothStep, animated: true }, eds)
      ),
    [setEdges]
  );

  const nodeTypes = useMemo(() => ({ replTopoNode: ReplTopoNode }), []);

  const getAllIncomers = (node: any, nodes: any[], edges: any[], prevIncomers: any[] = []) => {
    const incomers = getIncomers(node, nodes, edges);
    const result = incomers.reduce((memo, incomer) => {
      memo.push(incomer);

      if (prevIncomers.findIndex((n) => n.id == incomer.id) == -1) {
        prevIncomers.push(incomer);

        getAllIncomers(incomer, nodes, edges, prevIncomers).forEach((foundNode: any) => {
          memo.push(foundNode);

          if (prevIncomers.findIndex((n) => n.id == foundNode.id) == -1) {
            prevIncomers.push(incomer);
          }
        });
      }
      return memo;
    }, []);
    return result;
  };

  const getAllOutgoers = (node: any, nodes: any[], edges: any[], prevOutgoers: any[] = []) => {
    const outgoers = getOutgoers(node, nodes, edges);
    return outgoers.reduce((memo, outgoer) => {
      memo.push(outgoer);

      if (prevOutgoers.findIndex((n) => n.id == outgoer.id) == -1) {
        prevOutgoers.push(outgoer);

        getAllOutgoers(outgoer, nodes, edges, prevOutgoers).forEach((foundNode: any) => {
          memo.push(foundNode);

          if (prevOutgoers.findIndex((n) => n.id == foundNode.id) == -1) {
            prevOutgoers.push(foundNode);
          }
        });
      }
      return memo;
    }, []);
  };

  const getIncomers = (node: any, nodes: any[], edges: any[]): any[] => {
    if (!isNode(node)) {
      return [];
    }
  
    const incomersIds = edges.filter((e) => e.target === node.id).map((e) => e.source);
  
    return nodes.filter((e) =>
      incomersIds
        .map((id) => {
          const matches = /([\w-^]+)__([\w-]+)/.exec(id);
          if (matches === null) {
            return id;
          }
          return matches[1];
        })
        .includes(e.id)
    );
  };
  
  const getOutgoers = (node: any, nodes: any[], edges: any[]): any[] => {
    if (!isNode(node)) {
      return [];
    }
  
    const outgoerIds = edges.filter((e) => e.source === node.id).map((e) => e.target);
  
    return nodes.filter((n) =>
      outgoerIds
        .map((id) => {
          const matches = /([\w-^]+)__([\w-]+)/.exec(id);
          if (matches === null) {
            return id;
          }
          return matches[1];
        })
        .includes(n.id)
    );
  };

  const highlightPath = (node: any, nodes: any[], edges: any[], selection: any) => {
    if (node && [...nodes, ...edges]) {
      const allIncomers = getAllIncomers(node, nodes, edges);
      const allOutgoers = getAllOutgoers(node, nodes, edges);
  
      setNodes((prevElements) => {
        return prevElements?.map((elem) => {
          const incomerIds = allIncomers.map((i: any) => i.id);
          const outgoerIds = allOutgoers.map((o: any) => o.id);
  
          if (isNode(elem) && (allOutgoers.length > 0 || allIncomers.length > 0)) {
            const highlight = elem.id === node.id || incomerIds.includes(elem.id) || outgoerIds.includes(elem.id);
  
            elem.style = {
              ...elem.style,
              opacity: highlight ? 1 : 0.25
            };
          }
          return elem;
        });
      });

      setEdges((prevElements) => {
        return prevElements?.map((elem) => {
          if (isEdge(elem)) {
            if (selection) {
              const highlightIncoming = allIncomers.map((i: any) => i.id).includes(elem.source) && (allIncomers.map((i: any) => i.id).includes(elem.target) || node.id === elem.target);
              const highlightOutgoing = allOutgoers.map((o: any) => o.id).includes(elem.target) && (allOutgoers.map((o: any) => o.id).includes(elem.source) || node.id === elem.source);
              elem.animated = true;
  
              // eslint-disable-next-line @typescript-eslint/ban-ts-comment
              {/* @ts-ignore */}
              if (elem.edgeType === 'unidirectional') {
                elem.style = {stroke: '#fafafa'};
                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                {/* @ts-ignore */}
              } else if (elem.edgeType === 'bidirectional') {
                elem.style = {stroke: '#FF0072'};
              }

              elem.style = {
                ...elem.style,
                strokeWidth: highlightIncoming || highlightOutgoing ? 2 : 0.5,
                opacity: highlightIncoming || highlightOutgoing ? 1 : 0.25,
                transition: 'opacity 0.3s'
              };
            }
          }
  
          return elem;
        });
      });
    }
  };
  
  const resetNodeStyles = () => {
    setNodes((prevElements) => {
      return prevElements?.map((elem: any) => {
        if (isNode(elem)) {
          elem.style = {
            ...elem.style,
            opacity: 1,
            transition: 'opacity 0.3s'
          };
        } else {
          elem.animated = false;
          elem.style = {
            ...elem.style,
            stroke: '#b1b1b7',
            opacity: 1,
            transition: 'opacity 0.3s'
          };
        }
  
        return elem;
      });
    });

    setEdges((prevElements) => {
      return prevElements?.map((elem: any) => {
        if (isEdge(elem)) {
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          {/* @ts-ignore */}
          if (elem.edgeType === 'unidirectional' || elem.edgeType === 'dms') {
            elem.style = {
              ...elem.style,
              strokeWidth: 1.5,
              opacity: 1,
              transition: 'opacity 0.3s'
            };
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            {/* @ts-ignore */}
          } else if (elem.edgeType === 'bidirectional') {
            elem.style = {
              ...elem.style,
              strokeWidth: 2,
              stroke: '#FF0072',
              opacity: 1,
              transition: 'opacity 0.3s'
            };
          }
        }
  
        return elem;
      });
    });
  };

  const [selectedNode] = useState();

  return (
    <div style={{ height: 1050, backgroundColor: '#555555' }}>
      <ReactFlow
        nodeTypes={nodeTypes}
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        connectionLineType={ConnectionLineType.SmoothStep}
        fitView
        minZoom={0.2}
        maxZoom={4}
        preventScrolling
        onNodeMouseEnter={(_event, node) => !selectedNode && highlightPath(node, nodes, edges, true)}
        onNodeMouseLeave={() => !selectedNode && resetNodeStyles()}
      >
        <Background color="#aaa" gap={16} />
        <MiniMap />
        <Controls />
        <Panel className='cds--tile' position="top-right">
          <BsArrows size={25} color="rgb(255 0 114)"/> Bidirectional replication
          <br/>
          <br/>
          <TbArrowWaveRightDown color='#fafafa' size={25}/> Standard replication
          <br/>
          <br/>
          <TbDatabaseHeart color='#b44e4e' size={23}/> DB Primary source
          <br/>
          <br/>
          <TbDatabaseStar color='#24a148' size={23}/> DB replication source
          <br/>
          <br/>
          <TbDatabase color='#ff832b' size={23}/> Standard DB replica
        </Panel>
      </ReactFlow>
    </div>
  );
};