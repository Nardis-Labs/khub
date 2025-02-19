import { Loading, TableToolbarAction, TableToolbarMenu, Tag } from '@carbon/react';
import React from 'react';
import { useGetConfigMapsQuery } from '../../../service/khub';
import { InfoDrawer } from '../../components/InfoDrawer/InfoDrawer';
import { useSelector } from 'react-redux';
import { RootState, useAppDispatch } from '../../store';
import { updateTreeMapResourceDrawer } from '../../../service/resourceDrawerState';
import { InformationFilled, OverflowMenuVertical } from '@carbon/icons-react';
import { ResourceDataTable } from '../../components/ResourceDataTable/ResourceDataTable';

export const ConfigMaps = () => {

  const {data: configmaps = [], isLoading} = useGetConfigMapsQuery({});
  const [configMapsFilter, setConfigMapsFilter] = React.useState('');
  const filterConfigMaps = (args: any) => {
    setConfigMapsFilter(args.target.value);
  };

  const dispatch = useAppDispatch();

  const headers = [
    {
      header: 'Name',
      key: 'name',
      isSortable: true
    },
    {
      header: 'Namespace',
      key: 'namespace',
      isSortable: true
    },
    {
      header: 'Immutable',
      key: 'immutable'
    },
    {
      header: 'Files',
      key: 'dataKeys'
    },
    {
      header: '',
      key: 'controls'
    }
  ];

  const rows: any[] = configmaps.filter((configmap: any) => {
    return (
      (configmap.data.metadata?.name && configmap.data.metadata.name.toLowerCase().includes(configMapsFilter.toLowerCase())) ||
      (configmap.data.metadata?.namespace && configmap.data.metadata.namespace.toLowerCase().includes(configMapsFilter.toLowerCase())) ||
      (configmap.data?.data !== undefined && Object.keys(configmap.data.data).find((key: string) => key.toLowerCase().includes(configMapsFilter.toLowerCase()))));
  }).map((configmap: any) => {
    return {
      id: configmap.data?.metadata?.name + '-' + configmap.data?.metadata?.namespace,
      name: configmap.data?.metadata?.name,
      namespace: configmap.data?.metadata?.namespace,
      immutable: configmap.data?.immutable ? 'Yes' : 'No',
      dataKeys: configmap.data?.data !== undefined ? Object.keys(configmap.data.data).map((key: string) => {
        return <Tag key={key} type='purple'>{key}</Tag>;
      }) : "null",
      controls: <div style={{float: 'right'}}><TableToolbarMenu iconDescription='actions' renderIcon={OverflowMenuVertical} style={{float: 'right'}}>
                      <TableToolbarAction onClick={() => openDrawer(configmap.data)}>
                        <InformationFilled style={{marginRight: '10px'}}/>View Details
                      </TableToolbarAction> 
                  </TableToolbarMenu></div>
    };
  });
  
  const resourceDrawer = useSelector((state: RootState) => state.treeMapResourceDrawer);
  const openDrawer = (data: any) => {
    const jobsData = {resourceData: data, resourceType: 'configmap'};
    dispatch(updateTreeMapResourceDrawer({open: !resourceDrawer.open, data: jobsData}));
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
        filterFunction={filterConfigMaps} 
        filterPlaceholder="Filter config maps"
        filterValue={configMapsFilter}
        title={'ConfigMaps (' + configmaps.length + ')'}
        batchActions={[]}
      />
      <InfoDrawer open={resourceDrawer.open} onClose={closeDrawer} direction="right" style={{ padding: '75px 20px 20px 20px' }}/>
    </div>
  );
};