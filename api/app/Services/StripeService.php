<?php

namespace App\Services;

class StripeService
{
    protected $_stripe;

    public function __construct()
    {

        $this->_stripe = new \Stripe\StripeClient(env("STRIPE_SK", ""));
    }

    public function allProducts($admin, $language) {
        $products = $this->_stripe->products->all($admin ? [] : ['active' => true]);

        $lang_filtered_products = collect($products->data)
        ->filter(function ($product) use ($language) {
            return isset($product->metadata['lang']) && $product->metadata['lang'] === $language;
        });

        return $lang_filtered_products;
    }
}
