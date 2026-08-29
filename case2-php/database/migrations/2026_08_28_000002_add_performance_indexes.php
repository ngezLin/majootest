<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

/**
 * Performance Indexing Migration
 *
 * DML Documentation & Strategy:
 * 1. idx_transactions_merchant_date (merchant_id, created_at)
 *    - Accelerates monthly merchant reports: WHERE merchant_id = ? AND created_at BETWEEN ? AND ?
 * 2. idx_transactions_outlet_date (outlet_id, created_at)
 *    - Accelerates monthly outlet reports: WHERE outlet_id = ? AND created_at BETWEEN ? AND ?
 * 3. idx_outlets_merchant (merchant_id)
 *    - Accelerates tenant ownership validation in MultiTenancyMiddleware
 * 4. idx_merchants_user (user_id)
 *    - Accelerates user -> merchant resolution on JWT authentication
 */
return new class extends Migration
{
    public function up(): void
    {
        DB::statement('ALTER TABLE Transactions ADD INDEX idx_transactions_merchant_date (merchant_id, created_at)');
        DB::statement('ALTER TABLE Transactions ADD INDEX idx_transactions_outlet_date (outlet_id, created_at)');
        DB::statement('ALTER TABLE Outlets ADD INDEX idx_outlets_merchant (merchant_id)');
        DB::statement('ALTER TABLE Merchants ADD INDEX idx_merchants_user (user_id)');
    }

    public function down(): void
    {
        DB::statement('ALTER TABLE Transactions DROP INDEX idx_transactions_merchant_date');
        DB::statement('ALTER TABLE Transactions DROP INDEX idx_transactions_outlet_date');
        DB::statement('ALTER TABLE Outlets DROP INDEX idx_outlets_merchant');
        DB::statement('ALTER TABLE Merchants DROP INDEX idx_merchants_user');
    }
};
