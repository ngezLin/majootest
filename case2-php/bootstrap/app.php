<?php

use App\Http\Middleware\MultiTenancyMiddleware;
use Illuminate\Foundation\Application;
use Illuminate\Foundation\Configuration\Exceptions;
use Illuminate\Foundation\Configuration\Middleware;
use Illuminate\Http\Request;
use Tymon\JWTAuth\Exceptions\JWTException;
use Tymon\JWTAuth\Exceptions\TokenExpiredException;
use Tymon\JWTAuth\Exceptions\TokenInvalidException;

return Application::configure(basePath: dirname(__DIR__))
    ->withRouting(
        api: __DIR__.'/../routes/api.php',
        apiPrefix: 'api',
        web: __DIR__.'/../routes/web.php',
        commands: __DIR__.'/../routes/console.php',
        health: '/up',
    )
    ->withMiddleware(function (Middleware $middleware): void {
        $middleware->alias([
            'multi_tenancy' => MultiTenancyMiddleware::class,
        ]);
    })
    ->withExceptions(function (Exceptions $exceptions): void {
        // Return JSON errors for all API requests
        $exceptions->shouldRenderJsonWhen(
            fn(Request $request) => $request->is('api/*') || $request->expectsJson(),
        );

        // Handle JWT-specific exceptions with clean JSON responses
        $exceptions->render(function (TokenExpiredException $e, Request $request) {
            return response()->json([
                'success' => false,
                'error'   => ['code' => 'TOKEN_EXPIRED', 'message' => 'Token has expired.'],
            ], 401);
        });

        $exceptions->render(function (TokenInvalidException $e, Request $request) {
            return response()->json([
                'success' => false,
                'error'   => ['code' => 'TOKEN_INVALID', 'message' => 'Token is invalid.'],
            ], 401);
        });

        $exceptions->render(function (JWTException $e, Request $request) {
            return response()->json([
                'success' => false,
                'error'   => ['code' => 'TOKEN_ABSENT', 'message' => 'Token not provided.'],
            ], 401);
        });
    })->create();
