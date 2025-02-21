/*eslint-disable */
export interface IAppConfig {
  id: number;
  data: IAppConfigData;
}

export interface IAppConfigData {
  defaultReplicaScaleLimit: number;
  replicaScaleLimits: { [key: string]: number };
  enableK8sGlobalReadOnly: boolean;
  k8sClusterName: string;
  k8sClusterNamespaces: string[];
  k8sPodExecPlugins: IPodExecPlugin[];
}

export interface IPodExecPlugin {
  name: string;
  container: string
  command: string;
  labelFilter: string;
}