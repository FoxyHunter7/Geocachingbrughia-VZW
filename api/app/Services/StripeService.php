<?php

namespace App\Services;

class StripeService
{
    protected $_stripe;

    public function __construct()
    {

        $this->_stripe = new \Stripe\StripeClient(env("STRIPE_SK", ""));
    }

    public function allProducts($admin) {
        return $this->_stripe->products->all($admin ? [] : ['active' => true]);
    }
}
