<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

/**
 * Performance Indexing Migration
 *
 * DML Documentation & Justification:
 *
 * 1. idx_transactions_merchant_date (merchant_id, created_at)
 *    - Used by: GET /api/merchant/report
 *    - Query pattern: WHERE merchant_id = ? AND created_at BETWEEN ? AND ?
 *    - Without index: full scan of Transactions (grows unbounded with business volume)
 *    - With index: InnoDB B-tree seeks directly to merchant's rows then scans date range
 *
 * 2. idx_transactions_outlet_date (outlet_id, created_at)
 *    - Used by: GET /api/outlet/{id}/report
 *    - Same pattern but scoped to outlet_id
 *
 * 3. idx_outlets_merchant (merchant_id)
 *    - Used by: MultiTenancyMiddleware outlet ownership check
 *    - Query pattern: WHERE id = ? AND merchant_id = ?
 *    - Prevents table scan when verifying outlet belongs to authenticated merchant
 *
 * 4. idx_merchants_user (user_id)
 *    - Used by: MultiTenancyMiddleware resolving merchant from JWT user_id
 *    - Ensures O(log n) lookup instead of full Merchants scan on every request
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
