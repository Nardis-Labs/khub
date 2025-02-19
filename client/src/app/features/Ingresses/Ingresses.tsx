import { Loading, TableToolbarAction, TableToolbarMenu } from '@carbon/react';
import React from 'react';
import { useGetIngressesQuery } from '../../../service/khub';
import { InfoDrawer } from '../../components/InfoDrawer/InfoDrawer';
import { useSelector } from 'react-redux';
import { RootState, useAppDispatch } from '../../store';
import { updateTreeMapResourceDrawer } from '../../../service/resourceDrawerState';
import { InformationFilled, OverflowMenuVertical } from '@carbon/icons-react';
import { ResourceDataTable } from '../../components/ResourceDataTable/ResourceDataTable';

export const Ingresses = () => {

  const {data: ingresses = [], isLoading} = useGetIngressesQuery({});
  const [ingressesFilter, setIngressesFilter] = React.useState('');
  const filterIngresses = (args: any) => {
    setIngressesFilter(args.target.value);
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
      header: 'Rules',
      key: 'rules'
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

  const rows: any[] = ingresses.filter((ingress: any) => {
    return (
      (ingress.data.metadata?.name && ingress.data.metadata.name.toLowerCase().includes(ingressesFilter.toLowerCase())) ||
      (ingress.data.metadata?.namespace && ingress.data.metadata.namespace.toLowerCase().includes(ingressesFilter.toLowerCase()))
    );
  }).map((ingress: any) => {
    const age = Math.floor(Math.abs(new Date().getTime() - new Date(ingress.data.metadata.creationTimestamp).getTime()));
    const ageHours = age / 36e5;
    const ageMinutes = age / 60000;
    const ageSeconds = age / 1000;
    return {
      id: ingress.data?.metadata?.name !== undefined ? ingress.data.metadata.namespace + '-' + ingress.data.metadata.name : "undefined",
      name: ingress.data?.metadata?.name !== undefined ? ingress.data.metadata.name : "undefined",
      namespace: ingress.data?.metadata?.namespace !== undefined ? ingress.data.metadata.namespace : "undefined",
      rules: <div>
              {ingress.data.spec?.rules !== undefined && ingress.data.spec.rules.map((rule: any) => {
                const rls: string[] = rule.http.paths !== undefined && rule.http.paths.map((path: any) => {
                  let host = rule.host;
                  if (ingress.data.metadata?.annotations !== undefined && ingress.data.metadata.annotations['external-dns.alpha.kubernetes.io/hostname'] !== undefined) {
                    host = ingress.data.metadata.annotations['external-dns.alpha.kubernetes.io/hostname'];
                  }
                  if (path.backend.service.port.number === undefined) {
                    return host + ' → ' + path.backend.service.name;
                  }
                  return host + ' → ' + path.backend.service.name + ':' + path.backend.service.port.number;
                });

                return rls.join('\n');
              })}
            </div>,
      age: ageHours > 1 ? Math.floor(ageHours) + 'h' : ageMinutes > 1 ? Math.floor(ageMinutes) + 'm' : Math.floor(ageSeconds) + 's',
      controls: <div style={{float: 'right'}}><TableToolbarMenu iconDescription='actions' renderIcon={OverflowMenuVertical}>
                      <TableToolbarAction onClick={() => openDrawer(ingress.data)}>
                        <InformationFilled style={{marginRight: '10px'}}/>View Details
                      </TableToolbarAction> 
                  </TableToolbarMenu></div>
    };
  });
  
  const resourceDrawer = useSelector((state: RootState) => state.treeMapResourceDrawer);
  const openDrawer = (data: any) => {
    const ingressData = {resourceData: data, resourceType: 'ingress'};
    dispatch(updateTreeMapResourceDrawer({open: !resourceDrawer.open, data: ingressData}));
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
        filterFunction={filterIngresses} 
        filterPlaceholder="Filter ingresses"
        filterValue={ingressesFilter}
        title={'Ingresses (' + ingresses.length + ')'}
        batchActions={[]}
      />
      <InfoDrawer open={resourceDrawer.open} onClose={closeDrawer} direction="right" style={{ padding: '75px 20px 20px 20px' }}/>
    </div>
  );
};