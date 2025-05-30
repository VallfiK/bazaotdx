import React, { useState } from 'react';
import {
  Container,
  Paper,
  Typography,
  TextField,
  Button,
  Box,
  Grid,
  Card,
  CardContent,
  CardMedia,
  Divider,
} from '@mui/material';
import { useNavigate, useParams } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { bookCottage } from '../store/cottageSlice';

const Booking = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const { id } = useParams();
  const [bookingData, setBookingData] = useState({
    checkIn: '',
    checkOut: '',
    guests: 1,
  });
  const cottage = useSelector((state) => state.cottages.items.find(c => c._id === id));

  const handleChange = (e) => {
    setBookingData({
      ...bookingData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await dispatch(bookCottage({
        cottageId: id,
        ...bookingData
      }));
      navigate('/profile');
    } catch (error) {
      console.error('Booking error:', error);
    }
  };

  if (!cottage) {
    return <Typography variant="h6">Коттедж не найден</Typography>;
  }

  return (
    <Container component="main" maxWidth="md">
      <Box sx={{ mt: 4, mb: 4 }}>
        <Grid container spacing={3}>
          <Grid item xs={12} md={8}>
            <Paper elevation={3}>
              <Typography variant="h5" sx={{ p: 2, mb: 2 }}>
                Бронирование коттеджа
              </Typography>
              <Box component="form" onSubmit={handleSubmit} sx={{ p: 3 }}>
                <TextField
                  margin="normal"
                  required
                  fullWidth
                  type="date"
                  label="Дата заезда"
                  name="checkIn"
                  InputLabelProps={{ shrink: true }}
                  value={bookingData.checkIn}
                  onChange={handleChange}
                />
                <TextField
                  margin="normal"
                  required
                  fullWidth
                  type="date"
                  label="Дата выезда"
                  name="checkOut"
                  InputLabelProps={{ shrink: true }}
                  value={bookingData.checkOut}
                  onChange={handleChange}
                />
                <TextField
                  margin="normal"
                  required
                  fullWidth
                  type="number"
                  label="Количество гостей"
                  name="guests"
                  value={bookingData.guests}
                  onChange={handleChange}
                  inputProps={{ min: 1, max: cottage.capacity }}
                />
                <Button
                  type="submit"
                  fullWidth
                  variant="contained"
                  color="primary"
                  sx={{ mt: 3, mb: 2 }}
                >
                  Забронировать
                </Button>
              </Box>
            </Paper>
          </Grid>
          <Grid item xs={12} md={4}>
            <Card elevation={3}>
              <CardMedia
                component="img"
                height="200"
                image={cottage.images[0]}
                alt={cottage.name}
              />
              <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                  {cottage.name}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {cottage.description}
                </Typography>
                <Divider sx={{ my: 2 }} />
                <Typography variant="h6">
                  Цена: {cottage.price} ₽ за сутки
                </Typography>
                <Typography variant="body2">
                  Вместимость: {cottage.capacity} человек
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </Box>
    </Container>
  );
};

export default Booking;
