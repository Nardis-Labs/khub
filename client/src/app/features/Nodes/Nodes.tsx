import { Loading, ProgressBar, TableToolbarAction, TableToolbarMenu, Tag } from '@carbon/react';
import React from 'react';
import { useGetNodesQuery } from '../../../service/khub';
import { InfoDrawer } from '../../components/InfoDrawer/InfoDrawer';
import { useSelector } from 'react-redux';
import { RootState, useAppDispatch } from '../../store';
import { updateTreeMapResourceDrawer } from '../../../service/resourceDrawerState';
import { InformationFilled, OverflowMenuVertical } from '@carbon/icons-react';
import { ResourceDataTable } from '../../components/ResourceDataTable/ResourceDataTable';

export const Nodes = () => {

  const {data: nodes = [], isLoading} = useGetNodesQuery({});
  const [nodesFilter, setNodesFilter] = React.useState('');
  const filterNodes = (args: any) => {
    setNodesFilter(args.target.value);
  };

  const dispatch = useAppDispatch();

  const headers = [
    {
      header: 'Name',
      key: 'name',
      isSortable: true
    },
    {
      header: 'CPU (%)',
      key: 'cpu'
    },
    {
      header: 'Memory (%)',
      key: 'mem'
    },
    {
      header: 'Version',
      key: 'version',
      isSortable: true
    },
    {
      header: 'Age',
      key: 'age',
      isSortable: true
    },
    {
      header: 'Conditions',
      key: 'conditions'
    },
    {
      header: '',
      key: 'controls'
    }
  ];

  const rows: any[] = nodes.filter((node: any) => {
    return (
      (node?.data?.node?.metadata.name && node.data.node.metadata.name.toLowerCase().includes(nodesFilter.toLowerCase()))
    );
  }).map((node: any) => {
    const age = Math.floor(Math.abs(new Date().getTime() - new Date(node.data.node.metadata.creationTimestamp).getTime()));
    const ageDays = age / 864e5;
    const ageHours = age / 36e5;
    const ageMinutes = age / 60000;
    const ageSeconds = age / 1000;

    // This article helps with the calculations: https://hwchiu.medium.com/introduction-to-kubernetes-resources-capacity-and-allocatable-4dc1bfbd1caf
    const cpuUsage = node.data.metrics?.Usage?.cpu !== undefined ? (((parseInt(node.data.metrics.Usage.cpu) / 1000000000) / parseInt(node.data.node.status.capacity.cpu)) * 100) : 0;  
    const memUsage = node.data.metrics?.Usage?.memory !== undefined ? ((parseInt(node.data.metrics.Usage.memory) / parseInt(node.data.node.status.capacity.memory)) * 100) : 0;

    return {
      id: node.data.node.metadata.name,
      name: node.data.node.metadata.name,
      cpu: <ProgressBar label={cpuUsage.toFixed(1) + '%'} value={cpuUsage} />,
      mem: <ProgressBar label={memUsage.toFixed(1) + '%'} value={memUsage} />,
      version: node.data.node.status.nodeInfo.kubeletVersion,
      age: ageHours >= 24 ? Math.floor(ageDays) + 'd' : ageHours > 1 ? Math.floor(ageHours) + 'h' : ageMinutes > 1 ? Math.floor(ageMinutes) + 'm' : Math.floor(ageSeconds) + 's',
      conditions: node.data.node.status.conditions.map((condition: any) => {
        if (condition.type === 'Ready' && condition.status === 'True'){
          return (
            <Tag key={condition.type} type="green" title={condition.type}>{condition.type}</Tag>
          );
        }  
        if (condition.type === 'PIDPressure' && condition.status === 'True'){
          return (
            <Tag key={condition.type} type="red" title={condition.type}>{condition.type}</Tag>
          );
        } 
        if (condition.type === 'DiskPressure' && condition.status === 'True'){
          return (
            <Tag key={condition.type} type="red" title={condition.type}>{condition.type}</Tag>
          );
        }
         if (condition.type === 'MemoryPressure' && condition.status === 'True'){
          return (
            <Tag key={condition.type} type="red" title={condition.type}>{condition.type}</Tag>
          );
        }
        
      }),
      controls: <div style={{float: 'right'}}><TableToolbarMenu iconDescription='actions' renderIcon={OverflowMenuVertical}>
                      <TableToolbarAction onClick={() => openDrawer(node.data.node)}>
                        <InformationFilled style={{marginRight: '10px'}}/>View Details
                      </TableToolbarAction> 
                  </TableToolbarMenu></div>
    };
  });
  
  const resourceDrawer = useSelector((state: RootState) => state.treeMapResourceDrawer);
  const openDrawer = (data: any) => {
    const serviceData = {resourceData: data, resourceType: 'node'};
    dispatch(updateTreeMapResourceDrawer({open: !resourceDrawer.open, data: serviceData}));
  };
  const closeDrawer = () => {
    dispatch(updateTreeMapResourceDrawer({open: !resourceDrawer.open, data: null}));
  };
  

  return (
    <div>
      {isLoading &&
        <div>
          <Loading withOverlay={true}/>
        </div>
      }
      <ResourceDataTable 
        rows={rows} 
        headers={headers} 
        filterFunction={filterNodes} 
        filterPlaceholder="Filter nodes"
        filterValue={nodesFilter}
        title={'Nodes (' + nodes.length + ')'}
        batchActions={[]}
      />
      <InfoDrawer open={resourceDrawer.open} onClose={closeDrawer} direction="right" style={{ padding: '75px 20px 20px 20px' }}/>
    </div>
  );
};