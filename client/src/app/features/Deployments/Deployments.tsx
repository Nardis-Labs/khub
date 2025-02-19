import { Button, ButtonSet, ComposedModal, Loading, ModalBody, ModalHeader, NumberInput, TableToolbarAction, TableToolbarMenu, Tag } from '@carbon/react';
import React from 'react';
import { useGetDeploymentsQuery, useRolloutRestartMutation, useScaleDeploymentMutation } from '../../../service/khub';
import { InfoDrawer } from '../../components/InfoDrawer/InfoDrawer';
import { useSelector } from 'react-redux';
import { RootState, useAppDispatch } from '../../store';
import { updateTreeMapResourceDrawer } from '../../../service/resourceDrawerState';
import { InformationFilled, OverflowMenuVertical, Recycle } from '@carbon/icons-react';
import { updateNotifications } from '../../../service/notifications';
import { ResourceDataTable } from '../../components/ResourceDataTable/ResourceDataTable';
import { updateDeployFilter } from '../../../service/resource-filters';

import { HiMiniArrowDown, HiMiniArrowsUpDown, HiMiniArrowUp } from "react-icons/hi2";

export const Deployments = () => {

  const {data: deployments = [], isLoading} = useGetDeploymentsQuery({});
  const dispatch = useAppDispatch();

  const deployFilterState = useSelector((state: RootState) => state.deployFilter);
  const filterDeploys = (args: any) => {
    dispatch(updateDeployFilter({filter: args.target.value}));
  };

  const [rolloutRestart] = useRolloutRestartMutation();

  const handleRolloutRestart = (deploy: any) => {
    rolloutRestart({name: deploy.metadata.name, namespace: deploy.metadata.namespace, kind: 'deployment', labels: deploy.metadata.labels}).unwrap()
      .then((payload: any) => dispatch(updateNotifications({notifications: [{notif: deploy.metadata.name + ' restart initiated: ' + payload, status: 'success'}]})))
      .catch((error) => dispatch(updateNotifications({notifications: [{notif: 'Error restarting ' + deploy.metadata.name + ' deployment ' + JSON.stringify(error), status: 'error'}]})));
  };

  // Used when initiating a rolling restart of multiple deployments
  const handleBulkRolloutRestart = (name: string, namespace: string) => {
    const labels = deployments.filter((deploy: any) => deploy.data.metadata.name === name && deploy.data.metadata.namespace === namespace)[0].data.metadata.labels;
    rolloutRestart({name: name, namespace: namespace, kind: 'deployment', labels: labels}).unwrap()
      .then((payload: any) => dispatch(updateNotifications({notifications: [{notif: name + ' restart initiated: ' + payload, status: 'success'}]})))
      .catch((error) => dispatch(updateNotifications({notifications: [{notif: 'Error restarting ' + name + ' deployment ' + JSON.stringify(error), status: 'error'}]})));
  };

  const [selectedDeploymentForScale, setSelectedDeploymentForScale]: any = React.useState({});
  const [desiredReplicas, setDesiredReplicas] = React.useState(0);
  const [scaleModalOpen, setScaleModalOpen] = React.useState(false);
  const [scaleDeployments] = useScaleDeploymentMutation();

  const handleScaleDeploymentSubmit = () => {
    scaleDeployments({
        name: selectedDeploymentForScale.metadata.name, 
        namespace: selectedDeploymentForScale.metadata.namespace, 
        replicas: desiredReplicas,
        labels: selectedDeploymentForScale.metadata.labels
      }).unwrap()
      .then((payload: any) => dispatch(updateNotifications({notifications: [{notif: selectedDeploymentForScale.metadata.name + ' scaled to ' + desiredReplicas + ' replicas: ' + payload, status: 'success'}]})))
      .catch((error: any) => dispatch(updateNotifications({notifications: [{notif: 'Error scaling ' + selectedDeploymentForScale.metadata.name + ' deployment ' + JSON.stringify(error), status: 'error'}]})));
    setScaleModalOpen(false);
    setSelectedDeploymentForScale({});
  };

  const [desiredBulkScalingPercentage, setDesiredBulkScalingPercentage] = React.useState(0);
  const [bulkScaleDirection, setBulkScaleDirection] = React.useState("");
  const [bulkScaleDeploymentSelection, setBulkScaleDeploymentSelection] = React.useState([]);
  const [bulkScaleModalOpen, setBulkScaleModalOpen] = React.useState(false);
  const handleBulkScaleAction = (name: string, namespace: string, replicas: number) => {
    scaleDeployments({
        name: name, 
        namespace: namespace, 
        replicas: replicas,
        labels: deployments.filter((deploy: any) => deploy.data.metadata.name === name && deploy.data.metadata.namespace === namespace)[0].data.metadata.labels
      }).unwrap()
      .then((payload: any) => dispatch(updateNotifications({notifications: [{notif: name + ' scaled to ' + replicas + ' replicas: ' + payload, status: 'success'}]})))
      .catch((error: any) => dispatch(updateNotifications({notifications: [{notif: 'Error scaling ' + name + ' deployment ' + JSON.stringify(error), status: 'error'}]})));
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

  const rows: any[] = deployments.filter((deploy: any) => {
    return (
      (deploy.data.metadata?.name && deploy.data.metadata.name.toLowerCase().includes(deployFilterState.filter.toLowerCase())) ||
      (deploy.data.metadata?.namespace && deploy.data.metadata.namespace.toLowerCase().includes(deployFilterState.filter.toLowerCase())) ||
      ((deploy.data.status.availableReplicas === deploy.data.status.replicas) && deployFilterState.filter.toLowerCase() === 'running') ||
      ((deploy.data.status.availableReplicas !== deploy.data.status.replicas) && deployFilterState.filter.toLowerCase() === 'pending')
    );
  }).map((deploy: any) => {
    const age = Math.floor(Math.abs(new Date().getTime() - new Date(deploy.data.metadata.creationTimestamp).getTime()));
    const ageDays = age / 864e5;
    const ageHours = age / 36e5;
    const ageMinutes = age / 60000;
    const ageSeconds = age / 1000;
    return {
      id: deploy.data.metadata.namespace + '-' + deploy.data.metadata.name,
      name: deploy.data.metadata.name,
      namespace: deploy.data.metadata.namespace,
      pods: deploy.data.status.readyReplicas !== undefined ? deploy.data.status.readyReplicas + '/' + deploy.data.status.replicas : 0 + '/' + deploy.data.status.replicas,
      replicas: deploy.data.status.replicas,
      age: ageHours >= 24 ? Math.floor(ageDays) + 'd' : ageHours > 1 ? Math.floor(ageHours) + 'h' : ageMinutes > 1 ? Math.floor(ageMinutes) + 'm' : Math.floor(ageSeconds) + 's',
      conditions: deploy.data.status.conditions.map((condition: any) => {
        if (condition.type === 'Progressing' && condition.status === 'False'){
          return (
            <Tag key={condition.type} type="red" title={condition.type}>{condition.type}</Tag>
          );
        } else if (condition.type === 'Progressing' && condition.status === 'True'){
          return (
            <Tag key={condition.type} type="blue" title={condition.type}>{condition.type}</Tag>
          );
        } 
        else if (condition.type === 'Available'){
          return (
            <Tag key={condition.type} type="green" title={condition.type}>{condition.type}</Tag>
          );
        }
      }),
      controls: <div style={{float: 'right'}}><TableToolbarMenu iconDescription='actions' renderIcon={OverflowMenuVertical}>
                      <TableToolbarAction onClick={() => openDrawer(deploy.data)}>
                        <InformationFilled style={{marginRight: '10px'}}/>View Details
                      </TableToolbarAction> 
                      <TableToolbarAction onClick={() => handleRolloutRestart(deploy.data)}>
                         <Recycle style={{marginRight: '10px'}}/> Restart
                      </TableToolbarAction>
                      <TableToolbarAction onClick={() => {
                        setSelectedDeploymentForScale(deploy.data);
                        setScaleModalOpen(true);
                      }}>
                         <HiMiniArrowsUpDown style={{marginRight: '10px'}}/> Scale
                      </TableToolbarAction>
                  </TableToolbarMenu></div>
    };
  });
  
  const resourceDrawer = useSelector((state: RootState) => state.treeMapResourceDrawer);
  const openDrawer = (data: any) => {
    const deployData = {resourceData: data, resourceType: 'deployment'};
    dispatch(updateTreeMapResourceDrawer({open: !resourceDrawer.open, data: deployData}));
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
        filterFunction={filterDeploys} 
        filterPlaceholder={'Filter deployments'}
        filterValue={deployFilterState.filter}
        title={'Deployments (' + deployments.length + ')'}
        batchActions={[
          {
            actionFunc: (args: any) => {
              args.forEach((arg: any) => {
                handleBulkRolloutRestart(arg.cells[0].value, arg.cells[1].value);
              });
            },
            actionDescription: 'Initiate rolling restart of selected deployments',
            actionLabel: 'Rolling Restart',
            actionIcon: Recycle
          },
          {
            actionFunc: (args: any) => {
              setBulkScaleDirection("up");
              setBulkScaleDeploymentSelection(args);
              setDesiredBulkScalingPercentage(0);
              setBulkScaleModalOpen(true);
            },
            actionDescription: 'Scale selected deployments up by percentage',
            actionLabel: 'Scale Up (%)',
            actionIcon: HiMiniArrowUp
          },
          {
            actionFunc: (args: any) => {
              setBulkScaleDirection("down");
              setBulkScaleDeploymentSelection(args);
              setDesiredBulkScalingPercentage(0);
              setBulkScaleModalOpen(true);
            },
            actionDescription: 'Scale selected deployments up by percentage',
            actionLabel: 'Scale Down (%)',
            actionIcon: HiMiniArrowDown
          }
        ]}
      />
      <InfoDrawer open={resourceDrawer.open} onClose={closeDrawer} direction="right" style={{ padding: '75px 20px 20px 20px' }}/>
      
      <ComposedModal open={bulkScaleModalOpen} onClose={() => {
        setBulkScaleModalOpen(false);
        setBulkScaleDeploymentSelection([]);
        setDesiredBulkScalingPercentage(0);
        setBulkScaleDirection("");
      }}>
        <ModalHeader label="" title={"Bulk scale " + bulkScaleDirection + " by percentage"} />
        <ModalBody>
          {bulkScaleDeploymentSelection.map((deploy: any) => {
            return <div key={deploy.cells[0].value}>
                  <Tag key={deploy.cells[0].value} 
                        type="blue" 
                        title={deploy.cells[0].value}>
                          {deploy.cells[0].value} : {bulkScaleDirection === "up" ? Math.ceil(deploy.cells[3].value + (deploy.cells[3].value * (desiredBulkScalingPercentage / 100))) : Math.floor(deploy.cells[3].value - (deploy.cells[3].value * (desiredBulkScalingPercentage / 100)))}
                    </Tag>
                  </div>;
          })}
          <NumberInput id="carbon-number" 
            min={0} 
            max={200} 
            value={desiredBulkScalingPercentage} 
            label="By which percentage would you like to scale these deployments?" 
            invalidText="Desired percentage must be between 0 and 100" 
            onChange={(e: any, data: any) => setDesiredBulkScalingPercentage(data.value)}
            />
          <ButtonSet style={{marginTop: '20px'}}>
            <Button kind="primary" onClick={() => {
              bulkScaleDeploymentSelection.forEach((deploy: any) => {
                const currentDesiredReplicas = deploy.cells[3].value;
                const desiredReplicas = bulkScaleDirection === "up" ? Math.ceil(currentDesiredReplicas + (currentDesiredReplicas * (desiredBulkScalingPercentage / 100))) : Math.floor(currentDesiredReplicas - (currentDesiredReplicas * (desiredBulkScalingPercentage / 100)));
                handleBulkScaleAction(deploy.cells[0].value, deploy.cells[1].value, desiredReplicas);
              });
              setBulkScaleModalOpen(false);
              setBulkScaleDeploymentSelection([]);
              setDesiredBulkScalingPercentage(0);
              setBulkScaleDirection("");
            }}>
              Submit
            </Button>
            <Button kind="secondary" onClick={() => {
                setBulkScaleModalOpen(false);
                setBulkScaleDeploymentSelection([]);
                setDesiredBulkScalingPercentage(0);
                setBulkScaleDirection("");
              }}>
              Cancel
            </Button>
          </ButtonSet>
        </ModalBody>
      </ComposedModal>

      <ComposedModal open={scaleModalOpen} onClose={() => {
        setScaleModalOpen(false);
        setSelectedDeploymentForScale({});
      }}>
        <ModalHeader label="" title={"Scale " + selectedDeploymentForScale.metadata?.name} />
        <ModalBody>
          <NumberInput id="carbon-number" 
            min={0} 
            max={100} 
            value={selectedDeploymentForScale.spec?.replicas} 
            label="Scale this deployment up or down. Must be a positive integer between 0 and 100." 
            invalidText="Desired replicas must be between 0 and 100" 
            onChange={(e: any, data: any) => setDesiredReplicas(data.value)}
            />
          <ButtonSet style={{marginTop: '20px'}}>
            <Button kind="primary" onClick={() => handleScaleDeploymentSubmit()}>
              Submit
            </Button>
            <Button kind="secondary" onClick={() => {
                setScaleModalOpen(false);
                setSelectedDeploymentForScale({});
              }}>
              Cancel
            </Button>
          </ButtonSet>
        </ModalBody>
      </ComposedModal>
    </div>
  );
};

