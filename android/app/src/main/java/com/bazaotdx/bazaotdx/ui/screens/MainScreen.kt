package com.bazaotdx.bazaotdx.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.bazaotdx.bazaotdx.ui.components.*
import com.bazaotdx.bazaotdx.ui.viewmodels.*
import java.time.LocalDate

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun MainScreen(
    cottagesViewModel: CottagesViewModel,
    tariffsViewModel: TariffsViewModel,
    bookingsViewModel: BookingsViewModel,
    reportsViewModel: ReportsViewModel
) {
    val navController = rememberNavController()
    val scaffoldState = rememberScaffoldState()

    Scaffold(
        scaffoldState = scaffoldState,
        topBar = {
            TopAppBar(
                title = { Text("База Отдых") },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = MaterialTheme.colorScheme.primary
                )
            )
        },
        bottomBar = {
            BottomNavigation {
                val navBackStackEntry by navController.currentBackStackEntryAsState()
                val currentRoute = navBackStackEntry?.destination?.route

                BottomNavigationItem(
                    icon = { Icon(Icons.Default.Home, contentDescription = "Домики") },
                    label = { Text("Домики") },
                    selected = currentRoute == "cottages",
                    onClick = { navController.navigate("cottages") }
                )

                BottomNavigationItem(
                    icon = { Icon(Icons.Default.Money, contentDescription = "Тарифы") },
                    label = { Text("Тарифы") },
                    selected = currentRoute == "tariffs",
                    onClick = { navController.navigate("tariffs") }
                )

                BottomNavigationItem(
                    icon = { Icon(Icons.Default.Event, contentDescription = "Бронирования") },
                    label = { Text("Брони") },
                    selected = currentRoute == "bookings",
                    onClick = { navController.navigate("bookings") }
                )

                BottomNavigationItem(
                    icon = { Icon(Icons.Default.Analytics, contentDescription = "Статистика") },
                    label = { Text("Статистика") },
                    selected = currentRoute == "statistics",
                    onClick = { navController.navigate("statistics") }
                )
            }
        }
    ) { paddingValues ->
        NavHost(
            navController = navController,
            startDestination = "cottages",
            modifier = Modifier.padding(paddingValues)
        ) {
            composable("cottages") {
                CottagesScreen(
                    cottagesViewModel = cottagesViewModel,
                    onAddCottage = { navController.navigate("add_cottage") }
                )
            }

            composable("tariffs") {
                TariffsScreen(
                    tariffsViewModel = tariffsViewModel,
                    onAddTariff = { navController.navigate("add_tariff") }
                )
            }

            composable("bookings") {
                BookingsScreen(
                    bookingsViewModel = bookingsViewModel,
                    onAddBooking = { navController.navigate("add_booking") }
                )
            }

            composable("statistics") {
                StatisticsScreen(reportsViewModel = reportsViewModel)
            }

            composable("add_cottage") {
                AddCottageScreen(
                    cottagesViewModel = cottagesViewModel,
                    onBack = { navController.popBackStack() }
                )
            }

            composable("add_tariff") {
                AddTariffScreen(
                    tariffsViewModel = tariffsViewModel,
                    onBack = { navController.popBackStack() }
                )
            }

            composable("add_booking") {
                AddBookingScreen(
                    bookingsViewModel = bookingsViewModel,
                    onBack = { navController.popBackStack() }
                )
            }
        }
    }
}
