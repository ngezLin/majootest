<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Merchant extends Model
{
    protected $table = 'Merchants';

    protected $fillable = ['user_id', 'merchant_name', 'created_by', 'updated_by'];

    /** A merchant belongs to a user. */
    public function user()
    {
        return $this->belongsTo(User::class, 'user_id');
    }

    /** A merchant has many outlets. */
    public function outlets()
    {
        return $this->hasMany(Outlet::class, 'merchant_id');
    }

    /** A merchant has many transactions. */
    public function transactions()
    {
        return $this->hasMany(Transaction::class, 'merchant_id');
    }
}
