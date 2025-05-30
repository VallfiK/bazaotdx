package com.bazaotdx.bazaotdx

import com.bazaotdx.bazaotdx.data.database.AppDatabase
import com.bazaotdx.bazaotdx.data.models.Booking
import com.bazaotdx.bazaotdx.data.services.*
import com.bazaotdx.bazaotdx.ui.viewmodels.*
import org.koin.android.ext.koin.androidContext
import org.koin.dsl.module

val appModule = module {
    single { AppDatabase.getDatabase(androidContext()) }
    
    // Services
    single { CottageService(CottageRepository(get().cottageDao())) }
    single { TariffService(TariffRepository(get().tariffDao())) }
    single { BookingService(BookingRepository(get().bookingDao())) }
    single { ReportService(get(), get(), get()) }
    
    // ViewModels
    single { CottagesViewModel(get()) }
    single { TariffsViewModel(get()) }
    single { BookingsViewModel(get()) }
    single { ReportsViewModel(get()) }
}
