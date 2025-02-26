import React, { useEffect } from 'react';
import {
  Header,
  HeaderContainer,
  HeaderName,
  // HeaderNavigation,
  HeaderMenuButton,
  HeaderGlobalBar,
  HeaderGlobalAction,
  SkipToContent,
  SideNav,
  SideNavItems,
  SideNavLink,
  SideNavMenu,

  TableToolbarMenu,
  TableToolbarAction,
  PopoverContent,
  Popover,
  InlineNotification,
} from '@carbon/react';
import { CiMicrochip } from "react-icons/ci";

import { UserAvatar, BareMetalServer, DocumentEpdf, LogoKubernetes, KubernetesPod, KubernetesIpAddress, DocumentMultiple_01, WarningDiamond, ColorPalette, TreeView, CheckboxCheckedFilled, RuleLocked } from '@carbon/icons-react';
import { MdCheckBoxOutlineBlank } from "react-icons/md";
import { NavLink } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState, useAppDispatch } from '../store';
import { updateAppTheme } from '../../service/themeState';
import { dismissNotification } from '../../service/notifications';
import { useGetClusterNameQuery, useUpdateUserThemeMutation } from '../../service/khub';

export const Shell = ({userPreferredTheme, userInfo}: any) => {
  const appTheme = useSelector((state: RootState) => state.appTheme);

  const [updateUserTheme] = useUpdateUserThemeMutation();
  const dispatch = useAppDispatch();
  
  const setAppTheme = (theme: string) => {
    setLocalTheme(theme);
    const darkMode = theme === 'dark' ? true : false; 
    updateUserTheme({name: userInfo.name, darkMode: darkMode});
  };

  const setLocalTheme = (theme: string) => {
    dispatch(updateAppTheme({theme}));
  };

  useEffect(() => {
    setLocalTheme(userPreferredTheme);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const notificationsState = useSelector((state: RootState) => state.notifications);
  const [showNotifs, updateShowNotifs] = React.useState(false);
  const handleOpenNotifs = () => {
    updateShowNotifs(!showNotifs);
  };

  const dismissNotif = (arg: string) => {    
    dispatch(dismissNotification(arg));
  };

  const {data: clusterName} = useGetClusterNameQuery({});

  const logoutRedirect = () => {
    const logoutUrl = 
      process.env.REACT_APP_API_URL !== undefined && 
      process.env.REACT_APP_API_URL !== '' && 
      process.env.REACT_APP_API_URL !== null ? `${process.env.REACT_APP_API_URL}/logout` : `${window.location.origin.toString()}/logout`;
    window.location.replace(logoutUrl);
  };

  const userIsAdmin = useSelector((state: RootState) => state.userIsAdminState);

  return (
    <HeaderContainer
      render={({ isSideNavExpanded, onClickSideNavExpand }) => (
        <Header aria-label="Carbon Tutorial">
          <SkipToContent />
          <HeaderMenuButton aria-label={isSideNavExpanded ? 'Close menu' : 'Open menu'} onClick={onClickSideNavExpand} isActive={isSideNavExpanded} aria-expanded={isSideNavExpanded} />
          <CiMicrochip color='#0f62fe' size={20} style={{"marginLeft": "15px", "marginRight": "-15px"}}/>
          <HeaderName href="/" prefix="khub">
            [{clusterName}]
          </HeaderName>
          <SideNav aria-label="Side navigation" expanded={isSideNavExpanded} onSideNavBlur={onClickSideNavExpand} href="#main-content">
              <SideNavItems>
                <SideNavLink as={NavLink} to="/" end renderIcon={LogoKubernetes}>Cluster Overview</SideNavLink>
                <SideNavMenu defaultExpanded renderIcon={KubernetesPod} title="Resources">
                  <SideNavLink as={NavLink} to="/pods" end>Pods</SideNavLink>
                  <SideNavLink as={NavLink} to="/deployments" end>Deployments</SideNavLink>
                  <SideNavLink as={NavLink} to="/Daemonsets" end>Daemonsets</SideNavLink>
                  <SideNavLink as={NavLink} to="/statefulsets" end>Statefulsets</SideNavLink>
                  <SideNavLink as={NavLink} to="/jobs" end>Jobs</SideNavLink>
                  <SideNavLink as={NavLink} to="/cronjobs" end>CronJobs</SideNavLink>
                </SideNavMenu>
                <SideNavMenu defaultExpanded renderIcon={KubernetesIpAddress} title="Network">
                  <SideNavLink as={NavLink} to="/services" end>Services</SideNavLink>
                  <SideNavLink as={NavLink} to="/ingresses" end>Ingresses</SideNavLink>
                </SideNavMenu>
                <SideNavMenu defaultExpanded renderIcon={DocumentMultiple_01} title="Configuration">
                  <SideNavLink as={NavLink} to="/configmaps" end>ConfigMaps</SideNavLink>
                </SideNavMenu>
                <SideNavLink as={NavLink} to="/nodes" end renderIcon={BareMetalServer} >Nodes</SideNavLink>
                <SideNavMenu defaultExpanded renderIcon={TreeView} title="Data & Infrastructure">
                  <SideNavLink as={NavLink} to="/mysql-replication-topology" end>
                    MySQL Replication
                  </SideNavLink>
                </SideNavMenu>
                <SideNavLink as={NavLink} to="/reports" end renderIcon={DocumentEpdf} >Reports</SideNavLink>
                {userIsAdmin.isAdmin && <SideNavMenu defaultExpanded renderIcon={RuleLocked} title="Administration">
                  <SideNavLink as={NavLink} to="/general-settings" end>General</SideNavLink>
                  <SideNavLink as={NavLink} to="/access-control" end>Access Control</SideNavLink>
                </SideNavMenu>}
                
              </SideNavItems>
            </SideNav>
          <HeaderGlobalBar>
            <Popover open={showNotifs || notificationsState.notifications.length > 0} autoAlign> 
              <HeaderGlobalAction aria-label="Notifications" tooltipAlignment="center" onClick={() => handleOpenNotifs()}>
                <div style={{marginTop: '2px'}}>
                  <div>
                    {notificationsState.notifications.length > 0 &&
                    <div className='notif-tag'></div>
                    }
                    <WarningDiamond size={18} />
                  </div>
                </div>
              </HeaderGlobalAction>
              <PopoverContent className="p-3">
                      {notificationsState.notifications.map((notif: any) => {
                        return (<InlineNotification onCloseButtonClick={() => dismissNotif(notif.notif)} key={notif} title={notif.notif} kind={notif.status} role='log'/>);
                      })}
              </PopoverContent>
            </Popover>
            <HeaderGlobalAction aria-label="Theme Selector" tooltipAlignment={undefined}>
              <TableToolbarMenu iconDescription='theme selector' renderIcon={() => {return <ColorPalette/>;}}>
                  <TableToolbarAction onClick={() => setAppTheme('light')} className='color-theme-toolbar-actions'>
                  {appTheme.theme === 'light' && 
                      <div style={{display: 'inline-flex'}}>
                        <CheckboxCheckedFilled id='light' size={25}/>
                        <span style={{fontWeight: 'bold', color: '#0f62fe', marginTop: '3px', marginLeft: '5px'}}>Light</span>
                      </div>
                      
                    }
                    {appTheme.theme !== 'light' && 
                      <div style={{display: 'inline-flex'}}>
                        <MdCheckBoxOutlineBlank size={25}/>
                        <span style={{marginTop: '3px', marginLeft: '5px'}}> Light</span>
                      </div>
                    }
                  </TableToolbarAction>
                  <TableToolbarAction onClick={() => setAppTheme('dark')} className='color-theme-toolbar-actions'>
                    {appTheme.theme === 'dark' && 
                      <div style={{display: 'inline-flex'}}>
                        <CheckboxCheckedFilled id='light' size={25}/>
                        <span style={{fontWeight: 'bold', color: '#0f62fe', marginTop: '3px', marginLeft: '5px'}}>Dark</span>
                      </div>
                      
                    }
                    {appTheme.theme !== 'dark' && 
                      <div style={{display: 'inline-flex'}}>
                        <MdCheckBoxOutlineBlank size={25}/>
                        <span style={{marginTop: '3px', marginLeft: '5px'}}> Dark</span>
                      </div>
                    }
                  </TableToolbarAction>
              </TableToolbarMenu>
            </HeaderGlobalAction>
            <HeaderGlobalAction aria-label={userInfo.name} tooltipAlignment={undefined} >
              <TableToolbarMenu iconDescription='user info' renderIcon={() => {return <UserAvatar/>;}}>
                  <TableToolbarAction onClick={() => logoutRedirect()} className='color-theme-toolbar-actions'>
                      Refresh Session
                  </TableToolbarAction>
              </TableToolbarMenu>
            </HeaderGlobalAction>
            <span style={{fontSize: '12px', marginTop: '18px', marginRight: '10px', marginLeft: '5px'}}>{userInfo.name}</span>
          </HeaderGlobalBar>
        </Header>
      )}
    />
  );
};
