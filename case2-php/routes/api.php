<?php

use App\Http\Controllers\AuthController;
use App\Http\Controllers\MerchantReportController;
use App\Http\Controllers\OutletReportController;
use App\Http\Middleware\MultiTenancyMiddleware;
use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
*/

// Public routes
Route::prefix('auth')->group(function () {
    Route::post('/login', [AuthController::class, 'login']);
});

// Protected routes — require JWT + multi-tenancy middleware
Route::middleware([MultiTenancyMiddleware::class])->group(function () {
    Route::post('/auth/logout', [AuthController::class, 'logout']);

    // Merchant revenue report (default month: November)
    Route::get('/merchant/report', [MerchantReportController::class, 'index']);

    // Outlet revenue report (default month: August)
    Route::get('/outlet/{outlet_id}/report', [OutletReportController::class, 'index'])
        ->where('outlet_id', '[0-9]+');
});
