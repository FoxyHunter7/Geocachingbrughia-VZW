<?php

namespace App\Services;

use App\Models\StaticSiteProperty;
use App\Services\Service;
use Illuminate\Support\MessageBag;

class StaticSiteContentService extends Service
{
    protected $_rules = [
        'property' => 'required',
        'langcode' => 'required',
        'content' => 'requried'
    ];

    protected $_fields = ['property'];
    protected $_searchOn = 'property';
    protected $_defaultSortBy = 'property';
    protected $_paginate = false;

    public function __construct(StaticSiteProperty $model)
    {
        $this->_model = $model;
        $this->_errors = new MessageBag;
    }

    protected function getLanguageSpecificFields($language, $isAdmin)
    {
        return ['contents'];
    }
}
