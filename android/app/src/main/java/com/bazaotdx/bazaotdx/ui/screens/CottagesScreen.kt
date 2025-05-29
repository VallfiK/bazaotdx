package com.bazaotdx.bazaotdx.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.bazaotdx.bazaotdx.data.models.Cottage
import com.bazaotdx.bazaotdx.ui.components.CottageCard

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun CottagesScreen(
    cottages: List<Cottage>,
    onAddCottageClick: () -> Unit,
    onEditCottageClick: (Cottage) -> Unit,
    onDeleteCottageClick: (Cottage) -> Unit
) {
    Scaffold(
        floatingActionButton = {
            FloatingActionButton(
                onClick = onAddCottageClick,
                containerColor = MaterialTheme.colorScheme.primary
            ) {
                Icon(
                    imageVector = Icons.Default.Add,
                    contentDescription = "Add Cottage"
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
                text = "Домики",
                style = MaterialTheme.typography.headlineMedium
            )
            Spacer(modifier = Modifier.height(16.dp))
            
            LazyColumn {
                items(cottages) { cottage ->
                    CottageCard(
                        cottage = cottage,
                        onEditClick = { onEditCottageClick(cottage) },
                        onDeleteClick = { onDeleteCottageClick(cottage) }
                    )
                    Spacer(modifier = Modifier.height(8.dp))
                }
            }
        }
    }
}
