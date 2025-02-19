import { Loading, TableToolbarAction, TableToolbarMenu } from '@carbon/react';
import React from 'react';
import { useGetStatefulsetsQuery, useRolloutRestartMutation } from '../../../service/khub';
import { InfoDrawer } from '../../components/InfoDrawer/InfoDrawer';
import { useSelector } from 'react-redux';
import { RootState, useAppDispatch } from '../../store';
import { updateTreeMapResourceDrawer } from '../../../service/resourceDrawerState';
import { InformationFilled, OverflowMenuVertical, Recycle } from '@carbon/icons-react';
import { updateNotifications } from '../../../service/notifications';
import { ResourceDataTable } from '../../components/ResourceDataTable/ResourceDataTable';
import { updateStatefulsetFilter } from '../../../service/resource-filters';

export const Statefulsets = () => {

  const {data: statefulsets = [], isLoading} = useGetStatefulsetsQuery({});
  
  const dispatch = useAppDispatch();

  const statefulsetFilterState = useSelector((state: RootState) => state.statefulsetFilter);
  const filterStatefulsets = (args: any) => {
    dispatch(updateStatefulsetFilter({filter: args.target.value}));
  };

  const [rolloutRestart] = useRolloutRestartMutation();

  const handleRolloutRestart = (statefulset: any) => {
    rolloutRestart({name: statefulset.data?.metadata.name, namespace: statefulset.data?.metadata.namespace, kind: 'statefulset', labels: statefulset.data?.metadata.labels}).unwrap()
      .then((payload) => dispatch(updateNotifications({notifications: [{notif: statefulset.data?.metadata.name + ' restart initiated: ' + payload, status: 'success'}]})))
      .catch((error) => dispatch(updateNotifications({notifications: [{notif: 'Error restarting ' + statefulset.data?.metadata.name + ' statefulset ' + JSON.stringify(error), status: 'error'}]})));
  };

  // Used when initiating a rolling restart of multiple statefulsets
  const handleBulkRolloutRestart = (name: string, namespace: string) => {
    const labels = statefulsets.filter((s: any) => s.data.metadata.name === name && s.data.metadata.namespace === namespace)[0].data.metadata.labels;
    rolloutRestart({name: name, namespace: namespace, kind: 'statefulset', labels: labels}).unwrap()
      .then((payload: any) => dispatch(updateNotifications({notifications: [{notif: name + ' restart initiated: ' + payload, status: 'success'}]})))
      .catch((error) => dispatch(updateNotifications({notifications: [{notif: 'Error restarting ' + name + ' statefulset ' + JSON.stringify(error), status: 'error'}]})));
  };

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
      header: 'Pods',
      key: 'pods'
    },
    {
      header: 'Desired',
      key: 'replicas',
      isSortable: true
    },
    {
      header: 'Age',
      key: 'age'
    },
    {
      header: '',
      key: 'controls'
    }
  ];

  const rows: any[] = statefulsets.filter((statefulset: any) => {
    return (
      (statefulset.data?.metadata?.name && statefulset.data?.metadata.name.toLowerCase().includes(statefulsetFilterState.filter.toLowerCase())) ||
      (statefulset.data?.metadata?.namespace && statefulset.data?.metadata.namespace.toLowerCase().includes(statefulsetFilterState.filter.toLowerCase())) ||
      ((statefulset.data.status.availableReplicas === statefulset.data.status.replicas) && statefulsetFilterState.filter.toLowerCase() === 'running') ||
      ((statefulset.data.status.availableReplicas !== statefulset.data.status.replicas) && statefulsetFilterState.filter.toLowerCase() === 'pending')
    );
  }).map((statefulset: any) => {
    const age = Math.floor(Math.abs(new Date().getTime() - new Date(statefulset.data?.metadata.creationTimestamp).getTime()));
    const ageHours = age / 36e5;
    const ageMinutes = age / 60000;
    const ageSeconds = age / 1000;
    return {
      id: statefulset.data.metadata.namespace + '-' + statefulset.data?.metadata.name,
      name: statefulset.data?.metadata.name,
      namespace: statefulset.data?.metadata.namespace,
      pods: statefulset.data?.status.readyReplicas !== undefined ? statefulset.data?.status.readyReplicas + '/' + statefulset.data?.status.replicas : 0 + '/' + statefulset.data?.status.replicas,
      replicas: statefulset.data?.status.replicas,
      age: ageHours > 1 ? Math.floor(ageHours) + 'h' : ageMinutes > 1 ? Math.floor(ageMinutes) + 'm' : Math.floor(ageSeconds) + 's',
      controls: <div style={{float: 'right'}}><TableToolbarMenu iconDescription='actions' renderIcon={OverflowMenuVertical} style={{float: 'right'}}>
                      <TableToolbarAction onClick={() => openDrawer(statefulset.data)}>
                        <InformationFilled style={{marginRight: '10px'}}/>View Details
                      </TableToolbarAction> 
                      <TableToolbarAction onClick={() => handleRolloutRestart(statefulset)}>
                         <Recycle style={{marginRight: '10px'}}/> Restart
                      </TableToolbarAction>
                  </TableToolbarMenu></div>
    };
  });
  
  const resourceDrawer = useSelector((state: RootState) => state.treeMapResourceDrawer);
  const openDrawer = (data: any) => {
    const statefulsetData = {resourceData: data, resourceType: 'statefulset'};
    dispatch(updateTreeMapResourceDrawer({open: !resourceDrawer.open, data: statefulsetData}));
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
        filterFunction={filterStatefulsets} 
        filterPlaceholder={'Filter statefulsets'}
        filterValue={statefulsetFilterState.filter}
        title={'Statefulsets (' + statefulsets.length + ')'}
        batchActions={[
          {
            actionFunc: (args: any) => {
              args.forEach((arg: any) => {
                handleBulkRolloutRestart(arg.cells[0].value, arg.cells[1].value);
              });
            },
            actionDescription: 'Initiate rolling restart of selected statefulsets',
            actionLabel: 'Rolling Restart',
            actionIcon: Recycle
          },
        ]}
      />
      <InfoDrawer open={resourceDrawer.open} onClose={closeDrawer} direction="right" style={{ padding: '75px 20px 20px 20px' }}/>
    </div>
  );
};