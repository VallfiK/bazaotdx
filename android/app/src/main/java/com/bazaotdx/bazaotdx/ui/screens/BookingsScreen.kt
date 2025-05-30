package com.bazaotdx.bazaotdx.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bazaotdx.bazaotdx.ui.components.BookingCard
import com.bazaotdx.bazaotdx.ui.components.BookingCalendar
import java.time.LocalDate

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun BookingsScreen(
    bookings: List<Booking>,
    onAddBookingClick: () -> Unit,
    onEditBookingClick: (Booking) -> Unit,
    onDeleteBookingClick: (Booking) -> Unit,
    selectedDate: LocalDate,
    onDateSelected: (LocalDate) -> Unit
) {
    Scaffold(
        floatingActionButton = {
            FloatingActionButton(
                onClick = onAddBookingClick,
                containerColor = MaterialTheme.colorScheme.primary
            ) {
                Icon(
                    imageVector = Icons.Default.Add,
                    contentDescription = "Add Booking"
                )
            }
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(16.dp)
        ) {
            Text(
                text = "Бронирования",
                style = MaterialTheme.typography.headlineMedium
            )
            Spacer(modifier = Modifier.height(16.dp))
            
            BookingCalendar(
                selectedDate = selectedDate,
                onDateSelected = onDateSelected
            )
            
            Spacer(modifier = Modifier.height(16.dp))
            
            LazyColumn {
                items(bookings) { booking ->
                    BookingCard(
                        booking = booking,
                        onEditClick = { onEditBookingClick(booking) },
                        onDeleteClick = { onDeleteBookingClick(booking) }
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                }
            }
        }
    }
}
