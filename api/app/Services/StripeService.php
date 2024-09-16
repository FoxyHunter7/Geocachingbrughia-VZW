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
            return isset($product->metadata['lang']) && ($product->metadata['lang'] === $language || $product->metadata['lang'] === '*');
        });

        return $lang_filtered_products->map(function ($product) {
            $price = $this->_stripe->prices->retrieve($product->default_price);

            return [
                'id' => $product->id,
                'name' => $product->name,
                'description' => $product->description,
                'images' => $product->images,
                'price' => [
                    'currency' => $price->currency,
                    'amount' => $price->unit_amount
                ],
                'metadata' => $product->metadata
            ];
        });
    }
}
