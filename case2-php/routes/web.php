<?php

use Illuminate\Support\Facades\Route;

// Interactive Swagger UI documentation
Route::get('/', function () {
    return view('swagger');
});

Route::get('/docs', function () {
    return view('swagger');
});
