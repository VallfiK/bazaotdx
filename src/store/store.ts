import { configureStore } from '@reduxjs/toolkit';
import authReducer from './authSlice';
import cottageReducer from './cottageSlice';
import bookingReducer from './bookingSlice';

export const store = configureStore({
  reducer: {
    auth: authReducer,
    cottages: cottageReducer,
    bookings: bookingReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
