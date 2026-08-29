<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Outlet extends Model
{
    protected $table = 'Outlets';

    protected $fillable = ['merchant_id', 'outlet_name', 'created_by', 'updated_by'];

    public function merchant()
    {
        return $this->belongsTo(Merchant::class, 'merchant_id');
    }

    public function transactions()
    {
        return $this->hasMany(Transaction::class, 'outlet_id');
    }
}
