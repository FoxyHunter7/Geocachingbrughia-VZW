<?php

namespace App\Services;

use App\Models\Social;
use App\Services\Service;
use Illuminate\Support\MessageBag;

class SocialService extends Service
{
    protected $_model;
    protected $_errors;
    protected $_rules = [
        'name' => 'required | string | max:70',
        'url' => 'required | url:https',
        'image' => 'required | images | mimes:jpeg,png,jpg,gif,svg|max:4096'
    ];

    protected $_fields = ['id', 'name', 'url', 'imageUrl'];
    protected $_fieldsAdminOnly = [];
    protected $_searchOn = 'name';
    protected $_defaultSortBy = 'name';
    protected $_paginate = false;
    protected $_imageLocation = 'socials/';

    public function __construct(Social $model)
    {
        $this->_model = $model;
        $this->_errors = new MessageBag;
    }

    protected function modelHiddenoverwrites($results, $isAdmin)
    {
        if ($isAdmin) {
            $results->makeVisible(['id']);
        }
    }
}
