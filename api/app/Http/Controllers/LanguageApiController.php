<?php

namespace App\Http\Controllers;

use App\Services\LanguageService;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class LanguageApiController
{
    private $_service;

    public function __construct(LanguageService $service)
    {
        $this->_service = $service;
    }

    public function all(Request $request)
    {
        $search = $request->query('search', '');
        $sortBy = $request->query('sort_by', '');
        $sortDirection = $request->query('sort_direction', 'desc');

        return response()->json($this->_service->all(search: $search, sortBy: $sortBy, sortDirection: $sortDirection), Response::HTTP_OK);
    }

    public function add(Request $request)
    {
        $data = $request->all();
        $language = $this->_service->add($data);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["data" => $language];
    }

    public function update(Request $request, string $lang)
    {
        $data = $request->all();
        $language = $this->_service->update($lang, $data);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["data" => $language];
    }

    public function delete(string $lang)
    {
        // TODO: Get default lang from config
        if ($lang == "NL") {
            return ["warning", "You cannot delete the default language"];
        }

        $deletedLanguage = $this->_service->delete($lang);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["deleted" => $deletedLanguage];
    }
}
