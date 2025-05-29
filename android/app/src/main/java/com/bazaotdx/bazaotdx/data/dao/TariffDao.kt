package com.bazaotdx.bazaotdx.data.dao

import androidx.room.*
import com.bazaotdx.bazaotdx.data.models.Tariff
import kotlinx.coroutines.flow.Flow

@Dao
interface TariffDao {
    @Query("SELECT * FROM tariffs")
    fun getAllTariffs(): Flow<List<Tariff>>

    @Insert
    suspend fun insertTariff(tariff: Tariff)

    @Update
    suspend fun updateTariff(tariff: Tariff)

    @Delete
    suspend fun deleteTariff(tariff: Tariff)

    @Query("SELECT * FROM tariffs WHERE id = :id")
    suspend fun getTariffById(id: Long): Tariff?
}
