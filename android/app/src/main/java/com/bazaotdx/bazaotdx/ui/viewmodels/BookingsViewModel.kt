package com.bazaotdx.bazaotdx.ui.viewmodels

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bazaotdx.bazaotdx.data.models.Booking
import com.bazaotdx.bazaotdx.data.services.BookingService
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import java.time.LocalDate

class BookingsViewModel(private val service: BookingService) : ViewModel() {
    private val _bookings = MutableStateFlow<List<Booking>>(emptyList())
    val bookings: StateFlow<List<Booking>> = _bookings.asStateFlow()

    init {
        viewModelScope.launch {
            service.allBookings.collect { bookings ->
                _bookings.value = bookings
            }
        }
    }

    fun addBooking(booking: Booking) {
        viewModelScope.launch {
            service.addBooking(booking)
        }
    }

    fun updateBooking(booking: Booking) {
        viewModelScope.launch {
            service.updateBooking(booking)
        }
    }

    fun deleteBooking(booking: Booking) {
        viewModelScope.launch {
            service.deleteBooking(booking)
        }
    }

    fun getBookingsByDateRange(startDate: LocalDate, endDate: LocalDate) {
        viewModelScope.launch {
            service.getBookingsByDateRange(startDate, endDate).collect { bookings ->
                _bookings.value = bookings
            }
        }
    }

    suspend fun getBookingCountByDateRange(startDate: LocalDate, endDate: LocalDate): Int {
        return service.getBookingCountByDateRange(startDate, endDate)
    }

    suspend fun getRevenueByDateRange(startDate: LocalDate, endDate: LocalDate): Double {
        return service.getRevenueByDateRange(startDate, endDate)
    }

    suspend fun checkCottageAvailability(
        cottageId: Long,
        checkInDate: LocalDate,
        checkOutDate: LocalDate
    ): Boolean {
        return service.checkCottageAvailability(cottageId, checkInDate, checkOutDate)
    }
}
