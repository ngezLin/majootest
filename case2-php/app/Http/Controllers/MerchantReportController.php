<?php

namespace App\Http\Controllers;

use App\Services\RevenueReportService;
use Illuminate\Http\Request;

class MerchantReportController extends Controller
{
    public function __construct(private RevenueReportService $reportService) {}

    /**
     * GET /api/merchant/report
     * Daily revenue report for authenticated merchant.
     */
    public function index(Request $request)
    {
        $request->validate([
            'month'    => 'nullable|integer|min:1|max:12',
            'year'     => 'nullable|integer|min:2000|max:2100',
            'page'     => 'nullable|integer|min:1',
            'per_page' => 'nullable|integer|min:1|max:31',
        ]);

        $merchantId = $request->attributes->get('merchant_id');
        $month      = (int) $request->query('month', 11); // Defaults to November as per assessment
        $year       = (int) $request->query('year', now()->year);
        $page       = (int) $request->query('page', 1);
        $perPage    = (int) $request->query('per_page', 10);

        $result = $this->reportService->getMerchantReport(
            merchantId: $merchantId,
            month:      $month,
            year:       $year,
            page:       $page,
            perPage:    $perPage,
        );

        return response()->json([
            'success' => true,
            'message' => 'Merchant revenue report retrieved successfully.',
            'data'    => $result['data'],
            'meta'    => $result['meta'],
        ], 200);
    }
}
