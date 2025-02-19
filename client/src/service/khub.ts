import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import { wsConnect } from './websocketConnector';
import { IAppConfig } from './types/AppConfig';


const baseURL = 
  process.env.REACT_APP_API_URL !== undefined && 
  process.env.REACT_APP_API_URL !== '' && 
  process.env.REACT_APP_API_URL !== null ? `${process.env.REACT_APP_API_URL}/api` : `${window.location.protocol}//${window.location.hostname}/api`;

const wsProtocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
const wsBaseUrl = 
  process.env.REACT_APP_API_URL !== undefined && 
  process.env.REACT_APP_API_URL !== '' && 
  process.env.REACT_APP_API_URL !== null ? `${wsProtocol}://localhost:8080/api` : `${wsProtocol}://${window.location.hostname}/api`;
const baseQuery = fetchBaseQuery({ baseUrl: baseURL, credentials: 'include', mode: 'cors' });

export const khubApi = createApi({
  reducerPath: 'khubApi',
  baseQuery: baseQuery,
  tagTypes: ['Groups', 'Permissions', 'Reports', 'MySQLDBCatalog', 'DynamicAppConfig', 'ClusterName'],
  endpoints: (builder) => ({
    userInfo: builder.query<any, any>({
      query: () => ({
        url: `/users/me`,
        method: 'GET',
      }),
    }),
    getClusterName: builder.query<any, any>({
      query: () => ({
        url: `/k8s/name`,
        method: 'GET',
      }),
      providesTags: ['ClusterName']
    }),
    getDynamicAppConfig: builder.query<any, any>({
      query: () => ({
        url: `/appconfig`,
        method: 'GET',
      }),
      providesTags: ['DynamicAppConfig']
    }),
    updateDynamicAppConfig: builder.mutation<any, IAppConfig>({
      query: (arg) => ({
        url: `/appconfig`,
        method: 'PUT',
        body: arg
      }),
      invalidatesTags: ['DynamicAppConfig', 'ClusterName']
    }),
    getPods: builder.query<any, any>({
      query: () => ({
        url: `/k8s/pods`,
        method: 'GET',
      }),
      async onCacheEntryAdded(
        arg,
        { updateCachedData, cacheDataLoaded, cacheEntryRemoved }
      ) {
        wsConnect(`${wsBaseUrl}/k8s/pods`, cacheDataLoaded, updateCachedData, cacheEntryRemoved);
      }
    }),
    deletePod: builder.mutation<any, {podName: string, namespace: string}>({
      query: (arg) => ({
        url: `/k8s/pods?name=${arg.podName}&namespace=${arg.namespace}`,
        method: 'DELETE',
      }),
    }),
    getDeployments: builder.query<any, any>({
      query: () => ({
        url: `/k8s/deployments`,
        method: 'GET',
      }),
      async onCacheEntryAdded(
        arg,
        { updateCachedData, cacheDataLoaded, cacheEntryRemoved }
      ) {
        wsConnect(`${wsBaseUrl}/k8s/deployments`, cacheDataLoaded, updateCachedData, cacheEntryRemoved);
      }
    }),
    scaleDeployment: builder.mutation<any, { name: string, namespace: string, replicas: number, labels: any }>({
      query: (arg) => ({
        url: `/k8s/deployments/scale`,
        method: 'POST',
        body: {
          name: arg.name,
          namespace: arg.namespace,
          replicas: arg.replicas,
          resourceLabels: arg.labels
        },
      })
    }),
    getReplicasets: builder.query<any, any>({
      query: () => ({
        url: `/k8s/replicasets`,
        method: 'GET',
      }),
      async onCacheEntryAdded(
        arg,
        { updateCachedData, cacheDataLoaded, cacheEntryRemoved }
      ) {
        wsConnect(`${wsBaseUrl}/k8s/replicasets`, cacheDataLoaded, updateCachedData, cacheEntryRemoved);
      },
    }),
    getDaemonsets: builder.query<any, any>({
      query: () => ({
        url: `/k8s/daemonsets`,
        method: 'GET',
      }),
      async onCacheEntryAdded(
        arg,
        { updateCachedData, cacheDataLoaded, cacheEntryRemoved }
      ) {
        wsConnect(`${wsBaseUrl}/k8s/daemonsets`, cacheDataLoaded, updateCachedData, cacheEntryRemoved);
      },
    }),
    getStatefulsets: builder.query<any, any>({
      query: () => ({
        url: `/k8s/statefulsets`,
        method: 'GET',
      }),
      async onCacheEntryAdded(
        arg,
        { updateCachedData, cacheDataLoaded, cacheEntryRemoved }
      ) {
        wsConnect(`${wsBaseUrl}/k8s/statefulsets`, cacheDataLoaded, updateCachedData, cacheEntryRemoved);
      },
    }),
    getJobs: builder.query<any, any>({
      query: () => ({
        url: `/k8s/jobs`,
        method: 'GET',
      }),
      async onCacheEntryAdded(
        arg,
        { updateCachedData, cacheDataLoaded, cacheEntryRemoved }
      ) {
        wsConnect(`${wsBaseUrl}/k8s/jobs`, cacheDataLoaded, updateCachedData, cacheEntryRemoved);
      },
    }),
    getCronJobs: builder.query<any, any>({
      query: () => ({
        url: `/k8s/cronjobs`,
        method: 'GET',
      }),
      async onCacheEntryAdded(
        arg,
        { updateCachedData, cacheDataLoaded, cacheEntryRemoved }
      ) {
        wsConnect(`${wsBaseUrl}/k8s/cronjobs`, cacheDataLoaded, updateCachedData, cacheEntryRemoved);
      },
    }),
    getServices: builder.query<any, any>({
      query: () => ({
        url: `/k8s/services`,
        method: 'GET',
      }),
      async onCacheEntryAdded(
        arg,
        { updateCachedData, cacheDataLoaded, cacheEntryRemoved }
      ) {
        wsConnect(`${wsBaseUrl}/k8s/services`, cacheDataLoaded, updateCachedData, cacheEntryRemoved);
      },
    }),
    getIngresses: builder.query<any, any>({
      query: () => ({
        url: `/k8s/ingresses`,
        method: 'GET',
      }),
      async onCacheEntryAdded(
        arg,
        { updateCachedData, cacheDataLoaded, cacheEntryRemoved }
      ) {
        wsConnect(`${wsBaseUrl}/k8s/ingresses`, cacheDataLoaded, updateCachedData, cacheEntryRemoved);
      },
    }),
    getConfigMaps: builder.query<any, any>({
      query: () => ({
        url: `/k8s/configmaps`,
        method: 'GET',
      }),
      async onCacheEntryAdded(
        arg,
        { updateCachedData, cacheDataLoaded, cacheEntryRemoved }
      ) {
        wsConnect(`${wsBaseUrl}/k8s/configmaps`, cacheDataLoaded, updateCachedData, cacheEntryRemoved);
      },
    }),
    getNodes: builder.query<any, any>({
      query: () => ({
        url: `/k8s/nodes`,
        method: 'GET',
      }),
      async onCacheEntryAdded(
        arg,
        { updateCachedData, cacheDataLoaded, cacheEntryRemoved }
      ) {
        wsConnect(`${wsBaseUrl}/k8s/nodes`, cacheDataLoaded, updateCachedData, cacheEntryRemoved, 25000);
      },
    }),
    getClusterEvents: builder.query<any, any>({
      query: () => ({
        url: `/k8s/clusterevents`,
        method: 'GET',
      }),
      async onCacheEntryAdded(
        arg,
        { updateCachedData, cacheDataLoaded, cacheEntryRemoved }
      ) {
        wsConnect(`${wsBaseUrl}/k8s/clusterevents`, cacheDataLoaded, updateCachedData, cacheEntryRemoved);
      },
    }),
    rolloutRestart: builder.mutation<any, { kind: string, name: string, namespace: string, labels: any }>({
      query: (arg) => ({
        url: `/k8s/rolloutrestart`,
        method: 'POST',
        body: {
          kind: arg.kind,
          name: arg.name,
          namespace: arg.namespace,
          labels: arg.labels
        }
      })
    }),
    tomcatThreadDump: builder.mutation<any, { kind: string, name: string, namespace: string }>({
      query: (arg) => ({
        url: `/k8s/threaddump`,
        method: 'POST',
        body: {
          kind: arg.kind,
          name: arg.name,
          namespace: arg.namespace
        }
      })
    }),
    getUsers: builder.query<any, any>({
      query: () => ({
        url: `/users`,
        method: 'GET',
      })
    }),
    updateUserTheme: builder.mutation<any, {name: string, darkMode: boolean}>({
      query: (arg) => ({
        url: `/users/theme/${arg.name}?darkMode=${arg.darkMode}`,
        method: 'PUT',
      })
    }),
    getGroups: builder.query<any, any>({
      query: () => ({
        url: `/groups`,
        method: 'GET',
      }),
      providesTags: ['Groups']
    }),
    upsertGroup: builder.mutation<any, { id: string, name: string, users: any[], permissions: any[] }>({
      query: (arg) => ({
        url: `/groups`,
        method: 'PUT',
        body: {
          id: arg.id,
          name: arg.name,
          users: arg.users,
          permissions: arg.permissions
        }
      }),
      invalidatesTags: ['Groups']
    }),
    getPermissions: builder.query<any, any>({
      query: () => ({
        url: `/permissions`,
        method: 'GET',
      }),
      providesTags: ['Permissions']
    }),
    upsertPermission: builder.mutation<any, { id: string, name: string, grant: string}>({
      query: (arg) => ({
        url: `/permissions`,
        method: 'PUT',
        body: {
          id: arg.id,
          name: arg.name,
          appTag: arg.grant
        }
      }),
      invalidatesTags: ['Permissions']
    }),
    getMySQLDBCatalog: builder.query<any, any>({
      query: () => ({
        url: `/infra/mysql`,
        method: 'GET',
      }),
      providesTags: ['MySQLDBCatalog']
    }),
    upsertMySQLDBInfo: builder.mutation<any, { shortName: string, host: string, username: string, port: number, isPrimary: boolean }>({
      query: (arg) => ({
        url: `/infra/mysql`,
        method: 'PUT',
        body: {
          shortName: arg.shortName,
          host: arg.host,
          username: arg.username,
          port: arg.port,
          isPrimary: arg.isPrimary
        }
      }),
      invalidatesTags: ['MySQLDBCatalog']
    }),
    deleteMySQLDBInfo: builder.mutation<any, {host: string}>({
      query: (arg) => ({
        url: `/infra/mysql?dbHost=${arg.host}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['MySQLDBCatalog']
    }),
    getMySQLReplicationTopologyGraph: builder.query<any, any>({
      query: () => ({
        url: `/infra/mysql/topology`,
        method: 'GET',
      }),
    }),
    getReports: builder.query<any, any>({
      query: () => ({
        url: `/reports`,
        method: 'GET',
      }),
      providesTags: ['Reports']
    }),
    getReportDownloadURL: builder.query<any, {key: string}>({
      query: (arg) => ({
        url: `/reports/download/${arg.key}`,
        method: 'GET',
      }),
    }),
  })
});

// Export hooks for usage in functional components, which are
// auto-generated based on the defined endpoints
export const { 
  useUserInfoQuery,
  useGetClusterNameQuery,
  useGetDynamicAppConfigQuery,
  useUpdateDynamicAppConfigMutation,
  useGetPodsQuery,
  useDeletePodMutation,
  useGetDeploymentsQuery,
  useScaleDeploymentMutation,
  useGetReplicasetsQuery,
  useGetDaemonsetsQuery,
  useGetStatefulsetsQuery,
  useGetCronJobsQuery,
  useGetServicesQuery,
  useGetConfigMapsQuery,
  useGetIngressesQuery,
  useGetNodesQuery,
  useGetJobsQuery,
  useGetClusterEventsQuery,
  useRolloutRestartMutation,
  useTomcatThreadDumpMutation,
  useGetUsersQuery,
  useUpdateUserThemeMutation,
  useGetGroupsQuery,
  useUpsertGroupMutation,
  useGetPermissionsQuery,
  useUpsertPermissionMutation,
  useGetReportsQuery,
  useGetReportDownloadURLQuery,
  useGetMySQLDBCatalogQuery,
  useUpsertMySQLDBInfoMutation,
  useDeleteMySQLDBInfoMutation,
  useGetMySQLReplicationTopologyGraphQuery
} = khubApi;


