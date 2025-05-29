package com.bazaotdx.bazaotdx

import android.app.Application
import org.koin.android.ext.koin.androidContext
import org.koin.android.ext.koin.androidLogger
import org.koin.core.context.startKoin

class BazaOtdxApp : Application() {
    override fun onCreate() {
        super.onCreate()
        
        startKoin {
            androidLogger()
            androidContext(this@BazaOtdxApp)
            modules(KoinModule.appModule)
        }
    }
}
