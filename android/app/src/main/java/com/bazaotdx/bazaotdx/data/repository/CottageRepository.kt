package com.bazaotdx.bazaotdx.data.repository

import com.bazaotdx.bazaotdx.data.dao.CottageDao
import com.bazaotdx.bazaotdx.data.models.Cottage
import kotlinx.coroutines.flow.Flow

class CottageRepository(private val cottageDao: CottageDao) {
    val allCottages: Flow<List<Cottage>> = cottageDao.getAllCottages()

    suspend fun insertCottage(cottage: Cottage) {
        cottageDao.insertCottage(cottage)
    }

    suspend fun updateCottage(cottage: Cottage) {
        cottageDao.updateCottage(cottage)
    }

    suspend fun deleteCottage(cottage: Cottage) {
        cottageDao.deleteCottage(cottage)
    }

    suspend fun getCottageById(id: Long): Cottage? {
        return cottageDao.getCottageById(id)
    }
}
