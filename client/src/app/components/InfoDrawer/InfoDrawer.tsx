import React, { useEffect, useMemo, useRef } from 'react';
import type { CSSProperties } from 'react';
import './InfoDrawer.scss';
import { useSelector } from 'react-redux';
import { RootState, useAppDispatch } from '../../store';
import { Button, Table, TableBody, TableCell, TableContainer, TableRow, TableToolbar, TableToolbarContent, Tag } from '@carbon/react';
import { Recycle, Report } from "@carbon/icons-react";

import CodeMirror from '@uiw/react-codemirror';
import { yaml as yamlint } from '@codemirror/lang-yaml';
import yaml from 'js-yaml';
import { useDeletePodMutation, useExecPluginMutation, useRolloutRestartMutation } from '../../../service/khub';
import { updateNotifications } from '../../../service/notifications';
import { IPodExecPlugin } from '../../../service/types/AppConfig';

type InfoDrawerProps = {
  open: boolean;
  onClose?: (event: any) => void;
  direction: 'left' | 'right' | 'top' | 'bottom';
  lockBackgroundScroll?: boolean;
  children?: React.ReactNode;
  duration?: number;
  overlayOpacity?: number;
  overlayColor?: string;
  enableOverlay?: boolean;
  style?: React.CSSProperties;
  zIndex?: number;
  size?: number | string;
  className?: string;
  customIdSuffix?: string;
  overlayClassName?: string;
  podExecPlugins?: IPodExecPlugin[]
};

const getDirectionStyle = (dir: string, size?: number | string): Record<string, never> | React.CSSProperties => {
  switch (dir) {
    case 'left':
      return {
        top: 0,
        left: 0,
        transform: 'translate3d(-100%, 0, 0)',
        width: size,
        height: '100vh'
      };
    case 'right':
      return {
        top: 0,
        right: 0,
        transform: 'translate3d(100%, 0, 0)',
        width: size,
        height: '100vh'
      };
    case 'bottom':
      return {
        left: 0,
        right: 0,
        bottom: 0,
        transform: 'translate3d(0, 100%, 0)',
        width: '100%',
        height: size
      };
    case 'top':
      return {
        left: 0,
        right: 0,
        top: 0,
        transform: 'translate3d(0, -100%, 0)',
        width: '100%',
        height: size
      };

    default:
      return {};
  }
};

