<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class Transaction extends Model
{
    protected $table = 'Transactions';

    protected $fillable = [
        'merchant_id',
        'outlet_id',
        'bill_total',
        'created_by',
        'updated_by',
    ];

    protected $casts = [
        'bill_total' => 'float',
    ];

    public function merchant()
    {
        return $this->belongsTo(Merchant::class, 'merchant_id');
    }

    public function outlet()
    {
        return $this->belongsTo(Outlet::class, 'outlet_id');
    }
}
