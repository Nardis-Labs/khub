import React from "react";
import cx from 'classnames';
import { Column, Grid, InlineLoading, Loading, Table, TableBody, TableCell, TableContainer, TableHead, TableHeader, TableRow, TableToolbar, TableToolbarContent, TableToolbarSearch, Tile } from '@carbon/react';
import { useGetClusterEventsQuery, useGetCronJobsQuery, useGetDaemonsetsQuery, useGetDeploymentsQuery, useGetJobsQuery, useGetPodsQuery, useGetReplicasetsQuery, useGetStatefulsetsQuery } from "../../../service/khub";
import { PieData } from "../../components/PieChart/PieChart";
import { renderCronJobStatus, renderDeploymentStats, renderJobStatus, renderPodStats, useWindowDimensions } from "../../../Util";

import { useAppDispatch } from '../../store';
import { updateCronJobFilter, updateDaemonsetFilter, updateDeployFilter, updateJobFilter, updatePodFilter, updateStatefulsetFilter } from "../../../service/resource-filters";
import { useNavigate } from "react-router-dom";


export const ClusterOverview = () => {
  const dispatch = useAppDispatch();
  const { width } = useWindowDimensions();
  const useResponsiveOffset = width > 1055;
  const classNameFirstColumn = cx({
    'cds--col-lg-13': true,
    'cds--offset-lg-3': useResponsiveOffset
  });

  const {data: pods = [], isLoading: podsLoading} = useGetPodsQuery({});
  const {data: deployments = [], isLoading: deploysLoading} = useGetDeploymentsQuery({});
  const {data: daemonsets = [], isLoading: daemonsetsLoading} = useGetDaemonsetsQuery({});
  const {data: replicasets = [], isLoading: replicasetsLoading} = useGetReplicasetsQuery({});
  const {data: statefulsets = [], isLoading: statefulsetsLoading} = useGetStatefulsetsQuery({});
  const {data: jobs = [], isLoading: jobsLoading} = useGetJobsQuery({});
  const {data: cronJobs = [], isLoading: cronjobsLoading} = useGetCronJobsQuery({});
  const {data: clusterEvents = []} = useGetClusterEventsQuery({});

  const [eventFilter, setEventFilter] = React.useState('');
  const filterEvents = (args: any) => {
    setEventFilter(args.target.value);
  };

  const navigate = useNavigate();
  const filterPods = (filter: string) => {
    dispatch(updatePodFilter({filter: filter}));
    navigate('/pods', {replace: true});
  };

  const filterDeploys = (filter: string) => {
    dispatch(updateDeployFilter({filter: filter}));
    navigate('/deployments', {replace: true});
  };

  const filterReplicaSets = (filter: string) => {
    console.log('warning: replicaset filter shows all deployments' + ':' + filter);
    dispatch(updateDeployFilter({filter: ''}));
    navigate('/deployments', {replace: true});
  };

  const filterDaemonsets = (filter: string) => {
    dispatch(updateDaemonsetFilter({filter: filter}));
    navigate('/daemonsets', {replace: true});
  };

  const filterStatefulsets = (filter: string) => {
    dispatch(updateStatefulsetFilter({filter: filter}));
    navigate('/statefulsets', {replace: true});
  };

  const filterJobs = (filter: string) => {
    dispatch(updateJobFilter({filter: filter}));
    navigate('/jobs', {replace: true});
  };

  const filterCronjobs = (filter: string) => {
    dispatch(updateCronJobFilter({filter: filter}));
    navigate('/cronjobs', {replace: true});
  };

  const isLoading = podsLoading && deploysLoading && daemonsetsLoading && replicasetsLoading && statefulsetsLoading && jobsLoading && cronjobsLoading;
  return (
    <div>
      {isLoading &&
        <div>
          <Loading withOverlay={true}/>
        </div>
      }
      <div className="cds--row">
        <div>
          <Grid fullWidth style={{ marginBottom: '10px' }}>
            <Column sm="100%">
              <Tile id="tile-1" style={{ height: '65px', width: '100%' }}>
                Overview
              </Tile>
            </Column>
          </Grid>
        </div>
      </div>
      <div className="cds--row">
        <div className={classNameFirstColumn}>
          <Grid fullWidth style={{ marginBottom: '30px' }}>
              <Column sm={4} className="pie-column">
                <Tile id="tile-1" style={{ height: '300px' }}>
                  {podsLoading ? <InlineLoading status="active" iconDescription="Loading" description="Loading data..." /> : 'Pods (' + pods.length + ')'}
                  <PieData data={renderPodStats(pods, filterPods)} />
                </Tile>
              </Column>
              <Column sm={4} className="pie-column">
                <Tile id="tile-1" style={{ height: '300px' }}>
                  {deploysLoading ? <InlineLoading status="active" iconDescription="Loading" description="Loading data..." /> : 'Deployments (' + deployments.length + ')'}
                  <PieData data={renderDeploymentStats(deployments, filterDeploys, 'deployment')} />
                </Tile>
              </Column>
              <Column sm={4} className="pie-column">
                <Tile id="tile-1" style={{ height: '300px' }}>
                 {replicasetsLoading ? <InlineLoading status="active" iconDescription="Loading" description="Loading data..." /> : 'Replicasets (' + replicasets.length + ')'}
                  <PieData data={renderDeploymentStats(replicasets, filterReplicaSets, 'replicaset')} />
                </Tile>
              </Column>   
              <Column sm={4} className="pie-column">
                <Tile id="tile-1" style={{ height: '300px' }}>
                {statefulsetsLoading ? <InlineLoading status="active" iconDescription="Loading" description="Loading data..." /> : 'Statefulsets (' + statefulsets.length + ')'}
                  <PieData data={renderDeploymentStats(statefulsets, filterStatefulsets, 'statefulset')} />
                </Tile>
              </Column>
              <Column sm={4} className="pie-column">
                <Tile id="tile-1" style={{ height: '300px' }}>
                {daemonsetsLoading ? <InlineLoading status="active" iconDescription="Loading" description="Loading data..." /> : 'Daemonsets (' + daemonsets.length + ')'}
                  <PieData data={renderDeploymentStats(daemonsets, filterDaemonsets, 'daemonset')} />
                </Tile>
              </Column>
              <Column sm={4} className="pie-column">
                <Tile id="tile-1" style={{ height: '300px' }}>
                {cronjobsLoading ? <InlineLoading status="active" iconDescription="Loading" description="Loading data..." /> : 'CronJobs (' + cronJobs.length + ')'}
                  <PieData data={renderCronJobStatus(cronJobs, filterCronjobs)} />
                </Tile>
              </Column>
              <Column sm={4} className="pie-column">
                <Tile id="tile-1" style={{ height: '300px', marginBottom: '0px' }}>
                {jobsLoading ? <InlineLoading status="active" iconDescription="Loading" description="Loading data..." /> : 'Jobs (' + jobs.length + ')'}
                  <PieData data={renderJobStatus(jobs, filterJobs)} />
                </Tile>
              </Column>
          </Grid>
        </div>
      </div>
      <div className="cds--row">
        <div className={classNameFirstColumn}>
          <Grid fullWidth>
            <Column sm="100%">
              <TableContainer title="Events" description="all cluster events" className='events-table'>
                <TableToolbar aria-label="k8s cluster events">
                  <TableToolbarContent>
                    <TableToolbarSearch placeholder="Filter events" onChange={filterEvents} persistent/>
                  </TableToolbarContent>
                </TableToolbar>
                <Table aria-label="sample table">
                  <TableHead >
                    <TableRow className="theader">
                      <TableHeader key="message">Message</TableHeader>
                      <TableHeader key="reason">Reason</TableHeader>
                      <TableHeader key="object">Object</TableHeader>
                      <TableHeader key="type">Type</TableHeader>
                      <TableHeader key="lastseen">Last Seen</TableHeader>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {clusterEvents.length === 0 &&
                        <TableRow key="no-events">
                          <TableCell key="message">No events</TableCell>
                          <TableCell key="reason"></TableCell>
                          <TableCell key="object"></TableCell>
                          <TableCell key="type"></TableCell>
                          <TableCell key="lastseen"></TableCell>
                        </TableRow>
                    }
                    {clusterEvents
                      .filter((event: any) => {
                        return (
                          (event.data.message !== undefined && event.data.message.toLowerCase().includes(eventFilter.toLowerCase())) ||
                          (event.data.reason !== undefined && event.data.reason.toLowerCase().includes(eventFilter.toLowerCase())) ||
                          (event.data.object !== undefined && event.data.object?.toLowerCase().includes(eventFilter.toLowerCase())) ||
                          (event.data.type !== undefined && event.data.type.toLowerCase().includes(eventFilter.toLowerCase()))
                        );
                      })
                      .map((event: any) => {
                        return (
                          <TableRow key={event.data.metadata.uid}>
                            <TableCell key={event.data.metadata.uid + 'message'}>{event.data.message}</TableCell>
                            <TableCell key={event.data.metadata.uid + 'reason'}>{event.data.reason}</TableCell>
                            <TableCell key={event.data.metadata.uid + 'object'}>{event.data.object}</TableCell>
                            <TableCell key={event.data.metadata.uid + 'type'}>{event.data.type}</TableCell>
                            <TableCell key={event.data.metadata.uid + 'lastseen'}>{event.data.interval}</TableCell>
                          </TableRow>
                        );
                      })}
                  </TableBody>
                </Table>
              </TableContainer>
            </Column>
          </Grid>
        </div>
      </div>
    </div>
  );
};