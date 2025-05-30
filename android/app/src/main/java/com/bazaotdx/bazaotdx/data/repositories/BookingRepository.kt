package com.bazaotdx.bazaotdx.data.repositories

import com.bazaotdx.bazaotdx.data.dao.BookingDao
import com.bazaotdx.bazaotdx.data.models.Booking
import kotlinx.coroutines.flow.Flow

class BookingRepository(private val dao: BookingDao) {
    fun getAllBookings(): Flow<List<Booking>> = dao.getAllBookings()
    
    suspend fun insertBooking(booking: Booking) = dao.insertBooking(booking)
    
    suspend fun updateBooking(booking: Booking) = dao.updateBooking(booking)
    
    suspend fun deleteBooking(booking: Booking) = dao.deleteBooking(booking)
    
    suspend fun getBookingById(id: Long): Booking? = dao.getBookingById(id)
}
