import React, { CSSProperties } from "react";
import { Content } from '@carbon/react';
import { useWindowDimensions } from "../../Util";
import { Route, Routes } from "react-router-dom";
import { ClusterOverview } from "../features/ClusterOverview/ClusterOverview";
import { MySQLReplTopo } from "../features/MySQLTopology/MySQLTopology";
import { Pods } from "../features/Pods/Pods";
import { Deployments } from "../features/Deployments/Deployments";
import { Statefulsets } from "../features/Statefulsets/Statefulsets";
import { Daemonsets } from "../features/Daemonsets/Daemonsets";
import { Jobs } from "../features/Jobs/Jobs";
import { Services } from "../features/Services/Services";
import { Ingresses } from "../features/Ingresses/Ingresses";
import { CronJobs } from "../features/CronJobs/CronJobs";
import { ConfigMaps } from "../features/ConfigMaps/ConfigMaps";
import { Nodes } from "../features/Nodes/Nodes";
import { AccessControl } from "../features/AccessControl/AccessControl";
import { Reports } from "../features/Reports/Reports";
import { GeneralSettings } from "../features/GeneralSettings/GeneralSettings";
import { useGetDynamicAppConfigQuery } from "../../service/khub";



export const ContentArea = () => {
  // const { authState } = useOktaAuth();
  const { width } = useWindowDimensions();
  const useResponsiveOffset = width > 1055;

  const {data: appConfig} = useGetDynamicAppConfigQuery({});
  
  const content = (
    <div className="cds--grid">
      <Routes>
        <Route path="/" element={
          <ClusterOverview />
        }/>
        <Route path="/mysql-replication-topology" element={
          <MySQLReplTopo />
        }/>
        <Route path="/pods" element={
          <Pods appConfig={appConfig}/>
        }/>
        <Route path="/deployments" element={
          <Deployments />
        }/>
        <Route path="/statefulsets" element={
          <Statefulsets />
        }/>
        <Route path="/daemonsets" element={
          <Daemonsets />
        }/>
        <Route path="/jobs" element={
          <Jobs />
        }/>
        <Route path="/cronjobs" element={
          <CronJobs />
        }/>
        <Route path="/services" element={
          <Services />
        }/>
        <Route path="/ingresses" element={
          <Ingresses />
        }/>
        <Route path="/configmaps" element={
          <ConfigMaps />
        }/>
        <Route path="/nodes" element={
          <Nodes />
        }/>
        <Route path="/access-control" element={
          <AccessControl />
        }/>
        <Route path="/general-settings" element={
          <GeneralSettings appConfig={appConfig}/>
        }/>
        <Route path="/reports" element={
          <Reports />
        }/>
      </Routes>
    </div>
  );
  const style: CSSProperties = {
    height: '100%'
  };
  if (useResponsiveOffset) {
    style.marginLeft = '256px';
    style.marginTop = '48px';
  } else {
    style.marginLeft = '0px';
  }
  return <Content id="main-content" style={style}>
          {content}
         </Content>;
};