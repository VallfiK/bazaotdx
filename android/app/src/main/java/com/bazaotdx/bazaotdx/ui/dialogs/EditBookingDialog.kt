package com.bazaotdx.bazaotdx.ui.dialogs

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bazaotdx.bazaotdx.data.models.Booking
import java.time.LocalDate

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun EditBookingDialog(
    booking: Booking? = null,
    cottages: List<String>,
    tariffs: List<String>,
    onConfirm: (Booking) -> Unit,
    onDismiss: () -> Unit
) {
    var fullName by remember { mutableStateOf(booking?.fullName ?: "") }
    var email by remember { mutableStateOf(booking?.email ?: "") }
    var phone by remember { mutableStateOf(booking?.phone ?: "") }
    var selectedCottage by remember { mutableStateOf(booking?.cottageId?.toString() ?: "") }
    var selectedTariff by remember { mutableStateOf(booking?.tariffId?.toString() ?: "") }
    var checkInDate by remember { mutableStateOf(booking?.checkInDate?.toString() ?: "") }
    var checkOutDate by remember { mutableStateOf(booking?.checkOutDate?.toString() ?: "") }

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text(if (booking == null) "Добавить бронирование" else "Редактировать бронирование") },
        text = {
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp)
            ) {
                OutlinedTextField(
                    value = fullName,
                    onValueChange = { fullName = it },
                    label = { Text("ФИО") },
                    modifier = Modifier.fillMaxWidth()
                )
                Spacer(modifier = Modifier.height(16.dp))
                
                OutlinedTextField(
                    value = email,
                    onValueChange = { email = it },
                    label = { Text("Email") },
                    modifier = Modifier.fillMaxWidth()
                )
                Spacer(modifier = Modifier.height(16.dp))
                
                OutlinedTextField(
                    value = phone,
                    onValueChange = { phone = it },
                    label = { Text("Телефон") },
                    modifier = Modifier.fillMaxWidth()
                )
                Spacer(modifier = Modifier.height(16.dp))
                
                OutlinedTextField(
                    value = selectedCottage,
                    onValueChange = { selectedCottage = it },
                    label = { Text("Домик") },
                    modifier = Modifier.fillMaxWidth()
                )
                Spacer(modifier = Modifier.height(16.dp))
                
                OutlinedTextField(
                    value = selectedTariff,
                    onValueChange = { selectedTariff = it },
                    label = { Text("Тариф") },
                    modifier = Modifier.fillMaxWidth()
                )
                Spacer(modifier = Modifier.height(16.dp))
                
                OutlinedTextField(
                    value = checkInDate,
                    onValueChange = { checkInDate = it },
                    label = { Text("Дата заезда") },
                    modifier = Modifier.fillMaxWidth()
                )
                Spacer(modifier = Modifier.height(16.dp))
                
                OutlinedTextField(
                    value = checkOutDate,
                    onValueChange = { checkOutDate = it },
                    label = { Text("Дата выезда") },
                    modifier = Modifier.fillMaxWidth()
                )
            }
        },
        confirmButton = {
            Button(
                onClick = {
                    val checkIn = LocalDate.parse(checkInDate)
                    val checkOut = LocalDate.parse(checkOutDate)
                    val cottageId = selectedCottage.toLongOrNull() ?: 0
                    val tariffId = selectedTariff.toLongOrNull() ?: 0
                    
                    onConfirm(Booking(
                        id = booking?.id ?: 0,
                        fullName = fullName,
                        email = email,
                        phone = phone,
                        cottageId = cottageId,
                        documentScanPath = null,
                        checkInDate = checkIn,
                        checkOutDate = checkOut,
                        tariffId = tariffId
                    ))
                    onDismiss()
                }
            ) {
                Text("Сохранить")
            }
        },
        dismissButton = {
            Button(
                onClick = onDismiss,
                colors = ButtonDefaults.buttonColors(
                    containerColor = MaterialTheme.colorScheme.secondaryContainer
                )
            ) {
                Text("Отмена")
            }
        }
    )
}
