<?php

namespace App\Models;

use Illuminate\Foundation\Auth\User as Authenticatable;
use Tymon\JWTAuth\Contracts\JWTSubject;

class User extends Authenticatable implements JWTSubject
{
    protected $table = 'Users';

    protected $fillable = ['email', 'password'];

    protected $hidden = ['password'];

    protected $casts = [
        'password' => 'hashed',
    ];

    /**
     * Get the identifier that will be stored in the JWT subject claim.
     */
    public function getJWTIdentifier(): mixed
    {
        return $this->getKey();
    }

    /**
     * Return a key-value array with any custom claims to add to the JWT payload.
     */
    public function getJWTCustomClaims(): array
    {
        return [];
    }

    /**
     * A user has one merchant.
     */
    public function merchant()
    {
        return $this->hasOne(Merchant::class, 'user_id');
    }
}
