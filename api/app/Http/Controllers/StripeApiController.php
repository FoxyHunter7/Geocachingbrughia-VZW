<?php

namespace App\Http\Controllers;

use App\Services\StripeService;
use Illuminate\Http\Request;

class StripeApiController
{
    private $_service;

    public function __construct(StripeService $service)
    {
        $this->_service = $service;
    }

    public function allProductsPublic(Request $request)
    {
        $language = $request->get('lang');
        return $this->_service->allProducts(false, $language);
    }

    public function allProductsAdmin(Request $request)
    {
        $language = $request->get('lang');
        return $this->_service->allProducts(true, $language);
    }
}
