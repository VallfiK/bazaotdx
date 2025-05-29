package com.bazaotdx.bazaotdx.ui.components

import androidx.compose.foundation.layout.*
import com.bazaotdx.bazaotdx.ui.theme.BazaOtdxIcons
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import java.time.LocalDate
import java.time.YearMonth
import java.time.format.DateTimeFormatter

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun BookingCalendar(
    selectedDate: LocalDate,
    onDateSelected: (LocalDate) -> Unit
) {
    var currentMonth by remember { mutableStateOf(YearMonth.now()) }
    
    Column(
        modifier = Modifier
            .fillMaxWidth()
            .padding(8.dp)
    ) {
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween
        ) {
            IconButton(
                onClick = { currentMonth = currentMonth.minusMonths(1) }
            ) {
                Icon(
                    imageVector = BazaOtdxIcons.ArrowBack,
                    contentDescription = "Previous month"
                )
            }
            Text(
                text = currentMonth.format(DateTimeFormatter.ofPattern("MMMM yyyy")),
                style = MaterialTheme.typography.titleMedium
            )
            IconButton(
                onClick = { currentMonth = currentMonth.plusMonths(1) }
            ) {
                Icon(
                    imageVector = BazaOtdxIcons.ArrowForward,
                    contentDescription = "Next month"
                )
            }
        }
        
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceEvenly
        ) {
            Text("Пн")
            Text("Вт")
            Text("Ср")
            Text("Чт")
            Text("Пт")
            Text("Сб")
            Text("Вс")
        }
        
        val firstDayOfMonth = currentMonth.atDay(1)
        val firstDayOfWeek = firstDayOfMonth.dayOfWeek.value
        val daysInMonth = currentMonth.lengthOfMonth()
        
        var day = 1
        
        for (week in 0..5) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceEvenly
            ) {
                for (dayOfWeek in 1..7) {
                    if (week == 0 && dayOfWeek < firstDayOfWeek) {
                        Text("")
                        continue
                    }
                    
                    if (day > daysInMonth) {
                        Text("")
                        break
                    }
                    
                    val date = currentMonth.atDay(day)
                    val isSelected = date == selectedDate
                    
                    Button(
                        onClick = { onDateSelected(date) },
                        colors = ButtonDefaults.buttonColors(
                            containerColor = if (isSelected) MaterialTheme.colorScheme.primary else MaterialTheme.colorScheme.surface
                        )
                    ) {
                        Text(
                            text = day.toString(),
                            color = if (isSelected) MaterialTheme.colorScheme.onPrimary else MaterialTheme.colorScheme.onSurface
                        )
                    }
                    
                    day++
                }
            }
        }
    }
}
