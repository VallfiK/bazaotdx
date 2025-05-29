package com.bazaotdx.bazaotdx.data.services

import com.bazaotdx.bazaotdx.data.models.Booking
import com.bazaotdx.bazaotdx.data.models.Cottage
import com.bazaotdx.bazaotdx.data.models.Tariff
import com.bazaotdx.bazaotdx.data.reports.FinancialReport
import com.bazaotdx.bazaotdx.data.reports.OccupancyReport
import java.time.LocalDate
import java.time.temporal.ChronoUnit

class ReportService(
    private val bookingService: BookingService,
    private val cottageService: CottageService,
    private val tariffService: TariffService
) {
    suspend fun generateOccupancyReport(date: LocalDate): OccupancyReport {
        val bookings = bookingService.getBookingsByDateRange(date, date).first()
        val totalCottages = cottageService.allCottages.first().size
        val occupiedCottages = bookings.map { it.cottageId }.distinct().size
        val occupancyRate = (occupiedCottages.toDouble() / totalCottages) * 100
        
        return OccupancyReport(
            date = date,
            totalCottages = totalCottages,
            occupiedCottages = occupiedCottages,
            occupancyRate = occupancyRate,
            bookings = bookings
        )
    }

    suspend fun generateFinancialReport(startDate: LocalDate, endDate: LocalDate): FinancialReport {
        val bookings = bookingService.getBookingsByDateRange(startDate, endDate).first()
        val totalBookings = bookings.size
        val totalRevenue = bookingService.getRevenueByDateRange(startDate, endDate)
        val averageBookingValue = if (totalBookings > 0) totalRevenue / totalBookings else 0.0

        val tariffs = tariffService.allTariffs.first()
        val cottages = cottageService.allCottages.first()

        val tariffUsage = bookings.groupingBy { it.tariffId }
            .eachCount()
            .maxByOrNull { it.value }
            ?.let { tariffs.find { it.id == it.key }?.name } ?: ""

        val cottageUsage = bookings.groupingBy { it.cottageId }
            .eachCount()
            .maxByOrNull { it.value }
            ?.let { cottages.find { it.id == it.key }?.name } ?: ""

        return FinancialReport(
            startDate = startDate,
            endDate = endDate,
            totalBookings = totalBookings,
            totalRevenue = totalRevenue,
            averageBookingValue = averageBookingValue,
            mostPopularTariff = tariffUsage,
            mostBookedCottage = cottageUsage
        )
    }

    suspend fun generateMonthlyReport(year: Int, month: Int): List<OccupancyReport> {
        val firstDay = LocalDate.of(year, month, 1)
        val lastDay = firstDay.plusMonths(1).minusDays(1)
        
        return (firstDay..lastDay).map { date ->
            generateOccupancyReport(date)
        }
    }

    suspend fun generateYearlyReport(year: Int): List<FinancialReport> {
        return (1..12).map { month ->
            val firstDay = LocalDate.of(year, month, 1)
            val lastDay = firstDay.plusMonths(1).minusDays(1)
            generateFinancialReport(firstDay, lastDay)
        }
    }
}
