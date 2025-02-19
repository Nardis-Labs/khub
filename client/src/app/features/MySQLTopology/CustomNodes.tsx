import { Tile } from '@carbon/react';
import React from 'react';
import { Handle, Position } from 'reactflow';
import { TbDatabaseStar, TbDatabase, TbDatabaseHeart } from "react-icons/tb";


export const ReplTopoNode = ({data}: any) => {
  const innerData = () => {
    return (
      <div>
        <div>
          {data.replicas !== null && data.replicas.length > 0 && data.isPrimary &&
            <TbDatabaseHeart color='#b44e4e' size={50} style={{marginTop: '5px', marginBottom: '5px'}}/>
          }

          {data.replicas !== null && data.replicas.length > 0 && !data.isPrimary &&
            <TbDatabaseStar color='#24a148' size={50} style={{marginTop: '5px', marginBottom: '5px'}}/>
          }

          {(data.replicas === null || data.replicas.length === 0) && 
            <TbDatabase color='#ff832b' size={50} style={{marginTop: '5px', marginBottom: '5px'}}/>
          }
        </div>
        <strong style={{marginLeft: '10px', fontSize: '16px'}}>
            {data.shortName}
        </strong> <br/>
      </div>
    );
  };

  let borderColor = 'gray';
  if (data.replicas && data.replicas.length > 0 && data.isPrimary) {
    borderColor = '#b44e4e';
  } else if (data.replicas && data.replicas.length > 0 && !data.isPrimary) {
    borderColor = '#24a148';
  } else {
    borderColor = '#ff832b';
  }

  return (
    <>
      {(data.source !== null || data.source !== "") && <Handle style={{visibility: 'hidden'}} type="target" position={Position.Top} id="left"/>}
      {(data.source !== null || data.source !== "") && <Handle style={{visibility: 'hidden'}} type="target" position={Position.Top} id="top"/>}
        <Tile style={{
          borderRadius: '5%', 
          border: '1px solid ' + borderColor,
          borderRight: '3px solid ' + borderColor,
          borderLeft: '3px solid ' + borderColor, 
          maxWidth: '300px',
          minWidth: '200px',
          height: '115px',
          padding: '10px',
          fontSize: '10px',
        }}>
        {innerData()}
        </Tile>
        
      {data.replicas !== null && data.replicas.length > 0 && <Handle style={{visibility: 'hidden'}} type="source" position={Position.Right} id="right"/>}
      {data.replicas !== null && data.replicas.length > 0 && <Handle style={{visibility: 'hidden'}} type="source" position={Position.Bottom} id="bottom"/>}
    </>
  );
};