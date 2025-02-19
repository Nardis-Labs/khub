import { Loading, TableToolbarAction, TableToolbarMenu, Tag, Tooltip } from '@carbon/react';
import React from 'react';
import { useDeletePodMutation, useGetPodsQuery } from '../../../service/khub';
import { InfoDrawer } from '../../components/InfoDrawer/InfoDrawer';
import { useSelector } from 'react-redux';
import { RootState, useAppDispatch } from '../../store';
import { updateTreeMapResourceDrawer } from '../../../service/resourceDrawerState';
import { RiBox2Fill } from "react-icons/ri";
import { SiOpensearch } from "react-icons/si";
import { InformationFilled, OverflowMenuVertical, Recycle } from '@carbon/icons-react';
import { ResourceDataTable } from '../../components/ResourceDataTable/ResourceDataTable';
import { updatePodFilter } from '../../../service/resource-filters';
import { updateNotifications } from '../../../service/notifications';

export const Pods = () => {
  const dispatch = useAppDispatch();
  const {data: pods = [], isLoading} = useGetPodsQuery({});

  const filterPods = (args: any) => {
    dispatch(updatePodFilter({filter: args.target.value}));
  };

  const [deletePod] = useDeletePodMutation();
  const handleDeletePod = (podName: string, namespace: string) => {
    deletePod({podName: podName, namespace: namespace}).unwrap()
    .then(() => dispatch(updateNotifications({notifications: [{notif: podName + ' delete initiated', status: 'success'}]})))
    .catch((error) => dispatch(updateNotifications({notifications: [{notif: 'Error deleting ' + podName + ' ' + JSON.stringify(error), status: 'error'}]})));
  };

  const podFilterState = useSelector((state: RootState) => state.podFilter);

  const headers = [
    {
      header: 'Name',
      key: 'name'
    },
    {
      header: 'Namespace',
      key: 'namespace',
      isSortable: true
    },
    {
      header: 'Containers',
      key: 'containers'
    },
    {
      header: 'Restarts',
      key: 'restarts',
      isSortable: true
    },
    {
      header: 'Age',
      key: 'age',
      isSortable: true
    },
    {
      header: 'Controlled By',
      key: 'controlledBy'
    },
    {
      header: '',
      key: 'controls'
    }
  ];

  const rows: any[] = pods.filter((pod: any) => {
    return (
      (pod.data?.metadata?.name && pod.data?.metadata.name.toLowerCase().includes(podFilterState.filter.toLowerCase())) ||
      (pod.data?.metadata?.namespace && pod.data?.metadata.namespace.toLowerCase().includes(podFilterState.filter.toLowerCase())) ||
      (pod.data?.status?.phase && pod.data?.status.phase.toLowerCase().includes(podFilterState.filter.toLowerCase())) ||
      (pod.data?.metadata?.ownerReference !== undefined && pod.data?.metadata?.ownerReferences.map((ref: any) => {
        return ref.kind + "/" + ref.name;
      }).join(', ').toLowerCase().includes(podFilterState.filter.toLowerCase()))
    );
  }).map((pod: any) => {
    const age = Math.floor(Math.abs(new Date().getTime() - new Date(pod.data.metadata.creationTimestamp).getTime()));
    const ageDays = age / 864e5;
    const ageHours = age / 36e5;
    const ageMinutes = age / 60000;
    const ageSeconds = age / 1000;
    return {
      id: pod.data.metadata.name,
      name: <div>{pod.data.metadata.name}
              {(pod.data.status.phase === 'Running' || pod.data.status.phase === 'Succeeded') && <Tag key={pod.data.metadata.name} type="green" title={pod.data.status.phase}>{pod.data.status.phase}</Tag>}
              {(pod.data.metadata.deletionTimestamp !== undefined) && <Tag key={pod.data.metadata.name} type="cool-gray" title='Terminating'>Terminating</Tag>}
              {(pod.data.status.phase === 'Pending') && <Tag key={pod.data.metadata.name} type="blue" title={pod.data.status.phase}>{pod.data.status.phase}</Tag>}
              {(pod.data.status.phase === 'Failed') && <Tag key={pod.data.metadata.name} type="red" title={pod.data.status.phase}>{pod.data.status.phase}</Tag>}
              {(pod.data.status.phase === 'Unknown') && <Tag key={pod.data.metadata.name} type="cool-gray" title={pod.data.status.phase}>{pod.data.status.phase}</Tag>}
            </div>
      ,
      namespace: pod.data.metadata.namespace,
      containers: pod.data.status.containerStatuses !== undefined && pod.data.status.containerStatuses.map((c: any) => {
          return (
            <div key={pod.data.name} style={{display: 'inline-flex'}}>
              {(c.state.running !== undefined || c.state.succeeded !== undefined) &&
                <Tooltip align="right" label={ c.state.running !== undefined ? `(${c.name}) running` : `(${c.name}) succeeded`}>
                  <RiBox2Fill key={c.name} color={"green"}/>
               </Tooltip>
              }
              {(c.state.pending !== undefined || c.state.waiting !== undefined) &&
                <Tooltip align="right" label={ c.state.pending !== undefined ? `(${c.name}) pending` : `(${c.name}) waiting`}>
                  <RiBox2Fill key={c.name} color={"orange"}/>
               </Tooltip>
              }
              {(c.state.failed !== undefined) &&
                <Tooltip align="right" label={`failed`}>
                  <RiBox2Fill key={c.name} color={"red"}/>
                </Tooltip>
              }
              {(c.state.terminated !== undefined || c.state.unknown !== undefined) &&
                <Tooltip align="right" label={ c.state.terminated !== undefined ? `(${c.name}) terminated` : `(${c.name}) unkown`}>
                  <RiBox2Fill key={c.name} color={"gray"}/>
               </Tooltip>
              }
            </div>
          );
      }),
      restarts: pod.data.status.containerStatuses !== undefined && pod.data.status.containerStatuses.reduce((acc: number, container: any) => container.restartCount, 0),
      age: ageHours >= 24 ? Math.floor(ageDays) + 'd' : ageHours > 1 ? Math.floor(ageHours) + 'h' : ageMinutes > 1 ? Math.floor(ageMinutes) + 'm' : Math.floor(ageSeconds) + 's',
      controlledBy: pod.data.metadata.ownerReferences ? pod.data.metadata.ownerReferences.map((ref: any) => {
        return ref.kind + "/" + ref.name;
      }).join(', ') : '',
      controls: <div style={{float: 'right'}}><TableToolbarMenu iconDescription='actions' renderIcon={OverflowMenuVertical}>
                      <TableToolbarAction onClick={() => openDrawer(pod.data)}>
                        <InformationFilled style={{marginRight: '10px'}}/> View Details
                      </TableToolbarAction>
                      <TableToolbarAction onClick={
                          () => {
                            const logsURL = `https://logs.observability.prod.smar.cloud/_dashboards/app/discover?security_tenant=global#/?_g=(filters:!(),refreshInterval:(pause:!t,value:0),time:(from:now-15m,to:now))&_a=(columns:!(_source),filters:!(('$state':(store:appState),meta:(alias:!n,disabled:!f,index:'27991650-c7cd-11ec-8f18-cd1446f63700',key:kubernetes.pod_name,negate:!f,params:(query:${pod.data!.metadata!.name}),type:phrase),query:(match_phrase:(kubernetes.pod_name:${pod.data!.metadata!.name!})))),index:'27991650-c7cd-11ec-8f18-cd1446f63700',interval:auto,query:(language:kuery,query:''),sort:!())`;
                            const w = window.open(logsURL, '_blank');
                            if (w) {
                                w.focus();
                            };
                          }}>
                        <SiOpensearch style={{marginRight: '10px'}}/> View Logs
                      </TableToolbarAction>
                      <TableToolbarAction onClick={
                        () => handleDeletePod(pod.data.metadata.name, pod.data.metadata.namespace)
                      } style={{color: 'red'}}>
                        <Recycle style={{marginRight: '10px'}}/> Delete Pod
                      </TableToolbarAction>
                     
                  </TableToolbarMenu></div>
    };
  });
  
  const resourceDrawer = useSelector((state: RootState) => state.treeMapResourceDrawer);
  
  const openDrawer = (data: any) => {
    const podData = {resourceData: data, resourceType: 'pod'}; 
    dispatch(updateTreeMapResourceDrawer({open: !resourceDrawer.open, data: podData}));
  };
  const closeDrawer = () => {
    dispatch(updateTreeMapResourceDrawer({open: !resourceDrawer.open, data: null}));
  };
  

  return (
    <div style={{height: '100%'}}>
     {isLoading &&
        <div>
          <Loading withOverlay={true}/>
        </div>
      }
      <ResourceDataTable 
        rows={rows} 
        headers={headers} 
        filterFunction={filterPods} 
        filterPlaceholder={'Filter pods'}
        filterValue={podFilterState.filter}
        title={'Pods (' + pods.length + ')'}
        batchActions={[]}
      />
      
      <InfoDrawer open={resourceDrawer.open} onClose={closeDrawer} direction="right" style={{ padding: '75px 20px 20px 20px' }}/>
    </div>
  );
};