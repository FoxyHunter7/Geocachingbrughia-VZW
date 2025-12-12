<?php

namespace App\Services;

use App\Models\Language;
use App\Services\Service;
use Illuminate\Support\MessageBag;

class LanguageService extends Service
{
    protected $_model;
    protected $_errors;
    protected $_rules = [
        'code' => 'required | string | max:2',
        'name' => 'required | string',
        'image' => 'required | image | mimes:jpeg,png,jpg,gif,svg|max:4096'
    ];

    protected $_fields = ['code', 'name', 'imageUrl'];
    protected $_fieldsAdminOnly = [];
    protected $_searchOn = 'code';
    protected $_defaultSortBy = 'code';
    protected $_paginate = false;
    protected $_imageLocation = 'langFlags/';

    public function __construct(Language $model)
    {
        $this->_model = $model;
        $this->_errors = new MessageBag;
    }

    public function checkIfExists($language)
    {
        return $this->_model->select($this->_fields)->get()->contains('code', strtoupper($language));
    }

    protected function imageName($image, $data)
    {
        return $data['code'] . '.' . $image->extension();
    }
}
