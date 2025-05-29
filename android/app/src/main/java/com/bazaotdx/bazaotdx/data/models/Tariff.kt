package com.bazaotdx.bazaotdx.data.models

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "tariffs")
data class Tariff(
    @PrimaryKey(autoGenerate = true)
    val id: Long = 0,
    val name: String,
    val pricePerDay: Double
)
