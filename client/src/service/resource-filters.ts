import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';

interface FilterState {
  filter: string
}

// Define the initial state using that type
const initialState: FilterState = {
 filter: ''
};

export const podFilterSlice = createSlice({
  name: 'podFilter',
  initialState,
  reducers: {
    updatePodFilter: (state, action: PayloadAction<FilterState>) => {
      state.filter = action.payload.filter;
    }
  },
});

export const deployFilterSlice = createSlice({
  name: 'deployFilter',
  initialState,
  reducers: {
    updateDeployFilter: (state, action: PayloadAction<FilterState>) => {
      state.filter = action.payload.filter;
    }
  },
});

export const daemonsetFilterSlice = createSlice({
  name: 'daemonsetFilter',
  initialState,
  reducers: {
    updateDaemonsetFilter: (state, action: PayloadAction<FilterState>) => {
      state.filter = action.payload.filter;
    }
  },
});

export const statefulsetFilterSlice = createSlice({
  name: 'statefulsetFilter',
  initialState,
  reducers: {
    updateStatefulsetFilter: (state, action: PayloadAction<FilterState>) => {
      state.filter = action.payload.filter;
    }
  },
});

export const cronjobFilterSlice = createSlice({
  name: 'cronjobFilter',
  initialState,
  reducers: {
    updateCronJobFilter: (state, action: PayloadAction<FilterState>) => {
      state.filter = action.payload.filter;
    }
  },
});

export const jobFilterSlice = createSlice({
  name: 'jobFilter',
  initialState,
  reducers: {
    updateJobFilter: (state, action: PayloadAction<FilterState>) => {
      state.filter = action.payload.filter;
    }
  },
});

export const { updatePodFilter } = podFilterSlice.actions;
export const { updateCronJobFilter } = cronjobFilterSlice.actions;
export const { updateDaemonsetFilter } = daemonsetFilterSlice.actions;
export const { updateDeployFilter } = deployFilterSlice.actions;
export const { updateJobFilter } = jobFilterSlice.actions;
export const { updateStatefulsetFilter } = statefulsetFilterSlice.actions;

export const podFilterReducer = podFilterSlice.reducer;
export const deployFilterReducer = deployFilterSlice.reducer;
export const daemonsetFilterReducer = daemonsetFilterSlice.reducer;
export const statefulsetFilterReducer = statefulsetFilterSlice.reducer;
export const cronjobFilterReducer = cronjobFilterSlice.reducer;
export const jobFilterReducer = jobFilterSlice.reducer;
