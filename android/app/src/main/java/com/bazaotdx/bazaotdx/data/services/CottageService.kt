package com.bazaotdx.bazaotdx.data.services

import com.bazaotdx.bazaotdx.data.models.Cottage
import com.bazaotdx.bazaotdx.data.repository.CottageRepository
import kotlinx.coroutines.flow.Flow

class CottageService(private val repository: CottageRepository) {
    val allCottages: Flow<List<Cottage>> = repository.allCottages

    suspend fun addCottage(cottage: Cottage): Long {
        return repository.insertCottage(cottage)
    }

    suspend fun updateCottage(cottage: Cottage): Int {
        return repository.updateCottage(cottage)
    }

    suspend fun deleteCottage(cottage: Cottage): Int {
        return repository.deleteCottage(cottage)
    }

    suspend fun getCottageById(id: Long): Cottage? {
        return repository.getCottageById(id)
    }

    suspend fun getCottagesByStatus(status: String): List<Cottage> {
        return repository.allCottages.first().filter { it.status == status }
    }
}
