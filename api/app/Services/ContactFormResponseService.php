<?php

namespace App\Services;

use App\Models\ContactFormResponse;
use App\Services\Service;
use Illuminate\Support\MessageBag;

class ContactFormResponseService extends Service
{
    protected $_model;
    protected $_errors;
    protected $_rules = [
        'email' => 'required | email',
        'subject' => 'required | string | max:100',
        'message' => 'required | string | max:5000'
    ];

    protected $_fields = ['id', 'email', 'subject', 'message', 'created_at'];
    protected $_searchOn = 'subject';
    protected $_defaultSortBy = 'created_at';

    public function __construct(ContactFormResponse $model)
    {
        $this->_model = $model;
        $this->_errors = new MessageBag;
    }


}
