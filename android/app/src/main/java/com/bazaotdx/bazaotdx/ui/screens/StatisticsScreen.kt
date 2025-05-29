package com.bazaotdx.bazaotdx.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bazaotdx.bazaotdx.ui.components.StatisticsCard
import com.bazaotdx.bazaotdx.ui.viewmodels.ReportsViewModel
import java.time.LocalDate

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun StatisticsScreen(
    reportsViewModel: ReportsViewModel
) {
    var selectedStartDate by remember { mutableStateOf(LocalDate.now().minusMonths(1)) }
    var selectedEndDate by remember { mutableStateOf(LocalDate.now()) }
    var selectedTab by remember { mutableStateOf(0) }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp)
    ) {
        Text(
            text = "Статистика",
            style = MaterialTheme.typography.headlineMedium
        )
        Spacer(modifier = Modifier.height(16.dp))

        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween
        ) {
            OutlinedTextField(
                value = selectedStartDate.toString(),
                onValueChange = { selectedStartDate = LocalDate.parse(it) },
                label = { Text("Начальная дата") }
            )
            OutlinedTextField(
                value = selectedEndDate.toString(),
                onValueChange = { selectedEndDate = LocalDate.parse(it) },
                label = { Text("Конечная дата") }
            )
            Button(
                onClick = { 
                    reportsViewModel.generateFinancialReport(selectedStartDate, selectedEndDate)
                }
            ) {
                Text("Обновить")
            }
        }

        Spacer(modifier = Modifier.height(16.dp))

        TabRow(
            selectedTabIndex = selectedTab,
            containerColor = MaterialTheme.colorScheme.primary
        ) {
            Tab(
                selected = selectedTab == 0,
                onClick = { selectedTab = 0 },
                text = { Text("Финансовая") }
            )
            Tab(
                selected = selectedTab == 1,
                onClick = { selectedTab = 1 },
                text = { Text("Загрузка") }
            )
        }

        Spacer(modifier = Modifier.height(16.dp))

        when (selectedTab) {
            0 -> FinancialStats(reportsViewModel)
            1 -> OccupancyStats(reportsViewModel)
        }
    }
}

@Composable
fun FinancialStats(reportsViewModel: ReportsViewModel) {
    val reports = reportsViewModel.financialReports.value
    
    if (reports.isNotEmpty()) {
        val report = reports.first()
        
        StatisticsCard(
            title = "Общая выручка",
            value = "${String.format("%.2f", report.totalRevenue)} ₽",
            description = "За выбранный период"
        )
        
        Spacer(modifier = Modifier.height(16.dp))
        
        StatisticsCard(
            title = "Средняя стоимость бронирования",
            value = "${String.format("%.2f", report.averageBookingValue)} ₽",
            description = "За выбранный период"
        )
        
        Spacer(modifier = Modifier.height(16.dp))
        
        StatisticsCard(
            title = "Наиболее популярный тариф",
            value = report.mostPopularTariff,
            description = "Наиболее часто выбираемый"
        )
        
        Spacer(modifier = Modifier.height(16.dp))
        
        StatisticsCard(
            title = "Наиболее востребованный домик",
            value = report.mostBookedCottage,
            description = "Наиболее часто бронируемый"
        )
    }
}

@Composable
fun OccupancyStats(reportsViewModel: ReportsViewModel) {
    val reports = reportsViewModel.occupancyReports.value
    
    if (reports.isNotEmpty()) {
        val report = reports.first()
        
        StatisticsCard(
            title = "Общее количество домиков",
            value = report.totalCottages.toString(),
            description = "Всего домиков"
        )
        
        Spacer(modifier = Modifier.height(16.dp))
        
        StatisticsCard(
            title = "Занятых домиков",
            value = report.occupiedCottages.toString(),
            description = "На дату"
        )
        
        Spacer(modifier = Modifier.height(16.dp))
        
        StatisticsCard(
            title = "Процент загрузки",
            value = "${String.format("%.1f", report.occupancyRate)}%",
            description = "На дату"
        )
    }
}
