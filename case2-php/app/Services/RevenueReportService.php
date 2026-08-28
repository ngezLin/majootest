<?php

namespace App\Services;

use App\Models\Transaction;
use Carbon\Carbon;
use Carbon\CarbonPeriod;
use Illuminate\Support\Facades\DB;

/**
 * RevenueReportService
 *
 * Generates daily revenue (omzet) reports with zero-fill:
 * every date in the requested month appears in the result,
 * even if there were no transactions on that date (revenue = 0).
 */
class RevenueReportService
{
    /**
     * Merchant monthly revenue report.
     */
    public function getMerchantReport(
        int $merchantId,
        int $month,
        int $year,
        int $page,
        int $perPage,
    ): array {
        $dailyRevenue = $this->queryDailyRevenue(
            filters:  ['merchant_id' => $merchantId],
            month:    $month,
            year:     $year,
        );

        return $this->buildPaginatedReport(
            dailyRevenue: $dailyRevenue,
            month:        $month,
            year:         $year,
            page:         $page,
            perPage:      $perPage,
        );
    }

    /**
     * Outlet monthly revenue report.
     */
    public function getOutletReport(
        int $outletId,
        int $month,
        int $year,
        int $page,
        int $perPage,
    ): array {
        $dailyRevenue = $this->queryDailyRevenue(
            filters:  ['outlet_id' => $outletId],
            month:    $month,
            year:     $year,
        );

        return $this->buildPaginatedReport(
            dailyRevenue: $dailyRevenue,
            month:        $month,
            year:         $year,
            page:         $page,
            perPage:      $perPage,
        );
    }

    /**
     * Query aggregated daily revenue from the Transactions table.
     * Uses indexed columns (merchant_id/outlet_id + created_at) for performance.
     *
     * @param array $filters  e.g. ['merchant_id' => 1] or ['outlet_id' => 2]
     */
    private function queryDailyRevenue(array $filters, int $month, int $year): array
    {
        $start = Carbon::create($year, $month, 1)->startOfDay();
        $end   = $start->copy()->endOfMonth()->endOfDay();

        $rows = DB::table('Transactions')
            ->select(
                DB::raw('DATE(created_at) as date'),
                DB::raw('SUM(bill_total) as revenue'),
            )
            ->where($filters)
            ->whereBetween('created_at', [$start, $end])
            ->groupBy(DB::raw('DATE(created_at)'))
            ->orderBy('date')
            ->get()
            ->keyBy('date');

        // Return as plain array keyed by date string for O(1) lookup
        return $rows->map(fn($row) => (float) $row->revenue)->all();
    }

    /**
     * Build a zero-filled, paginated daily report for the entire month.
     *
     * For every day in the month:
     *   - If a transaction existed → use actual SUM(bill_total)
     *   - Otherwise              → use 0
     *
     * Then apply page/per_page slicing on the complete month array.
     */
    private function buildPaginatedReport(
        array $dailyRevenue,
        int   $month,
        int   $year,
        int   $page,
        int   $perPage,
    ): array {
        $start  = Carbon::create($year, $month, 1);
        $end    = $start->copy()->endOfMonth();
        $period = CarbonPeriod::create($start, $end);

        // Build full-month array with zero-fill
        $allDays = [];
        foreach ($period as $date) {
            $dateStr   = $date->toDateString();
            $allDays[] = [
                'date'    => $dateStr,
                'revenue' => $dailyRevenue[$dateStr] ?? 0,
            ];
        }

        $totalDays  = count($allDays);
        $totalPages = (int) ceil($totalDays / $perPage);
        $offset     = ($page - 1) * $perPage;

        $pageData = array_slice($allDays, $offset, $perPage);

        return [
            'data' => $pageData,
            'meta' => [
                'month'       => $month,
                'year'        => $year,
                'current_page'=> $page,
                'per_page'    => $perPage,
                'total_days'  => $totalDays,
                'total_pages' => $totalPages,
            ],
        ];
    }
}