export const InfoDrawer: React.FC<InfoDrawerProps> = (props) => {
  
  const appTheme = useSelector((state: RootState) => state.appTheme);
  const resourceDrawer: any = useSelector((state: RootState) => state.treeMapResourceDrawer);
  const {
    open,
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    onClose = () => {},
    children,
    style,
    enableOverlay = true,
    overlayColor = '#000',
    overlayOpacity = 0.4,
    zIndex = 100,
    duration = 500,
    direction,
    size = '40%',
    customIdSuffix,
    lockBackgroundScroll = false,
    overlayClassName = ''
  } = props;

  const bodyRef = useRef<HTMLBodyElement | null>(null);

  useEffect(() => {
    const updatePageScroll = () => {
      bodyRef.current = window.document.querySelector('body');

      if (bodyRef.current && lockBackgroundScroll) {
        if (open) {
          bodyRef.current.style.overflow = 'hidden';
        } else {
          bodyRef.current.style.overflow = '';
        }
      }
    };

    updatePageScroll();
  }, [open, lockBackgroundScroll]);

  const idSuffix = useMemo(() => {
    return customIdSuffix || (Math.random() + 1).toString(36).substring(7);
  }, [customIdSuffix]);

  const overlayStyles: CSSProperties = {
    backgroundColor: `${overlayColor}`,
    opacity: `${overlayOpacity}`,
    zIndex: zIndex
  };

  const drawerStyles: CSSProperties = {
    zIndex: zIndex + 1,
    transitionDuration: `${duration}ms`,
    ...getDirectionStyle(direction, size),
    ...style,
    backgroundColor: appTheme.theme === 'dark' ? 'var(--cds-background)' : '#fff'
  };

  const dispatch = useAppDispatch();
  const [rolloutRestart] = useRolloutRestartMutation();
  const [execPlugin] = useExecPluginMutation();

  const handleRolloutRestart = (res: any) => {
    rolloutRestart({name: res.name, namespace: res.namespace, kind: res.kind, labels: res.metadata.labels}).unwrap()
      .then((payload) => dispatch(updateNotifications({notifications: [{notif: res.name + ' restart initiated: ' + payload, status: 'success'}]})))
      .catch((error) => dispatch(updateNotifications({notifications: [{notif: 'Error restarting ' + res.name + ' ' + JSON.stringify(error), status: 'error'}]})));
  };

  const handleExecPlugin = (res: any) => {
    dispatch(updateNotifications({notifications: [{notif: res.name + ' exec plugin initiated ' , status: 'info'}]}));
    execPlugin({name: res.name, namespace: res.namespace, kind: res.kind, container: res.container, command: res.command, pluginName: res.pluginName}).unwrap()
      .then((payload) => dispatch(updateNotifications({notifications: [{notif: res.name + ' pod exec success: ' + payload, status: 'success'}]})))
      .catch((error) => dispatch(updateNotifications({notifications: [{notif: 'Error running pod exec plugin ' + res.name + ' ' + JSON.stringify(error), status: 'error'}]})));
  };
  
  const handleYamlRender = (data: any) => {
    if (data === null) {
      return '';
    }
    const yamlData = JSON.parse(JSON.stringify(data));
    if (yamlData.resourceData && yamlData.resourceData?.metadata?.managedFields) {
      yamlData.resourceData.metadata.managedFields = undefined;
    }
    return yaml.dump(yamlData.resourceData);
  };

  
  const [deletePod] = useDeletePodMutation();

  const handleDeletePod = (podName: string, namespace: string) => {
    deletePod({podName: podName, namespace: namespace}).unwrap()
    .then(() => dispatch(updateNotifications({notifications: [{notif: podName + ' delete initiated', status: 'success'}]})))
    .catch((error) => dispatch(updateNotifications({notifications: [{notif: 'Error deleting ' + podName + ' ' + JSON.stringify(error), status: 'error'}]})));
  };

  const getExecPlugins = (resource: any): any[] => {

    const plugins: any[] = [];
    if (resource?.resourceType === 'pod') {
      console.log('resource: ', JSON.stringify(props.podExecPlugins));
      props.podExecPlugins?.forEach((plugin: IPodExecPlugin) => {
        console.log('plugin: ', plugin.name);
        plugin.labelFilter.split(',').forEach((labelValue: string) => {
          console.log('labelValue: ', labelValue);
          const labelKeys = Object.keys(resource?.resourceData?.metadata?.labels); 
          labelKeys.forEach((labelKey: string) => {
            console.log('labelKey: ', labelKey);
            if (resource?.resourceData?.metadata?.labels[labelKey] === labelValue) {
              if (resource?.resourceData?.spec?.containers !== undefined) {
                resource?.resourceData?.spec?.containers.forEach((container: any) => {
                  if (container.name !== undefined && container.name === plugin.container) {
                    plugins.push(plugin);
                  }
                });
              }
            }
          });
        });
      });
    }
    return plugins;
  };

  return (
    <div id={'InfoDrawer' + idSuffix} className="InfoDrawer">
      <input type="checkbox" id={'InfoDrawer__checkbox' + idSuffix} className="InfoDrawer__checkbox" onChange={onClose} checked={open} />
      <nav role="navigation" id={'InfoDrawer__container' + idSuffix} style={drawerStyles} className={'InfoDrawer__container'}>
        {children}
        <div>
          <TableContainer title="Resource info">
            <TableToolbar>
              <TableToolbarContent>
                {(
                resourceDrawer.data?.resourceType === 'deployment' || 
                resourceDrawer.data?.resourceType === 'statefulset' || 
                resourceDrawer.data?.resourceType === 'daemonset' ) && 
                  <Button renderIcon={Recycle} kind="tertiary" onClick={() => handleRolloutRestart({name: resourceDrawer.data?.resourceData?.metadata?.name, namespace: resourceDrawer.data?.resourceData?.metadata?.namespace, kind: resourceDrawer.data.resourceType})}>Restart {resourceDrawer.data?.resourceType}</Button>
                }
                {(
                resourceDrawer.data?.resourceType === 'pod') &&
                  <div>
                    <Button 
                      renderIcon={Recycle} 
                      kind="danger--tertiary" 
                      onClick={
                        () => handleDeletePod(resourceDrawer.data?.resourceData?.metadata?.name, resourceDrawer.data?.resourceData?.metadata?.namespace)
                      }
                      >Delete {resourceDrawer.data?.resourceType}</Button>
                    {getExecPlugins(resourceDrawer.data).map((plugin: any) => {
                        return <Button key={plugin.name} renderIcon={Report} kind="tertiary" style={{marginLeft: '5px'}} onClick={() => {
                          handleExecPlugin({
                              name: resourceDrawer.data?.resourceData?.metadata?.name, 
                              namespace: resourceDrawer.data?.resourceData?.metadata?.namespace, 
                              kind: resourceDrawer.data.resourceType,
                              container: plugin.container,
                              command: plugin.command,
                              pluginName: plugin.name
                            });
                        }}>{plugin.name}</Button>;
                      })
                    }
                  </div>
                }
                {(resourceDrawer.data?.resourceType === 'node') &&
                  <div>
                    <Button disabled kind="tertiary" onClick={() => alert("yo duuude")}>Drain {resourceDrawer.data?.resourceType}</Button>
                    <Button disabled kind="tertiary" onClick={() => alert("yo duuude")}>Uncordon {resourceDrawer.data?.resourceType}</Button>
                  </div>
                  
                }
              </TableToolbarContent>
            </TableToolbar>
            <div id='infoDrawerTable'>
              <Table aria-label="sample table">
                {resourceDrawer.data?.resourceData !== null && (
                  <TableBody>
                    <TableRow key="resource-data">
                      <TableCell key="name">
                        <strong>name</strong>
                      </TableCell>
                      <TableCell key="nameData">{resourceDrawer.data?.resourceData?.metadata?.name}</TableCell>
                    </TableRow>

                    <TableRow key="resource-data2">
                      <TableCell key="namespace">
                        <strong>namespace</strong>
                      </TableCell>
                      <TableCell key="namespaceData">{resourceDrawer.data?.resourceData?.metadata?.namespace}</TableCell>
                    </TableRow>

                    <TableRow key="resource-data3">
                      <TableCell key="resourceVersion">
                        <strong>resourceVersion</strong>
                      </TableCell>
                      <TableCell key="resourceVersionData">{resourceDrawer.data?.resourceData?.metadata?.resourceVersion}</TableCell>
                    </TableRow>

                    {resourceDrawer.data?.resourceType === 'pod' && (
                      <TableRow key="resource-data4">
                      <TableCell key="containers">
                        <strong>containers</strong>
                      </TableCell>
                      <TableCell key="containersData">
                        {resourceDrawer.data?.resourceData?.spec?.containers?.map((container: any) => {
                          return <Tag key={container.name} type="cool-gray" title="pod container">
                                    {container.name}
                                </Tag>;
                        })}
                      </TableCell>
                    </TableRow>
                    )}
                  </TableBody>
                )}
              </Table>
            </div>
            
          </TableContainer>
          <CodeMirror
                placeholder={'resource yaml loading. . .'}
                value={handleYamlRender(resourceDrawer.data)}
                extensions={[yamlint()]}
                theme={'dark'}
                editable={false}
                style={{
                  minHeight: 550,
                  maxHeight: 800,
                  maxWidth: 1150,
                  overflowY: 'scroll'
                }}
              />
        </div>
      </nav>
      {enableOverlay && (
        <label
          htmlFor={'InfoDrawer__checkbox' + idSuffix}
          id={'InfoDrawer__overlay' + idSuffix}
          className={'InfoDrawer__overlay ' + overlayClassName}
          style={overlayStyles}
        />
      )}
    </div>
  );
};
