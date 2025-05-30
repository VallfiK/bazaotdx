package com.bazaotdx.bazaotdx.ui.viewmodels

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.bazaotdx.bazaotdx.data.models.Tariff
import com.bazaotdx.bazaotdx.data.services.TariffService
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

class TariffsViewModel(private val service: TariffService) : ViewModel() {
    private val _tariffs = MutableStateFlow<List<Tariff>>(emptyList())
    val tariffs: StateFlow<List<Tariff>> = _tariffs.asStateFlow()

    init {
        viewModelScope.launch {
            service.allTariffs.collect { tariffs ->
                _tariffs.value = tariffs
            }
        }
    }

    fun addTariff(tariff: Tariff) {
        viewModelScope.launch {
            service.addTariff(tariff)
        }
    }

    fun updateTariff(tariff: Tariff) {
        viewModelScope.launch {
            service.updateTariff(tariff)
        }
    }

    fun deleteTariff(tariff: Tariff) {
        viewModelScope.launch {
            service.deleteTariff(tariff)
        }
    }

    fun getTariffsByPriceRange(minPrice: Double, maxPrice: Double) {
        viewModelScope.launch {
            _tariffs.value = service.getTariffsByPriceRange(minPrice, maxPrice)
        }
    }
}
