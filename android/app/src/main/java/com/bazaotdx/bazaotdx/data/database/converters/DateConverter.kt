package com.bazaotdx.bazaotdx.data.database.converters

import androidx.room.TypeConverter
import java.time.LocalDate
import java.time.format.DateTimeFormatter

class DateConverter {
    private val formatter = DateTimeFormatter.ISO_LOCAL_DATE
    
    @TypeConverter
    fun fromLocalDate(date: LocalDate?): String? = date?.format(formatter)
    
    @TypeConverter
    fun toLocalDate(dateStr: String?): LocalDate? = dateStr?.let { LocalDate.parse(it, formatter) }
}
