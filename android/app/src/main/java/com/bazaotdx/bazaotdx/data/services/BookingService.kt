package com.bazaotdx.bazaotdx.data.services

import com.bazaotdx.bazaotdx.data.models.Booking
import com.bazaotdx.bazaotdx.data.repository.BookingRepository
import kotlinx.coroutines.flow.Flow
import java.time.LocalDate

class BookingService(private val repository: BookingRepository) {
    val allBookings: Flow<List<Booking>> = repository.allBookings

    suspend fun addBooking(booking: Booking): Long {
        return repository.insertBooking(booking)
    }

    suspend fun updateBooking(booking: Booking): Int {
        return repository.updateBooking(booking)
    }

    suspend fun deleteBooking(booking: Booking): Int {
        return repository.deleteBooking(booking)
    }

    suspend fun getBookingById(id: Long): Booking? {
        return repository.getBookingById(id)
    }

    fun getBookingsByDateRange(startDate: LocalDate, endDate: LocalDate): Flow<List<Booking>> {
        return repository.getBookingsByDateRange(startDate, endDate)
    }

    suspend fun getBookingCountByDateRange(startDate: LocalDate, endDate: LocalDate): Int {
        return repository.getBookingCountByDateRange(startDate, endDate)
    }

    suspend fun getRevenueByDateRange(startDate: LocalDate, endDate: LocalDate): Double {
        return repository.getRevenueByDateRange(startDate, endDate)
    }

    suspend fun checkCottageAvailability(
        cottageId: Long,
        checkInDate: LocalDate,
        checkOutDate: LocalDate
    ): Boolean {
        val bookings = repository.getBookingsByDateRange(checkInDate, checkOutDate).first()
        return bookings.none { it.cottageId == cottageId }
    }
}
