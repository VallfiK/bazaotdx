package com.bazaotdx.bazaotdx.data.repository

import com.bazaotdx.bazaotdx.data.dao.BookingDao
import com.bazaotdx.bazaotdx.data.models.Booking
import kotlinx.coroutines.flow.Flow
import java.time.LocalDate

class BookingRepository(private val bookingDao: BookingDao) {
    val allBookings: Flow<List<Booking>> = bookingDao.getAllBookings()

    suspend fun insertBooking(booking: Booking) {
        bookingDao.insertBooking(booking)
    }

    suspend fun updateBooking(booking: Booking) {
        bookingDao.updateBooking(booking)
    }

    suspend fun deleteBooking(booking: Booking) {
        bookingDao.deleteBooking(booking)
    }

    suspend fun getBookingById(id: Long): Booking? {
        return bookingDao.getBookingById(id)
    }

    fun getBookingsByDateRange(startDate: LocalDate, endDate: LocalDate): Flow<List<Booking>> {
        return bookingDao.getBookingsByDateRange(startDate, endDate)
    }

    suspend fun getBookingCountByDateRange(startDate: LocalDate, endDate: LocalDate): Int {
        return bookingDao.getBookingCountByDateRange(startDate, endDate)
    }

    suspend fun getRevenueByDateRange(startDate: LocalDate, endDate: LocalDate): Double {
        return bookingDao.getRevenueByDateRange(startDate, endDate)
    }
}
