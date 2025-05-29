package com.bazaotdx.bazaotdx.ui.dialogs

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bazaotdx.bazaotdx.data.models.Cottage

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun EditCottageDialog(
    cottage: Cottage? = null,
    onConfirm: (Cottage) -> Unit,
    onDismiss: () -> Unit
) {
    var name by remember { mutableStateOf(cottage?.name ?: "") }
    var status by remember { mutableStateOf(cottage?.status ?: "Свободно") }

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text(if (cottage == null) "Добавить домик" else "Редактировать домик") },
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
                    value = status,
                    onValueChange = { status = it },
                    label = { Text("Статус") },
                    modifier = Modifier.fillMaxWidth()
                )
            }
        },
        confirmButton = {
            Button(
                onClick = {
                    onConfirm(Cottage(
                        id = cottage?.id ?: 0,
                        name = name,
                        status = status
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
