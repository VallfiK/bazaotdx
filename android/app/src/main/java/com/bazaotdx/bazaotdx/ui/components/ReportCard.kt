package com.bazaotdx.bazaotdx.ui.components

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bazaotdx.bazaotdx.data.reports.FinancialReport
import com.bazaotdx.bazaotdx.data.reports.OccupancyReport
import java.time.format.DateTimeFormatter

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ReportCard(
    report: Any,
    modifier: Modifier = Modifier
) {
    Card(
        modifier = modifier.fillMaxWidth(),
        onClick = {}
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp)
        ) {
            when (report) {
                is OccupancyReport -> {
                    Text(
                        text = "Отчет по занятости на ${report.date.format(DateTimeFormatter.ISO_DATE)}",
                        style = MaterialTheme.typography.titleMedium
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(
                        text = "Всего домиков: ${report.totalCottages}",
                        style = MaterialTheme.typography.bodyMedium
                    )
                    Text(
                        text = "Занятых: ${report.occupiedCottages}",
                        style = MaterialTheme.typography.bodyMedium
                    )
                    Text(
                        text = "Процент занятости: ${"%.2f".format(report.occupancyRate)}%",
                        style = MaterialTheme.typography.bodyMedium
                    )
                }
                is FinancialReport -> {
                    Text(
                        text = "Финансовый отчет с ${report.startDate.format(DateTimeFormatter.ISO_DATE)} по ${report.endDate.format(DateTimeFormatter.ISO_DATE)}",
                        style = MaterialTheme.typography.titleMedium
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(
                        text = "Всего бронирований: ${report.totalBookings}",
                        style = MaterialTheme.typography.bodyMedium
                    )
                    Text(
                        text = "Общая выручка: ${"%.2f".format(report.totalRevenue)} ₽",
                        style = MaterialTheme.typography.bodyMedium
                    )
                    Text(
                        text = "Средняя стоимость бронирования: ${"%.2f".format(report.averageBookingValue)} ₽",
                        style = MaterialTheme.typography.bodyMedium
                    )
                    Text(
                        text = "Популярный тариф: ${report.mostPopularTariff}",
                        style = MaterialTheme.typography.bodyMedium
                    )
                    Text(
                        text = "Чаще всего бронировали: ${report.mostBookedCottage}",
                        style = MaterialTheme.typography.bodyMedium
                    )
                }
            }
        }
    }
}
