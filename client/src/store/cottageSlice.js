import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import axios from 'axios';

export const fetchCottages = createAsyncThunk(
  'cottages/fetchCottages',
  async () => {
    const response = await axios.get('http://localhost:5000/api/cottages');
    return response.data;
  }
);

export const bookCottage = createAsyncThunk(
  'cottages/bookCottage',
  async (bookingData) => {
    const token = localStorage.getItem('token');
    const response = await axios.post(
      'http://localhost:5000/api/bookings',
      bookingData,
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

const cottageSlice = createSlice({
  name: 'cottages',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchCottages.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchCottages.fulfilled, (state, action) => {
        state.loading = false;
        state.items = action.payload;
      })
      .addCase(fetchCottages.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message;
      })
      .addCase(bookCottage.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(bookCottage.fulfilled, (state, action) => {
        state.loading = false;
      })
      .addCase(bookCottage.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message;
      });
  },
});

export default cottageSlice.reducer;
