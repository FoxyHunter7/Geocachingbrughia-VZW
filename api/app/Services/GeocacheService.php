<?php

namespace App\Services;

use App\Models\Geocache;
use App\Services\Service;
use Illuminate\Support\MessageBag;

class GeocacheService extends Service
{
    protected $_model;
    protected $_errors;
    protected $_rules = [
        'state' => 'required | in:ONLINE,DRAFT, ARCHIVED',
        'title' => 'required | max:100 | string',
        'geolink' => 'nullable | url:https | starts_with:https://www.geocaching.com/geocache/| string',
        'type' => 'required | in:TRADITIONAL,MULTI,MYSTERY,EARTH,LETTERBOX,WHEREIGO,VIRTUAL,LAB,WEBCAM',
        'difficulty' => 'required | numeric | between:1,5',
        'terrain' => 'required | numeric | between:1,5',
        'placed_on' => 'date_format:Y-m-d'
    ];

    protected $_fields = ['id', 'title', 'geolink', 'type', 'difficulty', 'terrain'];
    protected $_fieldsAdminOnly = ['state', 'placed_on'];
    protected $_searchOn = 'title';
    protected $_defaultSortBy = 'placed_on';

    public function __construct(Geocache $model)
    {
        $this->_model = $model;
        $this->_errors = new MessageBag;
    }

    protected function nonAdminConditions($query)
    {
        $query->where('state', 'ONLINE');
    }

    protected function modelHiddenoverwrites($results, $isAdmin)
    {
        if ($isAdmin) {
            $results->makeVisible(['id']);
        }
    }
}
