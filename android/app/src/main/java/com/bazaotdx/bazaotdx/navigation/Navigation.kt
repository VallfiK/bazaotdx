package com.bazaotdx.bazaotdx.navigation

import androidx.compose.runtime.Composable
import androidx.navigation.NavHostController
import androidx.navigation.NavType
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.navArgument
import com.bazaotdx.bazaotdx.ui.screens.*
import com.bazaotdx.bazaotdx.ui.viewmodels.*

@Composable
fun Navigation(
    navController: NavHostController,
    cottagesViewModel: CottagesViewModel,
    tariffsViewModel: TariffsViewModel,
    bookingsViewModel: BookingsViewModel,
    reportsViewModel: ReportsViewModel
) {
    NavHost(
        navController = navController,
        startDestination = Screen.Cottages.route
    ) {
        composable(Screen.Cottages.route) {
            CottagesScreen(
                cottages = cottagesViewModel.cottages,
                onAddCottageClick = { navController.navigate(Screen.AddCottage.route) },
                onEditCottageClick = { cottage ->
                    navController.navigate(Screen.EditCottage.createRoute(cottage.id))
                },
                onDeleteCottageClick = { cottage ->
                    cottagesViewModel.deleteCottage(cottage)
                }
            )
        }

        composable(Screen.Tariffs.route) {
            TariffsScreen(
                tariffs = tariffsViewModel.tariffs,
                onAddTariffClick = { navController.navigate(Screen.AddTariff.route) },
                onEditTariffClick = { tariff ->
                    navController.navigate(Screen.EditTariff.createRoute(tariff.id))
                },
                onDeleteTariffClick = { tariff ->
                    tariffsViewModel.deleteTariff(tariff)
                }
            )
        }

        composable(Screen.Bookings.route) {
            BookingsScreen(
                bookings = bookingsViewModel.bookings,
                onAddBookingClick = { navController.navigate(Screen.AddBooking.route) },
                onEditBookingClick = { booking ->
                    navController.navigate(Screen.EditBooking.createRoute(booking.id))
                },
                onDeleteBookingClick = { booking ->
                    bookingsViewModel.deleteBooking(booking)
                },
                selectedDate = bookingsViewModel.selectedDate,
                onDateSelected = { date ->
                    bookingsViewModel.selectDate(date)
                }
            )
        }

        composable(Screen.Statistics.route) {
            StatisticsScreen(
                reportsViewModel = reportsViewModel
            )
        }

        composable(Screen.AddCottage.route) {
            // TODO: Implement AddCottageScreen
        }

        composable(
            route = Screen.EditCottage.route,
            arguments = listOf(
                navArgument("cottageId") { type = NavType.StringType }
            )
        ) { backStackEntry ->
            val cottageId = backStackEntry.arguments?.getString("cottageId")
            // TODO: Implement EditCottageScreen
        }

        composable(Screen.AddTariff.route) {
            // TODO: Implement AddTariffScreen
        }

        composable(
            route = Screen.EditTariff.route,
            arguments = listOf(
                navArgument("tariffId") { type = NavType.StringType }
            )
        ) { backStackEntry ->
            val tariffId = backStackEntry.arguments?.getString("tariffId")
            // TODO: Implement EditTariffScreen
        }

        composable(Screen.AddBooking.route) {
            // TODO: Implement AddBookingScreen
        }

        composable(
            route = Screen.EditBooking.route,
            arguments = listOf(
                navArgument("bookingId") { type = NavType.StringType }
            )
        ) { backStackEntry ->
            val bookingId = backStackEntry.arguments?.getString("bookingId")
            // TODO: Implement EditBookingScreen
        }
    }
}
