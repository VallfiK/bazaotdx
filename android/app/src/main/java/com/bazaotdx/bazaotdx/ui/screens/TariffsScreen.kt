package com.bazaotdx.bazaotdx.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bazaotdx.bazaotdx.ui.components.TariffCard

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TariffsScreen(
    tariffs: List<Tariff>,
    onAddTariffClick: () -> Unit,
    onEditTariffClick: (Tariff) -> Unit,
    onDeleteTariffClick: (Tariff) -> Unit
) {
    Scaffold(
        floatingActionButton = {
            FloatingActionButton(
                onClick = onAddTariffClick,
                containerColor = MaterialTheme.colorScheme.primary
            ) {
                Icon(
                    imageVector = Icons.Default.Add,
                    contentDescription = "Add Tariff"
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
                text = "Тарифы",
                style = MaterialTheme.typography.headlineMedium
            )
            Spacer(modifier = Modifier.height(16.dp))
            
            LazyColumn {
                items(tariffs) { tariff ->
                    TariffCard(
                        tariff = tariff,
                        onEditClick = { onEditTariffClick(tariff) },
                        onDeleteClick = { onDeleteTariffClick(tariff) }
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                }
            }
        }
    }
}
