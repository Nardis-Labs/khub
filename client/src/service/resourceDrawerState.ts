import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';

interface TreeMapResourceDrawerState {
  open: boolean
  data: any
}

// Define the initial state using that type
const initialState: TreeMapResourceDrawerState = {
  open: false,
  data: null
};

export const treeMapResourceDrawerSlice = createSlice({
  name: 'treeMapResourceDrawer',
  initialState,
  reducers: {
    updateTreeMapResourceDrawer: (state, action: PayloadAction<TreeMapResourceDrawerState>) => {
      state.open = action.payload.open;
      state.data = action.payload.data;
    }
  },
});

export const { updateTreeMapResourceDrawer } = treeMapResourceDrawerSlice.actions;

export default treeMapResourceDrawerSlice.reducer;
