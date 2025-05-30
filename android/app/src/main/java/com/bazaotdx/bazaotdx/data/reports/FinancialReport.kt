package com.bazaotdx.bazaotdx.data.reports

data class FinancialReport(
    val startDate: LocalDate,
    val endDate: LocalDate,
    val totalBookings: Int,
    val totalRevenue: Double,
    val averageBookingValue: Double,
    val mostPopularTariff: String,
    val mostBookedCottage: String
)
