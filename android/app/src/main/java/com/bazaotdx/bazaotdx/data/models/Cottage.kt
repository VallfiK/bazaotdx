package com.bazaotdx.bazaotdx.data.models

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "cottages")
data class Cottage(
    @PrimaryKey(autoGenerate = true)
    val id: Long = 0,
    val name: String,
    val status: String
)
