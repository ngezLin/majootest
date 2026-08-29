<?php

namespace App\Http\Controllers;

use App\Models\User;
use Illuminate\Http\Request;
use Tymon\JWTAuth\Facades\JWTAuth;
use Tymon\JWTAuth\Exceptions\JWTException;

class AuthController extends Controller
{
    /**
     * POST /api/auth/login
     * Authenticate user via JWT.
     */
    public function login(Request $request)
    {
        $request->validate([
            'email'    => 'required|email',
            'password' => 'required|string',
        ]);

        $credentials = $request->only('email', 'password');

        try {
            if (!$token = JWTAuth::attempt($credentials)) {
                return response()->json([
                    'success' => false,
                    'error'   => [
                        'code'    => 'INVALID_CREDENTIALS',
                        'message' => 'Invalid email or password.',
                    ],
                ], 401);
            }
        } catch (JWTException $e) {
            return response()->json([
                'success' => false,
                'error'   => [
                    'code'    => 'TOKEN_CREATION_FAILED',
                    'message' => 'Could not create JWT token.',
                ],
            ], 500);
        }

        $user     = auth()->user();
        $merchant = $user->merchant;

        return response()->json([
            'success' => true,
            'message' => 'Login successful.',
            'data'    => [
                'token'       => $token,
                'token_type'  => 'Bearer',
                'expires_in'  => JWTAuth::factory()->getTTL() * 60,
                'user'        => [
                    'id'          => $user->id,
                    'email'       => $user->email,
                    'merchant_id' => $merchant?->id,
                    'merchant'    => $merchant?->merchant_name,
                ],
            ],
        ], 200);
    }

    /**
     * POST /api/auth/logout
     * Invalidate current JWT token.
     */
    public function logout()
    {
        try {
            JWTAuth::invalidate(JWTAuth::getToken());
        } catch (JWTException $e) {
            // Already invalid
        }

        return response()->json([
            'success' => true,
            'message' => 'Logged out successfully.',
        ]);
    }
}
