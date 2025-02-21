import { Accordion, AccordionItem, Button, ButtonSet, ClickableTile, ComposedModal, Heading, ModalBody, ModalHeader, NumberInput, StructuredListBody, StructuredListCell, StructuredListHead, StructuredListRow, StructuredListWrapper, Tag, TextInput, Tile, Toggle, Tooltip } from "@carbon/react";
import React, { useEffect } from "react";
import { SiMysql } from "react-icons/si";
import { TbTrash, TbEdit } from "react-icons/tb";
import { PiPlusBold } from "react-icons/pi";
import { SiKubernetes } from "react-icons/si";
import { useDeleteMySQLDBInfoMutation, useGetMySQLDBCatalogQuery, useUpdateDynamicAppConfigMutation, useUpsertMySQLDBInfoMutation } from "../../../service/khub";
import { useAppDispatch } from "../../store";
import { updateNotifications } from "../../../service/notifications";
import { CheckmarkFilled, Misuse } from "@carbon/icons-react";
import { PiCaretCircleUpDownLight } from "react-icons/pi";
import { CgNametag } from "react-icons/cg";
import { IAppConfig } from "../../../service/types/AppConfig";

const DBSettingsTitle = () => {
  return (
    <>
    <div style={{ display: 'flex'}} >
      <SiMysql size={50}/>
      <h4 className='accordian-title'>Database Catalog</h4>
    </div>
    <span className='accordian-subtitle'>Manage MySQL databases represented in the replication topology</span>
    <br/>
    </>
  );
};

const KubernetesSettingsTile = () => {
  return (
    <>
    <div style={{ display: 'flex'}} >
      <SiKubernetes style={{marginTop: '8px'}} size={35}/>
      <h4 style={{marginTop: '13px', marginLeft: '10px'}}>Kubernetes Controls</h4>
    </div>
    <span className='accordian-subtitle'>Manage khub controls for kubernetes cluster</span>
    <br/>
    </>
  );
};


interface GeneralSettingsProps {
  appConfig: IAppConfig;
}

