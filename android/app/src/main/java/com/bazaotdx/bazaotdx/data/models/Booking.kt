package com.bazaotdx.bazaotdx.data.models

import androidx.room.Entity
import androidx.room.PrimaryKey
import java.time.LocalDate

@Entity(tableName = "bookings")
data class Booking(
    @PrimaryKey(autoGenerate = true)
    val id: Long = 0,
    val fullName: String,
    val email: String,
    val phone: String,
    val cottageId: Long,
    val documentScanPath: String?,
    val checkInDate: LocalDate,
    val checkOutDate: LocalDate,
    val tariffId: Long
)
