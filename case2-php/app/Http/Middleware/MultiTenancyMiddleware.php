<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Tymon\JWTAuth\Facades\JWTAuth;
use Tymon\JWTAuth\Exceptions\JWTException;
use App\Models\Merchant;

/**
 * MultiTenancyMiddleware
 *
 * Enforces strict data isolation per merchant.
 * Resolves the authenticated user's merchant_id from JWT and injects it
 * into the request attributes so controllers never need to trust user input
 * for scoping queries.
 */
class MultiTenancyMiddleware
{
    public function handle(Request $request, Closure $next)
    {
        try {
            $user = JWTAuth::parseToken()->authenticate();
        } catch (JWTException $e) {
            return response()->json([
                'success' => false,
                'error'   => [
                    'code'    => 'UNAUTHORIZED',
                    'message' => 'Token is invalid or expired.',
                ],
            ], 401);
        }

        if (!$user) {
            return response()->json([
                'success' => false,
                'error'   => [
                    'code'    => 'UNAUTHORIZED',
                    'message' => 'User not found.',
                ],
            ], 401);
        }

        // Resolve merchant for this authenticated user
        $merchant = Merchant::where('user_id', $user->id)->first();

        if (!$merchant) {
            return response()->json([
                'success' => false,
                'error'   => [
                    'code'    => 'FORBIDDEN',
                    'message' => 'No merchant associated with this account.',
                ],
            ], 403);
        }

        // Inject merchant_id into request attributes so controllers
        // can safely scope all queries without trusting URL/body inputs.
        $request->attributes->set('merchant_id', $merchant->id);
        $request->attributes->set('merchant', $merchant);

        return $next($request);
    }
}
