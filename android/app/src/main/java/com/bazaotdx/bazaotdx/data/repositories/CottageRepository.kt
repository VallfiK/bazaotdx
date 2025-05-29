package com.bazaotdx.bazaotdx.data.repositories

import com.bazaotdx.bazaotdx.data.dao.CottageDao
import com.bazaotdx.bazaotdx.data.models.Cottage
import kotlinx.coroutines.flow.Flow

class CottageRepository(private val dao: CottageDao) {
    fun getAllCottages(): Flow<List<Cottage>> = dao.getAllCottages()
    
    suspend fun insertCottage(cottage: Cottage) = dao.insertCottage(cottage)
    
    suspend fun updateCottage(cottage: Cottage) = dao.updateCottage(cottage)
    
    suspend fun deleteCottage(cottage: Cottage) = dao.deleteCottage(cottage)
    
    suspend fun getCottageById(id: Long): Cottage? = dao.getCottageById(id)
}
