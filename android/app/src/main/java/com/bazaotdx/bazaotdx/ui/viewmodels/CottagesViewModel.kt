package com.bazaotdx.bazaotdx.ui.viewmodels

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bazaotdx.bazaotdx.data.models.Cottage
import com.bazaotdx.bazaotdx.data.services.CottageService
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

class CottagesViewModel(private val service: CottageService) : ViewModel() {
    private val _cottages = MutableStateFlow<List<Cottage>>(emptyList())
    val cottages: StateFlow<List<Cottage>> = _cottages.asStateFlow()

    init {
        viewModelScope.launch {
            service.allCottages.collect { cottages ->
                _cottages.value = cottages
            }
        }
    }

    fun addCottage(cottage: Cottage) {
        viewModelScope.launch {
            service.addCottage(cottage)
        }
    }

    fun updateCottage(cottage: Cottage) {
        viewModelScope.launch {
            service.updateCottage(cottage)
        }
    }

    fun deleteCottage(cottage: Cottage) {
        viewModelScope.launch {
            service.deleteCottage(cottage)
        }
    }

    fun getCottagesByStatus(status: String) {
        viewModelScope.launch {
            _cottages.value = service.getCottagesByStatus(status)
        }
    }
}
