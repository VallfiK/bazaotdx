package com.bazaotdx.bazaotdx.data.services

import com.bazaotdx.bazaotdx.data.models.Tariff
import com.bazaotdx.bazaotdx.data.repository.TariffRepository
import kotlinx.coroutines.flow.Flow

class TariffService(private val repository: TariffRepository) {
    val allTariffs: Flow<List<Tariff>> = repository.allTariffs

    suspend fun addTariff(tariff: Tariff): Long {
        return repository.insertTariff(tariff)
    }

    suspend fun updateTariff(tariff: Tariff): Int {
        return repository.updateTariff(tariff)
    }

    suspend fun deleteTariff(tariff: Tariff): Int {
        return repository.deleteTariff(tariff)
    }

    suspend fun getTariffById(id: Long): Tariff? {
        return repository.getTariffById(id)
    }

    suspend fun getTariffsByPriceRange(minPrice: Double, maxPrice: Double): List<Tariff> {
        return repository.allTariffs.first().filter { it.pricePerDay in minPrice..maxPrice }
    }
}
