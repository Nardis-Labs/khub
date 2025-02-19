import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';

interface AdminModalState {
  openGroupsModal: boolean;
  openPermissionsModal: boolean;
}

// Define the initial state using that type
const initialState: AdminModalState = {
 openGroupsModal: false,
 openPermissionsModal: false
};

export const adminModalState = createSlice({
  name: 'adminModal',
  initialState,
  reducers: {
    updateAdminModalState: (state, action: PayloadAction<AdminModalState>) => {
      state.openGroupsModal = action.payload.openGroupsModal;
      state.openPermissionsModal = action.payload.openPermissionsModal;
    }
  },
});

export const { updateAdminModalState } = adminModalState.actions;

export default adminModalState.reducer;
