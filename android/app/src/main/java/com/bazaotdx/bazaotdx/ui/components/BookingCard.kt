package com.bazaotdx.bazaotdx.ui.components

import androidx.compose.foundation.layout.*
import com.bazaotdx.bazaotdx.ui.theme.BazaOtdxIcons
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bazaotdx.bazaotdx.data.models.Booking

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun BookingCard(
    booking: Booking,
    onEditClick: () -> Unit,
    onDeleteClick: () -> Unit
) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        onClick = onEditClick
    ) {
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp)
        ) {
            Text(
                text = booking.fullName,
                style = MaterialTheme.typography.titleMedium
            )
            Spacer(modifier = Modifier.height(8.dp))
            
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                Text(
                    text = "Дата заезда: ${booking.checkInDate}",
                    style = MaterialTheme.typography.bodyMedium
                )
                Text(
                    text = "Дата выезда: ${booking.checkOutDate}",
                    style = MaterialTheme.typography.bodyMedium
                )
            }
            
            Spacer(modifier = Modifier.height(8.dp))
            Text(
                text = "Телефон: ${booking.phone}",
                style = MaterialTheme.typography.bodyMedium
            )
            
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.End
            ) {
                IconButton(onClick = onEditClick) {
                    Icon(
                        imageVector = BazaOtdxIcons.Edit,
                        contentDescription = "Edit"
                    )
                }
                IconButton(onClick = onDeleteClick) {
                    Icon(
                        imageVector = BazaOtdxIcons.Delete,
                        contentDescription = "Delete"
                    )
                }
            }
        }
    }
}
