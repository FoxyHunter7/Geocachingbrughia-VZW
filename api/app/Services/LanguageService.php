<?php

namespace App\Services;

use App\Models\Language;
use App\Services\Service;
use Illuminate\Support\MessageBag;
use Illuminate\Support\Facades\Log;

class LanguageService extends Service
{
    protected $_model;
    protected $_errors;
    protected $_rules = [
        'state' => 'required | string | max:2'
    ];

    protected $_fields = ['code'];
    protected $_fieldsAdminOnly = [];
    protected $_searchOn = 'code';
    protected $_defaultSortBy = 'code';
    protected $_paginate = false;

    public function __construct(Language $model)
    {
        $this->_model = $model;
        $this->_errors = new MessageBag;
    }

    public function checkIfExists($language)
    {
        return $this->_model->select($this->_fields)->get()->contains('code', strtoupper($language));
    }
}
