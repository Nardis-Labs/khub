import { Loading, TableToolbarAction, TableToolbarMenu } from '@carbon/react';
import React from 'react';
import { useGetDaemonsetsQuery, useRolloutRestartMutation } from '../../../service/khub';
import { InfoDrawer } from '../../components/InfoDrawer/InfoDrawer';
import { useSelector } from 'react-redux';
import { RootState, useAppDispatch } from '../../store';
import { updateTreeMapResourceDrawer } from '../../../service/resourceDrawerState';
import { InformationFilled, OverflowMenuVertical, Recycle } from '@carbon/icons-react';
import { updateNotifications } from '../../../service/notifications';
import { ResourceDataTable } from '../../components/ResourceDataTable/ResourceDataTable';
import { updateDaemonsetFilter } from '../../../service/resource-filters';

export const Daemonsets = () => {

  const {data: daemonsets = [], isLoading} = useGetDaemonsetsQuery({});

  const dispatch = useAppDispatch();

  const daemonsetFilterState = useSelector((state: RootState) => state.daemonsetFilter);
  const filterDaemonsets = (args: any) => {
    dispatch(updateDaemonsetFilter({filter: args.target.value}));
  };


  const [rolloutRestart] = useRolloutRestartMutation();

  const handleRolloutRestart = (d: any) => {
    rolloutRestart({name: d.metadata.name, namespace: d.metadata.namespace, kind: 'daemonset', labels: d.metadata.labels}).unwrap()
      .then((payload) => dispatch(updateNotifications({notifications: [{notif: d.metadata.name + ' restart initiated: ' + payload, status: 'success'}]})))
      .catch((error) => dispatch(updateNotifications({notifications: [{notif: 'Error restarting ' + d.metadata.name + ' daemonset ' + JSON.stringify(error), status: 'error'}]})));
  };

  // Used when initiating a rolling restart of multiple daemonsets
  const handleBulkRolloutRestart = (name: string, namespace: string) => {
    const labels = daemonsets.filter((d: any) => d.data.metadata.name === name && d.data.metadata.namespace === namespace)[0].data.metadata.labels;
    rolloutRestart({name: name, namespace: namespace, kind: 'daemonset', labels}).unwrap()
      .then((payload: any) => dispatch(updateNotifications({notifications: [{notif: name + ' restart initiated: ' + payload, status: 'success'}]})))
      .catch((error) => dispatch(updateNotifications({notifications: [{notif: 'Error restarting ' + name + ' daemonset ' + JSON.stringify(error), status: 'error'}]})));
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
      key: 'desired',
      isSortable: true
    },
    {
      header: 'Age',
      key: 'age',
      isSortable: true
    },
    {
      header: '',
      key: 'controls'
    }
  ];

  const rows: any[] = daemonsets.filter((d: any) => {
    return (
      (d.data.metadata?.name && d.data.metadata.name.toLowerCase().includes(daemonsetFilterState.filter.toLowerCase())) ||
      (d.data.metadata?.namespace && d.data.metadata.namespace.toLowerCase().includes(daemonsetFilterState.filter.toLowerCase())) ||
      ((d.data.status.availableReplicas === d.data.status.replicas) && daemonsetFilterState.filter.toLowerCase() === 'running') ||
      ((d.data.status.availableReplicas !== d.data.status.replicas) && daemonsetFilterState.filter.toLowerCase() === 'pending')
    );
  }).map((daemonset: any) => {
    const age = Math.floor(Math.abs(new Date().getTime() - new Date(daemonset.data.metadata.creationTimestamp).getTime()));
    const ageDays = age / 864e5;
    const ageHours = age / 36e5;
    const ageMinutes = age / 60000;
    const ageSeconds = age / 1000;
    return {
      id: daemonset.data.metadata.namespace + '-' + daemonset.data.metadata.name,
      name: daemonset.data.metadata.name,
      namespace: daemonset.data.metadata.namespace,
      pods: daemonset.data.status.numberReady !== undefined ? daemonset.data.status.numberReady + '/' + daemonset.data.status.currentNumberScheduled : 0 + '/' + daemonset.data.status.currentNumberScheduled,
      desired: daemonset.data.status.desiredNumberScheduled,
      age: ageHours >= 24 ? Math.floor(ageDays) + 'd' : ageHours > 1 ? Math.floor(ageHours) + 'h' : ageMinutes > 1 ? Math.floor(ageMinutes) + 'm' : Math.floor(ageSeconds) + 's',
      controls: <div style={{float: 'right'}}><TableToolbarMenu iconDescription='actions' renderIcon={OverflowMenuVertical}>
                      <TableToolbarAction onClick={() => openDrawer(daemonset.data)}>
                        <InformationFilled style={{marginRight: '10px'}}/>View Details
                      </TableToolbarAction> 
                      <TableToolbarAction onClick={() => handleRolloutRestart(daemonset.data)}>
                         <Recycle style={{marginRight: '10px'}}/> Restart
                      </TableToolbarAction>
                  </TableToolbarMenu></div>
    };
  });
  
  const resourceDrawer = useSelector((state: RootState) => state.treeMapResourceDrawer);
  const openDrawer = (data: any) => {
    const daemonsetData = {resourceData: data, resourceType: 'daemonset'};
    dispatch(updateTreeMapResourceDrawer({open: !resourceDrawer.open, data: daemonsetData}));
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
        filterFunction={filterDaemonsets} 
        filterPlaceholder={'Filter daemonsets'}
        filterValue={daemonsetFilterState.filter}
        title={'Daemonsets (' + daemonsets.length + ')'}
        batchActions={[
          {
            actionFunc: (args: any) => {
              args.forEach((arg: any) => {
                handleBulkRolloutRestart(arg.cells[0].value, arg.cells[1].value);
              });
            },
            actionDescription: 'Initiate rolling restart of selected daemonsets',
            actionLabel: 'Rolling Restart',
            actionIcon: Recycle
          },
        ]}
      />
      <InfoDrawer open={resourceDrawer.open} onClose={closeDrawer} direction="right" style={{ padding: '75px 20px 20px 20px' }}/>
    </div>
  );
};