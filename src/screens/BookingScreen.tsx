import React, { useState } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  ScrollView,
  Image,
} from 'react-native';
import { useNavigation, useRoute } from '@react-navigation/native';
import { useDispatch, useSelector } from 'react-redux';
import { bookCottage } from '../store/cottageSlice';
import { format } from 'date-fns';
import ru from 'date-fns/locale/ru';

const BookingScreen = () => {
  const navigation = useNavigation();
  const route = useRoute();
  const dispatch = useDispatch();
  const cottage = route.params?.cottage;
  const [bookingData, setBookingData] = useState({
    checkIn: '',
    checkOut: '',
    guests: 1,
  });

  const handleChange = (name: string, value: string) => {
    setBookingData({
      ...bookingData,
      [name]: value,
    });
  };

  const handleSubmit = async () => {
    try {
      await dispatch(bookCottage({
        cottageId: cottage?._id,
        ...bookingData
      }));
      navigation.navigate('Profile');
    } catch (error) {
      console.error('Booking error:', error);
    }
  };

  if (!cottage) {
    return <Text>Коттедж не найден</Text>;
  }

  return (
    <ScrollView contentContainerStyle={styles.container}>
      <View style={styles.content}>
        <Image
          source={{ uri: cottage.images[0] }}
          style={styles.image}
          resizeMode="cover"
        />
        
        <Text style={styles.cottageName}>{cottage.name}</Text>
        <Text style={styles.description}>{cottage.description}</Text>
        
        <View style={styles.infoContainer}>
          <Text style={styles.infoLabel}>Цена:</Text>
          <Text style={styles.infoValue}>{cottage.price} ₽ за сутки</Text>
        </View>
        
        <View style={styles.infoContainer}>
          <Text style={styles.infoLabel}>Вместимость:</Text>
          <Text style={styles.infoValue}>{cottage.capacity} человек</Text>
        </View>
        
        <Text style={styles.sectionTitle}>Детали бронирования</Text>
        
        <TextInput
          style={styles.input}
          placeholder="Дата заезда"
          value={bookingData.checkIn}
          onChangeText={(value) => handleChange('checkIn', value)}
        />
        
        <TextInput
          style={styles.input}
          placeholder="Дата выезда"
          value={bookingData.checkOut}
          onChangeText={(value) => handleChange('checkOut', value)}
        />
        
        <TextInput
          style={styles.input}
          placeholder="Количество гостей"
          value={bookingData.guests.toString()}
          onChangeText={(value) => handleChange('guests', value)}
          keyboardType="numeric"
        />
        
        <TouchableOpacity style={styles.button} onPress={handleSubmit}>
          <Text style={styles.buttonText}>Забронировать</Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    flexGrow: 1,
    padding: 20,
    backgroundColor: '#fff',
  },
  content: {
    flex: 1,
  },
  image: {
    width: '100%',
    height: 200,
    borderRadius: 8,
    marginBottom: 20,
  },
  cottageName: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 10,
    color: '#4CAF50',
  },
  description: {
    fontSize: 16,
    marginBottom: 20,
    color: '#666',
  },
  infoContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 15,
  },
  infoLabel: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
  },
  infoValue: {
    fontSize: 16,
    color: '#666',
  },
  sectionTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 15,
    color: '#4CAF50',
  },
  input: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 8,
    padding: 15,
    marginBottom: 15,
    fontSize: 16,
  },
  button: {
    backgroundColor: '#4CAF50',
    padding: 15,
    borderRadius: 8,
    alignItems: 'center',
    marginTop: 20,
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
  },
});

export default BookingScreen;
