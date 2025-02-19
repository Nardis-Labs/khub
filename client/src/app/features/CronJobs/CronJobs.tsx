import { Loading, TableToolbarAction, TableToolbarMenu } from '@carbon/react';
import React from 'react';
import { useGetCronJobsQuery } from '../../../service/khub';
import { InfoDrawer } from '../../components/InfoDrawer/InfoDrawer';
import { useSelector } from 'react-redux';
import { RootState, useAppDispatch } from '../../store';
import { updateTreeMapResourceDrawer } from '../../../service/resourceDrawerState';
import { InformationFilled, OverflowMenuVertical } from '@carbon/icons-react';
import { ResourceDataTable } from '../../components/ResourceDataTable/ResourceDataTable';
import { updateCronJobFilter } from '../../../service/resource-filters';

export const CronJobs = () => {

  const dispatch = useAppDispatch();

  const {data: cronjobs = [], isLoading} = useGetCronJobsQuery({});
  
  const cronjobsFilterState = useSelector((state: RootState) => state.cronJobFilter);
  const filterCronjobs = (args: any) => {
    dispatch(updateCronJobFilter({filter: args.target.value}));
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
      header: 'Schedule',
      key: 'schedule'
    },
    {
      header: 'Last Schedule',
      key: 'lastScheduleTime'
    },
    {
      header: 'Last Successful',
      key: 'lastSuccessfulTime'
    },
    {
      header: 'Suspend',
      key: 'suspend'
    },
    {
      header: '',
      key: 'controls'
    }
  ];

  const rows: any[] = cronjobs.filter((cronjob: any) => {
    return (
      (cronjob.data.metadata?.name && cronjob.data.metadata.name.toLowerCase().includes(cronjobsFilterState.filter.toLowerCase())) ||
      (cronjob.data.metadata?.namespace && cronjob.data.metadata.namespace.toLowerCase().includes(cronjobsFilterState.filter.toLowerCase())) ||
      (cronjob.data.metadata?.ownerReferences.map((ref: any) => {
        return ref.kind + "/" + ref.name;
      }).join(', ').toLowerCase().includes(cronjobsFilterState.filter.toLowerCase()))
    );
  }).map((cronjob: any) => {
    const lastScheduleTime = Math.floor(Math.abs(new Date().getTime() - new Date(cronjob.data.status.lastScheduleTime || 0).getTime()));
    const lastScheduleTimeHours = lastScheduleTime / 36e5;
    const lastScheduleTimeMinutes = lastScheduleTime / 60000;
    const lastScheduleTimeSeconds = lastScheduleTime / 1000;

    const lastSuccessfulTime = Math.floor(Math.abs(new Date().getTime() - new Date(cronjob.data.status.lastSuccessfulTime || 0).getTime()));
    const lastSuccessfulTimeHours = lastSuccessfulTime / 36e5;
    const lastSuccessfulTimeMinutes = lastSuccessfulTime / 60000;
    const lastSuccessfulTimeSeconds = lastSuccessfulTime / 1000;
    return {
      id: cronjob.data.metadata.namespace + '-' + cronjob.data.metadata.name,
      name: cronjob.data.metadata.name,
      schedule: cronjob.data.spec.schedule,
      namespace: cronjob.data.metadata.namespace,
      lastScheduleTime: cronjob.data.status.startTime && lastScheduleTime > 1 ? Math.floor(lastScheduleTimeHours) + 'h' : lastScheduleTimeMinutes > 1 ? Math.floor(lastScheduleTimeMinutes) + 'm' : Math.floor(lastScheduleTimeSeconds) + 's',
      lastSuccessfulTime: cronjob.data.status.completionTime && lastSuccessfulTime > 1 ? Math.floor(lastSuccessfulTimeHours) + 'h' : lastSuccessfulTimeMinutes > 1 ? Math.floor(lastSuccessfulTimeMinutes) + 'm' : Math.floor(lastSuccessfulTimeSeconds) + 's',
      suspend: cronjob.data.spec.suspend ? 'True' : 'False',
      controls: <div style={{float: 'right'}}><TableToolbarMenu iconDescription='actions' renderIcon={OverflowMenuVertical} style={{float: 'right'}}>
                      <TableToolbarAction onClick={() => openDrawer(cronjob.data)}>
                        <InformationFilled style={{marginRight: '10px'}}/>View Details
                      </TableToolbarAction> 
                  </TableToolbarMenu></div>
    };
  });
  
  const resourceDrawer = useSelector((state: RootState) => state.treeMapResourceDrawer);
  const openDrawer = (data: any) => {
    const jobsData = {resourceData: data, resourceType: 'cronjob'};
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
        filterFunction={filterCronjobs} 
        filterPlaceholder={'Filter cronjobs'}
        filterValue={cronjobsFilterState.filter}
        title={'CronJobs (' + cronjobs.length + ')'}
        batchActions={[]}
      />
      <InfoDrawer open={resourceDrawer.open} onClose={closeDrawer} direction="right" style={{ padding: '75px 20px 20px 20px' }}/>
    </div>
  );
};