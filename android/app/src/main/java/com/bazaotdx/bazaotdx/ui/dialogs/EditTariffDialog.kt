package com.bazaotdx.bazaotdx.ui.dialogs

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bazaotdx.bazaotdx.data.models.Tariff

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun EditTariffDialog(
    tariff: Tariff? = null,
    onConfirm: (Tariff) -> Unit,
    onDismiss: () -> Unit
) {
    var name by remember { mutableStateOf(tariff?.name ?: "") }
    var pricePerDay by remember { mutableStateOf(tariff?.pricePerDay?.toString() ?: "0.0") }

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text(if (tariff == null) "Добавить тариф" else "Редактировать тариф") },
        text = {
            Column(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp)
            ) {
                OutlinedTextField(
                    value = name,
                    onValueChange = { name = it },
                    label = { Text("Название") },
                    modifier = Modifier.fillMaxWidth()
                )
                Spacer(modifier = Modifier.height(16.dp))
                
                OutlinedTextField(
                    value = pricePerDay,
                    onValueChange = { pricePerDay = it },
                    label = { Text("Цена за день (₽)") },
                    modifier = Modifier.fillMaxWidth()
                )
            }
        },
        confirmButton = {
            Button(
                onClick = {
                    val price = pricePerDay.toDoubleOrNull() ?: 0.0
                    onConfirm(Tariff(
                        id = tariff?.id ?: 0,
                        name = name,
                        pricePerDay = price
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
