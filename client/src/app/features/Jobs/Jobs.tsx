import { Loading, TableToolbarAction, TableToolbarMenu, Tag } from '@carbon/react';
import React from 'react';
import { useGetJobsQuery } from '../../../service/khub';
import { InfoDrawer } from '../../components/InfoDrawer/InfoDrawer';
import { useSelector } from 'react-redux';
import { RootState, useAppDispatch } from '../../store';
import { updateTreeMapResourceDrawer } from '../../../service/resourceDrawerState';
import { InformationFilled, OverflowMenuVertical } from '@carbon/icons-react';
import { ResourceDataTable } from '../../components/ResourceDataTable/ResourceDataTable';
import { updateJobFilter } from '../../../service/resource-filters';

export const Jobs = () => {

  const {data: jobs = [], isLoading} = useGetJobsQuery({});
  const dispatch = useAppDispatch();

  const jobsFilterState = useSelector((state: RootState) => state.jobFilter);
  const filterJobs = (args: any) => {
    dispatch(updateJobFilter({filter: args.target.value}));
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
      header: 'Conditions',
      key: 'conditions'
    },
    {
      header: 'Started',
      key: 'started'
    },
    {
      header: 'Completed',
      key: 'completed'
    },
    {
      header: '',
      key: 'controls'
    }
  ];

  const rows: any[] = jobs.filter((job: any) => {
    return (
      (job.data?.metadata?.name && job.data.metadata.name.toLowerCase().includes(jobsFilterState.filter.toLowerCase())) ||
      (job.data?.metadata?.namespace && job.data.metadata.namespace.toLowerCase().includes(jobsFilterState.filter.toLowerCase())) ||
      (job.data?.metadata?.ownerReferences !== undefined && job.data?.metadata?.ownerReferences.map((ref: any) => {
        return ref.kind + "/" + ref.name;
      }).join(', ').toLowerCase().includes(jobsFilterState.filter.toLowerCase()))
    );
  }).map((job: any) => {
    const started = Math.floor(Math.abs(new Date().getTime() - new Date(job.data.status.startTime || 0).getTime()));
    const startedHours = started / 36e5;
    const startedMinutes = started / 60000;
    const startedSeconds = started / 1000;

    const completed = Math.floor(Math.abs(new Date().getTime() - new Date(job.data.status.completionTime || 0).getTime()));
    const completedHours = completed / 36e5;
    const completedMinutes = completed / 60000;
    const completedSeconds = completed / 1000;
    return {
      id: job.data.metadata.namespace + '-' + job.data.metadata.name,
      name: <div>{job.data.metadata.name}
              {(job.data.status.phase === 'Running' || job.data.status.phase === 'Succeeded') && <Tag key={job.data.metadata.name} type="green" title={job.data.status.phase}>{job.data.status.phase}</Tag>}
              {(job.data.status.phase === 'Pending') && <Tag key={job.data.metadata.name} type="blue" title={job.data.status.phase}>{job.data.status.phase}</Tag>}
              {(job.data.status.phase === 'Failed') && <Tag key={job.data.metadata.name} type="red" title={job.data.status.phase}>{job.data.status.phase}</Tag>}
              {(job.data.status.phase === 'Unknown') && <Tag key={job.data.metadata.name} type="cool-gray" title={job.data.status.phase}>{job.data.status.phase}</Tag>}
            </div>
      ,
      namespace: job.data.metadata.namespace,
      started: job.data?.status?.startTime !== undefined && job.data.status.startTime && started > 1 ? Math.floor(startedHours) + 'h' : startedMinutes > 1 ? Math.floor(startedMinutes) + 'm' : Math.floor(startedSeconds) + 's',
      completed: job.data?.status?.completionTime !== undefined && job.data.status.completionTime && completed > 1 ? Math.floor(completedHours) + 'h' : completedMinutes > 1 ? Math.floor(completedMinutes) + 'm' : Math.floor(completedSeconds) + 's',
      conditions: job.data?.status?.conditions !== undefined && job.data.status.conditions.map((condition: any) => {
        if (condition.type === 'Failed' && condition.status === 'True'){
          return (
            <Tag key={condition.type} type="red" title={condition.type}>{condition.type}</Tag>
          );
        } else if (condition.type === 'Failed' && condition.status === 'False'){
          return (
            <Tag key={condition.type} type="cool-gray" title={condition.type}>{condition.type}</Tag>
          );
        } else if (condition.type === 'Complete' && condition.status === 'True'){
          return (
            <Tag key={condition.type} type="green" title={condition.type}>{condition.type}</Tag>
          );
        } else if (condition.type === 'Complete' && condition.status === 'False'){
          return (
            <Tag key={condition.type} type="cool-gray" title={condition.type}>{condition.type}</Tag>
          );
        }
      }),
      controls: <div style={{float: 'right'}}><TableToolbarMenu iconDescription='actions' renderIcon={OverflowMenuVertical} style={{float: 'right'}}>
                      <TableToolbarAction onClick={() => openDrawer(job.data)}>
                        <InformationFilled style={{marginRight: '10px'}}/>View Details
                      </TableToolbarAction> 
                  </TableToolbarMenu></div>
    };
  });
  
  const resourceDrawer = useSelector((state: RootState) => state.treeMapResourceDrawer);
  const openDrawer = (data: any) => {
    const jobsData = {resourceData: data, resourceType: 'job'};
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
        filterFunction={filterJobs} 
        filterPlaceholder={'Filter jobs'}
        filterValue={jobsFilterState.filter}
        title={'Jobs (' + jobs.length + ')'}
        batchActions={[]}
      />
      <InfoDrawer open={resourceDrawer.open} onClose={closeDrawer} direction="right" style={{ padding: '75px 20px 20px 20px' }}/>
    </div>
  );
};