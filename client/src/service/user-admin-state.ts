import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';

interface UserIsAdminState {
  isAdmin: boolean
}

// Define the initial state using that type
const initialState: UserIsAdminState = {
 isAdmin: false
};

export const userIsAdminSlice = createSlice({
  name: 'userIsAdmin',
  initialState,
  reducers: {
    updateUserIsAdmin: (state, action: PayloadAction<UserIsAdminState>) => {
      state.isAdmin = action.payload.isAdmin;
    }
  },
});

export const { updateUserIsAdmin } = userIsAdminSlice.actions;

export default userIsAdminSlice.reducer;
