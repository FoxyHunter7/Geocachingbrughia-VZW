<?php

namespace App\Services;

use App\Models\Message;
use App\Services\Service;
use Illuminate\Support\MessageBag;

class MessageService extends Service
{
    protected $_model;
    protected $_errors;
    protected $_rules = [
        'state' => 'required | in:ONLLINE,DRAFT, ARCHIVED',
        'title' => 'required | max:200 | string',
        'body' => 'nullable | max:20000 | string'
    ];

    protected $_fields = ['id', 'updated_at'];
    protected $_fieldsAdminOnly = ['state'];
    protected $_searchOn = 'title';
    protected $_defaultSortBy = 'updated_at';

    public function __construct(Message $model)
    {
        $this->_model = $model;
        $this->_errors = new MessageBag;
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
