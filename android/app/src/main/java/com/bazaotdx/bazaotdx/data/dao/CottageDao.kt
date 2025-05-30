package com.bazaotdx.bazaotdx.data.dao

import androidx.room.*
import com.bazaotdx.bazaotdx.data.models.Cottage
import kotlinx.coroutines.flow.Flow

@Dao
interface CottageDao {
    @Query("SELECT * FROM cottages")
    fun getAllCottages(): Flow<List<Cottage>>

    @Insert
    suspend fun insertCottage(cottage: Cottage)

    @Update
    suspend fun updateCottage(cottage: Cottage)

    @Delete
    suspend fun deleteCottage(cottage: Cottage)

    @Query("SELECT * FROM cottages WHERE id = :id")
    suspend fun getCottageById(id: Long): Cottage?
}
