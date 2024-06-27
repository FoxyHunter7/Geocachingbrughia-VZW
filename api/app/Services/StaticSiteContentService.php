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
    protected $_imageLocation = 'static/';
    protected $_translationsRelationName = 'contents';

    protected $_languages;

    public function __construct(StaticSiteProperty $model, LanguageService $languageService)
    {
        $this->_model = $model;
        $this->_errors = new MessageBag;
        $this->_languages = $languageService->all();
    }

    protected function getLanguageSpecificFields($language, $isAdmin)
    {
        return ['contents'];
    }

    protected function imageName($image, $data)
    {
        return 'slpash' . '.' . $image->extension();
    }

    protected function saveImageUrl(&$data, $imageName)
    {
        foreach ($this->_languages as $language) {
            $data['contents'][$language->code] = $this->_imageLocation.$imageName;
        }
    }
}
