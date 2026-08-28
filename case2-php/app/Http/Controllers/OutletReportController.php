<?php

namespace App\Http\Controllers;

use App\Models\Outlet;
use App\Services\RevenueReportService;
use Illuminate\Http\Request;

class OutletReportController extends Controller
{
    public function __construct(private RevenueReportService $reportService) {}

    /**
     * GET /api/outlet/{outlet_id}/report
     *
     * Returns daily revenue for a specific outlet.
     * The outlet must belong to the authenticated user's merchant (multi-tenancy).
     * Defaults to August if no month/year query params are provided.
     *
     * Query Params:
     *   - month    (int, 1-12, default: 8)
     *   - year     (int, default: current year)
     *   - page     (int, default: 1)
     *   - per_page (int, default: 10, max: 31)
     */
    public function index(Request $request, int $outletId)
    {
        $request->validate([
            'month'    => 'nullable|integer|min:1|max:12',
            'year'     => 'nullable|integer|min:2000|max:2100',
            'page'     => 'nullable|integer|min:1',
            'per_page' => 'nullable|integer|min:1|max:31',
        ]);

        // Multi-tenancy: merchant_id injected by MultiTenancyMiddleware
        $merchantId = $request->attributes->get('merchant_id');

        // Verify the outlet belongs to this merchant
        $outlet = Outlet::where('id', $outletId)
            ->where('merchant_id', $merchantId)
            ->first();

        if (!$outlet) {
            return response()->json([
                'success' => false,
                'error'   => [
                    'code'    => 'FORBIDDEN',
                    'message' => 'You do not have access to this outlet.',
                ],
            ], 403);
        }

        $month   = (int) $request->query('month', 8);
        $year    = (int) $request->query('year', now()->year);
        $page    = (int) $request->query('page', 1);
        $perPage = (int) $request->query('per_page', 10);

        $result = $this->reportService->getOutletReport(
            outletId: $outletId,
            month:    $month,
            year:     $year,
            page:     $page,
            perPage:  $perPage,
        );

        return response()->json([
            'success' => true,
            'message' => 'Outlet revenue report retrieved successfully.',
            'data'    => $result['data'],
            'meta'    => $result['meta'],
        ], 200);
    }
}
