package com.bazaotdx.bazaotdx.navigation

sealed class Screen(val route: String) {
    object Cottages : Screen("cottages")
    object Tariffs : Screen("tariffs")
    object Bookings : Screen("bookings")
    object Statistics : Screen("statistics")
    
    object AddCottage : Screen("add_cottage")
    object EditCottage : Screen("edit_cottage/{cottageId}") {
        fun createRoute(cottageId: String) = "edit_cottage/$cottageId"
    }
    
    object AddTariff : Screen("add_tariff")
    object EditTariff : Screen("edit_tariff/{tariffId}") {
        fun createRoute(tariffId: String) = "edit_tariff/$tariffId"
    }
    
    object AddBooking : Screen("add_booking")
    object EditBooking : Screen("edit_booking/{bookingId}") {
        fun createRoute(bookingId: String) = "edit_booking/$bookingId"
    }
}
