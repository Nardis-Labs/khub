import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';

interface AppThemeState {
  theme: string
}

// Define the initial state using that type
const initialState: AppThemeState = {
 theme: 'dark'
};

export const appThemeSlice = createSlice({
  name: 'appTheme',
  initialState,
  reducers: {
    updateAppTheme: (state, action: PayloadAction<AppThemeState>) => {
      state.theme = action.payload.theme;
    }
  },
});

export const { updateAppTheme } = appThemeSlice.actions;

export default appThemeSlice.reducer;
