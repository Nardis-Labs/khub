import { configureStore } from '@reduxjs/toolkit';
import { useDispatch } from 'react-redux';
import { setupListeners } from '@reduxjs/toolkit/query';

import { khubApi } from '../service/khub';
import treeMapResourceDrawerReducer from '../service/resourceDrawerState';
import appThemeReducer from '../service/themeState';
import userIsAdminReducer from '../service/user-admin-state';
import notificationsReducer from '../service/notifications';
import { cronjobFilterReducer, daemonsetFilterReducer, deployFilterReducer, jobFilterReducer, podFilterReducer, statefulsetFilterReducer } from '../service/resource-filters';
import adminModalReducer from '../service/adminmodals';

export const store = configureStore({
  reducer: {
    [khubApi.reducerPath]: khubApi.reducer,
    treeMapResourceDrawer: treeMapResourceDrawerReducer,
    appTheme: appThemeReducer,
    notifications: notificationsReducer,
    adminModalState: adminModalReducer,
    podFilter: podFilterReducer,
    deployFilter: deployFilterReducer,
    statefulsetFilter: statefulsetFilterReducer,
    daemonsetFilter: daemonsetFilterReducer,
    cronJobFilter: cronjobFilterReducer,
    jobFilter: jobFilterReducer,
    userIsAdminState: userIsAdminReducer,
  },
  middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(khubApi.middleware)
});

setupListeners(store.dispatch);

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>;
// Inferred type: {posts: PostsState, comments: CommentsState, users: UsersState}
export type AppDispatch = typeof store.dispatch;
export const useAppDispatch: () => AppDispatch = useDispatch;
