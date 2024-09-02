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

    public function allProductsPublic()
    {
        return $this->_service->allProducts(false);
    }

    public function allProductsAdmin()
    {
        return $this->_service->allProducts(true);
    }
}
