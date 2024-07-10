<?php

namespace App\Services;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Validator;
use Illuminate\Support\ItemNotFoundException;
use Illuminate\Support\MessageBag;
use Illuminate\Support\Str;

abstract class Service
{
    protected $_model;
    protected $_errors;
    protected $_rules;
    protected $_fields;
    protected $_fieldsAdminOnly;
    protected $_searchOn;
    protected $_defaultSortBy = 'created_at';
    protected $_paginate = true;
    protected $_imageLocation = 'app/images/default';
    protected $_translationsRelationName = 'translations';

    public function __construct(Model $model)
    {
        $this->_model = $model;
        $this->_errors = new MessageBag;
    }

    protected function validate($data)
    {
        $this->_errors = new MessageBag();

        $validator = Validator::make($data, $this->_rules);

        if ($validator->fails()) {
            $this->_errors = $validator->errors();
        }
    }

    public function hasErrors() {
        return $this->_errors->any();
    }

    public function getErrors() {
        return $this->_errors;
    }

    public function all($language = null, $perPage = 10, $search = null, $sortBy = null, $sortDirection = 'asc', $isAdmin = false)
    {
        $model = $this->_model;
        $sortBy ??= $this->_defaultSortBy;

        if ($search) {
            $model = $model->where($this->_searchOn, "like", "%$search%");
        }

        $results = $model
            ->select($this->getSelectFields($isAdmin))
            ->when($isAdmin, fn($query) => $this->adminConditions($query))
            ->when(!$isAdmin, fn($query) => $this->nonAdminConditions($query))
            ->orderBy($sortBy, $sortDirection)
            ->with($this->getLanguageSpecificFields($language, $isAdmin));

        if ($this->_paginate)
        {
            $results = $results->paginate($perPage)->withQueryString();
        } else {
            $results = $results->get();
        }

        $this->modelHiddenOverwrites($results, $isAdmin);

        return $results;
    }

    public function getById($id, $language = null, $isAdmin = false)
    {
        $model = $this->_model;
        $model = $model->where("id", $id);

        $results = $model
            ->select($this->getSelectFields($isAdmin))
            ->when($isAdmin, fn($query) => $this->adminConditions($query))
            ->when(!$isAdmin, fn($query) => $this->nonAdminConditions($query))
            ->with($this->getLanguageSpecificFields($language, $isAdmin))
            ->get();

        $this->modelHiddenOverwrites($results, $isAdmin);

        return $results;
    }

    public function add($data)
    {
        $this->validate($data);
        if ($this->hasErrors()) {
            return;
        }

        $translations = $data[$this->_translationsRelationName] ?? [];
        if (is_string($translations)) {
            $translations = json_decode($translations, true); // Decode JSON string to array
        }
        unset($data[$this->_translationsRelationName]);

        $image = $data['image'] ?? null;
        unset($data['image']);

        if ($image) {
            $imageName = $this->imageName($image, $data);
            $image->move(storage_path('app/images/'.$this->_imageLocation), $imageName);
            $this->saveImageUrl($data, $imageName);
        }

        $item = null;

        DB::transaction(function () use ($data, $translations, &$item) {
            $item = $this->_model->create($data);

            if (!empty($translations)) {
                $this->updateTranslations($item, $translations);
            }
        });

        return $item;
    }

    public function update($id, $data)
    {
        $this->validate($data);
        if ($this->hasErrors()) {
            return;
        }

        $translations = $data[$this->_translationsRelationName] ?? [];
        unset($data[$this->_translationsRelationName]);

        $image = $data['image'] ?? null;
        unset($data['image']);

        if ($image) {
            $imageName = $this->imageName($image, $data);
            $image->move(storage_path('app/images/'.$this->_imageLocation), $imageName);
            $this->saveImageUrl($data, $imageName);
        }

        DB::transaction(function () use ($id, $data, $translations) {
            $item = $this->_model->find($id);

            if (!$item) {
                $this->_errors = new MessageBag(["item_not_found", "No item with id: $id was found"]);
                DB::rollBack();
                return;
            }

            $item->update($data);

            if (!empty($translations)) {
                $this->updateTranslations($item, $translations);
            }
        });

        return $this->_model->with('translations')->find($id);
    }

    public function delete($id)
    {
        $item = $this->_model->find($id);

        if (!$item) {
            $this->_errors = new MessageBag(["item_not_found", "No item with id: $id was found"]);
            return;
        }

        $deleted = $item->delete();

        if (!$deleted) {
            $this->_errors = new MessageBag(["failed_to_delete", "Failed to delete item with it: $id"]);
            return;
        }

        return $deleted;
    }

    protected function updateTranslations($item, $translations)
    {
        foreach ($translations as $translation) {
            if (!isset($translation['lang_code'])) {
                continue;
            }

            $existing_translation = $item->translations()->where('lang_code', $translation['lang_code'])->first();
            if ($existing_translation) {
                $existing_translation->update($translation);
            } else {
                $item->translations()->create($translation);
            }
        }
    }

    protected function getSelectFields($isAdmin)
    {
        if ($isAdmin)
        {
            return [...$this->_fields, ...$this->_fieldsAdminOnly];
        } else {
            return $this->_fields;
        }
    }
    protected function getLanguageSpecificFields($language, $isAdmin)
    {
        return [];
    }

    protected function adminConditions($query)
    {}

    protected function nonAdminConditions($query)
    {}

    protected function modelHiddenoverwrites($results, $isAdmin)
    {}

    protected function imageName($image, $data)
    {
        return Str::uuid() . '.' . $image->extension();
    }

    protected function saveImageUrl(&$data, $imageName)
    {
        $data['imageUrl'] = $this->_imageLocation.$imageName;
    }
}
