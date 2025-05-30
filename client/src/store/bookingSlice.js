import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import axios from 'axios';

export const fetchBookings = createAsyncThunk(
  'bookings/fetchBookings',
  async () => {
    const token = localStorage.getItem('token');
    const response = await axios.get(
      'http://localhost:5000/api/bookings',
      {
        headers: { 'x-auth-token': token }
      }
    );
    return response.data;
  }
);

const initialState = {
  items: [],
  loading: false,
  error: null,
};

const bookingSlice = createSlice({
  name: 'bookings',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchBookings.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchBookings.fulfilled, (state, action) => {
        state.loading = false;
        state.items = action.payload;
      })
      .addCase(fetchBookings.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message;
      });
  },
});

export default bookingSlice.reducer;
