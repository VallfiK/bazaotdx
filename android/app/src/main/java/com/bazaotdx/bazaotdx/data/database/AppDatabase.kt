package com.bazaotdx.bazaotdx.data.database

import android.content.Context
import androidx.room.Database
import androidx.room.Room
import androidx.room.RoomDatabase
import com.bazaotdx.bazaotdx.data.dao.CottageDao
import com.bazaotdx.bazaotdx.data.models.Cottage
import com.bazaotdx.bazaotdx.data.models.Tariff
import com.bazaotdx.bazaotdx.data.models.Booking
import com.bazaotdx.bazaotdx.data.database.converters.DateConverter

@Database(
    entities = [Cottage::class, Tariff::class, Booking::class],
    version = 1,
    exportSchema = false
)
@TypeConverters(DateConverter::class)
abstract class AppDatabase : RoomDatabase() {
    abstract fun cottageDao(): CottageDao
    abstract fun tariffDao(): TariffDao
    abstract fun bookingDao(): BookingDao

    companion object {
        @Volatile
        private var INSTANCE: AppDatabase? = null

        fun getDatabase(context: Context): AppDatabase {
            return INSTANCE ?: synchronized(this) {
                val instance = Room.databaseBuilder(
                    context.applicationContext,
                    AppDatabase::class.java,
                    "bazaotdx_database"
                ).build()
                INSTANCE = instance
                instance
            }
        }
    }
}