export const GeneralSettings = ({appConfig}: GeneralSettingsProps) => {
  const dispatch = useAppDispatch();

  /* MySQL DB Info Catalog queries and mutations
  /
  /
  */
  const {data: mysqlDBCatalog = []} = useGetMySQLDBCatalogQuery({});
  const [upsertMySQLDBInfo] = useUpsertMySQLDBInfoMutation();
  const [deleteMySQLDBInfo] = useDeleteMySQLDBInfoMutation();

  const [selectedMySQLDBHost, setSelectedMySQLDBHost] = React.useState<string>('');
  const [selectedMySQLDBShortname, setSelectedMySQLDBShortname] = React.useState<string>('');
  const [selectedMySQLDBPort, setSelectedMySQLDBPort] = React.useState<number>(3306);
  const [selectedMySQLDBUsername, setSelectedMySQLDBUsername] = React.useState<string>('');
  const [selectedMySQLDBIsPrimary, setSelectedMySQLDBIsPrimary] = React.useState<boolean>(false);

  const [mySQLDBInfoModalOpen, setMySQLDBInfoModalOpen] = React.useState(false);

  const resetSelectedMySQLDB = () => {
    setSelectedMySQLDBHost('');
    setSelectedMySQLDBShortname('');
    setSelectedMySQLDBPort(3306);
    setSelectedMySQLDBUsername('');
    setMySQLDBInfoModalOpen(false);
    setSelectedMySQLDBIsPrimary(false);
  };

  const handleMySQLInfoUpsertModalOpen = (db: any) => {
    setMySQLDBInfoModalOpen(true);
    setSelectedMySQLDBHost(db.host);
    setSelectedMySQLDBShortname(db.shortName);
    setSelectedMySQLDBPort(db.port);
    setSelectedMySQLDBUsername(db.username);
    setSelectedMySQLDBIsPrimary(db.isPrimary);
  };

  const handleUpsertMySQLDBInfo = () => {
    upsertMySQLDBInfo({host: selectedMySQLDBHost, shortName: selectedMySQLDBShortname, port: selectedMySQLDBPort, username: selectedMySQLDBUsername, isPrimary: selectedMySQLDBIsPrimary}).unwrap()
    .then(() => dispatch(updateNotifications({notifications: [{notif: 'succesful mysql db info upsert', status: 'success'}]})))
    .catch((error) => dispatch(updateNotifications({notifications: [{notif: 'Error upserting db info: ' + JSON.stringify(error), status: 'error'}]})));
    
    setMySQLDBInfoModalOpen(false);
    resetSelectedMySQLDB();
  };

  const handleDeleteMySQLDBInfo = (db: any) => {
    deleteMySQLDBInfo({host: db.host}).unwrap()
    .then(() => dispatch(updateNotifications({notifications: [{notif: 'succesful mysql db info delete', status: 'success'}]})))
    .catch((error) => dispatch(updateNotifications({notifications: [{notif: 'Error deleting db info: ' + JSON.stringify(error), status: 'error'}]})));
  };

  /* Dynamic App config queries and mutations
  /
  /
  */
  const [updateAppConfig] = useUpdateDynamicAppConfigMutation();
  const handleUpdateAppConfig = (config: {id: number, data: any}, keyName: string) => {
    updateAppConfig(config).unwrap().then(() => {
      dispatch(updateNotifications({notifications: [{notif: `${keyName} updated successfully`, status: 'success'}]}));
    }).catch((error) => {
      dispatch(updateNotifications({notifications: [{notif: 'Error updating dynamic app config: ' + JSON.stringify(error), status: 'error'}]}));
    });
  };

  const [enableK8sGlobalReadOnly, setEnableK8sGlobalReadOnly] = React.useState<boolean>(false);
  const handleUpdateEnableK8sGlobalReadOnly = (val: boolean) => {
    setEnableK8sGlobalReadOnly(val);
    const updatedAppConfig: IAppConfig = {id: appConfig.id, data: {...appConfig.data, enableK8sGlobalReadOnly: val}};
    handleUpdateAppConfig(updatedAppConfig, 'enableK8sGlobalReadOnly');
  };

  const [defaultReplicaScaleLimit, setDefaultReplicaScaleLimit] = React.useState<number>(0);
  const handleUpdateDefaultReplicaScaleLimit = (val: number) => {
    setDefaultReplicaScaleLimit(val);
    const updatedAppConfig: IAppConfig = {id: appConfig.id, data: {...appConfig.data, defaultReplicaScaleLimit: val}};
    handleUpdateAppConfig(updatedAppConfig, 'defaultReplicaScaleLimit');
  };
  
  const [k8sClusterNamespaceModalOpen, setK8sClusterNamespaceModalOpen] = React.useState(false);
  const [k8sNamespaceName, setK8sNamespaceName] = React.useState('');
  const handleK8sClusterNamespaceModalOpen = (val: boolean) => {
    setK8sClusterNamespaceModalOpen(val);
  };

  const handleRemovek8sClusterNamespace = (namespace: string) => {
    const updatedAppConfig: IAppConfig = {id: appConfig.id, data: {
      ...appConfig.data, 
      k8sClusterNamespaces: appConfig.data?.k8sClusterNamespaces.filter((ns: string) => ns !== namespace)
    }};
    handleUpdateAppConfig(updatedAppConfig, 'k8sClusterNamespaces');
  };

  const handleAddK8sClusterNamespace = (namespace: string) => {
    if (namespace === '') {
      return;
    }
    const namespaces: string[] = [];
    if (appConfig.data?.k8sClusterNamespaces !== null) {
      namespaces.push(...appConfig.data?.k8sClusterNamespaces);
    }
    const updatedAppConfig: IAppConfig = {id: appConfig.id, data: {
      ...appConfig.data, 
      k8sClusterNamespaces: [...namespaces, namespace]
    }};
    handleUpdateAppConfig(updatedAppConfig, 'k8sClusterNamespaces');
  };

  const [k8sReplicaScaleLimitModalOpen, setK8sReplicaScaleLimitModalOpen] = React.useState(false);
  const [k8sReplicaScaleLimitLabel, setK8sReplicaScaleLimitLabel] = React.useState('');
  const [k8sReplicaScaleLimitValue, setK8sReplicaScaleLimitValue] = React.useState(0);
  const handleK8sReplicaScaleLimitModalOpen = (val: boolean) => {
    setK8sReplicaScaleLimitModalOpen(val);
  };

  const handleRemoveK8sReplicaScaleLimit = (label: string) => {
    const replicaScaleLimitMap: any = {};
    if (appConfig.data?.replicaScaleLimits !== null) {
      Object.assign(replicaScaleLimitMap, appConfig.data?.replicaScaleLimits);
    }
    delete replicaScaleLimitMap[label];

    const updatedAppConfig: IAppConfig = {id: appConfig.id, data: {
      ...appConfig.data,
      replicaScaleLimits: {...replicaScaleLimitMap}
    }};
    handleUpdateAppConfig(updatedAppConfig, 'k8sReplicaScaleLimits');
  };

  const handleAddK8sReplicaScaleLimit = (label: string, value: number) => {
    if (label === '') {
      return;
    }
    const replicaScaleLimitMap: any = {};
    if (appConfig.data?.replicaScaleLimits !== null) {
      Object.assign(replicaScaleLimitMap, appConfig.data?.replicaScaleLimits);
    }
    replicaScaleLimitMap[label] = value;

    const updatedAppConfig: IAppConfig = {id: appConfig.id, data: {
      ...appConfig.data,
      replicaScaleLimits: {...replicaScaleLimitMap}
    }};
    handleUpdateAppConfig(updatedAppConfig, 'k8sReplicaScaleLimits');
  };

  const [podExecPluginModalOpen, setPodExecPluginModalOpen] = React.useState(false);
  const [selectedPluginName, setSelectedPluginName] = React.useState('');
  const [selectedPluginContainerFilter, setSelectedPluginContainerFilter] = React.useState('');
  const [selectedPluginCommand, setSelectedPluginCommand] = React.useState('');
  const [selectedPluginLabelFilter, setSelectedPluginLabelFilter] = React.useState('');

  const handleUpsertPodExecPlugin = (name: string, container: string, command: string, labelFilter: string) => {
    if (!name || !command) return;

    const plugins = [...(appConfig.data?.k8sPodExecPlugins || [])];
    const existingIndex = plugins.findIndex(p => p.name === name);
    
    if (existingIndex >= 0) {
      plugins[existingIndex] = { name, container, command, labelFilter };
    } else {
      plugins.push({ name, container, command, labelFilter });
    }

    const updatedAppConfig: IAppConfig = {
      id: appConfig.id,
      data: {
        ...appConfig.data,
        k8sPodExecPlugins: plugins
      }
    };
    
    handleUpdateAppConfig(updatedAppConfig, 'podExecPlugins');
    setPodExecPluginModalOpen(false);
    setSelectedPluginName('');
    setSelectedPluginContainerFilter('');
    setSelectedPluginCommand('');
    setSelectedPluginLabelFilter('');
  };

  const handleRemovePodExecPlugin = (pluginName: string) => {
    const updatedPlugins = appConfig.data?.k8sPodExecPlugins.filter(
      (plugin: any) => plugin.name !== pluginName
    );

    const updatedAppConfig: IAppConfig = {
      id: appConfig.id,
      data: {
        ...appConfig.data,
        k8sPodExecPlugins: updatedPlugins
      }
    };

    handleUpdateAppConfig(updatedAppConfig, 'podExecPlugins');
  };
  const [clusterName, setClusterName] = React.useState('');

  useEffect(() => {
    setEnableK8sGlobalReadOnly(appConfig?.data?.enableK8sGlobalReadOnly);
    setDefaultReplicaScaleLimit(appConfig?.data?.defaultReplicaScaleLimit);
    setClusterName(appConfig?.data?.k8sClusterName);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [appConfig]);
  

  return (
    <>
      <Heading style={{marginBottom: '20px'}}>General Settings</Heading>
      <Tile>
        <Accordion size='lg' className='accordian-content-override'>
          <AccordionItem title={KubernetesSettingsTile()} className="accordian-border-0" open>
            <StructuredListWrapper>
              <StructuredListHead>
                <StructuredListRow head>
                  <StructuredListCell head></StructuredListCell>
                  <StructuredListCell head></StructuredListCell>
                </StructuredListRow>
              </StructuredListHead>
              <StructuredListBody>
                <StructuredListRow key='appConfig-clusterName'>
                  <StructuredListCell>
                    <strong>Cluster Name:</strong> <br/>
                    The name of the k8s cluster that khub is monitoring. This name appears in the application header. 
                    It should help identify the environment that the cluster is in.<br/><br/>
                    NOTE: This value is used as part of the cache key for the k8s resources. 
                    Changing it will result in new objects being created in your cache. 
                  </StructuredListCell>
                  <StructuredListCell>
                    <TextInput 
                      data-modal-primary-focus 
                      onChange={(e: any) => {setClusterName(e.target.value);}}
                      onKeyDown={(e: any) => {
                        if (e.key === 'Enter') {
                          const updatedAppConfig: IAppConfig = {id: appConfig.id, data: {...appConfig.data, k8sClusterName: clusterName}};
                          handleUpdateAppConfig(updatedAppConfig, 'clusterName');
                        }
                      }}
                      id="clusterName" 
                      labelText="Cluster Name" 
                      placeholder="e.g. app-core-staging" 
                      style={{marginBottom: '1rem'}}
                      value={clusterName} 
                    />
                    <em style={{fontSize: 11, float: 'right', color: 'light-gray'}}>press enter to save</em>
                  </StructuredListCell>
                </StructuredListRow>
                <StructuredListRow key='appConfig-enableGlobalReadOnly'>
                  <StructuredListCell>
                    <strong>Enable global read only:</strong> <br/>
                    Toggle this setting to allow all users to view k8s resources but not modify them, regardless of their permissions.
                  </StructuredListCell>
                  <StructuredListCell>
                    <Toggle 
                      labelText="Enable global read only of K8s resources" 
                      labelA="disabled" 
                      labelB="enabled" 
                      id="toggle-enableGlobalReadOnly"
                      toggled={enableK8sGlobalReadOnly}
                      onToggle={(val: any) => {handleUpdateEnableK8sGlobalReadOnly(val);}}
                    />
                  </StructuredListCell>
                </StructuredListRow>
                <StructuredListRow key='appConfig-K8sNamespaces'>
                  <StructuredListCell>
                    <strong>Kubernetes namespaces:</strong> <br/>
                    Configure the namespaces that are actively monitored by khub.
                  </StructuredListCell>
                  <StructuredListCell>
                    {appConfig?.data?.k8sClusterNamespaces !== null && appConfig?.data?.k8sClusterNamespaces.map((ns: string) => {
                      return (
                        <Tag 
                          renderIcon={CgNametag} 
                          key={ns} 
                          filter 
                          type="purple" 
                          title="Remove"
                          onClose={() => handleRemovek8sClusterNamespace(ns)} >
                          {ns}
                        </Tag>
                      );
                    })}
                    {(appConfig?.data?.k8sClusterNamespaces === null || 
                      appConfig?.data?.k8sClusterNamespaces.length === 0) && <Tag key={'allNS'} type="purple" title="">*</Tag>}
                    <div style={{float: 'right'}}>
                      <Tooltip label='add a namespace'>
                        <ClickableTile style={{display: 'flex', justifyContent: 'center', width: '20px'}} 
                          onClick={() => handleK8sClusterNamespaceModalOpen(true)}>
                            <PiPlusBold size={30}/>
                        </ClickableTile>
                      </Tooltip>
                    </div>
                  </StructuredListCell>
                </StructuredListRow>
                <StructuredListRow key='appConfig-k8sReplicaScaleLimits'>
                  <StructuredListCell>
                    <strong>Kubernetes replica scale limits:</strong> <br/>
                    Configure the scale limits of deployments based on label (such as an app-class label corresponding to a deployment).
                  </StructuredListCell>
                  <StructuredListCell>
                    {appConfig?.data?.replicaScaleLimits && Object.entries(appConfig?.data?.replicaScaleLimits).map(([key, value]: any) => {
                      return (
                        <Tag 
                          renderIcon={PiCaretCircleUpDownLight} 
                          key={key} filter 
                          type="green" 
                          title="Remove"
                          onClose={() => {
                            handleRemoveK8sReplicaScaleLimit(key);
                          }}>
                          {key}: {value}
                        </Tag>
                      );
                    })}
                    <div style={{float: 'right'}}>
                      <Tooltip label='add scale limit'>
                        <ClickableTile style={{display: 'flex', justifyContent: 'center', width: '20px'}} onClick={() => handleK8sReplicaScaleLimitModalOpen(true)}>
                            <PiPlusBold size={30}/>
                        </ClickableTile>
                      </Tooltip>
                    </div>
                  </StructuredListCell>
                </StructuredListRow>
                <StructuredListRow key='appConfig-defaultReplicaScaleLimit'>
                  <StructuredListCell>
                    <strong>Default replica scale limit:</strong> <br/>
                    Set the default replica scale limit for deployments. This is the scale limit that will be used for any deployments 
                    without a specific scale limit.
                  </StructuredListCell>
                  <StructuredListCell>
                    <NumberInput id="defaultReplicaScaleLimit" 
                      max={99999}
                      value={defaultReplicaScaleLimit} 
                      label="Default Replica Scale Limit" 
                      invalidText="Replica scale limit must be between 0 and 99999"
                      onKeyDown={(e: any) => {
                        if (e.key === 'Enter') {
                          handleUpdateDefaultReplicaScaleLimit(defaultReplicaScaleLimit);
                        }
                      }}
                      onChange={(e: any, data: any) => setDefaultReplicaScaleLimit(data.value)}
                      hideSteppers
                      />
                      <em style={{fontSize: 11, float: 'right', color: 'light-gray'}}>press enter to save</em>
                  </StructuredListCell>
                </StructuredListRow>
              </StructuredListBody>
            </StructuredListWrapper>
          </AccordionItem>
          <AccordionItem title="Pod Exec Plugins" className="accordian-border-0" open>
            <StructuredListWrapper selection>
              <StructuredListHead>
                <StructuredListRow head>
                  <StructuredListCell head>Plugin Name</StructuredListCell>
                  <StructuredListCell head>Container Filter</StructuredListCell>
                  <StructuredListCell head>Command</StructuredListCell>
                  <StructuredListCell head>Actions</StructuredListCell>
                </StructuredListRow>
              </StructuredListHead>
              <StructuredListBody>
                {appConfig?.data?.k8sPodExecPlugins && appConfig?.data?.k8sPodExecPlugins.map((plugin: any) => (
                  <StructuredListRow key={plugin.name}>
                    <StructuredListCell noWrap>{plugin.name}</StructuredListCell>
                    <StructuredListCell noWrap>{plugin.container}</StructuredListCell>
                    <StructuredListCell>{plugin.command}</StructuredListCell>
                    <StructuredListCell>
                      <ButtonSet stacked>
                        <Button 
                          tooltipPosition='right' 
                          size='sm' 
                          onClick={() => {
                            setSelectedPluginName(plugin.name);
                            setSelectedPluginContainerFilter(plugin.container);
                            setSelectedPluginCommand(plugin.command);
                            setPodExecPluginModalOpen(true);
                            setSelectedPluginLabelFilter(plugin.labelFilter);
                          }} 
                          iconDescription="Edit" 
                          style={{border: 'none'}} 
                          hasIconOnly 
                          renderIcon={TbEdit} 
                          kind="tertiary"
                        />
                        <Button 
                          tooltipPosition='right' 
                          size='sm' 
                          onClick={() => handleRemovePodExecPlugin(plugin.name)} 
                          iconDescription="Remove" 
                          style={{border: 'none'}} 
                          hasIconOnly 
                          renderIcon={TbTrash} 
                          kind="danger--tertiary"
                        />
                      </ButtonSet>
                    </StructuredListCell>
                  </StructuredListRow>
                ))}
              </StructuredListBody>
            </StructuredListWrapper>
            <ClickableTile style={{display: 'flex', justifyContent: 'center'}} onClick={() => setPodExecPluginModalOpen(true)}>
              <Tooltip label='add plugin'>
                <PiPlusBold size={30}/>
              </Tooltip>
            </ClickableTile>
          </AccordionItem>
        </Accordion>
      </Tile>

      <Tile style={{marginTop: '10px'}}>
        <Accordion size='lg' className='accordian-content-override'>
          <AccordionItem title={DBSettingsTitle()} className="accordian-border-0" open>
            <StructuredListWrapper selection>
              <StructuredListHead>
                <StructuredListRow head>
                  <StructuredListCell head>Name</StructuredListCell>
                  <StructuredListCell head>Host</StructuredListCell>
                  <StructuredListCell head>Username</StructuredListCell>
                  <StructuredListCell head>Primary</StructuredListCell>
                  <StructuredListCell head>Actions</StructuredListCell>
                </StructuredListRow>
              </StructuredListHead>
              <StructuredListBody>
                {mysqlDBCatalog.map((db: any) => (
                  <StructuredListRow key={db.host}>
                    <StructuredListCell noWrap>{db.shortName}</StructuredListCell>
                    <StructuredListCell>{db.host}:{db.port}</StructuredListCell>
                    <StructuredListCell>{db.username}</StructuredListCell>
                    <StructuredListCell>{db.isPrimary === true ? <CheckmarkFilled color="green" /> : <Misuse color="coral"/>}</StructuredListCell>
                    <StructuredListCell>
                      <ButtonSet stacked>
                        <Button tooltipPosition='right' size='sm' onClick={() => handleMySQLInfoUpsertModalOpen(db)} iconDescription="Edit" style={{border: 'none'}} hasIconOnly renderIcon={TbEdit} kind="tertiary"/>
                        <Button tooltipPosition='right' size='sm' onClick={() => handleDeleteMySQLDBInfo(db)} iconDescription="Remove" style={{border: 'none'}} hasIconOnly renderIcon={TbTrash} kind="danger--tertiary"/>
                      </ButtonSet>
                    </StructuredListCell>
                  </StructuredListRow>
                ))}
              </StructuredListBody>
            </StructuredListWrapper>
            <ClickableTile style={{display: 'flex', justifyContent: 'center'}} onClick={() => handleMySQLInfoUpsertModalOpen({port: 3306})}>
              <Tooltip label='add host'>
                <PiPlusBold size={30}/>
              </Tooltip>
            </ClickableTile>
          </AccordionItem>
        </Accordion>
      </Tile>
      
      {/* MySQL DB Info Upsert Modal */} 
      <ComposedModal open={mySQLDBInfoModalOpen} onClose={() => {resetSelectedMySQLDB();}}>
        <ModalHeader label="MySQL DB Catalog" title="Add a new db to the catalog" />
        <ModalBody>
          <p style={{marginBottom: '1rem'}}>
            MySQL DBs in this catalog will be represented in the MySQL replication topology plugin.
          </p>
          <TextInput 
            data-modal-primary-focus 
            onChange={(e: any) => {setSelectedMySQLDBHost(e.target.value);}}
            id="dbhost" 
            labelText="Host" 
            placeholder="e.g. db-core-007.platform-databases.staging.smar.cloud" 
            style={{marginBottom: '1rem'}}
            value={selectedMySQLDBHost} 
          />
          <TextInput 
            data-modal-primary-focus 
            onChange={(e: any) => {setSelectedMySQLDBShortname(e.target.value);}}
            id="dbshortname" 
            labelText="Shortname" 
            placeholder="e.g. db-core-007" 
            style={{marginBottom: '1rem'}}
            value={selectedMySQLDBShortname} 
          />
          <TextInput 
            data-modal-primary-focus 
            onChange={(e: any) => {setSelectedMySQLDBUsername(e.target.value);}}
            id="dbusername" 
            labelText="Username" 
            placeholder="e.g. khub" 
            style={{marginBottom: '1rem'}}
            value={selectedMySQLDBUsername} 
          />
          <NumberInput id="db-port" 
            max={99999} 
            value={selectedMySQLDBPort} 
            label="DB Port" 
            invalidText="Port must be between 0 and 99999" 
            onChange={(e: any, data: any) => setSelectedMySQLDBPort(data.value)}
            hideSteppers
            />
          <Toggle 
            labelText="Primary database" 
            labelA="no" 
            labelB="yes" 
            id="primary-toggle" 
            toggled={selectedMySQLDBIsPrimary}
            onToggle={(val: any) => {setSelectedMySQLDBIsPrimary(val);}}/>
          <ButtonSet style={{marginTop: '20px'}}>
            <Button kind="primary" onClick={() => handleUpsertMySQLDBInfo()}>
              Submit
            </Button>
            <Button kind="secondary" onClick={() => {
              resetSelectedMySQLDB();
              setMySQLDBInfoModalOpen(false);
            }}>
              Cancel
            </Button>
          </ButtonSet>
        </ModalBody>
      </ComposedModal>

      {/* K8s Namespace Modal */} 
      <ComposedModal open={k8sClusterNamespaceModalOpen} onClose={() => {handleK8sClusterNamespaceModalOpen(false); setK8sNamespaceName('');}}>
        <ModalHeader label="Kubernetes Namespace" title="Add a new namespace to the khub data sink" />
        <ModalBody>
          <TextInput 
            data-modal-primary-focus
            onChange={(e: any) => {setK8sNamespaceName(e.target.value);}} 
            onKeyDown={(e: any) => {
              if (e.key === 'Enter') {
                handleAddK8sClusterNamespace(k8sNamespaceName);
                setK8sNamespaceName('');
                setK8sClusterNamespaceModalOpen(false);
              }
            }}
            id="k8sNamespace" 
            labelText="Namespace name" 
            placeholder="e.g. default" 
            style={{marginBottom: '1rem'}}
            value={k8sNamespaceName} 
          />
        </ModalBody>
      </ComposedModal>

      {/* K8s Replica Scale Limits Modal */} 
      <ComposedModal open={k8sReplicaScaleLimitModalOpen} onClose={() => {
          handleK8sReplicaScaleLimitModalOpen(false);
          setK8sReplicaScaleLimitLabel('');
          setK8sReplicaScaleLimitValue(0);
        }}>
        <ModalHeader label="Kubernetes Replica Scale Limit" title="Add a new khub scale limit for a deployment." />
        <ModalBody>
          <TextInput 
            data-modal-primary-focus
            onChange={(e: any) => {setK8sReplicaScaleLimitLabel(e.target.value);}} 
            id="k8sReplicaLabel" 
            labelText="Label filter that corresponds with the deployment:" 
            placeholder="e.g. s2rapi or gridreader" 
            style={{marginBottom: '1rem'}}
            value={k8sReplicaScaleLimitLabel} 
          />
          <NumberInput id="k8sReplicaScaleLimitValue" 
            max={99999}
            value={k8sReplicaScaleLimitValue} 
            label="Replica scale limit:" 
            invalidText="Replica scale limit must be between 0 and 99999"
            onChange={(e: any, data: any) => setK8sReplicaScaleLimitValue(data.value)}
            hideSteppers
            />
            <ButtonSet style={{marginTop: '20px'}}>
              <Button kind="primary" onClick={() => {
                handleAddK8sReplicaScaleLimit(k8sReplicaScaleLimitLabel, k8sReplicaScaleLimitValue);
                setK8sReplicaScaleLimitLabel('');
                setK8sReplicaScaleLimitValue(0);
                setK8sReplicaScaleLimitModalOpen(false);
              }}>
                Submit
              </Button>
              <Button kind="secondary" onClick={() => {
                resetSelectedMySQLDB();
                setMySQLDBInfoModalOpen(false);
              }}>
                Cancel
            </Button>
          </ButtonSet>
        </ModalBody>
      </ComposedModal>

      {/* Pod Exec Plugin Modal */}
      <ComposedModal open={podExecPluginModalOpen} onClose={() => {
        setPodExecPluginModalOpen(false);
        setSelectedPluginName('');
        setSelectedPluginContainerFilter('');
        setSelectedPluginCommand('');
        setSelectedPluginLabelFilter('');
      }}>
        <ModalHeader label="Pod Exec Plugin" title="Add a new pod exec plugin" />
        <ModalBody>
          <TextInput 
            data-modal-primary-focus
            onChange={(e: any) => {setSelectedPluginName(e.target.value);}}
            id="pluginName"
            labelText="Plugin name"
            placeholder="e.g. list-files"
            style={{marginBottom: '1rem'}}
            value={selectedPluginName}
          />
          <TextInput 
            data-modal-primary-focus
            onChange={(e: any) => {setSelectedPluginContainerFilter(e.target.value);}}
            id="pluginContainerFilter"
            labelText="Plugin container filter"
            placeholder="e.g. busybox"
            style={{marginBottom: '1rem'}}
            value={selectedPluginContainerFilter}
          />
          <TextInput 
            data-modal-primary-focus
            onChange={(e: any) => {setSelectedPluginCommand(e.target.value);}}
            id="pluginCommand"
            labelText="Command"
            placeholder="e.g. ls -l"
            value={selectedPluginCommand}
          />
          <TextInput 
            data-modal-primary-focus
            onChange={(e: any) => {setSelectedPluginLabelFilter(e.target.value);}}
            id="pluginLabelFilter"
            labelText="Label filter"
            placeholder="e.g. busybox"
            value={selectedPluginLabelFilter}
          />
          <ButtonSet style={{marginTop: '20px'}}>
            <Button kind="primary" onClick={() => {
              handleUpsertPodExecPlugin(selectedPluginName, selectedPluginContainerFilter, selectedPluginCommand, selectedPluginLabelFilter);
            }}>
              Submit
            </Button>
            <Button kind="secondary" onClick={() => {
              setPodExecPluginModalOpen(false);
              setSelectedPluginName('');
              setSelectedPluginContainerFilter('');
              setSelectedPluginCommand('');
              setSelectedPluginLabelFilter('');
            }}>
              Cancel
            </Button>
          </ButtonSet>  
        </ModalBody>
      </ComposedModal>
    </>
  );
};