package com.bazaotdx.bazaotdx.data.reports

import com.bazaotdx.bazaotdx.data.models.Booking
import java.time.LocalDate

data class OccupancyReport(
    val date: LocalDate,
    val totalCottages: Int,
    val occupiedCottages: Int,
    val occupancyRate: Double,
    val bookings: List<Booking>
)
