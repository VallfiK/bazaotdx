package com.bazaotdx.bazaotdx.data.repository

import com.bazaotdx.bazaotdx.data.dao.TariffDao
import com.bazaotdx.bazaotdx.data.models.Tariff
import kotlinx.coroutines.flow.Flow

class TariffRepository(private val tariffDao: TariffDao) {
    val allTariffs: Flow<List<Tariff>> = tariffDao.getAllTariffs()

    suspend fun insertTariff(tariff: Tariff) {
        tariffDao.insertTariff(tariff)
    }

    suspend fun updateTariff(tariff: Tariff) {
        tariffDao.updateTariff(tariff)
    }

    suspend fun deleteTariff(tariff: Tariff) {
        tariffDao.deleteTariff(tariff)
    }

    suspend fun getTariffById(id: Long): Tariff? {
        return tariffDao.getTariffById(id)
    }
}
