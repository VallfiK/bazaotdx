package com.bazaotdx.bazaotdx.data.dao

import androidx.room.*
import com.bazaotdx.bazaotdx.data.models.Booking
import kotlinx.coroutines.flow.Flow
import java.time.LocalDate

@Dao
interface BookingDao {
    @Query("SELECT * FROM bookings")
    fun getAllBookings(): Flow<List<Booking>>

    @Query("SELECT * FROM bookings WHERE checkInDate BETWEEN :startDate AND :endDate")
    fun getBookingsByDateRange(startDate: LocalDate, endDate: LocalDate): Flow<List<Booking>>

    @Insert
    suspend fun insertBooking(booking: Booking)

    @Update
    suspend fun updateBooking(booking: Booking)

    @Delete
    suspend fun deleteBooking(booking: Booking)

    @Query("SELECT * FROM bookings WHERE id = :id")
    suspend fun getBookingById(id: Long): Booking?

    @Query("SELECT COUNT(*) FROM bookings WHERE checkInDate BETWEEN :startDate AND :endDate")
    suspend fun getBookingCountByDateRange(startDate: LocalDate, endDate: LocalDate): Int

    @Query("SELECT SUM(tariff.pricePerDay * (julianday(checkOutDate) - julianday(checkInDate))) FROM bookings JOIN tariffs ON bookings.tariffId = tariffs.id WHERE checkInDate BETWEEN :startDate AND :endDate")
    suspend fun getRevenueByDateRange(startDate: LocalDate, endDate: LocalDate): Double
}
