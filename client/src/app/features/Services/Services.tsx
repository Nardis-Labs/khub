import { Loading, TableToolbarAction, TableToolbarMenu, Tag } from '@carbon/react';
import React from 'react';
import { useGetServicesQuery } from '../../../service/khub';
import { InfoDrawer } from '../../components/InfoDrawer/InfoDrawer';
import { useSelector } from 'react-redux';
import { RootState, useAppDispatch } from '../../store';
import { updateTreeMapResourceDrawer } from '../../../service/resourceDrawerState';
import { InformationFilled, OverflowMenuVertical } from '@carbon/icons-react';
import { ResourceDataTable } from '../../components/ResourceDataTable/ResourceDataTable';

export const Services = () => {

  const {data: services = [], isLoading} = useGetServicesQuery({});
  const [servicesFilter, setServicesFilter] = React.useState('');
  const filterServices = (args: any) => {
    setServicesFilter(args.target.value);
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
      header: 'Type',
      key: 'type',
      isSortable: true
    },
    {
      header: 'Cluster IP',
      key: 'clusterIP'
    },
    {
      header: 'Ports',
      key: 'ports'
    },
    {
      header: 'External IP',
      key: 'externalIP'
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

  const rows: any[] = services.filter((service: any) => {
    return (
      (service.data?.metadata?.name && service.data?.metadata.name.toLowerCase().includes(servicesFilter.toLowerCase())) ||
      (service.data?.metadata?.namespace && service.data?.metadata.namespace.toLowerCase().includes(servicesFilter.toLowerCase())) ||
      (service.data?.spec?.type && service.data.spec.type.toLowerCase().includes(servicesFilter.toLowerCase()))
    );
  }).map((service: any) => {
    const age = Math.floor(Math.abs(new Date().getTime() - new Date(service.data?.metadata.creationTimestamp).getTime()));
    const ageDays = age / 864e5;
    const ageHours = age / 36e5;
    const ageMinutes = age / 60000;
    const ageSeconds = age / 1000;
    return {
      id: service.data.metadata.namespace + '-' + service.data?.metadata.name,
      name: service.data?.metadata.name,
      namespace: service.data?.metadata.namespace,
      type: <div>
              {service.data?.spec.type === 'LoadBalancer' && <Tag type='green'>{service.data?.spec.type}</Tag>} 
              {service.data?.spec.type === 'NodePort' && <Tag type='blue'>{service.data?.spec.type}</Tag>} 
              {service.data?.spec.type === 'ClusterIP' && <Tag type='magenta'>{service.data?.spec.type}</Tag>}
            </div>,
      clusterIP: service.data?.spec.clusterIP,
      ports: service.data?.spec.ports.map((port: any) => {
        return (
          <div key={port.port}>
              {port.port + '/' + port.targetPort + ' ' + port.protocol} <br/>
          </div>
          
        );
      }),
      externalIP: service.data?.spec.externalIPs ? service.data?.spec.externalIPs.join(', ') : 'None',
      age: ageHours >= 24 ? Math.floor(ageDays) + 'd' : ageHours > 1 ? Math.floor(ageHours) + 'h' : ageMinutes > 1 ? Math.floor(ageMinutes) + 'm' : Math.floor(ageSeconds) + 's',
      controls: <div style={{float: 'right'}}><TableToolbarMenu iconDescription='actions' renderIcon={OverflowMenuVertical}>
                      <TableToolbarAction onClick={() => openDrawer(service.data)}>
                        <InformationFilled style={{marginRight: '10px'}}/>View Details
                      </TableToolbarAction> 
                  </TableToolbarMenu></div>
    };
  });
  
  const resourceDrawer = useSelector((state: RootState) => state.treeMapResourceDrawer);
  const openDrawer = (data: any) => {
    const serviceData = {resourceData: data, resourceType: 'service'};
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
        filterFunction={filterServices} 
        filterPlaceholder="Filter services"
        filterValue={servicesFilter}
        title={'Services (' + services.length + ')'}
        batchActions={[]}
      />
      <InfoDrawer open={resourceDrawer.open} onClose={closeDrawer} direction="right" style={{ padding: '75px 20px 20px 20px' }}/>
    </div>
  );
};