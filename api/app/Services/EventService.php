<?php

namespace App\Services;

use App\Models\Event;
use App\Services\Service;
use Illuminate\Support\MessageBag;

class EventService extends Service
{
    protected $_rules = [
        'state' => 'required | in:ONLINE,DRAFT, ARCHIVED',
        'title' => 'required | max:100 | string',
        'geolink' => 'required | url:https | starts_with:https://www.geocaching.com/geocache/| string',
        'type' => 'required | in:REGULAR,CITO,MEGA,GIGA,BLOCK',
        'location' => 'nullable | regex:/^[NS]\s\d+°\s\d+\.\d+\s[EW]\s\d+°\s\d+\.\d+$/',
        'start_date' => 'required|date_format:Y-m-d H:i:s',
        'end_date' => 'required|date_format:Y-m-d H:i:s',
        'image' => 'required | image | mimes:jpeg,png,jpg,gif,svg|max:4096'
    ];

    protected $_fields = ['id', "on_home", 'title', 'geolink', 'type', 'location', 'start_date', 'end_date', 'imageUrl'];
    protected $_fieldsAdminOnly = ['state'];
    protected $_searchOn = 'title';
    protected $_defaultSortBy = 'start_date';
    protected $_imageLocation = 'events/';

    public function __construct(Event $model)
    {
        $this->_model = $model;
        $this->_errors = new MessageBag;
    }

    public function getHomePageEvents($language)
    {
        $model = $this->_model;
        $model = $model->where("on_home", "true");

        $results = $model
        ->select($this->_fields)
        ->with($this->getLanguageSpecificFields($language, false))
        ->get();

        return $results;
    }

    protected function getLanguageSpecificFields($language, $isAdmin)
    {
        return [
            'translations' => function($query) use ($language, $isAdmin) {
                if (!$isAdmin)
                {
                    $query->where('lang_code', 'like', $language);
                }
            }
        ];
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
