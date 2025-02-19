import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';

interface NotificationsState {
  notifications: {notif: string, status: string}[]
}

// Define the initial state using that type
const initialState: NotificationsState = {
 notifications: []
};

export const notificationsSlice = createSlice({
  name: 'notifications',
  initialState,
  reducers: {
    updateNotifications: (state, action: PayloadAction<NotificationsState>) => {
      state.notifications.push(...action.payload.notifications);
    },
    dismissNotification: (state, action: PayloadAction<string>) => {
      state.notifications = state.notifications.filter((notif) => notif.notif !== action.payload);
    }
  },
});

export const { updateNotifications, dismissNotification } = notificationsSlice.actions;

export default notificationsSlice.reducer;
