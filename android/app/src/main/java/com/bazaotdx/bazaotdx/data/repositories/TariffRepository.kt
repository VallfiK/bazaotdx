package com.bazaotdx.bazaotdx.data.repositories

import com.bazaotdx.bazaotdx.data.dao.TariffDao
import com.bazaotdx.bazaotdx.data.models.Tariff
import kotlinx.coroutines.flow.Flow

class TariffRepository(private val dao: TariffDao) {
    fun getAllTariffs(): Flow<List<Tariff>> = dao.getAllTariffs()
    
    suspend fun insertTariff(tariff: Tariff) = dao.insertTariff(tariff)
    
    suspend fun updateTariff(tariff: Tariff) = dao.updateTariff(tariff)
    
    suspend fun deleteTariff(tariff: Tariff) = dao.deleteTariff(tariff)
    
    suspend fun getTariffById(id: Long): Tariff? = dao.getTariffById(id)
}
