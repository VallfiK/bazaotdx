package com.bazaotdx.bazaotdx.ui.viewmodels

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bazaotdx.bazaotdx.data.models.Booking
import com.bazaotdx.bazaotdx.data.reports.FinancialReport
import com.bazaotdx.bazaotdx.data.reports.OccupancyReport
import com.bazaotdx.bazaotdx.data.services.ReportService
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import java.time.LocalDate

class ReportsViewModel(
    private val service: ReportService
) : ViewModel() {
    private val _occupancyReports = MutableStateFlow<List<OccupancyReport>>(emptyList())
    val occupancyReports: StateFlow<List<OccupancyReport>> = _occupancyReports.asStateFlow()

    private val _financialReports = MutableStateFlow<List<FinancialReport>>(emptyList())
    val financialReports: StateFlow<List<FinancialReport>> = _financialReports.asStateFlow()

    fun generateDailyReport(date: LocalDate) {
        viewModelScope.launch {
            _occupancyReports.value = listOf(service.generateOccupancyReport(date))
        }
    }

    fun generateMonthlyReport(year: Int, month: Int) {
        viewModelScope.launch {
            _occupancyReports.value = service.generateMonthlyReport(year, month)
        }
    }

    fun generateYearlyReport(year: Int) {
        viewModelScope.launch {
            _financialReports.value = service.generateYearlyReport(year)
        }
    }

    fun generateFinancialReport(startDate: LocalDate, endDate: LocalDate) {
        viewModelScope.launch {
            _financialReports.value = listOf(service.generateFinancialReport(startDate, endDate))
        }
    }
}
