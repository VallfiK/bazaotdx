package com.bazaotdx.bazaotdx

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.ui.Modifier
import androidx.navigation.compose.rememberNavController
import com.bazaotdx.bazaotdx.KoinModule
import com.bazaotdx.bazaotdx.navigation.Navigation
import com.bazaotdx.bazaotdx.ui.theme.BazaOtdxTheme
import org.koin.android.ext.koin.androidContext
import org.koin.android.ext.koin.androidLogger
import org.koin.core.context.startKoin
import org.koin.androidx.viewmodel.ext.android.viewModel

class MainActivity : ComponentActivity() {
    private val cottagesViewModel: CottagesViewModel by viewModel()
    private val tariffsViewModel: TariffsViewModel by viewModel()
    private val bookingsViewModel: BookingsViewModel by viewModel()
    private val reportsViewModel: ReportsViewModel by viewModel()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            BazaOtdxTheme {
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.background
                ) {
                    val navController = rememberNavController()
                    Navigation(
                        navController = navController,
                        cottagesViewModel = cottagesViewModel,
                        tariffsViewModel = tariffsViewModel,
                        bookingsViewModel = bookingsViewModel,
                        reportsViewModel = reportsViewModel
                    )
                }
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MainScreen(
    cottagesViewModel: CottagesViewModel,
    tariffsViewModel: TariffsViewModel,
    bookingsViewModel: BookingsViewModel
) {
    var selectedTab by remember { mutableStateOf(0) }
    val tabs = listOf("Домики", "Тарифы", "Бронирования", "Статистика")
    var showDialog by remember { mutableStateOf(false) }
    var selectedCottage by remember { mutableStateOf<Cottage?>(null) }
    var selectedTariff by remember { mutableStateOf<Tariff?>(null) }
    var selectedBooking by remember { mutableStateOf<Booking?>(null) }
    var selectedDate by remember { mutableStateOf(LocalDate.now()) }

    Column {
        // Верхняя панель
        TopAppBar(
            title = { Text("🌲 Звуки Леса") },
            colors = TopAppBarDefaults.topAppBarColors(
                containerColor = MaterialTheme.colorScheme.primary,
                titleContentColor = MaterialTheme.colorScheme.onPrimary
            )
        )

        // Табы
        TabRow(
            selectedTabIndex = selectedTab,
            containerColor = MaterialTheme.colorScheme.primary
        ) {
            tabs.forEachIndexed { index, title ->
                Tab(
                    selected = selectedTab == index,
                    onClick = { selectedTab = index },
                    text = { Text(title) }
                )
            }
        }

        // Основной контент
        when (selectedTab) {
            0 -> CottagesScreen(
                cottages = cottagesViewModel.cottages.value,
                onAddCottageClick = {
                    selectedCottage = null
                    showDialog = true
                },
                onEditCottageClick = { cottage ->
                    selectedCottage = cottage
                    showDialog = true
                },
                onDeleteCottageClick = { cottage ->
                    cottagesViewModel.deleteCottage(cottage)
                }
            )
            1 -> TariffsScreen(
                tariffs = tariffsViewModel.tariffs.value,
                onAddTariffClick = {
                    selectedTariff = null
                    showDialog = true
                },
                onEditTariffClick = { tariff ->
                    selectedTariff = tariff
                    showDialog = true
                },
                onDeleteTariffClick = { tariff ->
                    tariffsViewModel.deleteTariff(tariff)
                }
            )
            2 -> BookingsScreen(
                bookings = bookingsViewModel.bookings.value,
                onAddBookingClick = {
                    selectedBooking = null
                    showDialog = true
                },
                onEditBookingClick = { booking ->
                    selectedBooking = booking
                    showDialog = true
                },
                onDeleteBookingClick = { booking ->
                    bookingsViewModel.deleteBooking(booking)
                },
                selectedDate = selectedDate,
                onDateSelected = { date ->
                    selectedDate = date
                    bookingsViewModel.getBookingsByDateRange(date, date)
                }
            )
            3 -> StatisticsScreen(
                bookingsCount = bookingsViewModel.getBookingCountByDateRange(
                    LocalDate.now().minusMonths(1),
                    LocalDate.now()
                ),
                revenue = bookingsViewModel.getRevenueByDateRange(
                    LocalDate.now().minusMonths(1),
                    LocalDate.now()
                ),
                startDate = LocalDate.now().minusMonths(1),
                endDate = LocalDate.now(),
                onDateRangeSelected = { startDate, endDate ->
                    bookingsViewModel.getBookingsByDateRange(startDate, endDate)
                }
            )
        }

        // Диалоги
        if (showDialog) {
            when {
                selectedCottage != null -> EditCottageDialog(
                    cottage = selectedCottage,
                    onConfirm = { cottage ->
                        if (selectedCottage == null) {
                            cottagesViewModel.addCottage(cottage)
                        } else {
                            cottagesViewModel.updateCottage(cottage)
                        }
                        showDialog = false
                    },
                    onDismiss = { showDialog = false }
                )
                selectedTariff != null -> EditTariffDialog(
                    tariff = selectedTariff,
                    onConfirm = { tariff ->
                        if (selectedTariff == null) {
                            tariffsViewModel.addTariff(tariff)
                        } else {
                            tariffsViewModel.updateTariff(tariff)
                        }
                        showDialog = false
                    },
                    onDismiss = { showDialog = false }
                )
                selectedBooking != null -> EditBookingDialog(
                    booking = selectedBooking,
                    cottages = cottagesViewModel.cottages.value.map { it.name },
                    tariffs = tariffsViewModel.tariffs.value.map { it.name },
                    onConfirm = { booking ->
                        if (selectedBooking == null) {
                            bookingsViewModel.addBooking(booking)
                        } else {
                            bookingsViewModel.updateBooking(booking)
                        }
                        showDialog = false
                    },
                    onDismiss = { showDialog = false }
                )
            }
        }
    }
}
